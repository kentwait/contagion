package contagiongo

import (
	"sync"

	"github.com/segmentio/ksuid"
)

// StopCondition describes simulation conditions that must be
// satisfied in order for the simulation to continue.
// The Check method checks if the simulation still satisfies the
// imposed condition.
type StopCondition interface {
	Reason() string
	Check(sim Epidemic) bool
}

// AlleleExists is a stopping condition that checks if
// a particular allele in a given site still exists in the simulation.
type alleleExists struct {
	char uint8
	site int
}

// NewAlleleExistsCondition creates a new StopCondition that stops the
// simulation once the given char at a particular site becomes
// extinct.
func NewAlleleExistsCondition(char uint8, site int) StopCondition {
	cond := new(alleleExists)
	cond.char = uint8(char)
	cond.site = site
	return cond
}

func (cond *alleleExists) Reason() string {
	return "allele lost"
}

// Check looks at the
func (cond *alleleExists) Check(sim Epidemic) bool {
	c := make(chan bool)
	var wg sync.WaitGroup
	for _, host := range sim.HostMap() {
		wg.Add(1)
		go func(host Host, c chan<- bool, wg *sync.WaitGroup) {
			defer wg.Done()
			resultMap := make(map[ksuid.KSUID]bool)
			for _, node := range host.Pathogens() {
				genotypeUID := node.GenotypeUID()
				if _, exists := resultMap[genotypeUID]; !exists {
					if node.CurrentGenotype().Sequence()[cond.site] == cond.char {
						resultMap[genotypeUID] = true
						c <- true
					} else {
						resultMap[genotypeUID] = false
					}
				} else {
					if resultMap[genotypeUID] {
						c <- true
					}
				}
			}
		}(host, c, &wg)
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	exists := false
	for range c {
		exists = true
	}
	return exists
}

// GenotypeExists is a stopping condition that checks if
// a particular sequence still exists in the simulation.
type genotypeExists struct {
	sequence []uint8
}

// NewGenotypeExistsCondition creates a new StopCondition that stops the
// simulation once the given sequence genotype becomes
// extinct.
func NewGenotypeExistsCondition(sequence []uint8) StopCondition {
	cond := new(genotypeExists)
	cond.sequence = make([]uint8, len(sequence))
	copy(cond.sequence, sequence)
	return cond
}

func (cond *genotypeExists) Reason() string {
	return "genotype lost"
}

// Check looks in all infected hosts in the simulation to
// check if the genotype in question still exists.
// Return false if the genotype was not found in at least one host.
func (cond *genotypeExists) Check(sim Epidemic) bool {
	// This pipeline method does not work
	// getGenotypes := func(sim Epidemic) <-chan Genotype {
	// 	c := make(chan Genotype)
	// 	genotypeSet := make(map[ksuid.KSUID]bool)
	// 	for _, host := range sim.HostMap() {
	// 		for _, node := range host.Pathogens() {
	// 			uid := node.GenotypeUID()
	// 			if _, exists := genotypeSet[uid]; !exists {
	// 				c <- node.CurrentGenotype()
	// 			}
	// 		}
	// 	}
	// 	close(c)
	// 	return c
	// }
	// getMatches := func(sequence []uint8, c <-chan Genotype) <-chan bool {
	// 	d := make(chan bool)
	// 	var wg2 sync.WaitGroup
	// 	for genotype := range c {
	// 		go func(genotype Genotype, d chan<- bool, wg2 *sync.WaitGroup) {
	// 			defer wg2.Done()
	// 			matchCount := 0
	// 			for i, char := range genotype.Sequence() {
	// 				if char == cond.sequence[i] {
	// 					matchCount++
	// 				}
	// 			}
	// 			if matchCount == len(genotype.Sequence()) {
	// 				d <- true
	// 			}
	// 		}(genotype, d, &wg2)
	// 	}
	// 	wg2.Wait()
	// 	close(d)
	// 	return d
	// }
	// exists := false
	// for range getMatches(cond.sequence, getGenotypes(sim)) {
	// 	exists = true
	// }
	// return exists

	// TODO: Create more effective method that doesnt duplicate work on
	// reading genotypes
	c := make(chan bool)
	var wg sync.WaitGroup
	for _, host := range sim.HostMap() {
		wg.Add(1)
		go func(host Host, c chan<- bool, wg *sync.WaitGroup) {
			defer wg.Done()
			resultMap := make(map[ksuid.KSUID]bool)
			for _, node := range host.Pathogens() {
				genotype := node.CurrentGenotype()
				genotypeUID := node.GenotypeUID()
				if _, exists := resultMap[genotypeUID]; !exists {
					matchCount := 0
					for i, char := range genotype.Sequence() {
						if char == cond.sequence[i] {
							matchCount++
						}
					}
					if matchCount == len(genotype.Sequence()) {
						resultMap[genotypeUID] = true
						c <- true
					} else {
						resultMap[genotypeUID] = false
					}
				} else {
					if resultMap[genotypeUID] {
						c <- true
					}
				}
			}
		}(host, c, &wg)
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	exists := false
	for range c {
		exists = true
	}
	return exists
}

// AlleleFixedLost is a stopping condition that checks if
// a particular allele in a given site has fixed or been lost.
type alleleFixedLost struct {
	char   uint8
	site   int
	reason string
}

// NewAlleleFixedLostCondition creates a new StopCondition that stops the
// simulation once the particular allele in a given site has either
// been fixed or been lost.
func NewAlleleFixedLostCondition(char uint8, site int) StopCondition {
	cond := new(alleleFixedLost)
	cond.char = uint8(char)
	cond.site = site
	return cond
}

func (cond *alleleFixedLost) Reason() string {
	return cond.reason
}

// Check looks in all infected hosts in the simulation to
// check if the allele at a particular site has been fixed or lost.
// Check returns false when the allele is either fixed or lost across
// pathogens in all hosts.
// Check considers an allele fixed when all pathogens in all hosts
// have that allele. If the allele cannot be found on any pathogen
// sequence in any host, the allele is considered lost.
func (cond *alleleFixedLost) Check(sim Epidemic) bool {
	c := make(chan bool)
	var wg sync.WaitGroup
	for _, host := range sim.HostMap() {
		wg.Add(1)
		go func(host Host, c chan<- bool, wg *sync.WaitGroup) {
			defer wg.Done()
			resultMap := make(map[ksuid.KSUID]bool)
			for _, node := range host.Pathogens() {
				genotypeUID := node.GenotypeUID()
				if _, exists := resultMap[genotypeUID]; !exists {
					if node.CurrentGenotype().Sequence()[cond.site] == cond.char {
						resultMap[genotypeUID] = true
						c <- true
					} else {
						resultMap[genotypeUID] = false
						c <- false
					}
				} else {
					c <- resultMap[genotypeUID]
				}
			}
		}(host, c, &wg)
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	fixed := true
	lost := true
	for e := range c {
		if e {
			// One true output means it is not yet lost
			lost = false
		} else {
			// One false output means that it is not yet fixed
			fixed = false
		}
	}
	// Return true if not fixed, or not lost
	// fixed   lost    continue
	// true    false   false
	// false   true    false
	// false   false   true
	// true    true    false* (cannot happen)
	if fixed {
		cond.reason = "allele fixed"
	} else if lost {
		cond.reason = "allele lost"
	}
	return !fixed && !lost
}
