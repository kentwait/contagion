package contagiongo

// TODO: Update to use code of sir_simulation.go

import (
	"math"
	"strings"
	"sync"
)

// singleHostSimulation simulates pathogens infection within a single host.
type singleHostSimulation struct {
	host           Host
	status         int
	timer          int
	statusDuration map[int]int
	intrahostModel IntrahostModel
	fitnessModel   FitnessModel
	tree           GenotypeTree
}

func (sim *singleHostSimulation) Host() Host {
	return sim.host
}

func (sim *singleHostSimulation) HostStatus() int {
	return sim.status
}

func (sim *singleHostSimulation) SetHostStatus(status int) {
	sim.status = status
}

func (sim *singleHostSimulation) HostStatusDuration(status int) int {
	return sim.statusDuration[status]
}

func (sim *singleHostSimulation) HostTimer() int {
	return sim.timer
}

func (sim *singleHostSimulation) SetHostTimer(interval int) {
	sim.timer = interval
}

// SusceptibleProcess executes within-host processes that occurs when a host
// is in the susceptible state.
func (sim *singleHostSimulation) SusceptibleProcess(i, t int, host Host, wg *sync.WaitGroup) {
	defer wg.Done()
}

// ExposedProcess executes within-host processes that occurs when a host
// is in the exposed state. By default, it is same as InfectedProcess.
func (sim *singleHostSimulation) ExposedProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup) {
	// timer decrement is done within the InfectedProcess function
	// Done() signal also executed within the InfectedProcess function
	sim.InfectedProcess(i, t, host, c, wg)
	// TODO: Threshold to be considered infective instead of exposed
}

// InfectedProcess executes within-host processes that occurs when a host
// is in the infected state.
func (sim *singleHostSimulation) InfectedProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup) {
	defer wg.Done()
	pathogens := host.Pathogens()
	var replicatedC <-chan GenotypeNode
	switch strings.ToLower(sim.intrahostModel.ReplicationMethod()) {
	case "relative":
		// Get log fitness values for each pathogen
		logFitnesses := make([]float64, len(pathogens))
		var minLogFitness float64
		// Compute log total fitness and get max value
		for i, pathogen := range pathogens {
			logFitnesses[i] = pathogen.Fitness(host.GetFitnessModel())
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
		// TODO: Expose this in the interface
		nextPopSize := host.GetIntrahostModel().NextPathogenPopSize(currentPopSize)
		// Execute
		replicatedC = MultinomialReplication(pathogens, normedFitnesses, nextPopSize)
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
func (sim *singleHostSimulation) InfectiveProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup) {
	// timer decrement is done within the InfectedProcess function
	// Done() signal also executed within the InfectedProcess function
	sim.InfectedProcess(i, t, host, c, wg)
}

// RemovedProcess executes within-host processes that occurs when a host
// is in the removed state that is perpetually uninfectable.
func (sim *singleHostSimulation) RemovedProcess(i, t int, host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	host.RemoveAllPathogens()
}

// RecoveredProcess executes within-host processes that occurs when a host
// is in the recovered state that is perpetually uninfectable.
// This state is identically to Removed but is used to distinguish from
// a dead state.
func (sim *singleHostSimulation) RecoveredProcess(i, t int, host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	host.RemoveAllPathogens()
}

// DeadProcess executes within-host processes that occurs when a host
// is in the dead state state that is perpetually uninfectable.
// This state is identically to Removed but is used to distinguish from
// a recovered, but perpetually immune state.
func (sim *singleHostSimulation) DeadProcess(i, t int, host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	host.RemoveAllPathogens()
}

// VaccinatedProcess executes within-host processes that occurs when a host
// is in a globally immune state with the chance to become
// globally susceptible again.
func (sim *singleHostSimulation) VaccinatedProcess(i, t int, host Host, wg *sync.WaitGroup) {
	defer wg.Done()
	host.RemoveAllPathogens()
}
