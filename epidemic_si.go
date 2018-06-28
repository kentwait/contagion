package contagiongo

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/ksuid"
)

// SISimulation creates and runs an SI epidemiological simulation.
// Within this simulation, hosts may or may not run
// independent genetic evolution simulations.
type SISimulation struct {
	Epidemic
	DataLogger

	instanceID      int
	numGenerations  int
	t               int
	logFreq         int
	stopped         bool
	logTransmission bool
}

// NewSISimulation creates a new SI simulation.
func NewSISimulation(config Config, logger DataLogger) (*SISimulation, error) {
	epidemic, err := config.NewSimulation()
	if err != nil {
		return nil, err
	}
	sim := new(SISimulation)
	sim.Epidemic = epidemic
	sim.DataLogger = logger
	sim.numGenerations = config.NumGenerations()
	sim.logFreq = config.LogFreq()
	sim.logTransmission = config.LogTransmission()
	return sim, nil
}

// SetInstanceID sets the instance ID of the current realized simulation.
func (sim *SISimulation) SetInstanceID(i int) {
	sim.instanceID = i
}

// InstanceID returns the ID of the current realized simulation.
func (sim *SISimulation) InstanceID() int {
	return sim.instanceID
}

// SetTime sets the current internal time of the simulation.
// The simulation's internal time is based on the number of iterations
// that has taken place. This is equivalent to the number of pathogen
// generations.
func (sim *SISimulation) SetTime(t int) {
	sim.t = t
}

// Time returns the current internal time of the simulation.
// The simulation's internal time should be the number of iterations
// that has taken place. This is equivalent to the number of pathogen
// generations.
func (sim *SISimulation) Time() int {
	return sim.t
}

// SetGenerations sets the total number of pathogen generations
// the simulation will simulate. This is equivalent to the total
// number of iterations of the simulation.
func (sim *SISimulation) SetGenerations(n int) {
	sim.numGenerations = n
}

// NumGenerations returns the total number of pathogen generations
// the simulation will simulate. This is equivalent to the total
// number of iterations of the simulation.
func (sim *SISimulation) NumGenerations() int {
	return sim.numGenerations
}

// LogTransmission returns true is transmission events
// are saved to disk. If false, transmssion events occur but
// are not recorded.
func (sim *SISimulation) LogTransmission() bool {
	return sim.logTransmission
}

// LogFrequency returns the interval in number of pathogen
// generation between data recordings.
func (sim *SISimulation) LogFrequency() int {
	return sim.logFreq
}

// SetStopped sets the internal status of the current simulation.
// If set to true, this indicates that the simulation has stopped.
// If set to false, the current simulation has not yet stopped. By
// default, the value of internal status is false.
func (sim *SISimulation) SetStopped(b bool) {
	sim.stopped = b
}

// Stopped returns true if the current simulation has stopped.
// If it returns false, the current simulation has not yet stopped
func (sim *SISimulation) Stopped() bool {
	return sim.stopped
}

// Run instantiates, runs, and records the a new simulation.
func (sim *SISimulation) Run(i int) {
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
		stop := sim.StopSimulation()
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
		stop := sim.StopSimulation()
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

// Initialize initializes the simulation and accepts 0 or more parameters.
// For example, creating datbases etc.
func (sim *SISimulation) Initialize(params ...interface{}) {
	err := sim.DataLogger.Init()
	if err != nil {
		log.Fatal(err)
	}
}

// Update looks at the timer or internal state to decide if
// the status of the host remains the same of will change.
// After the status updates, each host's status is recorded to file.
func (sim *SISimulation) Update(t int) {
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
			status:     sim.HostStatus(hostID), // current host status before checking
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
					sim.SetHostTimer(host.ID(), newDuration)
					// Update status in pack and send
					pack.status = newStatus
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
func (sim *SISimulation) Process(t int) {
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

// Transmit facilitates the sampling and migration process of pathogens
// between hosts.
func (sim *SISimulation) Transmit(t int) {
	c := make(chan TransmissionEvent)
	d := make(chan TransmissionPackage)
	var wg sync.WaitGroup
	// Get hosts that are infected and determine pathogen pop size
	// Only hosts with an infected status can transmit
	var infectedHosts []Host
	var pathogenPopSizes []int
	for hostID, host := range sim.HostMap() {
		if sim.HostStatus(hostID) == InfectedStatusCode {
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
		numMigrants := host.GetTransmissionModel().TransmissionSize()
		transmissionProb := host.GetTransmissionModel().TransmissionProb()
		for _, neighbor := range sim.HostNeighbors(hostID) {
			status := sim.HostStatus(neighbor.ID())
			// Overrides default transmission prob set in the config file
			if t := sim.HostConnection(hostID, neighbor.ID()); t > 0 {
				transmissionProb = t
			}
			for _, infectableStatus := range sim.InfectableStatuses() {
				if status == infectableStatus {
					wg.Add(1)
					go TransmitPathogens(sim.InstanceID(), t, host, neighbor, numMigrants, transmissionProb, count, c, d, &wg)
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
		if sim.logTransmission {
			sim.WriteTransmission(d)
		} else {
			for range d {
			}
		}
		wg2.Done()
	}()
	wg2.Wait()
}

// Finalize performs processes to finish and close the simulation.
func (sim *SISimulation) Finalize() {
	// Record genotype tree
	var wg sync.WaitGroup
	c := make(chan GenotypeNode)
	d := make(chan Genotype)

	wg.Add(2)
	go func() {
		for _, node := range sim.GenotypeNodeMap() {
			c <- node
		}
		close(c)
		wg.Done()
	}()
	go func() {
		for _, genotype := range sim.GenotypeSet().Map() {
			d <- genotype
		}
		close(d)
		wg.Done()
	}()
	var wg2 sync.WaitGroup
	wg2.Add(2)
	go func() {
		sim.WriteGenotypeNodes(c)
		wg2.Done()
	}()
	go func() {
		sim.WriteGenotypes(d)
		wg2.Done()
	}()
	wg2.Wait()

	// Clear memory by deleting pathogens in host
	for _, host := range sim.HostMap() {
		host.RemoveAllPathogens()
	}
}
