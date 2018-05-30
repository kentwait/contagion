package contagiongo

import (
	"sync"

	"github.com/segmentio/ksuid"
)

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
	defer n.Unlock()
	n.Lock()
	n.children = append(n.children, child)
}

func (n *sequenceNode) Sequence() []int {
	return n.sequence
}

func (n *sequenceNode) Fitness(f SequenceFitness) float64 {
	id := f.ID()
	fitness, ok := n.fitness[id]
	if !ok {
		fitness, _ := f.Fitness(n.sequence...)
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
func NewSequenceTree(rootSequence []int) (SequenceTree, error) {
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
	// Count the initial number of states across all sites
	for _, s := range n.sequence {
		n.stateCounts[s]++
	}
	// Create new tree and initialize maps
	tree := new(sequenceTree)
	tree.pathogens = make(map[ksuid.KSUID]*sequenceNode)
	tree.roots = make(map[ksuid.KSUID]*sequenceNode)
	// Add the node to the tree and make it a root node
	tree.pathogens[n.UID()] = n
	tree.roots[n.UID()] = n
	return tree, nil
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
	// Copy state counts
	n.stateCounts = make(map[int]int)
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
	// Create fresh count of states
	n.stateCounts = make(map[int]int)
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
