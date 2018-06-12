package contagiongo

import (
	"fmt"
	"sync"
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
	// PickPathogens returns a random list of pathogens from the
	// current host.
	// Returns nil if no pathogen exists.
	PickPathogens(n int) []GenotypeNode
	// Pathogens returns a list of all pathogens present in the host.
	// This elements of the list are pointers to GenotypeNodes.
	Pathogens() []GenotypeNode
	// PathogenPopSize returns the number of pathogens inside the host.
	PathogenPopSize() int
	// AddPathogens appends a pathogen to the pathogen space of the host.
	// Returns the new pathogen population size.
	AddPathogens(p ...GenotypeNode) int
	// RemoveAllPathogens removes all the pathogens from the host.
	// Internally, this removes all the pointers that refer to GenotypeNodes.
	RemoveAllPathogens()
	// SetModel associates the current host to a given intrahost model.
	// The intrahost model governs intrahost processes by specifying the
	// mutation, replication, recombination, and infection modes and parameters
	// to be used.
	SetIntrahostModel(intrahostModel IntrahostModel) error
	SetFitnessModel(fitnessModel FitnessModel) error
	SetTransmissionModel(transmissionModel TransmissionModel) error

	GetIntrahostModel() IntrahostModel
	GetFitnessModel() FitnessModel
	GetTransmissionModel() TransmissionModel
}

type sequenceHost struct {
	sync.RWMutex
	IntrahostModel
	FitnessModel
	TransmissionModel

	id             int
	typeID         int
	pathogens      map[int]GenotypeNode
	lastPathogenID int
}

// EmptySequenceHost creates a new host without an intrahost model and
// no pathogens.
func EmptySequenceHost(ids ...int) Host {
	h := new(sequenceHost)
	h.id = ids[0]
	h.typeID = 0 // default
	if len(ids) > 1 {
		h.typeID = ids[1]
	}
	h.pathogens = make(map[int]GenotypeNode)
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

func (h *sequenceHost) PickPathogens(n int) []GenotypeNode {
	h.RLock()
	defer h.RUnlock()
	pathogens := make([]GenotypeNode, n)
	i := 0
	for _, node := range h.pathogens {
		pathogens[i] = node
		i++
		if i >= 3 {
			break
		}
	}
	return pathogens
}

func (h *sequenceHost) Pathogens() []GenotypeNode {
	h.RLock()
	defer h.RUnlock()
	pathogens := make([]GenotypeNode, len(h.pathogens))
	i := 0
	for _, p := range h.pathogens {
		pathogens[i] = p
		i++
	}
	return pathogens
}

func (h *sequenceHost) PathogenPopSize() int {
	return len(h.pathogens)
}

func (h *sequenceHost) AddPathogens(p ...GenotypeNode) int {
	h.Lock()
	defer h.Unlock()
	for _, node := range p {
		h.lastPathogenID++
		h.pathogens[h.lastPathogenID] = node
	}
	return len(h.pathogens)
}

func (h *sequenceHost) RemoveAllPathogens() {
	h.Lock()
	defer h.Unlock()
	for i := range h.pathogens {
		h.pathogens[i] = nil
	}
	h.pathogens = make(map[int]GenotypeNode)
}

func (h *sequenceHost) SetIntrahostModel(model IntrahostModel) error {
	if h.IntrahostModel != nil {
		return fmt.Errorf("intrahost "+ModelExistsError, h.IntrahostModel.ModelName(), h.IntrahostModel.ModelID())
	}
	h.IntrahostModel = model
	return nil
}

func (h *sequenceHost) SetFitnessModel(model FitnessModel) error {
	if h.FitnessModel != nil {
		return fmt.Errorf("fitness "+ModelExistsError, h.FitnessModel.ModelName(), h.FitnessModel.ModelID())
	}
	h.FitnessModel = model
	return nil
}

func (h *sequenceHost) SetTransmissionModel(model TransmissionModel) error {
	if h.TransmissionModel != nil {
		return fmt.Errorf("transmission "+ModelExistsError, h.TransmissionModel.ModelName(), h.TransmissionModel.ModelID())
	}
	h.TransmissionModel = model
	return nil
}

func (h *sequenceHost) GetIntrahostModel() IntrahostModel {
	return h.IntrahostModel
}
func (h *sequenceHost) GetFitnessModel() FitnessModel {
	return h.FitnessModel
}
func (h *sequenceHost) GetTransmissionModel() TransmissionModel {
	return h.TransmissionModel
}
