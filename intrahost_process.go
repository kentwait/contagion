package contagiongo

import (
	"fmt"
	"sync"

	rv "github.com/kentwait/randomvariate"
)

// MultinomialReplication replicates and selects sequences based on normalized fitness values used as probabilities.
func MultinomialReplication(pathogens []SequenceNode, normedFitnesses []float64, newPopSize int) <-chan SequenceNode {
	c := make(chan SequenceNode)
	var wg sync.WaitGroup
	wg.Add(len(pathogens))
	for i, count := range rv.Multinomial(newPopSize, normedFitnesses) {
		go func(i, count int, wg *sync.WaitGroup) {
			defer wg.Done()
			for x := 0; x < count; x++ {
				c <- pathogens[i]
			}
		}(i, count, &wg)
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	return c
}

// IntrinsicRateReplication replicates pathogens by considering their
// fitness value as the growth rate.
func IntrinsicRateReplication(pathogens []SequenceNode, replFitness []int, immuneSystem interface{}) <-chan SequenceNode {
	c := make(chan SequenceNode)
	var wg sync.WaitGroup
	wg.Add(len(pathogens))
	for i, pathogen := range pathogens {
		go func(pathogen SequenceNode, fitness int, wg *sync.WaitGroup) {
			defer wg.Done()
			growthRate := rv.PoissonXL(float64(fitness))
			for i := 0; i < growthRate; i++ {
				c <- pathogen
			}
		}(pathogen, replFitness[i], &wg)
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	return c
}

// MutateSite returns the new state of a site based on the
// given a set of transition probabilities.
func MutateSite(transitionProbs ...float64) int {
	// Get new state
	for i, v := range rv.Multinomial(1, transitionProbs) {
		if v == 1 {
			return i
		}
	}
	panic(fmt.Sprintf("improper transition probabilities %v", transitionProbs))
}
