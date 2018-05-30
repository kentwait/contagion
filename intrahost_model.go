package contagiongo

import "math"

// IntrahostModel is an interface for any type of intrahost model.
type IntrahostModel interface {
	ID() int
	Name() string

	// Mutation

	// MutationRate returns the mutation rate for this model.
	MutationRate() float64
	// States returns the list of possible states.
	States() []int
	// TransitionMatrix returns the conditioned mutation rate matrix
	// for this model.
	TransitionMatrix() [][]float64
	// TransitionProbs returns the conditioned transition probabilities
	// for the given state.
	TransitionProbs(char int) []float64

	// Replication

	// MaxPathogenPopSize returns the maximum number of pathogens allowed within
	// a single host of this particular host type.
	MaxPathogenPopSize() int
	// NextPathogenPopSize returns the pathogen population size for the next
	// generation of pathogens given the current population size.
	// This is used in conjunction with a population model under
	// relative fitness.
	NextPathogenPopSize(n int) int

	// Recombination

	// RecombinationRate returns the recombination rate for this model.
	RecombinationRate() float64
}

type mutationParams struct {
	mutationRate     float64
	transitionMatrix [][]float64
}

func (params *mutationParams) MutationRate() float64 {
	return params.mutationRate
}

func (params *mutationParams) TransitionMatrix() [][]float64 {
	return params.transitionMatrix
}

func (params *mutationParams) TransitionProbs(char int) []float64 {
	return params.transitionMatrix[char]
}

func (params *mutationParams) States() []int {
	if params.transitionMatrix == nil {
		return []int{}
	}
	states := make([]int, len(params.transitionMatrix))
	for i := range params.transitionMatrix {
		states[i] = i
	}
	return states
}

type recombinationParams struct {
	recombinationRate float64
}

func (params *recombinationParams) RecombinationRate() float64 {
	return params.recombinationRate
}

type constantPopModel struct {
	popSize int
}

func (m *constantPopModel) MaxPathogenPopSize() int {
	return m.popSize
}

func (m *constantPopModel) NextPathogenPopSize(n int) int {
	return m.popSize
}

type bhtPopModel struct {
	maxPopSize int
	growthRate float64
}

func (m *bhtPopModel) MaxPathogenPopSize() int {
	return m.maxPopSize
}

func (m *bhtPopModel) NextPathogenPopSize(n int) int {
	n64 := float64(n)
	k64 := float64(m.maxPopSize)
	res := (m.growthRate * n64 * k64) / (k64 + ((m.growthRate - 1.0) * n64))
	roundedRes := int(math.Ceil(res))
	if m.maxPopSize > roundedRes {
		return roundedRes
	}
	return m.maxPopSize
}

func (m *bhtPopModel) GrowthRate() float64 {
	return m.growthRate
}

type fitnessPopModel struct {
	maxPopSize int
}

func (m *fitnessPopModel) MaxPathogenPopSize() int {
	return m.maxPopSize
}

func (m *fitnessPopModel) NextPathogenPopSize(n int) int {
	// Not applicable
	return -1
}
