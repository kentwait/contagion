package contagiongo

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"
)

// MotifModel is a type of FitnessModel where the fitness of a sequence
// depends on the presence of the particular motifs.
type MotifModel interface {
	// ID returns the ID for this fitness model.
	ModelID() int
	// Name returns the name for this fitness model.
	ModelName() string
	SetModelID(id int)
	SetModelName(name string)
	// ComputeFitness returns the corresponding fitness value given
	// a set of sequences as integers.
	ComputeFitness(chars ...uint8) (fitness float64, err error)

	AddMotif(sequence []uint8, pos []int, value float64) error
}

// motif is a struct that describes an int-coded sequence motif.
type motif struct {
	motif   []uint8
	pos     []int
	fitness float64
}

// newMotif creates a new Motif.
func newMotif(sequence []uint8, pos []int, value float64) *motif {
	m := motif{
		motif:   sequence,
		pos:     pos,
		fitness: value,
	}
	return &m
}

// InSequence tells if the motif is present or absent in the given sequence.
func (m *motif) InSequence(seq []uint8) bool {
	if len(seq) == 0 || len(m.motif) == 0 {
		return false
	}
	for i, motifChar := range m.motif {
		pos := m.pos[i]
		seqChar := seq[pos]
		if motifChar != seqChar {
			return false
		}
	}
	return true
}

// UpdateFitness changes the value of the fitness value of the motif.
func (m *motif) UpdateFitness(v float64) {
	m.fitness = v
}

// MotifID returns the sequence and its corresponding positions interleaved
// as a string.
func (m *motif) MotifID() string {
	var b bytes.Buffer
	for i := 0; i < len(m.motif); i++ {
		s := int(m.motif[i])
		pos := m.pos[i]
		b.WriteString(strconv.Itoa(s))
		b.WriteString(strconv.Itoa(pos))
	}
	return b.String()
}

type motifModel struct {
	sync.RWMutex
	modelMetadata
	motifs         map[string]*motif
	positions      map[int]bool
	defaultFitness float64
}

// EmptyMotifModel returns a new motif model without any registered motifs.
func EmptyMotifModel(id int, name string) MotifModel {
	m := new(motifModel)
	m.id = id
	m.name = name
	m.motifs = make(map[string]*motif)
	m.positions = make(map[int]bool)
	return m
}

// ComputeFitness returns the corresponding fitness value given
// a set of sequences as integers.
func (m *motifModel) ComputeFitness(chars ...uint8) (fitness float64, err error) {
	m.RLock()
	defer m.RUnlock()
	var decFitness float64
	for _, motif := range m.motifs {
		if motif.InSequence(chars) {
			decFitness += motif.fitness
		}
	}
	return decFitness, nil
}

// AddMotif adds a new motif to the motif model. Returns an error if
// the given positions overlap with existing motifs in the model.
func (m *motifModel) AddMotif(sequence []uint8, pos []int, value float64) error {
	// Check positions
	m.RLock()
	for _, i := range pos {
		if _, exists := m.positions[i]; exists {
			m.RUnlock()
			return fmt.Errorf("site %d is already considered by another motif")
		}
	}
	m.RUnlock()

	m.Lock()
	defer m.Unlock()
	newMotif := newMotif(sequence, pos, value)
	if _, exists := m.motifs[newMotif.MotifID()]; exists {
		return fmt.Errorf("motif already exists")
	}
	m.motifs[newMotif.MotifID()] = newMotif
	for _, i := range pos {
		m.positions[i] = true
	}
	return nil
}
