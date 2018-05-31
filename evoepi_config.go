package contagiongo

import (
	"fmt"
	"strings"
)

// EvoEpiConfig contains parameters to create a simulated infection
// in a connected network of hosts.
type EvoEpiConfig struct {
	SimParams       epidemicSimConfig      `toml:"simulation"`
	LogParams       logConfig              `toml:"logging"`
	IntrahostModels []intrahostModelConfig `toml:"intrahost_model"`
	FitnessModels   []fitnessModelConfig   `toml:"fitness_model"`

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
	// TODO: validate file paths
	c.validated = true
	return nil
}

// NewSimulation creates a new SingleHostSimulation simulation.
func (c *EvoEpiConfig) NewSimulation() (Epidemic, error) {
	sim := new(evoEpiSimulation)
	// Create IntrahostModels
	for i, conf := range c.IntrahostModels {
		model, err := conf.CreateModel(i)
		if err != nil {
			return nil, err
		}
		model.SetModelID(i)
		sim.intrahostModels[i] = model
	}
	// Create FitnessModels
	for i, conf := range c.FitnessModels {
		model, err := conf.CreateModel(i)
		if err != nil {
			return nil, err
		}
		model.SetModelID(i)
		sim.fitnessModels[i] = model
	}
	// Create epidemic simulation
	switch c.SimParams.EpidemicModel {
	case "si":
	case "sis":
	case "sir":
	case "sirs":
	case "sei":
	case "seir":
	case "seirs":
	}
	return sim, nil
}

type epidemicSimConfig struct {
	NumGenerations int    `toml:"num_generations"`
	NumIntances    int    `toml:"num_instances"`
	HostPopSize    int    `toml:"host_popsize"`
	EpidemicModel  string `toml:"epidemic_model"` // si, sir, sirs, sei, seis, seirs

	PathogenSequencePath string `toml:"pathogen_sequence_path"` // fasta file for seeding infections
	HostNetworkPath      string `toml:"host_network_path"`
	validated            bool
}

func (c *epidemicSimConfig) Validate() error {
	// Check PathogenSequencePath
	exists, err := Exists(c.PathogenSequencePath)
	if err != nil {
		return fmt.Errorf("error checking if file in %s exists: %s", c.PathogenSequencePath, err)
	}
	if !exists {
		return fmt.Errorf("file in %s does not exist", c.PathogenSequencePath)
	}

	// Check HostNetworkPath
	exists, err = Exists(c.HostNetworkPath)
	if err != nil {
		return fmt.Errorf("error checking if file in %s exists: %s", c.HostNetworkPath, err)
	}
	if !exists {
		return fmt.Errorf("file in %s does not exist", c.HostNetworkPath)
	}

	// Check parameter values
	if c.NumGenerations < 1 {
		return fmt.Errorf(InvalidIntParameterError, "num_generations", c.NumGenerations, "must be greater than or equal to 1")
	}
	if c.NumIntances < 1 {
		return fmt.Errorf(InvalidIntParameterError, "num_instances", c.NumIntances, "must be greater than or equal to 1")
	}
	if c.HostPopSize < 1 {
		return fmt.Errorf(InvalidIntParameterError, "host_popsize", c.HostPopSize, "must be greater than or equal to 1")
	}

	switch c.EpidemicModel {
	case "si":
	case "sis":
	case "sir":
	case "sirs":
	case "sei":
	case "seir":
	case "seirs":
	default:
		return fmt.Errorf(UnrecognizedKeywordError, c.EpidemicModel, "epidemic_model")
	}
	c.validated = true
	return nil
}

type logConfig struct {
	LogFreq         uint   `toml:"log_freq"`
	PathogenLogPath string `toml:"pathogen_log_path"`
	validated       bool
}

func (c *logConfig) Validate() error {
	// Check parameter values
	if c.LogFreq < 1 {
		return fmt.Errorf(InvalidIntParameterError, "log_freq", c.LogFreq, "must be greater than or equal to 1")
	}
	c.validated = true
	return nil
}

