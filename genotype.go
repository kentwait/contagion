package contagiongo

import (
	"fmt"
	"sync"

	"github.com/segmentio/ksuid"
)

// Genotype represents a unique pathogen sequence.
type Genotype interface {
	GenotypeUID() ksuid.KSUID
	// Sequence returns the sequence of the current node.
	Sequence() []uint8
	// SetSequence changes the sequence of genotype.
	SetSequence(sequence []uint8)
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
	// StateCounts returns the number of sites by state.
	StateCounts() map[uint8]int
	// StatePositions returns the indexes of sites in a particular state.
	StatePositions(state uint8) []int
}

type genotype struct {
	sync.RWMutex
	uid      ksuid.KSUID
	sequence []uint8
	statePos map[uint8][]int // key is the state
	fitness  map[int]float64 // key is the fitness model id
}

// NewGenotype creates a new genotype from sequence.
func NewGenotype(s []uint8) Genotype {
	g := new(genotype)
	// Copy sequence
	g.sequence = make([]uint8, len(s))
	copy(g.sequence, s)
	// Initial count of states
	g.statePos = make(map[uint8][]int)
	for i, state := range g.sequence {
		g.statePos[state] = append(g.statePos[state], i)
	}
	// Initialize other maps
	g.fitness = make(map[int]float64)
	return g
}

func (n *genotype) GenotypeUID() ksuid.KSUID {
	return n.uid
}

func (n *genotype) Sequence() []uint8 {
	return n.sequence
}

func (n *genotype) SetSequence(sequence []uint8) {
	n.sequence = nil
	n.sequence = make([]uint8, len(sequence))
	copy(n.sequence, sequence)
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
		n.fitness[id] = fitness
		return fitness
	}
	return fitness
}

func (n *genotype) NumSites() int {
	return len(n.sequence)
}

func (n *genotype) StateCounts() map[uint8]int {
	stateCounts := make(map[uint8]int)
	for state, positions := range n.statePos {
		stateCounts[state] = len(positions)
	}
	return stateCounts
}

func (n *genotype) StatePositions(state uint8) []int {
	return n.statePos[state]
}

// GenotypeSet is a collection of genotypes.
type GenotypeSet interface {
	// Add adds the genotype to the set if the sequence does not exist yet.
	Add(g Genotype)
	// AddSequence creates a new genotype from the sequence if it is not present
	// in the set. Otherwise, returns the existing genotype in the set.
	AddSequence(s []uint8) Genotype
	// Remove removes genotype of a particular sequence from the set.
	Remove(s []uint8)
	// Size returns the size of the set.
	Size() int
	Map() map[string]Genotype
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

func (set *genotypeSet) AddSequence(s []uint8) Genotype {
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

func (set *genotypeSet) Remove(s []uint8) {
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
	return len(set.set)
}

func (set *genotypeSet) Map() map[string]Genotype {
	return set.set
}

// GenotypeNode represents a genotype together with its relationship to its parents and children.
type GenotypeNode interface {
	// UID returns the unique ID of the node. Uses KSUID to generate
	// random unique IDs with effectively no collision.
	UID() ksuid.KSUID
	GenotypeUID() ksuid.KSUID
	// Parents returns the parent of the node.
	Parents() []GenotypeNode
	// Children returns the children of the node.
	Children() []GenotypeNode
	// AddChild appends a child to the list of children.
	AddChild(child GenotypeNode)
	// Sequence returns the sequence of the current node.
	Sequence() []uint8
	// SetSequence changes the sequence of genotype.
	SetSequence(sequence []uint8)
	// StringSequence returns the string representation of the
	// integer-coded sequence of the current node.
	StringSequence() string
	// CurrentGenotype returns the current genotype of the current node.
	CurrentGenotype() Genotype
	// History returns the list of sequences that resulted into the extant
	// sequence.
	History(h [][]uint8) [][]uint8
	// Fitness returns the fitness value of this node based on its current
	// sequence and the given fitness model. If the fitness of the node has
	// been computed before using the same fitness model, then the value is
	// returned from memory and is not recomputed.
	Fitness(f FitnessModel) float64
	// NumSites returns the number of sites being modeled in this pathogen node.
	NumSites() int
	// StateCounts returns the number of sites by state.
	StateCounts() map[uint8]int
	// StatePositions returns the indexes of sites in a particular state.
	StatePositions(state uint8) []int
}

type genotypeNode struct {
	sync.RWMutex
	Genotype
	uid      ksuid.KSUID
	parents  []GenotypeNode
	children []GenotypeNode
}

// NewGenotypeNode creates a new genotype node from a sequence.
// This should not be used to create a new genotype. Use the NewNode method in
// GenotypeTree instead.
func newGenotypeNode(sequence []uint8, set GenotypeSet, parents ...GenotypeNode) GenotypeNode {
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

func (n *genotypeNode) Sequence() []uint8 {
	return n.Genotype.Sequence()
}

func (n *genotypeNode) CurrentGenotype() Genotype {
	return n.Genotype
}

func (n *genotypeNode) History(h [][]uint8) [][]uint8 {
	h = append(h, n.Genotype.Sequence())
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
	NewNode(sequence []uint8, parents ...GenotypeNode) GenotypeNode
	// Nodes returns the map of genotype node ID found in the tree to its
	// corresponding genotype.
	NodeMap() map[ksuid.KSUID]GenotypeNode
}

type genotypeTree struct {
	sync.RWMutex
	nodes map[ksuid.KSUID]GenotypeNode
	set   GenotypeSet
}

// EmptyGenotypeTree creates a new empty genotype tree.
func EmptyGenotypeTree() GenotypeTree {
	tree := new(genotypeTree)
	tree.nodes = make(map[ksuid.KSUID]GenotypeNode)
	tree.set = EmptyGenotypeSet()
	return tree
}

func (t *genotypeTree) Set() GenotypeSet {
	t.RLock()
	defer t.RUnlock()
	return t.set
}

func (t *genotypeTree) NewNode(sequence []uint8, parents ...GenotypeNode) GenotypeNode {
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

	// Add new sequence as child of its parent
	for _, parent := range parents {
		parent.AddChild(n)
	}

	// Add to tree map
	t.Lock()
	defer t.Unlock()
	t.nodes[n.uid] = n
	return n
}

func (t *genotypeTree) NodeMap() map[ksuid.KSUID]GenotypeNode {
	t.RLock()
	defer t.RUnlock()
	return t.nodes
}
