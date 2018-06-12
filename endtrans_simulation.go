package contagiongo

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/ksuid"
)

// import "sync"

// EndTransSimulation creates and runs a modified version of the
// SIR epidemiological simulation.
// In the endtrans model, transmissions are allowed to occur only
// at the last generation before pathogens are cleared from
// the host.
// To make transmission completely deterministic, set the transmission
// probability to 1.0, use the constant mode and set it to a constant
// size. This means that all paths connected to an infectable host gets
// infected at the end of the infection cycle of the current host.
// Endtrans assumes that the InfectedDuration is not zero.
type EndTransSimulation struct {
	EpidemicSimulation
}

// NewEndTransSimulation creates a new SI simulation.
func NewEndTransSimulation(config Config, logger DataLogger) (*EndTransSimulation, error) {
	sim := new(EndTransSimulation)
	var err error
	sim.EpidemicSimulation, err = NewSIRSimulation(config, logger)
	if err != nil {
		return nil, err
	}
	return sim, nil
}

// Run instantiates, runs, and records the a new simulation.
func (sim *EndTransSimulation) Run(i int) {
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

func (sim *EndTransSimulation) Update(t int) {
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
					sim.SetHostTimer(host.ID(), newDuration)
					// sim.SetHostTimer(host.ID(), newDuration)
					// Update status in pack and send
					pack.status = newStatus
				}
			case InfectedStatusCode:
				if timer == 0 || host.PathogenPopSize() == 0 {
					// Set new host status
					newStatus := RemovedStatusCode
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

// Transmit facilitates the sampling and migration process of pathogens
// between hosts.
func (sim *EndTransSimulation) Transmit(t int) {
	c := make(chan TransmissionEvent)
	d := make(chan TransmissionPackage)
	var wg sync.WaitGroup
	// Get hosts that are infected and at the end of their infection cycle,
	// and then determine pathogen pop size
	// Only hosts with an infected status can transmit
	var infectedHosts []Host
	var pathogenPopSizes []int
	for hostID, host := range sim.HostMap() {
		if sim.HostStatus(hostID) == InfectedStatusCode && sim.HostTimer(hostID) == 0 {
			infectedHosts = append(infectedHosts, host)
			pathogenPopSizes = append(pathogenPopSizes, host.PathogenPopSize())
		}
	}
	// Iterate using pre-assembled list of infected hosts
	for i, host := range infectedHosts {
		// Iterate over host's neighbors and create a new goroutine
		// that determines whether pathogens transmit or not
		hostID := host.ID()
		count := pathogenPopSizes[i]
		for _, neighbor := range sim.HostNeighbors(hostID) {
			status := sim.HostStatus(neighbor.ID())
			for _, infectableStatus := range sim.InfectableStatuses() {
				if status == infectableStatus {
					wg.Add(1)
					go TransmitPathogens(sim.InstanceID(), t, host, neighbor, count, c, d, &wg)
				}
			}

		}
	}
	go func() {
		wg.Wait()
		close(c)
		close(d)
	}()
	// Add the new pathogen to the destination host
	// and record
	var wg2 sync.WaitGroup
	wg2.Add(2)
	go func() {
		for t := range c {
			t.destination.AddPathogens(t.pathogen)
		}
		wg2.Done()
	}()
	go func() {
		sim.WriteTransmission(d)
		wg2.Done()
	}()
	wg2.Wait()
}
