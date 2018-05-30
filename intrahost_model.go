package contagiongo

// IntrahostModel is an interface for any type of intrahost model.
type IntrahostModel interface {
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
