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

// SequenceNode represents a sequence genotype in the sequence tree.
type SequenceNode interface {
	// UID returns the unique ID of the pathogen node. Uses KSUID to generate random unique IDs with effectively no collision.
	UID() ksuid.KSUID
	// Parent returns the parent of the node.
	Parents() []SequenceNode
	// Children returns the children of the node.
	Children() []SequenceNode
	// AddChild appends a child to the list of children.
	AddChild(n SequenceNode)
	// Sequence returns the sequence of the current node.
	Sequence() []int
	// Fitness returns the fitness value of this node based on its current
	// sequence and the given fitness model. If the fitness of the node has
	// been computed before using the same fitness model, then the value is
	// returned from memory and is not recomputed.
	Fitness(landscape FitnessModel) float64
	// NumSites returns the number of sites being modeled in this pathogen node.
	NumSites() int
	// StateCounts returns the number of sites by state, postion corresponds to the state from 0 to n.
	StateCounts() map[int]int
	// History returns the list of sequences that resulted into the extant
	// sequence.
	History(h [][]int) [][]int
}

type sequenceNode struct {
	sync.RWMutex
	uid      ksuid.KSUID
	parents  []SequenceNode
	children []SequenceNode
	// subs        map[int][]int             // key the position of the substitution, value is list of new states. Essentially tracks if multiple hits occured
	// recombs     map[int][][]*sequenceNode // key is the position of recombination, value is a list of sequence node pairs representing its parents
	sequence    []int
	stateCounts map[int]int     // key is the state
	fitness     map[int]float64 // key is the fitness model id
}

func (n *sequenceNode) UID() ksuid.KSUID {
	return n.uid
}

func (n *sequenceNode) Parents() []SequenceNode {
	return n.parents
}

func (n *sequenceNode) Children() []SequenceNode {
	return n.children
}

func (n *sequenceNode) AddChild(child SequenceNode) {
	n.Lock()
	defer n.Unlock()
	n.children = append(n.children, child)
}

func (n *sequenceNode) Sequence() []int {
	return n.sequence
}

func (n *sequenceNode) Fitness(f FitnessModel) float64 {
	id := f.ID()
	fitness, ok := n.fitness[id]
	if !ok {
		fitness, _ := f.ComputeFitness(n.sequence...)
		return fitness
	}
	return fitness
}

func (n *sequenceNode) NumSites() int {
	return len(n.sequence)
}

func (n *sequenceNode) StateCounts() map[int]int {
	return n.stateCounts
}

func (n *sequenceNode) History(h [][]int) [][]int {
	h = append(h, n.sequence)
	if len(n.parents) == 0 {
		return h
	}
	// TODO: Assumes no recombination. Only follows the first parent
	return n.parents[0].History(h)
}

// SequenceTree represents the pathogen as a series of one-hit differences
// from its ancestor.
type SequenceTree interface {
	// Sequence returns a node in the tree. If uid is not found, returns nil.
	Sequence(uid ksuid.KSUID) *sequenceNode
	// NewSub creates a new node based on the parent and the
	// position and new state of the substitution.
	NewSub(parent SequenceNode, position, state int) *sequenceNode
	// NewRecomb creates a new node given two parents and the position at which recombination occurred.
	NewRecomb(parent1, parent2 SequenceNode, position int) *sequenceNode
	// NewRoot adds a new root sequence node.
	NewRoot(sequence []int) *sequenceNode
}

// TODO: Create inf sites, fixed sites, multiple hit trees
// Currently, sequence tree is a multiple hits tree
// Depend on NewHit method
type sequenceTree struct {
	sync.RWMutex
	roots      map[ksuid.KSUID]*sequenceNode
	pathogens  map[ksuid.KSUID]*sequenceNode
	subHits    map[int]int
	recombHits map[int]int
}

