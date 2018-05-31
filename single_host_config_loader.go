package contagiongo

import (
	"fmt"
	"strings"
)

// SingleHostConfig contains parameters to create a simulated infection
// in a single host.
type SingleHostConfig struct {
	NumGenerations uint `toml:"num_generations"`
	NumReplicates  uint `toml:"num_replicates"`

	ModelName         string      `toml:"model_name"`
	MutationRate      float64     `toml:"mutation_rate"`
	TransitionMatrix  [][]float64 `toml:"transition_matrix"`
	RecombinationRate float64     `toml:"recombination_rate"`
	ReplicationModel  string      `toml:"replication_model"` // constant, bht, fitness
	ConstantPopSize   int         `toml:"constant_pop_size"` // only for constant
	MaxPopSize        int         `toml:"max_pop_size"`      // only for bht and fitness
	GrowthRate        float64     `toml:"growth_rate"`       // only for bht

	FitnessModel     string `toml:"fitness_model"` // multiplicative, additive, additive_motif
	FitnessModelPath string `toml:"fitness_model_path"`

	PathogenSequencePath string `toml:"pathogen_sequence_path"` // fasta file for seeding infections

	LogFreq         uint   `toml:"log_freq"`
	PathogenLogPath string `toml:"pathogen_log_path"`

	validated bool
}

// Validate checks the validity of the configuration.
func (c *SingleHostConfig) Validate() error {
	// check keywords
	// replication_model
	switch strings.ToLower(c.ReplicationModel) {
	case "constant":
	case "bht":
	case "fitness":
	default:
		return fmt.Errorf(UnrecognizedKeywordError, c.ReplicationModel, "replication_model")
	}
	// fitness_model
	switch strings.ToLower(c.FitnessModel) {
	case "multiplicative":
	case "additive":
	case "additive_motif":
	default:
		return fmt.Errorf(UnrecognizedKeywordError, c.FitnessModel, "fitness_model")
	}
	c.validated = true
	return nil
}

// NewSimulation creates a new SingleHostSimulation simulation.
func (c *SingleHostConfig) NewSimulation() (Infection, error) {
	sim := new(singleHostSimulation)
	// Create empty tree
	sim.tree = EmptyGenotypeTree()
	// Create empty host
	host := NewEmptySequenceHost(0, 0)

	// Create IntrahostModel
	switch c.ReplicationModel {
	case "constant":
		model := new(ConstantPopModel)
		model.popSize = c.ConstantPopSize
		model.mutationRate = c.MutationRate
		model.recombinationRate = c.RecombinationRate
		model.transitionMatrix = make([][]float64, len(c.TransitionMatrix))
		for i := 0; i < len(c.TransitionMatrix); i++ {
			model.transitionMatrix[i] = make([]float64, len(c.TransitionMatrix))
			copy(model.transitionMatrix[i], c.TransitionMatrix[i])
		}
		host.SetIntrahostModel(model)
		sim.intrahostModel = model
	case "bht":
		model := new(BevertonHoltThresholdPopModel)
		model.maxPopSize = c.MaxPopSize
		model.growthRate = c.GrowthRate
		model.mutationRate = c.MutationRate
		model.recombinationRate = c.RecombinationRate
		model.transitionMatrix = make([][]float64, len(c.TransitionMatrix))
		for i := 0; i < len(c.TransitionMatrix); i++ {
			model.transitionMatrix[i] = make([]float64, len(c.TransitionMatrix))
			copy(model.transitionMatrix[i], c.TransitionMatrix[i])
		}
		host.SetIntrahostModel(model)
		sim.intrahostModel = model
	case "fitness":
		model := new(FitnessDependentPopModel)
		model.maxPopSize = c.MaxPopSize
		model.mutationRate = c.MutationRate
		model.recombinationRate = c.RecombinationRate
		model.transitionMatrix = make([][]float64, len(c.TransitionMatrix))
		for i := 0; i < len(c.TransitionMatrix); i++ {
			model.transitionMatrix[i] = make([]float64, len(c.TransitionMatrix))
			copy(model.transitionMatrix[i], c.TransitionMatrix[i])
		}
		host.SetIntrahostModel(model)
		sim.intrahostModel = model
	}

	// Create FitnessModel
	switch c.FitnessModel {
	case "multiplicative":
		matrix, err := LoadFitnessMatrix(c.FitnessModelPath)
		if err != nil {
			return nil, err
		}
		fm := NewMultiplicativeFM(0, "multiplicative", matrix)
		sim.fitnessModel = fm
	case "additive":
		matrix, err := LoadFitnessMatrix(c.FitnessModelPath)
		if err != nil {
			return nil, err
		}
		fm := NewAdditiveFM(0, "additive", matrix)
		sim.fitnessModel = fm
	case "additive_motif":
		return nil, fmt.Errorf("additive_motif not yet implemented")
	}

	// Parse fitness model file
	pathogenHostMap, err := LoadSequences(c.PathogenSequencePath)
	if err != nil {
		return nil, err
	}
	// Assumes that target host ID is 0
	// Adds sequences to the tree
	for _, sequence := range pathogenHostMap[0] {
		// Each starting sequence is a root node
		sim.tree.NewNode(sequence)
	}
	// Initialize durations
	sim.statusDuration = make(map[int]int)
	sim.statusDuration[InfectedStatusCode] = int(c.NumGenerations)
	// Initialize status
	if len(pathogenHostMap[0]) > 0 {
		sim.status = InfectedStatusCode
	}
	return sim, nil
}
