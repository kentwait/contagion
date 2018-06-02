package contagiongo

import (
	"fmt"
	"sync"
	"testing"
)

func TestEvoEpiConfig_Validate(t *testing.T) {
	conf := sampleEvoEpiConfig()
	err := conf.Validate()
	if err != nil {
		t.Error(err)
	}
	// TODO: Add errors
}

func TestEvoEpiConfig_NewSimulation(t *testing.T) {
	conf := sampleEvoEpiConfig()
	err := conf.Validate()
	if err != nil {
		t.Error(err)
	}
	sim0, err := conf.NewSimulation()
	sim := sim0.(*evoEpiSimulation)
	if err != nil {
		t.Error(err)
	}

	// Expected values
	numHosts := 10
	infectedHostID := 0
	pathogenPopSize := 100
	numNeighbors := 7

	// Tests
	// Host population
	if l := len(sim.HostMap()); l != numHosts {
		t.Errorf(UnequalIntParameterError, "number of hosts", numHosts, l)
	}
	// Infected host should only be host 0
	for i, host := range sim.HostMap() {
		if host.PathogenPopSize() > 0 && infectedHostID != i {
			t.Errorf(
				UnequalIntParameterError,
				fmt.Sprintf("nubmer of pathogens in host %d", i),
				0, host.PathogenPopSize(),
			)
		}
	}
	// Number of pathogens in host 0
	if l := sim.Host(0).PathogenPopSize(); l != pathogenPopSize {
		t.Errorf(
			UnequalIntParameterError,
			fmt.Sprintf("nubmer of pathogens in host %d", 0),
			pathogenPopSize, l,
		)
	}
	// Number of neighbors in host 0
	if l := len(sim.HostNeighbors(0)); l != numNeighbors {
		t.Errorf(
			UnequalIntParameterError,
			fmt.Sprintf("nubmer of neighbors for host %d", 0),
			numNeighbors, l,
		)
	}
	// Number of intrahost models
	if l := len(sim.intrahostModels); l != 1 {
		t.Errorf(UnequalIntParameterError, "number of intrahost models", 1, l)
	}
	// Number of fitness models
	if l := len(sim.fitnessModels); l != 1 {
		t.Errorf(UnequalIntParameterError, "number of fitness models", 1, l)
	}
}

func TestEvoEpiConfig_NewSimulation_NewInstance(t *testing.T) {
	conf := sampleEvoEpiConfig()
	err := conf.Validate()
	if err != nil {
		t.Error(err)
	}
	sim0, err := conf.NewSimulation()
	if err != nil {
		t.Error(err)
	}

	sim1, err := sim0.NewInstance()
	if err != nil {
		t.Error(err)
	}

	// Check is configs have the same pointer
	sim0ConfigPtr := fmt.Sprintf("%p", sim0.(*evoEpiSimulation).config)
	sim1ConfigPtr := fmt.Sprintf("%p", sim1.(*evoEpiSimulation).config)
	if sim0ConfigPtr != sim1ConfigPtr {
		t.Errorf(
			NotIdenticalPointerError,
			"sim0 host 0", sim0.(*evoEpiSimulation).config,
			"sim1 host 0", sim1.(*evoEpiSimulation).config,
		)
	}
	// Check is simulations have the same pointer
	sim0Ptr := fmt.Sprintf("%p", sim0)
	sim1Ptr := fmt.Sprintf("%p", sim1)
	if sim0Ptr == sim1Ptr {
		t.Errorf(IdenticalPointerError, "sim0", sim0, "sim1", sim1)
	}
	// Check whether hosts have the same pointer
	sim0HostPtr := fmt.Sprintf("%p", sim0.Host(0))
	sim1HostPtr := fmt.Sprintf("%p", sim1.Host(0))
	if sim0HostPtr == sim1HostPtr {
		t.Errorf(
			IdenticalPointerError,
			"sim0 host 0", sim0.Host(0),
			"sim1 host 0", sim1.Host(0),
		)
	}
	// Check whether trees have the same pointer
	sim0TreePtr := fmt.Sprintf("%p", sim0.(*evoEpiSimulation).tree)
	sim1TreePtr := fmt.Sprintf("%p", sim1.(*evoEpiSimulation).tree)
	if sim0TreePtr == sim1TreePtr {
		t.Errorf(
			IdenticalPointerError,
			"sim0 tree", sim0.(*evoEpiSimulation).tree,
			"sim1 tree", sim1.(*evoEpiSimulation).tree,
		)
	}
	// Check whether adjacency matrices have the same pointer
	sim0NetPtr := fmt.Sprintf("%p", sim0.(*evoEpiSimulation).hostNetwork)
	sim1NetPtr := fmt.Sprintf("%p", sim1.(*evoEpiSimulation).hostNetwork)
	if sim0NetPtr == sim1NetPtr {
		t.Errorf(
			IdenticalPointerError,
			"sim0 host network", sim0.(*evoEpiSimulation).hostNetwork,
			"sim1 host network", sim1.(*evoEpiSimulation).hostNetwork,
		)
	}
}

