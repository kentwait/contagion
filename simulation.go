package contagiongo

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/segmentio/ksuid"
)

// The following are status codes for different preset compartments that
// describe the current epidemiological status of a host in the simulation.
const (
	SusceptibleStatusCode = 1
	ExposedStatusCode     = 2
	InfectedStatusCode    = 3
	InfectiveStatusCode   = 4
	RemovedStatusCode     = 5
	RecoveredStatusCode   = 6
	DeadStatusCode        = 7
	VaccinatedStatusCode  = 8
)

// Epidemic encapsulates the set of hosts, its connections,
// the pathogen tree lineage and the host types used to create
// a simulated epidemic.
type Epidemic interface {
	// Host returns the selected host in the simulation.
	Host(id int) Host

	// HostStatus retrieves the current status of the selected host.
	HostStatus(id int) int
	// SetHostStatus sets the current status of the selected host
	// to a given status code.
	SetHostStatus(id, status int)
	// HostTime returns the current number of generations remaining
	// before the host changes status.
	HostTimer(id int) int
	// SetHostTimer sets the number of generations for the host to
	// remain in its current status.
	SetHostTimer(id, interval int)
	// InfectableStatuses returns the list of statuses that infected
	// hosts can transmit to.
	InfectableStatuses() []int

	// HostMap returns the hosts in the simulation in the form of a map.
	// The key is the host's ID and the value is the pointer to the host.
	HostMap() map[int]Host
	// HostNeighbors retrieves the directly connected hosts to the current
	// host based on the supplied adjacency matrix.
	HostNeighbors(id int) []Host

	// NewInstance creates a NewInstancete a new instance from
	// the stored configuration
	NewInstance() (Epidemic, error)

	GenotypeNodeMap() map[ksuid.KSUID]GenotypeNode
	GenotypeSet() GenotypeSet
	CheckConditions() bool

	// The following methods perform intrahost processes associated with
	// the status. For every generation, one of the following is called for
	// each host.

	// SusceptibleProcess performs intrahost processes while the host is in
	// the susceptible status.
	SusceptibleProcess(i, t int, host Host, wg *sync.WaitGroup)
	// ExposedProcess performs intrahost processes while the host is in
	// the exposed status.
	ExposedProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup)
	// InfectedProcess performs intrahost processes while the host is in
	// the infected status.
	InfectedProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup)
	// InfectiveProcess performs intrahost processes while the host is in
	// the infective status.
	InfectiveProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup)
	// RemovedProcess performs intrahost processes while the host is in
	// the removed status.
	RemovedProcess(i, t int, host Host, wg *sync.WaitGroup)
	// RecoveredProcess performs intrahost processes while the host is in
	// the recovered status.
	RecoveredProcess(i, t int, host Host, wg *sync.WaitGroup)
	// DeadProcess performs intrahost processes while the host is in
	// the dead status.
	DeadProcess(i, t int, host Host, wg *sync.WaitGroup)
	// DeadProcess performs intrahost processes while the host is in
	// the dead status.
	VaccinatedProcess(i, t int, host Host, wg *sync.WaitGroup)
}

// SequenceNodeEpidemic is a type of Epidemic that uses a SequenceNode
// to represent pathogens.
type SequenceNodeEpidemic struct {
	sync.RWMutex
	hosts              map[int]Host
	statuses           map[int]int
	timers             map[int]int
	intrahostModels    map[int]IntrahostModel
	fitnessModels      map[int]FitnessModel
	transModels        map[int]TransmissionModel
	hostNeighborhoods  map[int][]Host
	hostNetwork        HostNetwork
	infectableStatuses []int
	tree               GenotypeTree
	config             Config

	stopConditions []StopCondition
}

func (sim *SequenceNodeEpidemic) Host(id int) Host {
	return sim.hosts[id]
}

func (sim *SequenceNodeEpidemic) HostStatus(id int) int {
	sim.RLock()
	defer sim.RUnlock()
	return sim.statuses[id]
}

func (sim *SequenceNodeEpidemic) SetHostStatus(id, status int) {
	sim.Lock()
	defer sim.Unlock()
	sim.statuses[id] = status
}

func (sim *SequenceNodeEpidemic) HostTimer(id int) int {
	sim.RLock()
	defer sim.RUnlock()
	return sim.timers[id]
}

func (sim *SequenceNodeEpidemic) SetHostTimer(id, interval int) {
	sim.Lock()
	defer sim.Unlock()
	sim.timers[id] = interval
}

func (sim *SequenceNodeEpidemic) InfectableStatuses() []int {
	return sim.infectableStatuses
}

func (sim *SequenceNodeEpidemic) HostMap() map[int]Host {
	return sim.hosts
}

