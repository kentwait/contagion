package contagiongo

import "math/rand"

func sampleSequence(sites int) []uint8 {
	sequence := make([]uint8, sites)
	for i := 0; i < sites; i++ {
		sequence[i] = uint8(rand.Intn(2))
	}
	return sequence
}

func sampleGenotype() Genotype {
	sequence := []uint8{0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1}
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
	host := EmptySequenceHost(hostID)
	tree := sampleGenotypeTree(pathogens, sites)
	for _, n := range tree.NodeMap() {
		host.AddPathogens(n)
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
		tree.NewNode(pathogen, 0)
	}
	return tree
}

func sampleSequenceNodeEpidemic() *SequenceNodeEpidemic {
	sites := 100
	mu := 0.01
	constPopSize := 1000
	initPopSize := 1
	multiplicative := true

	sim := new(SequenceNodeEpidemic)
	sim.tree = sampleGenotypeTree(initPopSize, sites)
	sim.intrahostModels = map[int]IntrahostModel{
		0: sampleIntrahostModel(mu, constPopSize),
	}
	sim.fitnessModels = map[int]FitnessModel{
		0: sampleFitnessModel(multiplicative, sites),
	}
	sim.hosts = map[int]Host{
		0: EmptySequenceHost(0),
		1: EmptySequenceHost(1),
	}
	sim.hosts[0].SetIntrahostModel(sim.intrahostModels[0])
	sim.hosts[0].SetFitnessModel(sim.fitnessModels[0])
	for _, n := range sim.tree.NodeMap() {
		sim.hosts[0].AddPathogens(n)
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
	sim.hostNetwork = adjacencyMatrix{
		0: map[int]float64{1: 1.0},
		1: map[int]float64{0: 1.0},
	}
	sim.infectableStatuses = []int{SusceptibleStatusCode}

	return sim
}

func sampleEpidemicSimConfig() *epidemicSimConfig {
	conf := new(epidemicSimConfig)
	conf.NumGenerations = 10
	conf.NumIntances = 10
	conf.HostPopSize = 10
	conf.EpidemicModel = "sir"
	conf.PathogenSequencePath = "examples/test1.sir.pathogens.fa"
	conf.HostNetworkPath = "examples/test1.sir.network.txt"
	return conf
}

func sampleLogConfig() *logConfig {
	conf := new(logConfig)
	conf.LogFreq = 1
	conf.LogPath = "examples/test1.sir.log"
	return conf
}

func sampleIntrahostModelConfig() *intrahostModelConfig {
	conf := new(intrahostModelConfig)
	conf.ModelName = "constant-high_mutation"
	conf.HostIDs = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	conf.MutationRate = 0.01
	conf.TransitionMatrix = [][]float64{
		[]float64{0.0, 1.0},
		[]float64{1.0, 0.0},
	}
	conf.RecombinationRate = 0.0
	conf.ReplicationModel = "constant"
	conf.ConstantPopSize = 100
	conf.InfectedDuration = 10
	return conf
}

func sampleFitnessModelConfig() *fitnessModelConfig {
	conf := new(fitnessModelConfig)
	conf.HostIDs = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	conf.ModelName = "multiplicative"
	conf.FitnessModel = "multiplicative"
	conf.FitnessModelPath = "examples/test1.sir.fm.txt"
	conf.validated = false
	return conf
}

func sampleTransModelConfig() *transModelConfig {
	conf := new(transModelConfig)
	conf.HostIDs = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	conf.ModelName = "constant"
	conf.Mode = "constant"
	conf.TransmissionProb = 1.0
	conf.TransmissionSize = 1
	return conf
}

func sampleEvoEpiConfig() *EvoEpiConfig {
	conf := new(EvoEpiConfig)
	conf.SimParams = sampleEpidemicSimConfig()
	conf.LogParams = sampleLogConfig()
	conf.IntrahostModels = []*intrahostModelConfig{sampleIntrahostModelConfig()}
	conf.FitnessModels = []*fitnessModelConfig{sampleFitnessModelConfig()}
	conf.TransmissionModels = []*transModelConfig{sampleTransModelConfig()}
	return conf
}
