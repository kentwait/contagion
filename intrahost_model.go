package contagiongo

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
