package contagiongo

import (
	"fmt"
	"log"
	"strings"
	"sync"

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
	SISimulation

	instanceID     int
	numGenerations int
	logFreq        int
}

// NewEndTransSimulation creates a new SI simulation.
func NewEndTransSimulation(config Config, logger DataLogger) (*EndTransSimulation, error) {
	epidemic, err := config.NewSimulation()
	if err != nil {
		return nil, err
	}
	sim := new(EndTransSimulation)
	sim.Epidemic = epidemic
	sim.DataLogger = logger
	sim.numGenerations = config.NumGenerations()
	sim.logFreq = config.LogFreq()
	return sim, nil
}

// Run instantiates, runs, and records the a new simulation.
func (sim *EndTransSimulation) Run(i int) {
	sim.Init()
	sim.instanceID = i
	sim.Update(0)
	t := 0
	for t < sim.numGenerations {
		t++
		fmt.Printf("instance %04d\tgeneration %05d\n", i, t)
		sim.Process(t)
		sim.Transmit(t)
		sim.Update(t)
	}
	fmt.Println(strings.Repeat("-", 80))
	sim.Finalize()
}

// Init initializes the simulation and accepts 0 or more parameters.
// For example, creating datbases etc.
func (sim *EndTransSimulation) Init(params ...interface{}) {
	err := sim.DataLogger.Init()
	if err != nil {
		log.Fatal(err)
	}
}

// Update looks at the timer or internal state to decide if
// the status of the host remains the same of will change.
// After the status updates, each host's status is recorded to file.
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
			instanceID: sim.instanceID,
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
			case InfectedStatusCode:
				if timer == 0 || host.PathogenPopSize() > 0 {
					// Set new host status
					newStatus := RemovedStatusCode
					newDuration := -1 // Host is permanently removed from the population of infectables
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
		}(sim.instanceID, t, host, timer, pack, c, d, &wg)
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
		sim.WriteStatus(c)
		wg2.Done()
	}()
	go func() {
		sim.WriteGenotypeFreq(d)
		wg2.Done()
	}()
	wg2.Wait()
}

// Process runs the internal evolution simulation in each host.
// During intrahost evolution, if new mutations appear, the new sequence
// and ancestry is recorded to file.
func (sim *EndTransSimulation) Process(t int) {
	c := make(chan MutationPackage)
	var wg sync.WaitGroup
	// Read all hosts and process based on the current status of the host
	for hostID, host := range sim.HostMap() {
		// Run the intrahost process asynchronously depending on the
		// current status of the host.
		wg.Add(1)
		switch sim.HostStatus(hostID) {
		case SusceptibleStatusCode:
			go sim.SusceptibleProcess(sim.instanceID, t, host, &wg)
		case InfectedStatusCode:
			go sim.InfectedProcess(sim.instanceID, t, host, c, &wg)
		case RemovedStatusCode:
			go sim.RemovedProcess(sim.instanceID, t, host, &wg)
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
					go TransmitPathogens(sim.instanceID, t, host, neighbor, count, c, d, &wg)
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
			t.destination.AddPathogen(t.pathogen)
		}
		wg2.Done()
	}()
	go func() {
		sim.WriteTransmission(d)
		wg2.Done()
	}()
	wg2.Wait()
}
