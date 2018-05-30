package contagiongo

import (
	"math"
	"strings"
	"sync"
)

// SingleHostSimulation simulates pathogens infection within a single host.
type SingleHostSimulation struct {
	Host           Host
	Status         int
	StatusDuration map[int]int
	IntrahostModel IntrahostModel
	FitnessModel   FitnessModel
	Tree           GenotypeTree
}

// SusceptibleProcess executes within-host processes that occurs when a host
// is in the susceptible state.
func (sim *SingleHostSimulation) SusceptibleProcess(host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	// Decrement to -1 if pathogens exist
	if host.PathogenPopSize() > 0 {
		host.DecrementTimer()
	}
}

// InfectedProcess executes within-host processes that occurs when a host
// is in the infected state.
func (sim *SingleHostSimulation) InfectedProcess(host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	pathogens := host.Pathogens()
	var replicatedC <-chan GenotypeNode
	switch strings.ToLower(sim.IntrahostModel.ReplicationMethod()) {
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
	mutatedC := MutateSequence(replicatedC, sim.Tree, host.(*sequenceHost).IntrahostModel)
	// Clear current set of pathogens and get new set from the channel
	host.RemoveAllPathogens()
	for pathogen := range mutatedC {
		host.AddPathogen(pathogen)
	}
	host.DecrementTimer()
}

// InfectiveProcess executes within-host processes that occurs when a host
// is in the infective state. By default, it is same as InfectedProcess.
func (sim *SingleHostSimulation) InfectiveProcess(host Host, wg *sync.WaitGroup) {
	// timer decrement is done within the InfectedProcess function
	// Done() signal also executed within the InfectedProcess function
	sim.InfectedProcess(host, wg)
}

// RemovedProcess executes within-host processes that occurs when a host
// is in the removed state that is perpetually uninfectable.
func (sim *SingleHostSimulation) RemovedProcess(host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	host.RemoveAllPathogens()
}

// RecoveredProcess executes within-host processes that occurs when a host
// is in the recovered state that is perpetually uninfectable.
// This state is identically to Removed but is used to distinguish from
// a dead state.
func (sim *SingleHostSimulation) RecoveredProcess(host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	host.RemoveAllPathogens()
}

// DeadProcess executes within-host processes that occurs when a host
// is in the dead state state that is perpetually uninfectable.
// This state is identically to Removed but is used to distinguish from
// a recovered, but perpetually immune state.
func (sim *SingleHostSimulation) DeadProcess(host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	host.RemoveAllPathogens()
}

// VaccinatedProcess executes within-host processes that occurs when a host
// is in a globally immune state with the chance to become
// globally susceptible again.
func (sim *SingleHostSimulation) VaccinatedProcess(host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	host.RemoveAllPathogens()
	host.DecrementTimer()
}
