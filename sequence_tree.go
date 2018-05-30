package contagiongo

import (
	"fmt"
	"sync"

	"github.com/segmentio/ksuid"
)

// SequenceNode represents a sequence genotype in the sequence tree.
type SequenceNode interface {
	// ID returns the unique ID of the pathogen node.
	ID() int
	// Parent returns the parent of the node.
	Parent() SequenceNode
	// Children returns the children of the node.
	Children() []SequenceNode
	// Sequence returns the sequence of the current node.
	Sequence() []int
	// Fitness returns the fitness value of this node based on its current
	// sequence and the given fitness model. If the fitness of the node has
	// been computed before using the same fitness model, then the value is
	// returned from memory and is not recomputed.
	Fitness(landscape SequenceFitness) float64
	// NumSites returns the number of sites being modeled in this pathogen node.
	NumSites() int
	// StateCounts returns the number of sites by state, postion corresponds to the state from 0 to n.
	StateCounts() map[int]int
	// History returns the list of sequences that resulted into the extant
	// sequence.
	History(h [][]int) [][]int
}

type sequenceNode struct {
	uid      ksuid.KSUID
	parents  []SequenceNode
	children []SequenceNode
	// subs        map[int][]int             // key the position of the substitution, value is list of new states. Essentially tracks if multiple hits occured
	// recombs     map[int][][]*sequenceNode // key is the position of recombination, value is a list of sequence node pairs representing its parents
	sequence    []int
	stateCounts map[int]int     // key is the state
	fitness     map[int]float64 // key is the fitness model id
}

func (n *sequenceNode) ID() ksuid.KSUID {
	return n.uid
}

func (n *sequenceNode) Parents() []SequenceNode {
	return n.parents
}

func (n *sequenceNode) Children() []SequenceNode {
	return n.children
}

func (n *sequenceNode) Sequence() []int {
	return n.sequence
}

func (n *sequenceNode) Fitness(f SequenceFitness) float64 {
	id := f.ID()
	if _, ok := n.fitness[id]; !ok {
		fitness, _ := f.Fitness(n.sequence...)
		return fitness
	}
	return n.fitness[id]
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
	// Sequence returns a node in the tree.
	Sequence(id int) (SequenceNode, error)
	// NewSub creates a new node based on the parent and the
	// position and new state of the substitution.
	NewSub(parent SequenceNode, position, state int) SequenceNode
	// NewRecomb creates a new node given two parents and the position at which recombination occurred.
	NewRecomb(parent1, parent2 SequenceNode, position int) SequenceNode
	// NewRoot adds a new root sequence node.
	NewRoot(id int, sequence []int) (SequenceNode, error)
}

// TODO: Create inf sites, fixed sites, multiple hit trees
// Currently, sequence tree is a multiple hits tree
// Depend on NewHit method
type sequenceTree struct {
	sync.RWMutex
	roots     map[int]*sequenceNode
	pathogens map[int]*sequenceNode
	lastID    int
	hits      map[int]int
}

// NewSequenceTree creates a new pathogen tree with a single root.
func NewSequenceTree(rootID int, rootSequence []int) (SequenceTree, error) {
	// Create new root node
	// A root node does not have a parent.
	n := new(sequenceNode)
	n.parent = nil
	n.children = []SequenceNode{}
	n.hits = make(map[int]int)
	n.sequence = rootSequence
	n.numStates = make(map[int]int)
	// Count the initial number of states across all sites
	for _, s := range n.sequence {
		n.numStates[s]++
	}
	// Assign ID
	n.id = rootID

	// Create new tree
	tree := new(sequenceTree)
	tree.pathogens = make(map[int]*sequenceNode)
	tree.roots = make(map[int]*sequenceNode)
	tree.lastID = 0
	// Add the node to the tree and make it a root node
	tree.pathogens[rootID] = n
	tree.roots[rootID] = n
	return tree, nil
}

// Sequence returns a node in the tree.
func (t *sequenceTree) Sequence(id int) (SequenceNode, error) {
	t.RLock()
	defer t.RUnlock()
	p, ok := t.pathogens[id]
	if !ok {
		return nil, fmt.Errorf("sequence tree "+IntKeyNotFoundError, id)
	}
	return p, nil
}

// NewHit creates a new node based on the parent and the position and the
// new state of the substitution.
func (t *sequenceTree) NewHit(parent SequenceNode, position, state int) SequenceNode {
	// Create new node
	n := new(sequenceNode)
	// Assign its parent
	n.parent = parent.(*sequenceNode)
	// Change sequence from parent
	n.numStates = make(map[int]int)
	n.sequence = make([]int, len(n.parent.sequence))
	copy(n.sequence, n.parent.sequence)
	for i, s := range n.sequence {
		if i == position {
			n.sequence[i] = state
			n.numStates[state]++
		} else {
			n.numStates[s]++
		}
	}

	// Assign ID and update pathogen map
	t.Lock()
	t.lastID++
	n.id = t.lastID
	t.pathogens[n.id] = n
	// Add new pathogen to parent
	n.parent.children = append(n.parent.children, n)
	t.Unlock()
	return n
}

func (t *sequenceTree) NewRoot(id int, sequence []int) (SequenceNode, error) {
	t.Lock()
	defer t.Unlock()
	_, exists := t.pathogens[id]
	if exists {
		return nil, fmt.Errorf("sequence tree "+IntKeyExists, id)
	}
	// Create new node in tree
	n := new(sequenceNode)
	n.parent = nil
	n.children = []SequenceNode{}
	n.hits = make(map[int]int)
	n.sequence = sequence
	n.numStates = make(map[int]int)
	for _, s := range n.sequence {
		n.numStates[s]++
	}
	// Assign ID
	n.id = id

	// Add to pathogen and root maps
	t.pathogens[id] = n
	t.roots[id] = n
	return n, nil
}
