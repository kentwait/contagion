package contagiongo

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
