package contagiongo

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/segmentio/ksuid"
)

func sampleGenotypeNode(numSites int, sets ...GenotypeSet) *genotypeNode {
	sequence := make([]int, numSites)
	for i := 0; i < numSites; i++ {
		s := rand.Intn(4)
		sequence[i] = s
	}
	var set *genotypeSet
	if len(sets) < 1 {
		set = EmptyGenotypeSet().(*genotypeSet)
	} else {
		set = sets[0].(*genotypeSet)
	}
	n := NewGenotypeNode(sequence, set).(*genotypeNode)
	return n
}

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
	growthRates := []int{0, 0, 0, 10}
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
		if pathogenCounter[i] != growthRate {
			t.Errorf(UnequalIntParameterError, fmt.Sprintf("number of pathogen %d", i), growthRate, pathogenCounter[i])
		}
	}
}

// func TestSequenceMutate(t *testing.T) {
// 	// Create a mock root
// 	root := sampleGenotypeNode(100)
// 	// Create mock IntrahostModel
// 	model := new(ConstantPopModel)
// 	model.mutationRate = 0.1
// 	model.transitionMatrix = [][]float64{
// 		[]float64{0, 1},
// 		[]float64{1, 0},
// 	}
// 	model.recombinationRate = 0
// 	model.popSize = 5
// 	// Send root 4 times to simulate population of 4
// 	c := make(chan GenotypeNode)
// 	go func() {
// 		for _, pathogen := range []GenotypeNode{root, root, root, root, root, root, root, root, root, root} {
// 			c <- pathogen
// 		}
// 		close(c)
// 	}()
// 	pathogens := MutateSequence(c, tree, model)
// 	pathogenCounter := make(map[ksuid.KSUID]int)
// 	counter := 0
// 	fmt.Println(root.UID())
// 	for pathogen := range pathogens {
// 		pathogenCounter[pathogen.UID()]++
// 		if pathogen.UID() != root.UID() {
// 			fmt.Print(pathogen.UID())
// 			fmt.Print(" ")
// 			cnt := 0
// 			for pathogen.UID() != root.UID() {
// 				pathogen = pathogen.Parents()[0]
// 				cnt++
// 			}
// 			fmt.Println(cnt)
// 			if cnt > 120 || cnt < 80 {
// 				t.Errorf(IntNotBetweenError, "number of mutations", 80, 120, cnt)
// 			}
// 		} else {
// 			fmt.Println(pathogen.UID())
// 		}
// 		counter++
// 	}

// 	if counter != 10 {
// 		t.Errorf(UnequalIntParameterError, "number of pathogens", 4, counter)
// 	}
// 	// TODO: test whether mutation took place, and if the number is correct
// 	// TODO: Add scenarios for binomial hits
// }
