package contagiongo

import (
	"fmt"
	"sync"

	"github.com/segmentio/ksuid"
)

// Genotype represents a unique pathogen sequence.
type Genotype struct {
	sync.RWMutex
	sequence    []int
	stateCounts map[int]int     // key is the state
	fitness     map[int]float64 // key is the fitness model id
}

// NewGenotype returns a new Genotype given a sequence.
func NewGenotype(s []int) *Genotype {
	g := new(Genotype)
	// Copy sequence
	g.sequence = make([]int, len(s))
	copy(g.sequence, s)
	// Initial count of states
	g.stateCounts = make(map[int]int)
	for _, state := range g.sequence {
		g.stateCounts[state]++
	}
	// Initialize other maps
	g.fitness = make(map[int]float64)
	return g
}

// Sequence returns the sequence of the current node.
func (n *Genotype) Sequence() []int {
	return n.sequence
}

// Fitness returns the fitness value of this node based on its current
// sequence and the given fitness model. If the fitness of the node has
// been computed before using the same fitness model, then the value is
// returned from memory and is not recomputed.
func (n *Genotype) Fitness(f FitnessModel) float64 {
	id := f.ID()
	fitness, ok := n.fitness[id]
	if !ok {
		fitness, _ := f.ComputeFitness(n.sequence...)
		return fitness
	}
	return fitness
}

// NumSites returns the number of sites being modeled in this pathogen node.
func (n *Genotype) NumSites() int {
	return len(n.sequence)
}

// StateCounts returns the number of sites by state, postion corresponds to the state from 0 to n.
func (n *Genotype) StateCounts() map[int]int {
	return n.stateCounts
}

// GenotypeSet is a collection of Genotypes
type GenotypeSet struct {
	sync.RWMutex
	set map[string]*Genotype
}

// EmptyGenotypeSet creates a new empty set.
func EmptyGenotypeSet() *GenotypeSet {
	set := new(GenotypeSet)
	set.set = make(map[string]*Genotype)
	return set
}

// Add adds the genotype to the set if the sequence does not exist yet.
func (set *GenotypeSet) Add(g *Genotype) {
	key := fmt.Sprintf("%v", g.Sequence())
	key = key[1 : len(key)-1]
	set.Lock()
	defer set.Unlock()
	if _, exists := set.set[key]; !exists {
		set.set[key] = g
	}
}

// AddSequence creates a new Genotype from the sequence if it is not present
// in the set. Otherwise, returns the existing Genotype in the set.
func (set *GenotypeSet) AddSequence(s []int) *Genotype {
	key := fmt.Sprintf("%v", s)
	key = key[1 : len(key)-1]
	set.Lock()
	defer set.Unlock()
	g, exists := set.set[key]
	if !exists {
		g := NewGenotype(s)
		set.set[key] = g
		return g
	}
	return g
}

// Remove removes Genotype of a particular sequence from the set.
func (set *GenotypeSet) Remove(key string) {
	set.Lock()
	defer set.Unlock()
	set.set[key] = nil
	delete(set.set, key)
}

// GenotypeNode represents a genotype together with its relationship to its parents and children.
type GenotypeNode struct {
	sync.RWMutex
	uid      ksuid.KSUID
	sequence []int
	genotype *Genotype
	parents  []*GenotypeNode
	children []*GenotypeNode
}

// NewGenotypeNode creates a new Genotype node from a sequence.
// Automatically adds sequence to the GenotypeSet if it is not yet present.
func NewGenotypeNode(sequence []int, set *GenotypeSet, parents ...*GenotypeNode) *GenotypeNode {
	genotype := set.AddSequence(sequence)

	// Create new node
	n := new(GenotypeNode)
	n.uid = ksuid.New()
	// Assign its parent
	n.parents = make([]*GenotypeNode, len(parents))
	copy(n.parents, parents)
	n.children = []*GenotypeNode{}
	n.genotype = genotype
	// Copy sequence
	n.sequence = make([]int, len(sequence))
	copy(n.sequence, sequence)

	// Add new sequence as child of its parent
	for _, parent := range parents {
		parent.AddChild(n)
	}
	return n
}

// UID returns the unique ID of the node. Uses KSUID to generate random unique IDs with effectively no collision.
func (n *GenotypeNode) UID() ksuid.KSUID {
	return n.uid
}

// Parents returns the parent of the node.
func (n *GenotypeNode) Parents() []*GenotypeNode {
	return n.parents
}

// Children returns the children of the node.
func (n *GenotypeNode) Children() []*GenotypeNode {
	n.RLock()
	defer n.RUnlock()
	return n.children
}

// AddChild appends a child to the list of children.
func (n *GenotypeNode) AddChild(child *GenotypeNode) {
	n.Lock()
	defer n.Unlock()
	n.children = append(n.children, child)
}

// Sequence returns the sequence of the current node.
func (n *GenotypeNode) Sequence() []int {
	return n.sequence
}

// History returns the list of sequences that resulted into the extant
// sequence.
func (n *GenotypeNode) History(h [][]int) [][]int {
	h = append(h, n.sequence)
	if len(n.parents) == 0 {
		return h
	}
	// TODO: Assumes no recombination. Only follows the first parent
	return n.parents[0].History(h)
}