// IntrahostModelConfig contains parameters to create an IntrahostModel.
type intrahostModelConfig struct {
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
func (c *intrahostModelConfig) Validate() error {
	// check keywords and associated values
	// replication_model
	switch strings.ToLower(c.ReplicationModel) {
	case "constant":
		if c.ConstantPopSize < 1 {
			return fmt.Errorf(InvalidIntParameterError, "constant_pop_size", c.ConstantPopSize, "must be greater than or equal to 1")
		}
	case "bht":
		if c.MaxPopSize < 1 {
			return fmt.Errorf(InvalidIntParameterError, "max_pop_size", c.MaxPopSize, "must be greater than or equal to 1")
		}
		if c.GrowthRate < 0 {
			return fmt.Errorf(InvalidFloatParameterError, "max_pop_size", c.GrowthRate, "must be greater than or equal to 0")
		}
	case "fitness":
		if c.MaxPopSize < 1 {
			return fmt.Errorf(InvalidIntParameterError, "max_pop_size", c.MaxPopSize, "must be greater than or equal to 1")
		}
	default:
		return fmt.Errorf(UnrecognizedKeywordError, c.ReplicationModel, "replication_model")
	}
	// Check mutation rate
	if c.MutationRate < 0 {
		return fmt.Errorf(InvalidFloatParameterError, "mutation_rate", c.MutationRate, "cannot be negative")
	}
	// Check recombination rate
	if c.RecombinationRate < 0 {
		return fmt.Errorf(InvalidFloatParameterError, "recombination_rate", c.RecombinationRate, "cannot be negative")
	}
	// Checks values of TransitionMatrix
	for i, row := range c.TransitionMatrix {
		for j := range row {
			if c.TransitionMatrix[i][j] < 0 {
				return fmt.Errorf(InvalidFloatParameterError, "transition rate", c.TransitionMatrix[i][j], "cannot be negative")
			}
		}
	}
	c.validated = true
	return nil
}

// CreateModel creates an IntrahostModel based on the configuration.
func (c *intrahostModelConfig) CreateModel(id int) (IntrahostModel, error) {
	if !c.validated {
		return nil, fmt.Errorf("validate model parameters first")
	}
	switch c.ReplicationModel {
	case "constant":
		model := new(ConstantPopModel)
		model.id = id
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
		model.id = id
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
	model.id = id
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
type fitnessModelConfig struct {
	ModelName        string `toml:"model_name"`
	FitnessModel     string `toml:"fitness_model"` // multiplicative, additive, additive_motif
	FitnessModelPath string `toml:"fitness_model_path"`
	validated        bool
}

// Validate checks the validity of the FitnessModelConfig configuration.
func (c *fitnessModelConfig) Validate() error {
	// check keywords
	// fitness_model
	switch strings.ToLower(c.FitnessModel) {
	case "multiplicative":
	case "additive":
	case "additive_motif":
	default:
		return fmt.Errorf(UnrecognizedKeywordError, c.FitnessModel, "fitness_model")
	}

	// Check FitnessModelPath
	exists, err := Exists(c.FitnessModelPath)
	if err != nil {
		return fmt.Errorf("error checking if file in %s exists: %s", c.FitnessModelPath, err)
	}
	if !exists {
		return fmt.Errorf("file in %s does not exist", c.FitnessModelPath)
	}

	c.validated = true
	return nil
}

// CreateModel creates an FitnessModel based on the configuration.
func (c *fitnessModelConfig) CreateModel(id int) (FitnessModel, error) {
	if !c.validated {
		return nil, fmt.Errorf("validate model parameters first")
	}
	// Create FitnessModel
	switch c.FitnessModel {
	case "multiplicative":
		matrix, err := LoadFitnessMatrix(c.FitnessModelPath)
		if err != nil {
			return nil, err
		}
		fm := NewMultiplicativeFM(id, "multiplicative", matrix)
		return fm, nil
	case "additive":
		matrix, err := LoadFitnessMatrix(c.FitnessModelPath)
		if err != nil {
			return nil, err
		}
		fm := NewAdditiveFM(id, "additive", matrix)
		return fm, nil
	}
	// additive_motif
	return nil, fmt.Errorf("additive_motif not yet implemented")
}
