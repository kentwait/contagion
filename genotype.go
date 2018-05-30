package contagiongo

import (
	"fmt"
	"sync"

	"github.com/segmentio/ksuid"
)

// Genotype represents a unique pathogen sequence.
type Genotype interface {
	// Sequence returns the sequence of the current node.
	Sequence() []int
	// StringSequence returns the string representation of the
	// integer-coded sequence of the current node.
	StringSequence() string
	// Fitness returns the fitness value of this node based on its current
	// sequence and the given fitness model. If the fitness of the node has
	// been computed before using the same fitness model, then the value is
	// returned from memory and is not recomputed.
	Fitness(f FitnessModel) float64
	// NumSites returns the number of sites being modeled in this pathogen node.
	NumSites() int
	// StateCounts returns the number of sites by state, postion corresponds
	// to the state from 0 to n.
	StateCounts() map[int]int
}

type genotype struct {
	sync.RWMutex
	sequence    []int
	stateCounts map[int]int     // key is the state
	fitness     map[int]float64 // key is the fitness model id
}

// NewGenotype creates a new genotype from sequence.
func NewGenotype(s []int) Genotype {
	g := new(genotype)
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

func (n *genotype) Sequence() []int {
	return n.sequence
}

func (n *genotype) StringSequence() string {
	key := fmt.Sprintf("%v", n.sequence)
	key = key[1 : len(key)-1]
	return key
}

func (n *genotype) Fitness(f FitnessModel) float64 {
	id := f.ModelID()
	fitness, ok := n.fitness[id]
	if !ok {
		fitness, _ := f.ComputeFitness(n.sequence...)
		return fitness
	}
	return fitness
}

func (n *genotype) NumSites() int {
	return len(n.sequence)
}

func (n *genotype) StateCounts() map[int]int {
	return n.stateCounts
}

// GenotypeSet is a collection of genotypes.
type GenotypeSet interface {
	// Add adds the genotype to the set if the sequence does not exist yet.
	Add(g Genotype)
	// AddSequence creates a new genotype from the sequence if it is not present
	// in the set. Otherwise, returns the existing genotype in the set.
	AddSequence(s []int) Genotype
	// Remove removes genotype of a particular sequence from the set.
	Remove(s []int)
	// Size returns the size of the set.
	Size() int
}

type genotypeSet struct {
	sync.RWMutex
	set map[string]Genotype
}

// EmptyGenotypeSet creates a new empty set.
func EmptyGenotypeSet() GenotypeSet {
	set := new(genotypeSet)
	set.set = make(map[string]Genotype)
	return set
}

func (set *genotypeSet) Add(g Genotype) {
	key := fmt.Sprintf("%v", g.Sequence())
	key = key[1 : len(key)-1]
	set.Lock()
	defer set.Unlock()
	if _, exists := set.set[key]; !exists {
		set.set[key] = g
	}
}

func (set *genotypeSet) AddSequence(s []int) Genotype {
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

func (set *genotypeSet) Remove(s []int) {
	key := fmt.Sprintf("%v", s)
	key = key[1 : len(key)-1]
	set.Lock()
	defer set.Unlock()
	if _, exists := set.set[key]; exists {
		set.set[key] = nil
		delete(set.set, key)
	}
}

func (set *genotypeSet) Size() int {
	set.RLock()
	defer set.RUnlock()
	return len(set.set)
}

// GenotypeNode represents a genotype together with its relationship to its parents and children.
type GenotypeNode interface {
	// UID returns the unique ID of the node. Uses KSUID to generate
	// random unique IDs with effectively no collision.
	UID() ksuid.KSUID
	// Parents returns the parent of the node.
	Parents() []GenotypeNode
	// Children returns the children of the node.
	Children() []GenotypeNode
	// AddChild appends a child to the list of children.
	AddChild(child GenotypeNode)
	// Sequence returns the sequence of the current node.
	Sequence() []int
	// StringSequence returns the string representation of the
	// integer-coded sequence of the current node.
	StringSequence() string
	// CurrentGenotype returns the current genotype of the current node.
	CurrentGenotype() Genotype
	// History returns the list of sequences that resulted into the extant
	// sequence.
	History(h [][]int) [][]int
	// Fitness returns the fitness value of this node based on its current
	// sequence and the given fitness model. If the fitness of the node has
	// been computed before using the same fitness model, then the value is
	// returned from memory and is not recomputed.
	Fitness(f FitnessModel) float64
	// NumSites returns the number of sites being modeled in this pathogen node.
	NumSites() int
	// StateCounts returns the number of sites by state, postion corresponds
	// to the state from 0 to n.
	StateCounts() map[int]int
}

type genotypeNode struct {
	sync.RWMutex
	Genotype
	uid      ksuid.KSUID
	sequence []int
	parents  []GenotypeNode
	children []GenotypeNode
}

// NewGenotypeNode creates a new genotype node from a sequence.
func NewGenotypeNode(sequence []int, set GenotypeSet, parents ...GenotypeNode) GenotypeNode {
	genotype := set.AddSequence(sequence)

	// Create new node
	n := new(genotypeNode)
	n.uid = ksuid.New()
	// Assign its parent
	if len(parents) > 0 {
		n.parents = make([]GenotypeNode, len(parents))
		copy(n.parents, parents)
	} else {
		n.parents = []GenotypeNode{}
	}
	// Initialize children
	n.children = []GenotypeNode{}
	// Assign genotype
	n.Genotype = genotype
	// Copy sequence
	n.sequence = make([]int, len(sequence))
	copy(n.sequence, sequence)

	// Add new sequence as child of its parent
	for _, parent := range parents {
		parent.AddChild(n)
	}
	return n
}

func (n *genotypeNode) UID() ksuid.KSUID {
	return n.uid
}

func (n *genotypeNode) Parents() []GenotypeNode {
	return n.parents
}

func (n *genotypeNode) Children() []GenotypeNode {
	n.RLock()
	defer n.RUnlock()
	return n.children
}

func (n *genotypeNode) AddChild(child GenotypeNode) {
	n.Lock()
	defer n.Unlock()
	n.children = append(n.children, child)
}

func (n *genotypeNode) Sequence() []int {
	return n.sequence
}

func (n *genotypeNode) CurrentGenotype() Genotype {
	return n.Genotype
}

func (n *genotypeNode) History(h [][]int) [][]int {
	h = append(h, n.sequence)
	if len(n.parents) == 0 {
		return h
	}
	// TODO: Assumes no recombination. Only follows the first parent
	return n.parents[0].History(h)
}

// GenotypeTree represents the genotypes as a series of differences
// from its ancestor.
type GenotypeTree interface {
	// Set returns the GenotypeSet associated with this tree.
	Set() GenotypeSet
	// NewNode creates a new genotype node from a given sequence.
	// Automatically adds sequence to the genotypeSet if it is not yet present.
	NewNode(sequence []int, parents ...GenotypeNode) GenotypeNode
}

type genotypeTree struct {
	sync.RWMutex
	genotypes map[ksuid.KSUID]GenotypeNode
	set       GenotypeSet
}

// EmptyGenotypeTree creates a new empty genotype tree.
func EmptyGenotypeTree() GenotypeTree {
	tree := new(genotypeTree)
	tree.genotypes = make(map[ksuid.KSUID]GenotypeNode)
	tree.set = EmptyGenotypeSet()
	return tree
}

func (t *genotypeTree) Set() GenotypeSet {
	return t.set
}

func (t *genotypeTree) NewNode(sequence []int, parents ...GenotypeNode) GenotypeNode {
	genotype := t.set.AddSequence(sequence)

	// Create new node
	n := new(genotypeNode)
	n.uid = ksuid.New()
	// Assign its parent
	if len(parents) > 0 {
		n.parents = make([]GenotypeNode, len(parents))
		copy(n.parents, parents)
	} else {
		n.parents = []GenotypeNode{}
	}
	// Initialize children
	n.children = []GenotypeNode{}
	// Assign genotype
	n.Genotype = genotype
	// Copy sequence
	n.sequence = make([]int, len(sequence))
	copy(n.sequence, sequence)

	// Add new sequence as child of its parent
	for _, parent := range parents {
		parent.AddChild(n)
	}

	// Add to tree map
	t.Lock()
	defer t.Unlock()
	t.genotypes[n.uid] = n
	return n
}
