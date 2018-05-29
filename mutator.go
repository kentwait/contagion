package contagiongo

import (
	"fmt"
)

// Mutator is an interface for mutation rate related values.
type Mutator interface {
	// MutatorDescription returns a short description about the mutator type
	MutatorDescription() string

	// States returns the list of possible states
	States() []int

	// MutationRate returns the mutation rate for the mutator
	MutationRate() float64

	// TransitionMatrix returns the conditioned mutation rate matrix for the mutator
	TransitionMatrix() [][]float64

	// TransitionProbs returns the conditioned transition probabilities for a given
	// state
	TransitionProbs(char int) []float64

	// RecombinationRate returns the recombination rate for the mutator
	RecombinationRate() float64
}

type mutator struct {
	desc   string
	states []int
	mu     float64
	matrix [][]float64
}

func (m *mutator) MutatorDescription() string {
	return m.desc
}

func (m *mutator) States() []int {
	return m.states
}

func (m *mutator) MutationRate() float64 {
	return m.mu
}
func (m *mutator) TransitionMatrix() [][]float64 {
	return m.matrix
}

func (m *mutator) TransitionProbs(i int) []float64 {
	return m.matrix[i]
}

func (m *mutator) RecombinationRate() float64 {
	return 0
}

// NewUniformRateMutator creates a new Mutator based on a uniform single
// mutation rate for all states.
func NewUniformRateMutator(mu float64, states int, desc ...string) (Mutator, error) {
	m := new(mutator)
	if mu < 0 {
		return nil, fmt.Errorf(InvalidFloatParameterError, "mutation rate", mu, "mu < 0")
	}
	m.mu = mu

	if len(desc) > 0 {
		m.desc = desc[0]
	}

	// Initialize states
	if states < 2 {
		return nil, fmt.Errorf(InvalidIntParameterError, "number of states", states, "s < 2")
	}
	m.states = make([]int, states)

	for i := 0; i < states; i++ {
		m.states[i] = i
	}

	// Initialize uniform rate matrix
	count := states - 1
	p := 1.0 / float64(count)
	m.matrix = make([][]float64, states)
	for i := 0; i < count; i++ {
		m.matrix[i] = make([]float64, states)
		for j := 0; j < count; j++ {
			if i != j {
				m.matrix[i][j] = p
			} else {
				m.matrix[i][j] = 0
			}
		}
	}
	return m, nil
}

// NewRateMatrixMutator creates a new Mutator where the mutation rate
// of states are unequal.
func NewRateMatrixMutator(mu float64, rateMatrix [][]float64, desc ...string) (Mutator, error) {
	m := new(mutator)
	if mu < 0 {
		return nil, fmt.Errorf(InvalidFloatParameterError, "mutation rate", mu, "mu < 0")
	}
	m.mu = mu

	if len(desc) > 0 {
		m.desc = desc[0]
	}
	// Test if rate matrix is not empty and is square
	if len(rateMatrix) < 2 {
		return nil, fmt.Errorf(InvalidIntParameterError, "number of states", len(rateMatrix), "s < 2")
	}
	for _, row := range rateMatrix {
		if len(row) != len(rateMatrix) {
			return nil, fmt.Errorf(UnequalIntParameterError, "number of rates", len(rateMatrix), len(row))
		}
	}
	// Initialize states
	m.states = make([]int, len(rateMatrix))
	for i := 0; i < len(rateMatrix); i++ {
		m.states[i] = i
	}
	// Copy rate matrix
	m.matrix = make([][]float64, len(rateMatrix))
	for i, row := range rateMatrix {
		// Check if row sums to 1.0
		var total float64
		for _, v := range row {
			total += v
		}
		if total < 0.999 || total > 1.001 {
			return nil, fmt.Errorf("row values must sum up to 1.0")
		}

		m.matrix[i] = make([]float64, len(row))
		for j, p := range row {
			if i == j && p != 0 {
				return nil, fmt.Errorf("diagonal values must be 0")
			}
			m.matrix[i][j] = p
		}
	}
	return m, nil
}
