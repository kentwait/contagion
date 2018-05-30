package contagiongo

import (
	"fmt"
	"sort"
)

// Host encapsulates pathogens together and ties its evolution to a particular
// set of parameters given by its assigned host type.
type Host interface {
	// ID returns the unique ID of the host.
	ID() int
	// TypeID returns the ID representing the host's type in a multi-host
	// simulation. Generally, the host type ID is used to identify hosts
	// belonging to the same group that share the same properties.
	TypeID() int
	// Pathogen returns one pathogen from the current host based on
	// its given position in the list of pathogens.
	// Returns nil if no pathogen exists in the specified position.
	Pathogen(i int) GenotypeNode
	// Pathogens returns a list of all pathogens present in the host.
	// This elements of the list are pointers to GenotypeNodes.
	Pathogens() []GenotypeNode
	// PathogenPopSize returns the number of pathogens inside the host.
	PathogenPopSize() int
	// AddPathogen appends a pathogen to the pathogen space of the host.
	// Returns the new pathogen population size.
	AddPathogen(p GenotypeNode) int
	// RemovePathogens removes pathogens based on the list of positions given.
	// Returns the number of pathogens remaining and any errors encountered.
	RemovePathogens(ids ...int) (n int, err error)
	// RemoveAllPathogens removes all the pathogens from the host.
	// Internally, this removes all the pointers that refer to GenotypeNodes.
	RemoveAllPathogens()
	// DecrementTimer decreases the internal timer by 1.
	DecrementTimer()
	// SetModel associates the current host to a given intrahost model.
	// The intrahost model governs intrahost processes by specifying the
	// mutation, replication, recombination, and infection modes and parameters
	// to be used.
	SetIntrahostModel(intrahostModel IntrahostModel) error
}

type sequenceHost struct {
	IntrahostModel
	FitnessModel

	id            int
	typeID        int
	internalTimer int
	pathogens     []GenotypeNode
}

// NewEmptySequenceHost creates a new host without an intrahost model and
// no pathogens.
func NewEmptySequenceHost(ids ...int) Host {
	h := new(sequenceHost)
	h.id = ids[0]
	h.typeID = 0 // default
	if len(ids) > 1 {
		h.typeID = ids[1]
	}
	h.internalTimer = 0
	h.pathogens = []GenotypeNode{}
	h.IntrahostModel = nil
	h.FitnessModel = nil
	return h
}

func (h *sequenceHost) ID() int {
	return h.id
}

func (h *sequenceHost) TypeID() int {
	return h.typeID
}

func (h *sequenceHost) Pathogen(i int) GenotypeNode {
	return h.pathogens[i]
}

func (h *sequenceHost) Pathogens() []GenotypeNode {
	// pathogens := make([]GenotypeNode, len(h.pathogens))
	// for i, p := range h.pathogens {
	// 	pathogens[i] = p
	// }
	// return pathogens
	return h.pathogens
}

func (h *sequenceHost) PathogenPopSize() int {
	return len(h.pathogens)
}

func (h *sequenceHost) AddPathogen(p GenotypeNode) int {
	h.pathogens = append(h.pathogens, p)
	return len(h.pathogens)
}

func (h *sequenceHost) RemovePathogens(ids ...int) (n int, err error) {
	sort.Ints(ids)
	// Check if the largest ID is less than the number of pathogens
	lastID := ids[len(ids)-1]
	if len(ids) > 0 && lastID >= len(h.pathogens) {
		return 0, fmt.Errorf("pathogen "+IntKeyNotFoundError, lastID)
	}
	for offset, i := range ids {
		pos := i - offset
		lastID := len(h.pathogens) - 1
		// Remove hit from list
		copy(h.pathogens[pos:], h.pathogens[pos+1:])
		h.pathogens[lastID] = nil // or the zero value of T
		h.pathogens = h.pathogens[:lastID]
	}
	return len(h.pathogens), nil
}

func (h *sequenceHost) RemoveAllPathogens() {
	for i := range h.pathogens {
		h.pathogens[i] = nil
	}
	h.pathogens = h.pathogens[:0]
}

func (h *sequenceHost) DecrementTimer() {
	h.internalTimer--
}

func (h *sequenceHost) SetIntrahostModel(intrahostModel IntrahostModel) error {
	if h.IntrahostModel != nil {
		return fmt.Errorf(IntrahostModelExistsError, h.IntrahostModel.ModelName(), h.IntrahostModel.ModelID())
	}
	h.IntrahostModel = intrahostModel
	return nil
}
