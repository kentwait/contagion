package contagiongo

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"

	rv "github.com/kentwait/randomvariate"
)

// MultinomialReplication replicates and selects sequences based on normalized fitness values used as probabilities.
func MultinomialReplication(pathogens []GenotypeNode, normedFitnesses []float64, newPopSize int) <-chan GenotypeNode {
	c := make(chan GenotypeNode)
	var wg sync.WaitGroup
	wg.Add(len(pathogens))
	for i, count := range rv.MultinomialA(newPopSize, normedFitnesses) {
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
func MutateSite(transitionProbs ...float64) uint8 {
	// Get new state
	for i, v := range rv.Multinomial(1, transitionProbs) {
		if v == 1 {
			return uint8(i)
		}
	}
	panic(fmt.Sprintf("improper transition probabilities %v", transitionProbs))
}

// MutateSequence adds substitution mutations to sequenceNode.
func MutateSequence(sequences <-chan GenotypeNode, tree GenotypeTree, model IntrahostModel) (<-chan GenotypeNode, <-chan GenotypeNode) {
	c := make(chan GenotypeNode) // all the sequences, whether mutated or untouched
	d := make(chan GenotypeNode) // new mutants
	var wg sync.WaitGroup
	for sequence := range sequences {
		wg.Add(1)
		go func(n GenotypeNode, model IntrahostModel, wg *sync.WaitGroup) {
			defer wg.Done()
			mu := model.MutationRate()
			// Copy sequence to make changes atomic
			sequence := make([]uint8, len(n.Sequence()))
			copy(sequence, n.Sequence())
			// Add mutations by state to account for unequal rates
			totalHits := 0
			if mu > 0 {
				for state, numSites := range n.StateCounts() {
					probs := model.TransitionProbs(int(state))
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
					// Returns empty list if hits == 0
					hitPositions := pickSites(hits, numSites, n.StatePositions(state))
					// Create new node per hit
					for _, pos := range hitPositions {
						sequence[pos] = MutateSite(probs...)
					}
					totalHits += hits
				}
			}
			if totalHits > 0 {
				newNode := tree.NewNode(sequence, totalHits, n)
				c <- newNode
				d <- newNode
			} else {
				c <- n
			}
		}(sequence, model, &wg)
	}
	go func() {
		wg.Wait()
		close(c)
		close(d)
	}()
	return c, d
}

// RecombineSequencePairs recombines two sequences at random positions
// similar to the behavior of diploid chromosomes.
func RecombineSequencePairs(numSeqs, numRecSites int, sequences <-chan GenotypeNode, tree GenotypeTree, model IntrahostModel) (<-chan GenotypeNode, <-chan GenotypeNode) {
	c := make(chan GenotypeNode) // all sequences regardless whether it recombined or not
	d := make(chan GenotypeNode) // only sequences that recombined

	// assumes all sequences have the same length
	// which may not be true if indels are implemented in the future
	// Expected number of recombinations over the entire sequence
	rate := model.RecombinationRate()

	// permute and determine pairs
	// if odd, the last 3 sequences form a triad
	permIdx := rand.Perm(numSeqs)
	seqGroupLookup := make(map[int]int)      // key is order of sequence in channel, value is the pair ID it belongs to
	seqGroup := make(map[int][]GenotypeNode) // key is the pair ID, values are list of genotype nodes
	seqTriadID := -1                         // is not negative if a triad exists

	// Create map of sequence pair ids
	groupID := 0
	for i := 1; i < len(permIdx); i += 2 {
		seqGroupLookup[permIdx[i]] = groupID
		seqGroupLookup[permIdx[i-1]] = groupID
		groupID++
	}
	if len(permIdx)%2 != 0 {
		seqGroupLookup[permIdx[len(permIdx)-1]] = groupID
		seqGroupLookup[permIdx[len(permIdx)-2]] = groupID
		seqGroupLookup[permIdx[len(permIdx)-3]] = groupID
		seqTriadID = groupID
	}

	x := 0
	var wg sync.WaitGroup
	for sequence := range sequences {
		groupID := seqGroupLookup[x]
		seqGroup[groupID] = append(seqGroup[groupID], sequence)
		switch {
		case len(seqGroup[groupID]) == 2 && (groupID != seqTriadID):
			fallthrough
		case len(seqGroup[groupID]) == 3 && (groupID == seqTriadID):
			go func(rate float64, numRecSites int, wg *sync.WaitGroup, seqs ...GenotypeNode) {
				defer wg.Done()
				// assumes all sequences have the same length
				// which may not be true if indels are implemented in the future
				// Expected number of recombinations over the entire sequence
				nrate := float64(numRecSites) * rate
				var hits int
				if nrate < 1/float64(numRecSites) {
					hits = rv.Poisson(nrate)
				} else {
					hits = rv.Binomial(numRecSites, rate)
				}

				// Create empty sequences
				recombinantSeqs := make([][]uint8, len(seqs))
				for i := range recombinantSeqs {
					recombinantSeqs[i] = make([]uint8, seqs[i].NumSites())
					copy(recombinantSeqs[i], seqs[i].Sequence())
				}

				// Determine positions
				hittablePositions := make([]int, numRecSites)
				for i := 0; i < numRecSites; i++ {
					hittablePositions[i] = i
				}
				hitPositions := pickSites(hits, numRecSites, hittablePositions)
				prevOrder := 0
				prevPos := 0
				if hitPositions[len(hitPositions)-1] < numRecSites-1 {
					hitPositions = append(hitPositions, numRecSites-1)
				}
				totalHits := 0
				for _, pos := range hitPositions {
					var idx0, idx1, idx2 int
					if len(seqs) == 3 {
						if prevPos == 0 {
							idx0, idx1, idx2 = 0, 1, 2
						} else {
							if prevOrder == 0 { // 0, 1, 2
								if r := rand.Intn(2); r == 0 {
									idx0, idx1, idx2 = 1, 2, 0
									prevOrder = 1
								} else {
									idx0, idx1, idx2 = 2, 0, 1
									prevOrder = 2
								}
							} else if prevOrder == 1 { // 1, 2, 0
								if r := rand.Intn(2); r == 0 {
									idx0, idx1, idx2 = 0, 1, 2
									prevOrder = 0
								} else {
									idx0, idx1, idx2 = 2, 0, 1
									prevOrder = 2
								}
							} else if prevOrder == 2 { // 2, 0, 1
								if r := rand.Intn(2); r == 0 {
									idx0, idx1, idx2 = 1, 2, 0
									prevOrder = 1
								} else {
									idx0, idx1, idx2 = 0, 1, 2
									prevOrder = 0
								}
							}
						}
						copy(recombinantSeqs[0][prevPos:pos], seqs[idx0].Sequence()[prevPos:pos])
						copy(recombinantSeqs[1][prevPos:pos], seqs[idx1].Sequence()[prevPos:pos])
						copy(recombinantSeqs[2][prevPos:pos], seqs[idx2].Sequence()[prevPos:pos])
					} else {
						if prevOrder%2 == 0 {
							idx0, idx1 = 0, 1
						} else if prevOrder%2 == 1 {
							idx0, idx1 = 1, 0
						}
						copy(recombinantSeqs[0][prevPos:pos], seqs[idx0].Sequence()[prevPos:pos])
						copy(recombinantSeqs[1][prevPos:pos], seqs[idx1].Sequence()[prevPos:pos])
						prevOrder++
					}
					prevPos = pos
					totalHits++
				}

				if totalHits > 0 {
					for _, recombinantSeq := range recombinantSeqs {
						newNode := tree.NewRecombinantNode(recombinantSeq, totalHits, seqs...)
						c <- newNode
						d <- newNode
					}
				} else {
					for _, n := range seqs {
						c <- n
					}
				}
			}(rate, numRecSites, &wg, seqGroup[groupID]...)
		}
		x++
	}
	go func() {
		wg.Wait()
		close(c)
		close(d)
	}()
	return c, d
}

// RecombineAnySequence recombines any two sequences at a random position
// similar to the behavior of template switching.
func RecombineAnySequence(numSeqs, numRecSites int, sequences <-chan GenotypeNode, tree GenotypeTree, model IntrahostModel) (<-chan GenotypeNode, <-chan GenotypeNode) {
	c := make(chan GenotypeNode) // all sequences regardless whether it recombined or not
	d := make(chan GenotypeNode) // only sequences that recombined

	// Collects all sequences into a list
	nodeMap := make(map[int]GenotypeNode)
	i := 0
	for sequence := range sequences {
		nodeMap[i] = sequence
		i++
	}

	// assumes all sequences have the same length
	// which may not be true if indels are implemented in the future
	rate := model.RecombinationRate()

	var wg sync.WaitGroup
	for x, node := range nodeMap {
		wg.Add(1)
		go func(x int, node GenotypeNode, nodes map[int]GenotypeNode, rate float64, numRecSites int, wg *sync.WaitGroup) {
			defer wg.Done()
			// Return immediately if nodes only has one node
			if len(nodes) == 1 {
				c <- node
				return
			}

			// Get recombination hits
			nrate := float64(numRecSites) * rate
			var hits int
			if nrate < 1/float64(numRecSites) {
				hits = rv.Poisson(nrate)
			} else {
				hits = rv.Binomial(numRecSites, rate)
			}

			// Create empty sequence
			recombinantSeq := make([]uint8, node.NumSites())
			copy(recombinantSeq, node.Sequence())

			// Determine positions
			hittablePositions := make([]int, numRecSites)
			for i := 0; i < numRecSites; i++ {
				hittablePositions[i] = i
			}
			hitPositions := pickSites(hits, numRecSites, hittablePositions)
			if hitPositions[len(hitPositions)-1] < numRecSites-1 {
				hitPositions = append(hitPositions, numRecSites-1)
			}

			prevNodePos := x
			prevPos := 0
			totalHits := 0
			var parents []GenotypeNode
			for _, pos := range hitPositions {
				repeat := true
				repeatCount := 0
				maxRepeat := 3
				for repeat {
					// iterate over the set of nodes until the selected node
					// is not the previous node
					for nodePos, node := range nodes {
						if prevNodePos != nodePos {
							copy(recombinantSeq[prevPos:pos], node.Sequence()[prevPos:pos])
							prevNodePos = nodePos
							repeat = false
							parents = append(parents, node)
							totalHits++
							break
						}
					}
					if !repeat {
						break
					}
					// In case after 3 tries it fails
					// fail on the current recombination position
					if repeatCount > maxRepeat {
						break
					}
					repeatCount++
				}
			}

			if totalHits > 0 {
				newNode := tree.NewRecombinantNode(recombinantSeq, totalHits, parents...)
				c <- newNode
				d <- newNode
			} else {
				c <- node
			}
		}(x, node, nodeMap, rate, numRecSites, &wg)
	}
	go func() {
		wg.Wait()
		close(c)
		close(d)
	}()
	return c, d
}

// // RecombineSequences recombines two sequences at a random position.
// func RecombineSequences(sequences <-chan GenotypeNode, tree GenotypeTree, model IntrahostModel) (<-chan GenotypeNode, <-chan GenotypeNode) {
// 	var sequenceList []GenotypeNode
// 	for sequence := range sequences {
// 		sequenceList = append(sequenceList, sequence)
// 	}
// 	// If no sequence received
// 	if len(sequenceList) == 0 {
// 		c := make(chan GenotypeNode)
// 		d := make(chan GenotypeNode)
// 		close(c)
// 		close(d)
// 		return c, d
// 	}
// 	rate := model.RecombinationRate()
// 	// assumes all sequences have the same length
// 	// which may not be true if indels are implemented in the future
// 	numSites := sequenceList[0].CurrentGenotype().NumSites() - 1
// 	// Expected number of recombinations over the entire sequence
// 	nrate := float64(numSites) * rate
// 	// Get number of hits in the sequence
// 	var hits int
// 	if nrate < 1/float64(numSites) {
// 		hits = rv.Poisson(nrate)
// 	} else {
// 		hits = rv.Binomial(numSites, rate)
// 	}
// 	// Get position of hits
// 	// Returns if hits == 0
// 	if hits == 0 {
// 		c := make(chan GenotypeNode)
// 		d := make(chan GenotypeNode)
// 		close(c)
// 		close(d)
// 		return c, d
// 	}
// 	hittablePositions := make([]int, numSites)
// 	for i := 0; i < numSites; i++ {
// 		hittablePositions[i] = i
// 	}
// 	hitPositions := pickSites(hits, numSites, hittablePositions)
// 	// Get pathogens to recombine
// 	// Recombine at site
// 	indices := make([]int, len(sequenceList))
// 	c := make(chan GenotypeNode)
// 	d := make(chan GenotypeNode)
// 	var wg sync.WaitGroup
// 	for _, pos := range hitPositions {
// 		go func(pos int, indices []int, sequenceList []GenotypeNode, wg *sync.WaitGroup) {
// 			pickedIndices := pickSites(2, len(indices), indices)
// 			node1 := sequenceList[pickedIndices[0]]
// 			node2 := sequenceList[pickedIndices[1]]
// 			// Recombine 2 sequences
// 			// a-a, b-b -> a-b
// 			newSequence1 := append(node1.Sequence()[:pos], node2.Sequence()[pos:]...)
// 			newNode1 := tree.NewNode(newSequence1, 0, node1, node2)
// 			c <- newNode1
// 			d <- newNode1
// 			// a-a, b-b -> b-a
// 			newSequence2 := append(node2.Sequence()[:pos], node1.Sequence()[pos:]...)
// 			newNode2 := tree.NewNode(newSequence2, 0, node2, node1)
// 			c <- newNode2
// 			d <- newNode2
// 		}(pos, indices, sequenceList, &wg)
// 	}
// 	go func() {
// 		wg.Wait()
// 		close(c)
// 		close(d)
// 	}()
// 	return c, d
// }

func pickSites(hitsNeeded, numSites int, positions []int) []int {
	if hitsNeeded == 0 {
		return []int{}
	}
	// Create hittable list of positions
	hittable := make([]int, len(positions))
	copy(hittable, positions)
	// Get position of hits
	hitPositions := make([]int, hitsNeeded)
	for i := 0; i < hitsNeeded; i++ {
		x := rand.Intn(numSites - i)
		hitPositions[i] = hittable[x]
		hittable = append(hittable[:x], hittable[x+1:]...)
	}
	sort.Ints(hitPositions)
	return hitPositions
}