func (sim *SequenceNodeEpidemic) HostNeighbors(id int) []Host {
	return sim.hostNeighborhoods[id]
}

func (sim *SequenceNodeEpidemic) GenotypeNodeMap() map[ksuid.KSUID]GenotypeNode {
	return sim.tree.NodeMap()
}

func (sim *SequenceNodeEpidemic) GenotypeSet() GenotypeSet {
	return sim.tree.Set()
}

func (sim *SequenceNodeEpidemic) NewInstance() (Epidemic, error) {
	return sim.config.NewSimulation()
}

func (sim *SequenceNodeEpidemic) CheckConditions() bool {
	continueSim := true
	for _, cond := range sim.stopConditions {
		continueSim = cond.Check(sim)
		if !continueSim {
			fmt.Printf(" \t\t- %s -\n", cond.Reason())
			break
		}
	}
	return continueSim
}

// The following methods are used as goroutines that performs tasks within
// each host when the host is in a particular state. Tasks performed are
// assumed to affect only data encapsulated within the host.

// SusceptibleProcess executes within-host processes that occurs when a host
// is in the susceptible state.
func (sim *SequenceNodeEpidemic) SusceptibleProcess(i, t int, host Host, wg *sync.WaitGroup) {
	defer wg.Done()
}

// ExposedProcess executes within-host processes that occurs when a host
// is in the exposed state. By default, it is same as InfectedProcess.
func (sim *SequenceNodeEpidemic) ExposedProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup) {
	// timer decrement is done within the InfectedProcess function
	// Done() signal also executed within the InfectedProcess function
	sim.InfectedProcess(i, t, host, c, wg)
	// TODO: Threshold to be considered infective instead of exposed
}

// InfectedProcess executes within-host processes that occurs when a host
// is in the infected state.
func (sim *SequenceNodeEpidemic) InfectedProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup) {
	defer wg.Done()
	pathogens := host.Pathogens()
	if host.PathogenPopSize() == 0 {
		return
	}
	var replicatedC <-chan GenotypeNode
	switch strings.ToLower(host.GetIntrahostModel().ReplicationMethod()) {
	case "relative":
		// Get log fitness values for each pathogen
		logFitnesses := make([]float64, len(pathogens))

		var maxLogFitness float64
		// Compute log total fitness and get max value
		for i, pathogen := range pathogens {
			logFitnesses[i] = pathogen.Fitness(host.GetFitnessModel())
			if maxLogFitness < logFitnesses[i] {
				maxLogFitness = logFitnesses[i]
			}
		}
		// exp-normalize algorithm
		// get normalizing constant by summing all elements
		var c float64
		for _, logF := range logFitnesses {
			c += math.Exp(logF - maxLogFitness)
		}
		// normalize
		normedDecFitnesses := make([]float64, len(pathogens))
		for i, logF := range logFitnesses {
			normedDecFitnesses[i] = math.Exp(logF-maxLogFitness) / c
		}
		// get current and next pop size based on popsize function
		currentPopSize := host.PathogenPopSize()
		// TODO: Expose this in interface
		nextPopSize := host.GetIntrahostModel().NextPathogenPopSize(currentPopSize)
		// Execute
		replicatedC = MultinomialReplication(pathogens, normedDecFitnesses, nextPopSize)
	case "absolute":
		// Get decimal fitness values. Each value is the expected number of
		// offspring
		replicativeFitnesses := make([]float64, len(pathogens))
		for i, pathogen := range pathogens {
			replicativeFitnesses[i] = pathogen.Fitness(host.GetFitnessModel())
		}
		// Execute
		replicatedC = IntrinsicRateReplication(pathogens, replicativeFitnesses, nil)
	}
	// Mutate replicated pathogens
	mutatedC, newMutantsC := MutateSequence(replicatedC, sim.tree, host.GetIntrahostModel())
	// Clear current set of pathogens and get new set from the channel
	host.RemoveAllPathogens()
	var wg2 sync.WaitGroup
	wg2.Add(2)
	go func() {
		for node := range mutatedC {
			host.AddPathogens(node)
		}
		wg2.Done()
	}()
	go func() {
		for node := range newMutantsC {
			for _, parent := range node.Parents() {
				c <- MutationPackage{
					instanceID:   i,
					genID:        t,
					hostID:       host.ID(),
					nodeID:       node.UID(),
					parentNodeID: parent.UID(),
				}
			}
		}
		wg2.Done()
	}()
	wg2.Wait()
}

