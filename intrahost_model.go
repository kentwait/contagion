package contagiongo

import "math"

// IntrahostModel is an interface for any type of intrahost model.
type IntrahostModel interface {
	ModelID() int
	ModelName() string
	SetModelID(id int)
	SetModelName(name string)

	// Mutation

	// MutationRate returns the mutation rate for this model.
	MutationRate() float64
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
	// ReplicationMethod returns whether relative, or absolute is used.
	ReplicationMethod() string

	// Recombination

	// RecombinationRate returns the recombination rate for this model.
	RecombinationRate() float64

	// Infection duration

	// StatusDuration
	StatusDuration(status int) int
}

// ConstantPopModel models a constant pathogen population size within the host.
type ConstantPopModel struct {
	intrahostMetadata
	mutationParams
	recombinationParams
	durationParams
	constantIntrahostPopModel
}

// BevertonHoltThresholdPopModel uses the Beverton-Holt population model
// modified to have a constant threshold population size.
type BevertonHoltThresholdPopModel struct {
	intrahostMetadata
	mutationParams
	recombinationParams
	durationParams
	bhtIntrahostPopModel
}

// FitnessDependentPopModel does not use a population model to determine the
// population of pathogens. Instead population size is dependent on fitness
// which is implemented outside of this model.
// The NextPathogenPopSize method for this model always returns -1 regardless
// of the input value.
type FitnessDependentPopModel struct {
	intrahostMetadata
	mutationParams
	recombinationParams
	durationParams
	fitnessIntrahostPopModel
}

type intrahostMetadata struct {
	id   int
	name string
}

func (meta *intrahostMetadata) ModelID() int {
	return meta.id
}

func (meta *intrahostMetadata) ModelName() string {
	return meta.name
}

func (meta *intrahostMetadata) SetModelID(id int) {
	meta.id = id
}
func (meta *intrahostMetadata) SetModelName(name string) {
	meta.name = name
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

type recombinationParams struct {
	recombinationRate float64
}

func (params *recombinationParams) RecombinationRate() float64 {
	return params.recombinationRate
}

type durationParams struct {
	statusDuration map[int]int
}

func (params *durationParams) StatusDuration(status int) int {
	return params.statusDuration[status]
}

type constantIntrahostPopModel struct {
	popSize int
}

func (m *constantIntrahostPopModel) MaxPathogenPopSize() int {
	return m.popSize
}

func (m *constantIntrahostPopModel) NextPathogenPopSize(n int) int {
	return m.popSize
}

func (m *constantIntrahostPopModel) ReplicationMethod() string {
	return "relative"
}

type bhtIntrahostPopModel struct {
	maxPopSize int
	growthRate float64
}

func (m *bhtIntrahostPopModel) MaxPathogenPopSize() int {
	return m.maxPopSize
}

func (m *bhtIntrahostPopModel) NextPathogenPopSize(n int) int {
	n64 := float64(n)
	k64 := float64(m.maxPopSize)
	res := (m.growthRate * n64 * k64) / (k64 + ((m.growthRate - 1.0) * n64))
	roundedRes := int(math.Ceil(res))
	if m.maxPopSize > roundedRes {
		return roundedRes
	}
	return m.maxPopSize
}

func (m *bhtIntrahostPopModel) GrowthRate() float64 {
	return m.growthRate
}

func (m *bhtIntrahostPopModel) ReplicationMethod() string {
	return "relative"
}

type fitnessIntrahostPopModel struct {
	maxPopSize int
}

func (m *fitnessIntrahostPopModel) MaxPathogenPopSize() int {
	return m.maxPopSize
}

func (m *fitnessIntrahostPopModel) NextPathogenPopSize(n int) int {
	// Not applicable
	return -1
}

func (m *fitnessIntrahostPopModel) ReplicationMethod() string {
	return "absolute"
}