// NewSequenceTree creates a new pathogen tree with a single root.
func NewSequenceTree(rootSequence []int) SequenceTree {
	// Create new root node
	// Assign unique ID
	n := new(sequenceNode)
	n.uid = ksuid.New()
	// A root node does not have any parents.
	n.parents = []SequenceNode{}
	n.children = []SequenceNode{}
	// Copy sequence
	n.sequence = make([]int, len(rootSequence))
	copy(n.sequence, rootSequence)
	// Initialize maps
	n.stateCounts = make(map[int]int)
	n.fitness = make(map[int]float64)
	// Count the initial number of states across all sites
	for _, s := range n.sequence {
		n.stateCounts[s]++
	}
	// Create new tree and initialize maps
	tree := new(sequenceTree)
	tree.pathogens = make(map[ksuid.KSUID]*sequenceNode)
	tree.roots = make(map[ksuid.KSUID]*sequenceNode)
	tree.subHits = make(map[int]int)
	tree.recombHits = make(map[int]int)
	// Add the node to the tree and make it a root node
	tree.pathogens[n.UID()] = n
	tree.roots[n.UID()] = n
	return tree
}

func (t *sequenceTree) Sequence(uid ksuid.KSUID) *sequenceNode {
	t.RLock()
	defer t.RUnlock()
	n, ok := t.pathogens[uid]
	if !ok {
		return nil
	}
	return n
}

func (t *sequenceTree) NewSub(parent SequenceNode, position, state int) *sequenceNode {
	// Create new node
	n := new(sequenceNode)
	n.uid = ksuid.New()
	// Assign its parent
	n.parents = []SequenceNode{parent.(*sequenceNode)}
	// Copy sequence from parent, then change at specified position
	n.sequence = make([]int, len(n.parents[0].Sequence()))
	copy(n.sequence, n.parents[0].Sequence())
	n.sequence[position] = state
	// Initialize maps
	n.stateCounts = make(map[int]int)
	n.fitness = make(map[int]float64)
	// Copy state counts
	for s, cnt := range n.parents[0].(*sequenceNode).stateCounts {
		n.stateCounts[s] = cnt
	}
	// Decrement count of original state, increase count in new state
	n.stateCounts[n.parents[0].Sequence()[position]]--
	n.stateCounts[state]++

	// Add new sequence as child of its parent
	parent.AddChild(n)
	// Add new sequence to map of tree
	t.Lock()
	defer t.Unlock()
	t.pathogens[n.uid] = n
	return n
}

func (t *sequenceTree) NewRecomb(parent1, parent2 SequenceNode, position int) *sequenceNode {
	// Create new node
	n := new(sequenceNode)
	n.uid = ksuid.New()
	// Assign its parent
	n.parents = []SequenceNode{parent1.(*sequenceNode), parent2.(*sequenceNode)}
	// Copy sequence from parent1, then append sequence of parent2 at the given
	// position
	n.sequence = make([]int, len(n.parents[0].Sequence()))
	copy(n.sequence, n.parents[0].Sequence())
	n.sequence = append(n.sequence[0:position], parent2.Sequence()[position:len(parent2.Sequence())]...)
	// Initialize maps
	n.stateCounts = make(map[int]int)
	n.fitness = make(map[int]float64)
	// Create fresh count of states
	for _, s := range n.sequence {
		n.stateCounts[s]++
	}

	// Add new sequence as child of both parents
	parent1.AddChild(n)
	parent2.AddChild(n)
	// Add new sequence to map of tree
	t.Lock()
	defer t.Unlock()
	t.pathogens[n.uid] = n
	return n
}

func (t *sequenceTree) NewRoot(rootSequence []int) *sequenceNode {
	// Create new root node
	// Assign unique ID
	n := new(sequenceNode)
	n.uid = ksuid.New()
	// A root node does not have any parents.
	n.parents = []SequenceNode{}
	n.children = []SequenceNode{}
	// Copy sequence
	n.sequence = make([]int, len(rootSequence))
	copy(n.sequence, rootSequence)
	// Initialize maps
	n.stateCounts = make(map[int]int)
	n.fitness = make(map[int]float64)
	// Count the initial number of states across all sites
	for _, s := range n.sequence {
		n.stateCounts[s]++
	}

	// Add to pathogen and root maps
	t.Lock()
	defer t.Unlock()
	t.pathogens[n.uid] = n
	t.roots[n.uid] = n
	return n
}
