package contagiongo

import (
	"fmt"
)

// SequenceFitness represents a fitness landscape where the sequence
// is used to look up the fitness value.
type SequenceFitness interface {
	// Fitness returns the corresponding fitness value given
	// a set of sequences as integers.
	Fitness(chars ...int) (fitness float64, err error)
}

// FitnessMatrix
type FitnessMatrix interface {
	Fitness(chars ...int) (fitness float64, err error)
	SiteFitness(position, state int) (fitness float64, err error)
	Log() bool
}

type multiplicativeFM struct {
	matrix map[int]map[int]float64
}

// NewMultiplicativeFM create a new multiplicative fitness matrix using
// a map of maps.
// Assumes that the values are in log form.
func NewMultiplicativeFM(matrix map[int]map[int]float64) FitnessMatrix {
	// Copy map of maps
	fm := new(multiplicativeFM)
	fm.matrix = make(map[int]map[int]float64)
	for k1, row := range matrix {
		fm.matrix[k1] = make(map[int]float64)
		for k2, v := range row {
			fm.matrix[k1][k2] = v
		}
	}
	return fm
}

func (fm *multiplicativeFM) Fitness(chars ...int) (fitness float64, err error) {
	// Assume coords are sequence of ints representing a sequence
	// Matrix values are in log
	// Returns log fitness total
	if len(chars) < 0 {
		return 0, fmt.Errorf("multiplicativeFM Fitness requires at least 1 coordinate (1 site character as an integer")
	}
	var logFitness float64
	for i, v := range chars {
		logFitness += fm.matrix[i][v]
	}
	return logFitness, nil
}

func (fm *multiplicativeFM) SiteFitness(position, state int) (fitness float64, err error) {
	return fm.matrix[position][state], nil
}

func (fm *multiplicativeFM) Log() bool {
	return true
}

type additiveFM struct {
	matrix map[int]map[int]float64
}

// NewAdditiveFM create a new additive fitness matrix using a map of maps.
// Assumes that the values are in decimal form.
func NewAdditiveFM(matrix map[int]map[int]float64) FitnessMatrix {
	// Copy map of maps
	fm := new(additiveFM)
	fm.matrix = make(map[int]map[int]float64)
	for k1, row := range matrix {
		fm.matrix[k1] = make(map[int]float64)
		for k2, v := range row {
			fm.matrix[k1][k2] = v
		}
	}
	return fm
}

func (fm *additiveFM) Fitness(chars ...int) (fitness float64, err error) {
	// Assume coords are sequence of ints representing a sequence
	// Matrix values are in decimal
	// Returns decimal fitness total
	if len(chars) < 0 {
		return 0, fmt.Errorf("additiveFM Fitness requires at least 1 coordinate (1 site character as an integer")
	}
	var decFitness float64
	for i, v := range chars {
		decFitness += fm.matrix[i][v]
	}
	return decFitness, nil
}

func (fm *additiveFM) SiteFitness(position, state int) (fitness float64, err error) {
	return fm.matrix[position][state], nil
}

func (fm *additiveFM) Log() bool {
	return false
}