// InfectiveProcess executes within-host processes that occurs when a host
// is in the infective state. By default, it is same as InfectedProcess.
func (sim *SequenceNodeEpidemic) InfectiveProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup) {
	// timer decrement is done within the InfectedProcess function
	// Done() signal also executed within the InfectedProcess function
	sim.InfectedProcess(i, t, host, c, wg)
}

// RemovedProcess executes within-host processes that occurs when a host
// is in the removed state that is perpetually uninfectable.
func (sim *SequenceNodeEpidemic) RemovedProcess(i, t int, host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	host.RemoveAllPathogens()
}

// RecoveredProcess executes within-host processes that occurs when a host
// is in the recovered state that is perpetually uninfectable.
// This state is identically to Removed but is used to distinguish from
// a dead state.
func (sim *SequenceNodeEpidemic) RecoveredProcess(i, t int, host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	host.RemoveAllPathogens()
}

// DeadProcess executes within-host processes that occurs when a host
// is in the dead state state that is perpetually uninfectable.
// This state is identically to Removed but is used to distinguish from
// a recovered, but perpetually immune state.
func (sim *SequenceNodeEpidemic) DeadProcess(i, t int, host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	host.RemoveAllPathogens()
}

// VaccinatedProcess executes within-host processes that occurs when a host
// is in a globally immune state with the chance to become
// globally susceptible again.
func (sim *SequenceNodeEpidemic) VaccinatedProcess(i, t int, host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	host.RemoveAllPathogens()
}

// EpidemicSimulation is a simulation environment that simulates
// the spread of disease between hosts in a connected host network.
type EpidemicSimulation interface {
	Epidemic
	DataLogger

	// Run runs the whole simulation
	Initialize(params ...interface{})
	Run(i int)
	Update(t int)
	Process(t int)
	Transmit(t int)
	Finalize()

	// Metadata
	SetInstanceID(i int)
	InstanceID() int
	SetTime(t int)
	Time() int
	SetGenerations(n int)
	NumGenerations() int
	LogTransmission() bool
	LogFrequency() int
	SetStopped(b bool)
	Stopped() bool
}

// Infection encapsulates a single host and the pathogen tree lineage
// to trace the evolution of pathogens within one host. This is useful
// to strudy intrahost evolutionary dynamics especially in chronic diseases.
type Infection interface {
	// Host returns the selected host in the simulation.
	Host() Host

	// HostStatus retrieves the current status of the selected host.
	HostStatus() int
	// SetHostStatus sets the current status of the selected host
	// to a given status code.
	SetHostStatus(status int)
	// HostStatusDuration returns the number of generations a host
	// remains in a given status.
	HostStatusDuration(status int) int
	// HostTime returns the current number of generations remaining
	// before the host changes status.
	HostTimer() int
	// SetHostTimer sets the number of generations for the host to
	// remain in its current status.
	SetHostTimer(interval int)

	// The following methods perform intrahost processes associated with
	// the status. For every generation, one of the following is called for
	// each host.

	// SusceptibleProcess performs intrahost processes while the host is in
	// the susceptible status.
	SusceptibleProcess(i, t int, host Host, wg *sync.WaitGroup)
	// ExposedProcess performs intrahost processes while the host is in
	// the exposed status.
	ExposedProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup)
	// InfectedProcess performs intrahost processes while the host is in
	// the infected status.
	InfectedProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup)
	// InfectiveProcess performs intrahost processes while the host is in
	// the infective status.
	InfectiveProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup)
	// RemovedProcess performs intrahost processes while the host is in
	// the removed status.
	RemovedProcess(i, t int, host Host, wg *sync.WaitGroup)
	// RecoveredProcess performs intrahost processes while the host is in
	// the recovered status.
	RecoveredProcess(i, t int, host Host, wg *sync.WaitGroup)
	// DeadProcess performs intrahost processes while the host is in
	// the dead status.
	DeadProcess(i, t int, host Host, wg *sync.WaitGroup)
	// DeadProcess performs intrahost processes while the host is in
	// the dead status.
	VaccinatedProcess(i, t int, host Host, wg *sync.WaitGroup)
}

// InfectionSimulation is a simulation environment that simulates
// the infection within a single hosts or in an evironment where a network
// configuration is not necessary.
type InfectionSimulation interface {
	Infection
	// Run runs the whole simulation
	Init(params ...interface{})
	Run(i int)
	Update(t int)
	Process(t int)
	Transmit(t int)
}

// TransmissionEvent is a struct for sending and receiving
// transmission event information.
type TransmissionEvent struct {
	destination Host
	pathogen    GenotypeNode
}

// ExchangeEvent is a struct for sending and receiving
// exchange event information.
type ExchangeEvent struct {
	source      Host
	destination Host
	// pathogenIndex int
	pathogen GenotypeNode
}
