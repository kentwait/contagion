package contagiongo

import (
	"fmt"
	"strings"
	"sync"
	"time"

	rv "github.com/kentwait/randomvariate"
	"github.com/segmentio/ksuid"
)

// SISSimulation creates and runs an SIR epidemiological simulation.
// Within this simulation, hosts may or may not run
// independent genetic evolution simulations.
type SISSimulation struct {
	EpidemicSimulation
}

// NewSISSimulation creates a new SIS simulation.
func NewSISSimulation(config Config, logger DataLogger) (*SISSimulation, error) {
	sim := new(SISSimulation)
	var err error
	sim.EpidemicSimulation, err = NewSISimulation(config, logger)
	if err != nil {
		return nil, err
	}
	return sim, nil
}

// Run instantiates, runs, and records the a new simulation.
func (sim *SISSimulation) Run(i int) {
	sim.Initialize()
	sim.SetInstanceID(i)
	// Initial state
	sim.Update(0)

	sim.SetTime(0)
	var maxElapsed int64
	// First five generations generation initializes time
	for sim.Time() < 6 {
		sim.SetTime(sim.Time() + 1)
		fmt.Printf(" instance %04d\tgeneration %05d\n", i, sim.Time())
		start := time.Now()
		sim.Process(sim.Time())
		sim.Transmit(sim.Time())
		// Check conditions before update
		stop := !sim.CheckConditions()
		if stop {
			sim.SetStopped(true)
		}
		// Update after condition. If stop, will override logging setting
		// and log last generation
		sim.Update(sim.Time())
		// Check time elapsed
		if elapsed := time.Since(start).Nanoseconds(); elapsed > maxElapsed {
			maxElapsed = elapsed
		}
		// Feedback that simulation is stopping
		if stop {
			fmt.Printf(" [stop]       \tgeneration %05d\tstop condition triggered\n", sim.Time())
			break
		}
	}
	if !sim.Stopped() {
		fmt.Printf(" \t\texpected time: %fms per generation\n", float64(maxElapsed)/1e6)
	}
	for sim.Time() < sim.NumGenerations() && !sim.Stopped() {
		sim.SetTime(sim.Time() + 1)
		// Print only every ten steps is time is short
		if maxElapsed < 0.02e9 {
			if sim.Time()%100 == 0 {
				fmt.Printf(" instance %04d\tgeneration %05d\n", i, sim.Time())
			}
		} else if maxElapsed < 0.2e9 {
			if sim.Time()%10 == 0 {
				fmt.Printf(" instance %04d\tgeneration %05d\n", i, sim.Time())
			}
		} else {
			fmt.Printf(" instance %04d\tgeneration %05d\n", i, sim.Time())
		}
		sim.Process(sim.Time())
		sim.Transmit(sim.Time())
		// Check conditions before update
		stop := !sim.CheckConditions()
		if stop {
			sim.SetStopped(true)
		}
		// Update after condition. If stop, will override logging setting
		// and log last generation
		sim.Update(sim.Time())
		// Feedback that simulation is stopping
		if stop {
			fmt.Printf(" [stop]       \tgeneration %05d\tstop condition triggered\n", sim.Time())
			break
		}
	}
	fmt.Println(strings.Repeat("-", 80))
	sim.Finalize()
}

