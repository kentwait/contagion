package contagiongo

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/segmentio/ksuid"
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
	// HostTimer returns the current number of generations remaining
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
	// HostConnection returns the weight of a connection between two hosts
	// if it exists, returns 0 otherwise.
	HostConnection(a, b int) float64
	// HostNeighbors retrieves the directly connected hosts to the current
	// host based on the supplied adjacency matrix.
	HostNeighbors(id int) []Host

	// NewInstance creates a new instance from the stored configuration
	NewInstance() (Epidemic, error)

	// GenotypeNodeMap returns the set of all GenotypeNodes seen since the
	// start of the simulation.
	GenotypeNodeMap() map[ksuid.KSUID]GenotypeNode

	// GenotypeSet returns the set of all Genotypes seen since the
	// start of the simulation.
	GenotypeSet() GenotypeSet

	// StopSimulation check whether the simulation has satisfied at least one
	// of the conditions that will halt the simulation in the current interation.
	StopSimulation() bool

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

// Host returns the selected host in the simulation.
func (sim *SequenceNodeEpidemic) Host(id int) Host {
	return sim.hosts[id]
}

// HostStatus retrieves the current status of the selected host.
func (sim *SequenceNodeEpidemic) HostStatus(id int) int {
	sim.RLock()
	defer sim.RUnlock()
	return sim.statuses[id]
}

// SetHostStatus sets the current status of the selected host
// to a given status code.
func (sim *SequenceNodeEpidemic) SetHostStatus(id, status int) {
	sim.Lock()
	defer sim.Unlock()
	sim.statuses[id] = status
}

// HostTimer returns the current number of generations remaining
// before the host changes status.
func (sim *SequenceNodeEpidemic) HostTimer(id int) int {
	sim.RLock()
	defer sim.RUnlock()
	return sim.timers[id]
}

// SetHostTimer sets the number of generations for the host to
// remain in its current status.
func (sim *SequenceNodeEpidemic) SetHostTimer(id, interval int) {
	sim.Lock()
	defer sim.Unlock()
	sim.timers[id] = interval
}

// InfectableStatuses returns the list of statuses that infected
// hosts can transmit to.
func (sim *SequenceNodeEpidemic) InfectableStatuses() []int {
	return sim.infectableStatuses
}

// HostMap returns the hosts in the simulation in the form of a map.
// The key is the host's ID and the value is the pointer to the host.
func (sim *SequenceNodeEpidemic) HostMap() map[int]Host {
	return sim.hosts
}

func (sim *SequenceNodeEpidemic) HostConnection(a, b int) float64 {
	return sim.hostNetwork.Connection(a, b)
}

// HostNeighbors retrieves the directly connected hosts to the current
// host based on the supplied adjacency matrix.
func (sim *SequenceNodeEpidemic) HostNeighbors(id int) []Host {
	return sim.hostNeighborhoods[id]
}

// NewInstance creates a new instance from the stored configuration
func (sim *SequenceNodeEpidemic) NewInstance() (Epidemic, error) {
	return sim.config.NewSimulation()
}

// GenotypeNodeMap returns the set of all GenotypeNodes seen since the
// start of the simulation.
func (sim *SequenceNodeEpidemic) GenotypeNodeMap() map[ksuid.KSUID]GenotypeNode {
	return sim.tree.NodeMap()
}

// GenotypeSet returns the set of all Genotypes seen since the
// start of the simulation.
func (sim *SequenceNodeEpidemic) GenotypeSet() GenotypeSet {
	return sim.tree.Set()
}

// StopSimulation check whether the simulation has satisfied at least one
// of the conditions that will halt the simulation in the current
// interation. Returns true is the simulation should stop, false otherwise.
func (sim *SequenceNodeEpidemic) StopSimulation() bool {
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
