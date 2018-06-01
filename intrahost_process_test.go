package contagiongo

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/segmentio/ksuid"
)

func TestSequenceMultinomialReplication(t *testing.T) {
	set := EmptyGenotypeSet()
	p1 := sampleGenotypeNode(100, set)
	p2 := sampleGenotypeNode(100, set)
	p3 := sampleGenotypeNode(100, set)
	p4 := sampleGenotypeNode(100, set)
	uids := []ksuid.KSUID{p1.UID(), p2.UID(), p3.UID(), p4.UID()}
	// Only p4 should be present
	normedFitness := []float64{0, 0, 0, 1.0}
	popSize := 4
	pathogens := MultinomialReplication([]GenotypeNode{p1, p2, p3, p4}, normedFitness, popSize)

	pathogenCounter := make(map[int]int)
	for pathogen := range pathogens {
		i := -1
		for idx, uid := range uids {
			if uid == pathogen.UID() {
				i = idx
				break
			}
		}
		pathogenCounter[i]++
	}
	for i, fitness := range normedFitness {
		expected := int(fitness) * 4
		if pathogenCounter[i] != expected {
			t.Errorf(UnequalIntParameterError, fmt.Sprintf("number of pathogen %d", i), expected, pathogenCounter[i])
		}
	}
}

func TestIntrinsicRateReplication(t *testing.T) {
	rand.Seed(1)
	set := EmptyGenotypeSet()
	p1 := sampleGenotypeNode(100, set)
	p2 := sampleGenotypeNode(100, set)
	p3 := sampleGenotypeNode(100, set)
	p4 := sampleGenotypeNode(100, set)
	uids := []ksuid.KSUID{p1.UID(), p2.UID(), p3.UID(), p4.UID()}
	// Only p4 should be present
	growthRates := []float64{0., 0., 0., 10.}
	pathogens := IntrinsicRateReplication([]GenotypeNode{p1, p2, p3, p4}, growthRates, nil)

	pathogenCounter := make(map[int]int)
	for pathogen := range pathogens {
		i := -1
		for idx, uid := range uids {
			if uid == pathogen.UID() {
				i = idx
				break
			}
		}
		pathogenCounter[i]++
	}
	for i, growthRate := range growthRates {
		if pathogenCounter[i] != int(growthRate) {
			t.Errorf(UnequalIntParameterError, fmt.Sprintf("number of pathogen %d", i), int(growthRate), pathogenCounter[i])
		}
	}
}

func TestSequenceMutate(t *testing.T) {
	rand.Seed(0)
	// Create a mock tree
	tree := EmptyGenotypeTree()
	rootSeq := []uint8{
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	}
	root := tree.NewNode(rootSeq, 0)
	// Create mock IntrahostModel
	model := new(ConstantPopModel)
	model.mutationRate = 0.1
	model.transitionMatrix = [][]float64{
		[]float64{0, 1},
		[]float64{1, 0},
	}
	model.recombinationRate = 0
	model.popSize = 10
	// Send root 4 times to simulate population of 4
	c := make(chan GenotypeNode)
	go func() {
		for i := 0; i < model.popSize; i++ {
			r := make([]uint8, len(rootSeq))
			copy(r, rootSeq)
			pathogen := tree.NewNode(r, 0)
			c <- pathogen
		}
		close(c)
	}()
	pathogens, newMutantsC := MutateSequence(c, tree, model)
	go func() {
		for range newMutantsC {
		}
	}()
	counter := 0
	diffMean := 0.0
	for pathogen := range pathogens {
		diff := 0
		// fmt.Println(pathogen.StringSequence())
		for i := 0; i < len(pathogen.StringSequence()); i++ {
			if pathogen.StringSequence()[i] != root.StringSequence()[i] {
				diff++
			}
		}
		diffMean += float64(diff)
		counter++
	}
	diffMean = diffMean / float64(counter)
	if counter != model.popSize {
		t.Errorf(UnequalIntParameterError, "number of pathogens", model.popSize, counter)
	}
	if diffMean < 8 || diffMean > 12 {
		t.Errorf(FloatNotBetweenError, "average number of mutations", 8., 12., diffMean)
	}
	// TODO: test whether mutation took place, and if the number is correct
	// TODO: Add scenarios for binomial hits
}

func TestSequenceMutate2(t *testing.T) {
	rand.Seed(0)
	// Create a mock tree
	tree := EmptyGenotypeTree()
	rootSeq := []uint8{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	}
	root := tree.NewNode(rootSeq, 0)
	// Create mock IntrahostModel
	model := new(ConstantPopModel)
	model.mutationRate = 0.1
	model.transitionMatrix = [][]float64{
		[]float64{0, 1},
		[]float64{1, 0},
	}
	model.recombinationRate = 0
	model.popSize = 10
	// Send root 4 times to simulate population of 4
	c := make(chan GenotypeNode)
	go func() {
		for i := 0; i < model.popSize; i++ {
			r := make([]uint8, len(rootSeq))
			copy(r, rootSeq)
			pathogen := tree.NewNode(r, 0)
			c <- pathogen
		}
		close(c)
	}()
	pathogens, newMutantsC := MutateSequence(c, tree, model)
	go func() {
		for range newMutantsC {
		}
	}()
	counter := 0
	diffMean := 0.0
	for pathogen := range pathogens {
		diff := 0
		// fmt.Println(pathogen.StringSequence())
		for i := 0; i < len(pathogen.StringSequence()); i++ {
			if pathogen.StringSequence()[i] != root.StringSequence()[i] {
				diff++
			}
		}
		diffMean += float64(diff)
		counter++
	}
	diffMean = diffMean / float64(counter)
	if counter != model.popSize {
		t.Errorf(UnequalIntParameterError, "number of pathogens", model.popSize, counter)
	}
	if diffMean < 8 || diffMean > 12 {
		t.Errorf(FloatNotBetweenError, "average number of mutations", 8., 12., diffMean)
	}
	// TODO: test whether mutation took place, and if the number is correct
	// TODO: Add scenarios for binomial hits
}

func TestPickSites(t *testing.T) {
	hitsNeeded := 10
	numSites := 10
	positions := make([]int, numSites)
	for i := range positions {
		positions[i] = i
	}
	hitPositions := pickSites(hitsNeeded, numSites, positions)
	sort.Slice(hitPositions, func(i, j int) bool { return hitPositions[i] < hitPositions[j] })
	if fmt.Sprintf("%v", positions) != fmt.Sprintf("%v", hitPositions) {
		t.Errorf(UnequalStringParameterError, "hit positions", fmt.Sprintf("%v", positions), fmt.Sprintf("%v", hitPositions))
	}
}