// Update looks at the timer or internal state to decide if
// the status of the host remains the same of will change.
// After the status updates, each host's status is recorded to file.
func (sim *SISSimulation) Update(t int) {
	// Update status first
	c := make(chan StatusPackage)
	d := make(chan GenotypeFreqPackage)
	var wg sync.WaitGroup
	// Read all hosts and process hosts concurrently
	// These succeeding steps connects the simulation's record of
	// each host's status and timer with the host's internal state.
	for hostID, host := range sim.HostMap() {
		// Simulation-level record of status and timer of particular host
		timer := sim.HostTimer(hostID)
		pack := StatusPackage{
			instanceID: sim.InstanceID(),
			genID:      t,
			hostID:     hostID,
			status:     sim.HostStatus(hostID), // invalid status
		}
		wg.Add(1)
		go func(i, t int, host Host, timer int, pack StatusPackage, c chan<- StatusPackage, d chan<- GenotypeFreqPackage, wg *sync.WaitGroup) {
			defer wg.Done()
			// Add cases depending on the compartmental model being used
			// In this case, SI only uses susceptible and infected statuses
			switch pack.status {
			case SusceptibleStatusCode:
				// Use timer or number of pathogens
				if timer == 0 || host.PathogenPopSize() > 0 {
					// Set new host status
					newStatus := InfectedStatusCode
					newDuration := host.GetIntrahostModel().StatusDuration(newStatus)
					sim.SetHostStatus(host.ID(), newStatus)
					// Makes the duration poisson
					sim.SetHostTimer(host.ID(), rv.Poisson(float64(newDuration)))
					// sim.SetHostTimer(host.ID(), newDuration)
					// Update status in pack and send
					pack.status = newStatus
				}
			case InfectedStatusCode:
				if timer == 0 || host.PathogenPopSize() == 0 {
					// Set new host status
					newStatus := SusceptibleStatusCode
					newDuration := -1 // Host goes back to being susceptible
					sim.SetHostStatus(host.ID(), newStatus)
					sim.SetHostTimer(host.ID(), newDuration)
					// Update status in pack and send
					pack.status = newStatus
					host.RemoveAllPathogens()
				}
			}
			// Send pack after all changes
			c <- pack
			// Record pathogen frequencies
			counts := make(map[ksuid.KSUID]int)
			for _, p := range host.Pathogens() {
				counts[p.GenotypeUID()]++
			}
			for uid, freq := range counts {
				d <- GenotypeFreqPackage{
					instanceID: i,
					genID:      t,
					hostID:     host.ID(),
					genotypeID: uid,
					freq:       freq,
				}
			}
		}(sim.InstanceID(), t, host, timer, pack, c, d, &wg)
	}
	go func() {
		wg.Wait()
		close(c)
		close(d)
	}()
	// Write status  and genotype frequencies using DataLogger
	var wg2 sync.WaitGroup
	wg2.Add(2)
	go func() {
		if sim.Time() == 0 || sim.Time()%sim.LogFrequency() == 0 || sim.Stopped() {
			sim.WriteStatus(c)
		} else {
			for range c {
			}
		}
		wg2.Done()
	}()
	go func() {
		if sim.Time() == 0 || sim.Time()%sim.LogFrequency() == 0 || sim.Stopped() {
			sim.WriteGenotypeFreq(d)
		} else {
			for range d {
			}
		}
		wg2.Done()
	}()
	wg2.Wait()
}

// Process runs the internal evolution simulation in each host.
// During intrahost evolution, if new mutations appear, the new sequence
// and ancestry is recorded to file.
func (sim *SISSimulation) Process(t int) {
	c := make(chan MutationPackage)
	var wg sync.WaitGroup
	// Read all hosts and process based on the current status of the host
	for hostID, host := range sim.HostMap() {
		// Run the intrahost process asynchronously depending on the
		// current status of the host.
		wg.Add(1)
		switch sim.HostStatus(hostID) {
		case SusceptibleStatusCode:
			go sim.SusceptibleProcess(sim.InstanceID(), t, host, &wg)
		case InfectedStatusCode:
			go sim.InfectedProcess(sim.InstanceID(), t, host, c, &wg)
		}
		// Decrement host timer.
		// If status depends on timer, then timer will go positive integer
		// to 0.
		// If status is not dependent on timer, then will just decrement
		// from -1 to more negative values
		sim.SetHostTimer(hostID, sim.HostTimer(hostID)-1)
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	// Write mutations to DataLogger
	sim.WriteMutations(c)
}
