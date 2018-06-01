package contagiongo

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

func TestEvoEpiSimulation_Getters(t *testing.T) {
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
	for _, n := range sim.tree.NodeMap() {
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
	sim.hostNetwork = adjacencyMatrix{
		0: map[int]float64{1: 1.0},
		1: map[int]float64{0: 1.0},
	}
	sim.infectableStatuses = []int{SusceptibleStatusCode}

	// Tests
	if sim.Host(0).ID() != sim.hosts[0].ID() {
		t.Errorf(UnequalIntParameterError, "host ID", sim.hosts[0].ID(), sim.Host(0).ID())
	}
	if sim.HostStatus(0) != sim.statuses[0] {
		t.Errorf(UnequalIntParameterError, "host 0 status", sim.statuses[0], sim.HostStatus(0))
	}
	if sim.HostTimer(0) != sim.timers[0] {
		t.Errorf(UnequalIntParameterError, "host 0 timer", sim.timers[0], sim.HostTimer(0))
	}
	if len(sim.HostMap()) != len(sim.hosts) {
		t.Errorf(UnequalIntParameterError, "number of hosts", len(sim.hosts), len(sim.HostMap()))
	}
	if len(sim.HostNeighbors(0)) != len(sim.hostNeighborhoods[0]) {
		t.Errorf(UnequalIntParameterError, "number of host 0 neighbors", len(sim.hostNeighborhoods[0]), len(sim.HostNeighbors(0)))
	}
	if sim.InfectableStatuses()[0] != SusceptibleStatusCode {
		t.Errorf(UnequalIntParameterError, "infectable status", SusceptibleStatusCode, sim.InfectableStatuses()[0])
	}
}

func TestEvoEpiSimulation_SusceptibleProcess(t *testing.T) {
	rand.Seed(0)
	sites := 100
	mu := 0.5
	constPopSize := 100
	initPopSize := 10
	multiplicative := true

	sim := new(evoEpiSimulation)
	// Creates a tree with n roots of all identical sequences (same genotype)
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
	var originalSequence string
	for _, n := range sim.tree.NodeMap() {
		sim.hosts[0].AddPathogen(n)
		originalSequence = n.StringSequence()
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

	// Test setup
	var wg sync.WaitGroup
	wg.Add(2)
	go sim.SusceptibleProcess(0, 0, sim.Host(0), &wg)
	go sim.SusceptibleProcess(0, 0, sim.Host(1), &wg)
	wg.Wait()

	// Expectations
	// population size should not change
	// no mutation should occur
	counter := 0
	diffMean := 0.0
	for _, n := range sim.hosts[0].Pathogens() {
		// compare sequences
		diff := 0
		for i := 0; i < len(originalSequence); i++ {
			if originalSequence[i] != n.StringSequence()[i] {
				diff++
			}
		}
		diffMean += float64(diff)
		counter++
	}
	diffMean = diffMean / float64(counter)
	if counter != initPopSize {
		t.Errorf(UnequalIntParameterError, "number of pathogens", initPopSize, counter)
	}
	if diffMean != 0 {
		t.Errorf(UnequalFloatParameterError, "average number of mutations", 0., diffMean)
	}
}

func TestEvoEpiSimulation_InfectedProcess_Relative(t *testing.T) {
	rand.Seed(0)
	sites := 100
	mu := 0.1
	constPopSize := 100
	initPopSize := 10
	multiplicative := true

	sim := new(evoEpiSimulation)
	// Creates a tree with n roots of all identical sequences (same genotype)
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
	var originalSequence string
	for _, n := range sim.tree.NodeMap() {
		sim.hosts[0].AddPathogen(n)
		originalSequence = n.StringSequence()
		fmt.Println(n.StringSequence())
	}
	fmt.Println("")

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

	// Test setup
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
	for pack := range c {
		fmt.Println(pack)
		newNodeCnt++
	}

	// Expectations
	// population size should increase from 10 to 100 after 1 step
	// on average, half of sites should be mutated compared to the original
	counter := 0
	diffMean := 0.0
	fmt.Println(originalSequence)
	for _, n := range sim.hosts[0].Pathogens() {
		// compare sequences
		diff := 0
		for i := 0; i < len(originalSequence); i++ {
			if originalSequence[i] != n.StringSequence()[i] {
				diff++
			}
		}
		fmt.Println(n.StringSequence(), diff)
		diffMean += float64(diff)
		counter++
	}
	diffMean = diffMean / float64(counter)
	if counter != constPopSize {
		t.Errorf(UnequalIntParameterError, "number of pathogens", constPopSize, counter)
	}
	if diffMean < 9 || diffMean > 11 {
		t.Errorf(FloatNotBetweenError, "average number of mutations", 9., 11., diffMean)
	}
	if newNodeCnt < counter {
		t.Errorf(UnequalIntParameterError, "number of new mutants", counter, newNodeCnt)
	}
}

func TestEvoEpiSimulation_InfectedProcess_Additive(t *testing.T) {
	rand.Seed(0)
	sites := 100
	mu := 0.1
	maxPopSize := 1000
	initPopSize := 10
	growthRate := 2
	multiplicative := false

	sim := new(evoEpiSimulation)
	// Creates a tree with n roots of all identical sequences (same genotype)
	sim.tree = sampleGenotypeTree(initPopSize, sites)
	sim.intrahostModels = map[int]IntrahostModel{
		0: sampleFitnessIntrahostModel(mu, maxPopSize),
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
	var originalSequence string
	for _, n := range sim.tree.NodeMap() {
		sim.hosts[0].AddPathogen(n)
		originalSequence = n.StringSequence()
		fmt.Println(n.StringSequence())
	}
	fmt.Println("")

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

	// Test setup
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
	for pack := range c {
		fmt.Println(pack)
		newNodeCnt++
	}

	// Expectations
	// population size should increase from 10 to 20 after 1 step
	// (initPopSize * growthRate)
	// on average, half of sites should be mutated compared to the original
	counter := 0
	diffMean := 0.0
	fmt.Println(originalSequence)
	for _, n := range sim.hosts[0].Pathogens() {
		// compare sequences
		diff := 0
		for i := 0; i < len(originalSequence); i++ {
			if originalSequence[i] != n.StringSequence()[i] {
				diff++
			}
		}
		fmt.Println(n.StringSequence(), diff)
		diffMean += float64(diff)
		counter++
	}
	diffMean = diffMean / float64(counter)
	expPopSize := initPopSize * growthRate
	errVal := 5
	if counter < expPopSize-errVal || counter > expPopSize+errVal {
		t.Errorf(IntNotBetweenError, "number of pathogens", expPopSize-errVal, expPopSize+errVal, counter)
	}
	if diffMean < 9 || diffMean > 11 {
		t.Errorf(FloatNotBetweenError, "average number of mutations", 9., 11., diffMean)
	}
	if newNodeCnt < counter {
		t.Errorf(UnequalIntParameterError, "number of new mutants", counter, newNodeCnt)
	}
}

// TODO: Create separate test for testing contents of MutationPackage channel