func TestEvoEpiConfig_NewSimulation_InfectedProcess(t *testing.T) {
	conf := sampleEvoEpiConfig()
	err := conf.Validate()
	if err != nil {
		t.Error(err)
	}
	sim0, err := conf.NewSimulation()
	sim := sim0.(*evoEpiSimulation)
	if err != nil {
		t.Error(err)
	}
	// Record parameters
	originalSequence := sim.Host(0).Pathogen(0).StringSequence()
	constPopSize := sim.config.(*EvoEpiConfig).IntrahostModels[0].ConstantPopSize
	// Run infected process on the simulation
	var wg sync.WaitGroup
	c := make(chan MutationPackage)
	wg.Add(2)
	go sim.InfectedProcess(0, 0, sim.Host(0), c, &wg)
	go sim.InfectedProcess(0, 0, sim.Host(1), c, &wg)
	go func() {
		wg.Wait()
		close(c)
	}()
	newNodeCnt := 0
	for range c {
		// fmt.Println(pack)
		newNodeCnt++
	}

	// Expectations
	// population size should increase from 10 to 100 after 1 step
	// on average, half of sites should be mutated compared to the original
	counter := 0
	diffMean := 0.0
	// fmt.Println(originalSequence)
	for _, n := range sim.hosts[0].Pathogens() {
		// compare sequences
		diff := 0
		for i := 0; i < len(originalSequence); i++ {
			if originalSequence[i] != n.StringSequence()[i] {
				diff++
			}
		}
		// fmt.Println(n.StringSequence(), diff)
		diffMean += float64(diff)
		if diff > 0 {
			counter++
		}
	}
	diffMean = diffMean / float64(counter)
	if counter != constPopSize {
		t.Errorf(UnequalIntParameterError, "number of pathogens", constPopSize, counter)
	}
	if diffMean < 0.8 || diffMean > 1.2 {
		t.Errorf(FloatNotBetweenError, "average number of mutations", 0.8, 1.2, diffMean)
	}
	if newNodeCnt < counter {
		t.Errorf(UnequalIntParameterError, "number of new mutants", counter, newNodeCnt)
	}
}

func TestIntrahostModelConfig_Validate(t *testing.T) {
	conf := sampleIntrahostModelConfig()
	err := conf.Validate()
	if err != nil {
		t.Error(err)
	}
	// TODO: Add errors
}

func TestIntrahostModelConfig_CreateModel(t *testing.T) {
	conf := sampleIntrahostModelConfig()
	err := conf.Validate()
	if err != nil {
		t.Error(err)
	}
	modelID := 0
	model, err := conf.CreateModel(modelID)
	if err != nil {
		t.Error(err)
	}

	// Tests
	if i := model.ModelID(); i != 0 {
		t.Errorf(UnequalIntParameterError, "model ID", 0, i)
	}
	if s := model.NextPathogenPopSize(1); s != conf.ConstantPopSize {
		t.Errorf(UnequalIntParameterError, "next pathogen population size", conf.ConstantPopSize, s)
	}
	if s := model.ReplicationMethod(); s != "relative" {
		t.Errorf(UnequalStringParameterError, "replication method", "relative", s)
	}
	// TODO: Add cases
}

func TestFitnessModelConfig_Validate(t *testing.T) {
	conf := sampleFitnessModelConfig()
	err := conf.Validate()
	if err != nil {
		t.Error(err)
	}

	// Add errors
	// Invalid fitness_model keyword
	conf = sampleFitnessModelConfig()
	conf.FitnessModel = "intrinsic"
	err = conf.Validate()
	if err == nil {
		t.Errorf(ExpectedErrorWhileError, "validating fitness_model keyword")
	}
	// Invalid fitness_model_path
	conf = sampleFitnessModelConfig()
	conf.FitnessModelPath = "examples/mock.txt"
	err = conf.Validate()
	if err == nil {
		t.Errorf(ExpectedErrorWhileError, "validating fitness_model_path")
	}
}

func TestFitnessModelConfig_CreateModel(t *testing.T) {
	// Multiplicative
	conf := sampleFitnessModelConfig()
	err := conf.Validate()
	if err != nil {
		t.Error(err)
	}
	modelID := 0
	model, err := conf.CreateModel(modelID)
	if err != nil {
		t.Error(err)
	}

	if i := model.ModelID(); i != 0 {
		t.Errorf(UnequalIntParameterError, "model ID", 0, i)
	}
	fitness, err := model.ComputeFitness(sampleSequence(100)...)
	if err != nil {
		t.Error(err)
	}
	if fitness != 0 {
		t.Errorf(UnequalFloatParameterError, "log fitness value", 0.0, fitness)
	}

	// Additive
	conf = sampleFitnessModelConfig()
	conf.FitnessModelPath = "examples/test1.sir.fm.additive.txt"
	conf.FitnessModel = "additive"
	err = conf.Validate()
	if err != nil {
		t.Error(err)
	}
	modelID = 0
	model, err = conf.CreateModel(modelID)
	if err != nil {
		t.Error(err)
	}

	if i := model.ModelID(); i != 0 {
		t.Errorf(UnequalIntParameterError, "model ID", 0, i)
	}
	fitness, err = model.ComputeFitness(sampleSequence(100)...)
	if err != nil {
		t.Error(err)
	}
	if fmt.Sprintf("%.6f", fitness) != fmt.Sprintf("%.6f", 2.) {
		t.Errorf(UnequalFloatParameterError, "decimal fitness value", 2.0, fitness)
	}
}
