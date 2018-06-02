package contagiongo

import (
	"fmt"
	"log"
	"sync"

	"github.com/segmentio/ksuid"
)

// import "sync"

// SISimulation creates and runs an SI epidemiological simulation.
// Within this simulation, hosts may or may not run
// independent genetic evolution simulations.
type SISimulation struct {
	Epidemic
	DataLogger

	instanceID     int
	numGenerations int
	logFreq        int
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
	return sim, nil
}

// Run instantiates, runs, and records the a new simulation.
func (sim *SISimulation) Run(i int) {
	sim.Init()
	sim.instanceID = i
	sim.Update(0)
	t := 1
	for t <= sim.numGenerations {
		fmt.Printf("instance %04d\tgeneration %05d\t", i, t)
		sim.Process(t)
		sim.Transmit(t)
		sim.Update(t)
		t++
		fmt.Printf("Done.\n")
	}
	sim.Update(t)
	sim.Finalize()
}

// Init initializes the simulation and accepts 0 or more parameters.
// For example, creating datbases etc.
func (sim *SISimulation) Init(params ...interface{}) {
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
			go sim.SusceptibleProcess(sim.instanceID, t, host, &wg)
		case InfectedStatusCode:
			go sim.InfectedProcess(sim.instanceID, t, host, c, &wg)
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
