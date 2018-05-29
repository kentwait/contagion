package contagiongo

import (
	"fmt"
	"sort"
	"strings"
)

// Host encapsulates pathogens together and ties its evolution to a particular
// set of parameters given by its assigned host type.
type Host interface {
	// HostID returns the unique ID of the host.
	HostID() int

	HostTypeID() int

	// Pathogen returns one pathogen from a specific host. Returns nil
	// if no pathogen exists in the specified position.
	Pathogen(i int) interface{}

	// Pathogens returns all the pathogens present in the host.
	Pathogens() []interface{}

	// PathogenPopSize returns the number of pathogens inside the host.
	PathogenPopSize() int

	// AddPathogen appends a pathogen to the pathogen space of the host.
	// Returns the new pathogen population size.
	AddPathogen(p interface{}) int

	// RemovePathogensByID removes pathogens based on the list of IDs given.
	// Returns the number of pathogens remaining and any errors encountered.
	RemovePathogensByID(ids ...int) (n int, err error)

	// // Returns a deep copy of the host, including all the pathogens
	// // present.
	// Clone() Host

	// ClearPathogens removes all the pathogens from the pathogen space
	// of the host, leaving an empty list.
	ClearPathogens()

	// DecrementTimer decreases the internal timer by 1.
	DecrementTimer()

	// SetModel sets the type of spreader, Replicator, mutator, and fitness to
	// be used by the host.
	SetModel(name string, model interface{}) error
}

type EpidemicHost interface {
	Host
	Spreader
	Replicator
	Mutator
}

// SequenceHost is type of host that implements the Host interface and embeds
// a Spreader, Replicator, Mutator, and SequenceFitness.
// The SequenceHost is identified by its hostID and can be accessed using the
// HostID method.
// Pathogens associated with the host are listed under the pathogen property
// and can be accessed using the Pathogen or Pathogens methods.
// Pathogens can be added into the host using the AddPathogen method and can
// be removed by pathogenID using the RemovePathogensByID method. All pathogens
// can be removed fromt the host using the ClearPathogens method.
type SequenceHost struct {
	Spreader
	Replicator
	Mutator
	SequenceFitness

	hostID        int
	pathogens     []SequenceNode
	internalTimer int
}

// NewEmptySequenceHost creates a new host without spreader, replicator, and
// mutator features. It also has no status time intervals and
// lacks any pathogens.
func NewEmptySequenceHost(id int) *SequenceHost {
	h := new(SequenceHost)
	h.hostID = id
	return h
}

// HostID returns the unique ID of the host.
func (h *SequenceHost) HostID() int {
	return h.hostID
}

// Pathogen returns one pathogen from a specific host. Returns nil
// if no pathogen exists in the specified position.
func (h *SequenceHost) Pathogen(i int) interface{} {
	return h.pathogens[i]
}

// Pathogens returns all the pathogens present in the host.
func (h *SequenceHost) Pathogens() []interface{} {
	var pathogens []interface{}
	for _, p := range h.pathogens {
		pathogens = append(pathogens, p)
	}
	return pathogens
}

// PathogenPopSize returns the number of pathogens inside the host.
func (h *SequenceHost) PathogenPopSize() int {
	return len(h.pathogens)
}

// AddPathogen appends a pathogen to the pathogen space of the host.
// Returns the new pathogen population size.
func (h *SequenceHost) AddPathogen(p interface{}) int {
	h.pathogens = append(h.pathogens, p.(SequenceNode))
	return len(h.pathogens)
}

// RemovePathogensByID removes pathogens based on the list of IDs given.
// Returns the number of pathogens remaining and any errors encountered.
func (h *SequenceHost) RemovePathogensByID(ids ...int) (n int, err error) {
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

// ClearPathogens removes all the pathogens from the pathogen space
// of the host, leaving an empty list.
func (h *SequenceHost) ClearPathogens() {
	h.pathogens = []SequenceNode{}
}

// DecrementTimer decreases the internal timer by 1.
func (h *SequenceHost) DecrementTimer() {
	h.internalTimer--
}

// SetModel sets the type of spreader, Replicator, mutator, and fitness to
// be used by the host.
func (h *SequenceHost) SetModel(keyword string, model interface{}) error {
	switch strings.ToLower(keyword) {
	case "spreader":
		h.Spreader = model.(Spreader)
	case "Replicator":
		h.Replicator = model.(Replicator)
	case "mutator":
		h.Mutator = model.(Mutator)
	case "landscape":
		h.SequenceFitness = model.(SequenceFitness)
	default:
		return fmt.Errorf("unknown keyword")
	}
	return nil
}
