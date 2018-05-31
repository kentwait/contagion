package contagiongo

import "math/rand"

func sampleSequence(sites int) []int {
	sequence := make([]int, sites)
	for i := 0; i < sites; i++ {
		sequence[i] = rand.Intn(2)
	}
	return sequence
}

func sampleGenotype() Genotype {
	sequence := []int{0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1}
	g := NewGenotype(sequence)
	return g
}

func sampleGenotypeNode(sites int, sets ...GenotypeSet) *genotypeNode {
	var set *genotypeSet
	if len(sets) < 1 {
		set = EmptyGenotypeSet().(*genotypeSet)
	} else {
		set = sets[0].(*genotypeSet)
	}
	n := newGenotypeNode(sampleSequence(sites), set).(*genotypeNode)
	return n
}

func sampleInfectedHost(hostID, pathogens, sites int) Host {
	host := NewEmptySequenceHost(hostID)
	tree := sampleGenotypeTree(pathogens, sites)
	for _, n := range tree.Nodes() {
		host.AddPathogen(n)
	}
	return host
}

func sampleIntrahostModel(mutationRate float64, popSize int) IntrahostModel {
	model := new(ConstantPopModel)
	model.mutationRate = mutationRate
	model.transitionMatrix = [][]float64{
		[]float64{0, 1},
		[]float64{1, 0},
	}
	model.recombinationRate = 0
	model.popSize = popSize
	model.statusDuration = map[int]int{
		InfectedStatusCode: 10,
	}
	return model
}

func sampleFitnessIntrahostModel(mutationRate float64, popSize int) IntrahostModel {
	model := new(FitnessDependentPopModel)
	model.mutationRate = mutationRate
	model.transitionMatrix = [][]float64{
		[]float64{0, 1},
		[]float64{1, 0},
	}
	model.recombinationRate = 0
	model.maxPopSize = popSize
	model.statusDuration = map[int]int{
		InfectedStatusCode: 10,
	}
	return model
}

func sampleFitnessModel(multiplicative bool, sites int) FitnessModel {
	if multiplicative {
		return NeutralMultiplicativeFM(0, "neutral", sites, 2)
	}
	return NeutralAdditiveFM(0, "additive", sites, 2, 2)
}

func sampleGenotypeTree(roots, sites int) GenotypeTree {
	tree := EmptyGenotypeTree()
	pathogen := sampleSequence(sites)
	for i := 0; i < roots; i++ {
		tree.NewNode(pathogen)
	}
	return tree
}

func sampleEvoEpiSimulation() *evoEpiSimulation {
	sites := 100
	mu := 10e-5
	constPopSize := 1000
	initPopSize := 1
	multiplicative := true

	sim := new(evoEpiSimulation)
	sim.tree = sampleGenotypeTree(initPopSize, sites)
	sim.intrahostModels = map[int]IntrahostModel{
		0: sampleIntrahostModel(mu, constPopSize),
	}
	sim.fitnessModels = map[int]FitnessModel{
		0: sampleFitnessModel(multiplicative, sites),
	}
	sim.hosts = map[int]Host{
		0: NewEmptySequenceHost(0),
		1: NewEmptySequenceHost(1),
	}
	sim.hosts[0].SetIntrahostModel(sim.intrahostModels[0])
	sim.hosts[0].SetFitnessModel(sim.fitnessModels[0])
	for _, n := range sim.tree.Nodes() {
		sim.hosts[0].AddPathogen(n)
	}
	sim.hosts[1].SetIntrahostModel(sim.intrahostModels[0])
	sim.hosts[1].SetFitnessModel(sim.fitnessModels[0])
	sim.statuses = map[int]int{
		0: InfectedStatusCode,
		1: SusceptibleStatusCode,
	}
	sim.timers = map[int]int{
		0: 10,
		1: 0,
	}
	sim.hostNeighborhoods = map[int][]Host{
		0: []Host{sim.hosts[1]},
		1: []Host{sim.hosts[0]},
	}
	sim.hostNetwork = map[int]map[int]float64{
		0: map[int]float64{1: 1.0},
		1: map[int]float64{0: 1.0},
	}
	sim.infectableStatuses = []int{SusceptibleStatusCode}

	return sim
}
