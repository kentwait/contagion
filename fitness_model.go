package contagiongo

import "fmt"

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
	ComputeFitness(chars ...int) (fitness float64, err error)
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
	ComputeFitness(chars ...int) (fitness float64, err error)
	// SiteFitness returns the fitness value associated for a particular
	// character at the given site.
	SiteCharFitness(position, state int) (fitness float64, err error)
	// Log tells whether the fitness values are decimal or log.
	// Usually fitness is in log.
	Log() bool
}

type multiplicativeFM struct {
	modelMetadata
	matrix map[int]map[int]float64
}

// NewMultiplicativeFM create a new multiplicative fitness matrix using a map of maps.
// Assumes that the values are in log form.
func NewMultiplicativeFM(id int, name string, matrix map[int]map[int]float64) FitnessMatrix {
	// Copy map of maps
	fm := new(multiplicativeFM)
	fm.id = id
	fm.name = name
	fm.matrix = make(map[int]map[int]float64)
	// Each row lists the fitness of alternative characters for that site
	for k1, row := range matrix {
		fm.matrix[k1] = make(map[int]float64)
		for k2, v := range row {
			fm.matrix[k1][k2] = v
		}
	}
	return fm
}

func (fm *multiplicativeFM) ComputeFitness(chars ...int) (fitness float64, err error) {
	// Assume coords are sequence of ints representing a sequence
	// Matrix values are in log
	// Returns log fitness total
	if len(chars) < 0 {
		return 0, fmt.Errorf(ZeroItemsError)
	}
	var logFitness float64
	for i, v := range chars {
		logFitness += fm.matrix[i][v]
	}
	return logFitness, nil
}

func (fm *multiplicativeFM) SiteCharFitness(position, state int) (fitness float64, err error) {
	return fm.matrix[position][state], nil
}

func (fm *multiplicativeFM) Log() bool {
	return true
}

type additiveFM struct {
	modelMetadata
	matrix map[int]map[int]float64
}

// NewAdditiveFM create a new additive fitness matrix using a map of maps.
// Assumes that the values are in decimal form.
func NewAdditiveFM(id int, name string, matrix map[int]map[int]float64) FitnessMatrix {
	// Copy map of maps
	fm := new(additiveFM)
	fm.id = id
	fm.name = name
	fm.matrix = make(map[int]map[int]float64)
	for k1, row := range matrix {
		fm.matrix[k1] = make(map[int]float64)
		for k2, v := range row {
			fm.matrix[k1][k2] = v
		}
	}
	return fm
}

func (fm *additiveFM) ComputeFitness(chars ...int) (fitness float64, err error) {
	// Assume coords are sequence of ints representing a sequence
	// Matrix values are in decimal
	// Returns decimal fitness total
	if len(chars) < 0 {
		return 0, fmt.Errorf(ZeroItemsError)
	}
	var decFitness float64
	for i, v := range chars {
		decFitness += fm.matrix[i][v]
	}
	return decFitness, nil
}

func (fm *additiveFM) SiteCharFitness(position, state int) (fitness float64, err error) {
	return fm.matrix[position][state], nil
}

func (fm *additiveFM) Log() bool {
	return false
}

// NeutralMultiplicativeFM returns a multiplicative fitness matrix where all the values are 0 (ln 1)
func NeutralMultiplicativeFM(id int, name string, sites, alleles int) FitnessMatrix {
	fm := new(multiplicativeFM)
	fm.id = id
	fm.name = name
	fm.matrix = make(map[int]map[int]float64)
	for i := 0; i < sites; i++ {
		fm.matrix[i] = make(map[int]float64)
		for j := 0; j < alleles; j++ {
			fm.matrix[i][j] = 0.0
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
	fm.matrix = make(map[int]map[int]float64)
	for i := 0; i < sites; i++ {
		fm.matrix[i] = make(map[int]float64)
		for j := 0; j < alleles; j++ {
			fm.matrix[i][j] = float64(growthRate) / float64(sites)
		}
	}
	return fm
}
