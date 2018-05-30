package contagiongo

import (
	"fmt"
	"strings"
)

// EvoEpiConfig contains parameters to create a simulated infection
// in a connected network of hosts.
type EvoEpiConfig struct {
	NumGenerations uint `toml:"num_generations"`
	NumReplicates  uint `toml:"num_replicates"`

	IntrahostModels []IntrahostModelConfig `toml:"intrahost_models"`
	FitnessModels   []FitnessModelConfig   `toml:"fitness_models"`

	PathogenSequencePath string `toml:"pathogen_sequence_path"` // fasta file for seeding infections

	LogFreq         uint   `toml:"log_freq"`
	PathogenLogPath string `toml:"pathogen_log_path"`

	validated bool
}

// Validate checks the validity of the configuration.
func (c *EvoEpiConfig) Validate() error {
	// Validate each intrahost model
	for _, model := range c.IntrahostModels {
		err := model.Validate()
		if err != nil {
			return err
		}
	}
	// Validate each fitness model
	for _, model := range c.FitnessModels {
		err := model.Validate()
		if err != nil {
			return err
		}
	}
	c.validated = true
	return nil
}

// NewSimulation creates a new SingleHostSimulation simulation.
func (c *EvoEpiConfig) NewSimulation() (*EvoEpiSimulation, error) {
	sim := new(EvoEpiSimulation)
	// Create IntrahostModels
	for i, conf := range c.IntrahostModels {
		model, err := conf.CreateModel()
		if err != nil {
			return nil, err
		}
		model.SetModelID(i)
		sim.IntrahostModels[i] = model
	}
	// Create FitnessModels
	for i, conf := range c.FitnessModels {
		model, err := conf.CreateModel()
		if err != nil {
			return nil, err
		}
		model.SetModelID(i)
		sim.FitnessModels[i] = model
	}
	return sim, nil
}

// IntrahostModelConfig contains parameters to create an IntrahostModel.
type IntrahostModelConfig struct {
	ModelName         string      `toml:"model_name"`
	MutationRate      float64     `toml:"mutation_rate"`
	TransitionMatrix  [][]float64 `toml:"transition_matrix"`
	RecombinationRate float64     `toml:"recombination_rate"`
	ReplicationModel  string      `toml:"replication_model"` // constant, bht, fitness
	ConstantPopSize   int         `toml:"constant_pop_size"` // only for constant
	MaxPopSize        int         `toml:"max_pop_size"`      // only for bht and fitness
	GrowthRate        float64     `toml:"growth_rate"`       // only for bht
	validated         bool
}

// Validate checks the validity of the IntrahostModelConfig configuration.
func (c *IntrahostModelConfig) Validate() error {
	// check keywords
	// replication_model
	switch strings.ToLower(c.ReplicationModel) {
	case "constant":
	case "bht":
	case "fitness":
	default:
		return fmt.Errorf(UnrecognizedKeywordError, c.ReplicationModel, "replication_model")
	}
	c.validated = true
	return nil
}

// CreateModel creates an IntrahostModel based on the configuration.
func (c *IntrahostModelConfig) CreateModel() (IntrahostModel, error) {
	switch c.ReplicationModel {
	case "constant":
		model := new(ConstantPopModel)
		model.name = c.ModelName
		model.popSize = c.ConstantPopSize
		model.mutationRate = c.MutationRate
		model.recombinationRate = c.RecombinationRate
		model.transitionMatrix = make([][]float64, len(c.TransitionMatrix))
		for i := 0; i < len(c.TransitionMatrix); i++ {
			model.transitionMatrix[i] = make([]float64, len(c.TransitionMatrix))
			copy(model.transitionMatrix[i], c.TransitionMatrix[i])
		}
		return model, nil
	case "bht":
		model := new(BevertonHoltThresholdPopModel)
		model.name = c.ModelName
		model.maxPopSize = c.MaxPopSize
		model.growthRate = c.GrowthRate
		model.mutationRate = c.MutationRate
		model.recombinationRate = c.RecombinationRate
		model.transitionMatrix = make([][]float64, len(c.TransitionMatrix))
		for i := 0; i < len(c.TransitionMatrix); i++ {
			model.transitionMatrix[i] = make([]float64, len(c.TransitionMatrix))
			copy(model.transitionMatrix[i], c.TransitionMatrix[i])
		}
		return model, nil
	}
	// fitness
	model := new(FitnessDependentPopModel)
	model.name = c.ModelName
	model.maxPopSize = c.MaxPopSize
	model.mutationRate = c.MutationRate
	model.recombinationRate = c.RecombinationRate
	model.transitionMatrix = make([][]float64, len(c.TransitionMatrix))
	for i := 0; i < len(c.TransitionMatrix); i++ {
		model.transitionMatrix[i] = make([]float64, len(c.TransitionMatrix))
		copy(model.transitionMatrix[i], c.TransitionMatrix[i])
	}
	return model, nil
}

// FitnessModelConfig contains parameters to create an FitnessModel.
type FitnessModelConfig struct {
	ModelName        string `toml:"model_name"`
	FitnessModel     string `toml:"fitness_model"` // multiplicative, additive, additive_motif
	FitnessModelPath string `toml:"fitness_model_path"`
	validated        bool
}

// Validate checks the validity of the FitnessModelConfig configuration.
func (c *FitnessModelConfig) Validate() error {
	// check keywords
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

// CreateModel creates an FitnessModel based on the configuration.
func (c *FitnessModelConfig) CreateModel() (FitnessModel, error) {
	// Create FitnessModel
	switch c.FitnessModel {
	case "multiplicative":
		matrix, err := LoadFitnessMatrix(c.FitnessModelPath)
		if err != nil {
			return nil, err
		}
		fm := NewMultiplicativeFM(0, "multiplicative", matrix)
		return fm, nil
	case "additive":
		matrix, err := LoadFitnessMatrix(c.FitnessModelPath)
		if err != nil {
			return nil, err
		}
		fm := NewAdditiveFM(0, "additive", matrix)
		return fm, nil
	}
	// additive_motif
	return nil, fmt.Errorf("additive_motif not yet implemented")
}
