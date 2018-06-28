package contagiongo

import (
	"bytes"
	"strconv"
	"sync"

	"github.com/pkg/errors"
)

// FitnessModel represents a general method to determine the fitness value
// associated to a particular genotype.
type FitnessModel interface {
	// ID returns the ID for this fitness model.
	ModelID() int
	// Name returns the name for this fitness model.
	ModelName() string
	SetModelID(id int)
	SetModelName(name string)
	// ComputeFitness returns the corresponding fitness value given
	// a set of sequences as integers.
	ComputeFitness(chars ...uint8) (fitness float64, err error)
}

// FitnessMatrix is a type of FitnessModel where the fitness of each individual
// character at every site is specified.
type FitnessMatrix interface {
	// ID returns the ID for this fitness model.
	ModelID() int
	// Name returns the name for this fitness model.
	ModelName() string
	SetModelID(id int)
	SetModelName(name string)
	// ComputeFitness returns the corresponding fitness value given
	// a set of sequences as integers.
	ComputeFitness(chars ...uint8) (fitness float64, err error)
	// SiteFitness returns the fitness value associated for a particular
	// character at the given site.
	SiteCharFitness(position int, state uint8) (fitness float64, err error)
	// Log tells whether the fitness values are decimal or log.
	// Usually fitness is in log.
	Log() bool
}

// multiplicativeFM is a multiplicative fitness matrix that computes
// the fitness of a pathogen by getting the product of each site fitness.
// Values are assumed to be in log space such that the product is the
// sum of log-space fitness contributions.
type multiplicativeFM struct {
	modelMetadata
	matrix map[int]map[uint8]float64
}

// NewMultiplicativeFM create a new multiplicative fitness matrix using
// a map of maps. Assumes that the values are in log form.
func NewMultiplicativeFM(id int, name string, matrix map[int]map[uint8]float64) FitnessMatrix {
	// Copy map of maps
	fm := new(multiplicativeFM)
	fm.id = id
	fm.name = name
	fm.matrix = make(map[int]map[uint8]float64)
	// Each row lists the fitness of alternative characters for that site
	for k1, row := range matrix {
		fm.matrix[k1] = make(map[uint8]float64)
		for k2, v := range row {
			fm.matrix[k1][uint8(k2)] = v
		}
	}
	return fm
}

// ComputeFitness takes a sequence of characters and gets the log-sum of
// their corresponding fitness values. The log sum of log fitness values
// is equivalent to the product of base 10 fitnesses.
func (fm *multiplicativeFM) ComputeFitness(chars ...uint8) (fitness float64, err error) {
	// Assume coords are sequence of ints representing a sequence
	// Matrix values are in log
	// Returns log fitness total
	if len(chars) < 0 {
		return 0, errors.Wrap(ZeroItemsError(), "computing multiplicative fitness failed")
	}
	var logFitness float64
	for i, v := range chars {
		logFitness += fm.matrix[i][v]
	}
	return logFitness, nil
}

func (fm *multiplicativeFM) SiteCharFitness(position int, state uint8) (fitness float64, err error) {
	return fm.matrix[position][state], nil
}

func (fm *multiplicativeFM) Log() bool {
	return true
}

// additiveFM is an additive fitness matrix that computes the fitness
// of a pathogen by getting the sum of each site fitness. Values are
// assumed to be in base 10 form.
type additiveFM struct {
	modelMetadata
	matrix map[int]map[uint8]float64
}

// NewAdditiveFM create a new additive fitness matrix using a map of maps.
// Assumes that the values are in decimal form.
func NewAdditiveFM(id int, name string, matrix map[int]map[uint8]float64) FitnessMatrix {
	// Copy map of maps
	fm := new(additiveFM)
	fm.id = id
	fm.name = name
	fm.matrix = make(map[int]map[uint8]float64)
	for k1, row := range matrix {
		fm.matrix[k1] = make(map[uint8]float64)
		for k2, v := range row {
			fm.matrix[k1][uint8(k2)] = v
		}
	}
	return fm
}

func (fm *additiveFM) ComputeFitness(chars ...uint8) (fitness float64, err error) {
	// Assume coords are sequence of ints representing a sequence
	// Matrix values are in decimal
	// Returns decimal fitness total
	if len(chars) < 0 {
		return 0, errors.Wrap(ZeroItemsError(), "computing additive fitness failed")
	}
	var decFitness float64
	for i, v := range chars {
		decFitness += fm.matrix[i][uint8(v)]
	}
	if decFitness < 0 {
		decFitness = 0.0
	}
	return decFitness, nil
}

func (fm *additiveFM) SiteCharFitness(position int, state uint8) (fitness float64, err error) {
	return fm.matrix[position][state], nil
}

func (fm *additiveFM) Log() bool {
	return false
}

// NeutralMultiplicativeFM returns a multiplicative fitness matrix
// where all the values are 0 (ln 1) such that all changes have no effect
// and are therefore neutral.
func NeutralMultiplicativeFM(id int, name string, sites, alleles int) FitnessMatrix {
	fm := new(multiplicativeFM)
	fm.id = id
	fm.name = name
	fm.matrix = make(map[int]map[uint8]float64)
	for i := 0; i < sites; i++ {
		fm.matrix[i] = make(map[uint8]float64)
		for j := 0; j < alleles; j++ {
			fm.matrix[i][uint8(j)] = 0.0
		}
	}
	return fm
}

// NeutralAdditiveFM returns a additive fitness matrix where the sum of
// all sites using any allele combination is equal to the growth rate.
func NeutralAdditiveFM(id int, name string, sites, alleles, growthRate int) FitnessMatrix {
	fm := new(additiveFM)
	fm.id = id
	fm.name = name
	fm.matrix = make(map[int]map[uint8]float64)
	for i := 0; i < sites; i++ {
		fm.matrix[i] = make(map[uint8]float64)
		for j := 0; j < alleles; j++ {
			fm.matrix[i][uint8(j)] = float64(growthRate) / float64(sites)
		}
	}
	return fm
}

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
			return OverlappingMotifError(i)
		}
	}
	m.RUnlock()

	m.Lock()
	defer m.Unlock()
	newMotif := newMotif(sequence, pos, value)
	newMotifID := newMotif.MotifID()
	if _, exists := m.motifs[newMotifID]; exists {
		return MotifExistsError(newMotifID)
	}
	m.motifs[newMotifID] = newMotif
	for _, i := range pos {
		m.positions[i] = true
	}
	return nil
}
