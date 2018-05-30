package contagiongo

import (
	"fmt"
	"math/rand"
	"sync"

	rv "github.com/kentwait/randomvariate"
)

// MultinomialReplication replicates and selects sequences based on normalized fitness values used as probabilities.
func MultinomialReplication(pathogens []GenotypeNode, normedFitnesses []float64, newPopSize int) <-chan GenotypeNode {
	c := make(chan GenotypeNode)
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
func IntrinsicRateReplication(pathogens []GenotypeNode, replFitness []float64, immuneSystem interface{}) <-chan GenotypeNode {
	c := make(chan GenotypeNode)
	var wg sync.WaitGroup
	wg.Add(len(pathogens))
	for i, pathogen := range pathogens {
		go func(pathogen GenotypeNode, fitness float64, wg *sync.WaitGroup) {
			defer wg.Done()
			growthRate := rv.Poisson(fitness)
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

// MutateSequence adds substitution mutations to sequenceNode.
func MutateSequence(sequences <-chan GenotypeNode, tree GenotypeTree, model IntrahostModel) <-chan GenotypeNode {
	c := make(chan GenotypeNode)
	var wg sync.WaitGroup
	for sequence := range sequences {
		wg.Add(1)
		go func(n GenotypeNode, model IntrahostModel, wg *sync.WaitGroup) {
			defer wg.Done()
			mu := model.MutationRate()
			stateCounts := n.StateCounts()
			sequence := n.Sequence()
			// Add mutations by state to account for unequal rates
			for state, numSites := range stateCounts {
				probs := model.TransitionProbs(state)
				// Expected number of mutations over the entire sequence
				nmu := float64(numSites) * mu

				// Get number of hits in the sequence
				var hits int
				if nmu < 1.0 {
					hits = rv.Poisson(nmu)
				} else {
					hits = rv.Binomial(numSites, mu)
				}
				// Get position of hits
				hitPositions := pickSites(hits, numSites)

				// Create new node per hit
				for _, pos := range hitPositions {
					sequence[pos] = MutateSite(probs...)
				}
			}
			c <- tree.NewNode(sequence, n)
		}(sequence, model, &wg)
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	return c
}

func pickSites(hitsNeeded, numSites int) []int {
	// Create hittable list of positions
	hittable := make([]int, numSites)
	for i := range hittable {
		hittable[i] = i
	}
	// Get position of hits
	hitPositions := make([]int, hitsNeeded)
	for i := 0; i < hitsNeeded; i++ {
		x := rand.Intn(numSites - i)
		hittable = append(hittable[:x], hittable[x+1:]...)
		hitPositions[i] = x
	}
	return hitPositions
}
