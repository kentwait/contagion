package contagiongo

import (
	"math"
	"strings"
	"sync"
)

// evoEpiSimulation is a type simulation that uses a SequenceNode
// to represent pathogens.
type evoEpiSimulation struct {
	hosts             map[int]Host
	statuses          map[int]int
	timers            map[int]int
	intrahostModels   map[int]IntrahostModel
	fitnessModels     map[int]FitnessModel
	hostNeighborhoods map[int][]Host
	hostNetwork       map[int]map[int]float64
	tree              GenotypeTree
}

func (sim *evoEpiSimulation) Host(id int) Host {
	return sim.hosts[id]
}

func (sim *evoEpiSimulation) HostStatus(id int) int {
	return sim.statuses[id]
}

func (sim *evoEpiSimulation) SetHostStatus(id, status int) {
	sim.statuses[id] = status
}

func (sim *evoEpiSimulation) HostTimer(id int) int {
	return sim.timers[id]
}

func (sim *evoEpiSimulation) SetHostTimer(id, interval int) {
	sim.timers[id] = interval
}

func (sim *evoEpiSimulation) HostMap() map[int]Host {
	return sim.hosts
}

func (sim *evoEpiSimulation) HostNeighbors(id int) []Host {
	return sim.hostNeighborhoods[id]
}

// The following methods are used as goroutines that performs tasks within
// each host when the host is in a particular state. Tasks performed are
// assumed to affect only data encapsulated within the host.

// SusceptibleProcess executes within-host processes that occurs when a host
// is in the susceptible state.
func (sim *evoEpiSimulation) SusceptibleProcess(host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	// Decrement to -1 if pathogens exist
	if host.PathogenPopSize() > 0 {
		host.DecrementTimer()
	}
}

// ExposedProcess executes within-host processes that occurs when a host
// is in the exposed state. By default, it is same as InfectedProcess.
func (sim *evoEpiSimulation) ExposedProcess(host Host, wg *sync.WaitGroup) {
	// timer decrement is done within the InfectedProcess function
	// Done() signal also executed within the InfectedProcess function
	sim.InfectedProcess(host, wg)
	// TODO: Threshold to be considered infective instead of exposed
}

// InfectedProcess executes within-host processes that occurs when a host
// is in the infected state.
func (sim *evoEpiSimulation) InfectedProcess(host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	pathogens := host.Pathogens()
	var replicatedC <-chan GenotypeNode
	switch strings.ToLower(host.(*sequenceHost).IntrahostModel.ReplicationMethod()) {
	case "relative":
		// Get log fitness values for each pathogen
		logFitnesses := make([]float64, len(pathogens))
		var minLogFitness float64
		// Compute log total fitness and get max value
		for i, pathogen := range pathogens {
			logFitnesses[i] = pathogen.Fitness(host.(*sequenceHost).FitnessModel)
			if minLogFitness > logFitnesses[i] {
				minLogFitness = logFitnesses[i]
			}
		}
		// exp-normalize algorithm
		// Get normalizing constant by summing all elements
		var c float64
		for _, logF := range logFitnesses {
			c += math.Exp(logF - minLogFitness)
		}
		// Normalize such that fitnesses sum to 1.0
		normedFitnesses := make([]float64, len(pathogens))
		for i, logF := range logFitnesses {
			normedFitnesses[i] = math.Exp(logF-minLogFitness) / c
		}
		// get current and next pop size based on popsize function
		currentPopSize := host.PathogenPopSize()
		nextPopSize := host.(*sequenceHost).NextPathogenPopSize(currentPopSize)
		// Execute
		replicatedC = MultinomialReplication(pathogens, normedFitnesses, nextPopSize)
	case "absolute":
		// Get decimal fitness values. Each value is the expected number of
		// offspring
		replicativeFitnesses := make([]float64, len(pathogens))
		for i, pathogen := range pathogens {
			replicativeFitnesses[i] = pathogen.Fitness(host.(*sequenceHost).FitnessModel)
		}
		// Execute
		replicatedC = IntrinsicRateReplication(pathogens, replicativeFitnesses, nil)
	}
	// Mutate replicated pathogens
	mutatedC := MutateSequence(replicatedC, sim.tree, host.(*sequenceHost).IntrahostModel)
	// Clear current set of pathogens and get new set from the channel
	host.RemoveAllPathogens()
	for pathogen := range mutatedC {
		host.AddPathogen(pathogen)
	}
	host.DecrementTimer()
}

// InfectiveProcess executes within-host processes that occurs when a host
// is in the infective state. By default, it is same as InfectedProcess.
func (sim *evoEpiSimulation) InfectiveProcess(host Host, wg *sync.WaitGroup) {
	// timer decrement is done within the InfectedProcess function
	// Done() signal also executed within the InfectedProcess function
	sim.InfectedProcess(host, wg)
}

// RemovedProcess executes within-host processes that occurs when a host
// is in the removed state that is perpetually uninfectable.
func (sim *evoEpiSimulation) RemovedProcess(host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	host.RemoveAllPathogens()
}

// RecoveredProcess executes within-host processes that occurs when a host
// is in the recovered state that is perpetually uninfectable.
// This state is identically to Removed but is used to distinguish from
// a dead state.
func (sim *evoEpiSimulation) RecoveredProcess(host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	host.RemoveAllPathogens()
}

// DeadProcess executes within-host processes that occurs when a host
// is in the dead state state that is perpetually uninfectable.
// This state is identically to Removed but is used to distinguish from
// a recovered, but perpetually immune state.
func (sim *evoEpiSimulation) DeadProcess(host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	host.RemoveAllPathogens()
}

// VaccinatedProcess executes within-host processes that occurs when a host
// is in a globally immune state with the chance to become
// globally susceptible again.
func (sim *evoEpiSimulation) VaccinatedProcess(host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	host.RemoveAllPathogens()
	host.DecrementTimer()
}
