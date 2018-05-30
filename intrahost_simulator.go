package contagiongo

// import (
// 	"strings"
// 	"sync"
// )

// type SingleHostSimulation struct {
// 	host           Host
// 	hostStatus     int
// 	statusDuration map[int]int
// 	intrahostModel IntrahostModel
// 	fitnessModel   FitnessModel
// 	tree           SequenceTree
// }

// func NewSingleHostSimulationFromConfig(config )

// // InfectedProcess executes within-host processes that occurs when a host
// // is in the infected state.
// func (sim *SequenceSimulation) InfectedProcess(host EpidemicHost, wg *sync.WaitGroup) {
// 	// Convert []interface{} to []PathogenNode
// 	var pathogens []SequenceNode
// 	for _, p := range host.Pathogens() {
// 		pathogens = append(pathogens, p.(SequenceNode))
// 	}
// 	var replicatedC <-chan SequenceNode
// 	switch strings.ToLower(sim.ReplicationMethod) {
// 	case "multinomial":
// 		// TODO: compute normalized fitness of all pathogens
// 		normedFitnesses := []float64{}
// 		// get current and next pop size based on popsize function
// 		currentPopSize := host.PathogenPopSize()
// 		nextPopSize := host.NextPathogenPopSize(currentPopSize)
// 		// Execute
// 		replicatedC = SequenceMultinomialReplication(pathogens, normedFitnesses, nextPopSize)
// 	case "intrinsic":
// 		// TODO: compute growth rate from fitness value for each pathogen
// 		growthRates := []int{}
// 		// Execute
// 		replicatedC = SequenceIntrinsicRateReplication(pathogens, growthRates, nil)
// 	}
// 	// Mutate replicated pathogens
// 	mutatedC := SequenceMutate(replicatedC, sim.Tree, host.(*PathogenHost))
// 	// Clear current set of pathogens and get new set from the channel
// 	host.ClearPathogens()
// 	for pathogen := range mutatedC {
// 		host.AddPathogen(pathogen)
// 	}
// 	host.DecrementTimer()
// 	wg.Done()
// }
