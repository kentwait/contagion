package contagiongo

import (
	"fmt"
	"strings"
)

// Config represents any top level TOML configuration
// that can create a new simulation.
type Config interface {
	Validate() error
	NewSimulation() (Epidemic, error)
	NumInstances() int
	NumGenerations() int
	LogFreq() int
	LogPath() string
}

// EvoEpiConfig contains parameters to create a simulated infection
// in a connected network of hosts.
type EvoEpiConfig struct {
	SimParams          *epidemicSimConfig      `toml:"simulation"`
	LogParams          *logConfig              `toml:"logging"`
	IntrahostModels    []*intrahostModelConfig `toml:"intrahost_model"`
	FitnessModels      []*fitnessModelConfig   `toml:"fitness_model"`
	TransmissionModels []*transModelConfig     `toml:"transmission_model"`

	validated bool
}

// Validate checks the validity of the configuration.
func (c *EvoEpiConfig) Validate() error {
	// Validate sections
	err := c.SimParams.Validate()
	if err != nil {
		return err
	}
	err = c.LogParams.Validate()
	if err != nil {
		return err
	}
	// Validate each intrahost model
	// Check if host_ids are unique
	hostIDSet := make(map[int]bool)
	for _, model := range c.IntrahostModels {
		err := model.Validate()
		if err != nil {
			return err
		}
		// Check if durations match EpidemicModel
		switch strings.ToLower(c.SimParams.EpidemicModel) {
		case "si":
			if model.InfectedDuration != 0 && model.InfectedDuration < c.SimParams.NumGenerations {
				return fmt.Errorf("cannot create %s model if %s (%d) is less than the number of generations (%d)",
					c.SimParams.EpidemicModel,
					"infected_duration",
					model.InfectedDuration,
					c.SimParams.NumGenerations,
				)
			}
			// Assign default value
			if model.InfectedDuration == 0 {
				model.InfectedDuration = c.SimParams.NumGenerations + 1
			}
		case "sis":
			if model.InfectedDuration > c.SimParams.NumGenerations {
				return fmt.Errorf("cannot create %s model if %s (%d) is greater than the number of generations (%d)",
					c.SimParams.EpidemicModel,
					"infected_duration",
					model.InfectedDuration,
					c.SimParams.NumGenerations,
				)
			}
		case "sir":
			if model.RemovedDuration != 0 && model.RemovedDuration < c.SimParams.NumGenerations {
				return fmt.Errorf("cannot create %s model if %s (%d) is less than the number of generations (%d)",
					c.SimParams.EpidemicModel,
					"removed_duration",
					model.RemovedDuration,
					c.SimParams.NumGenerations,
				)
			}
			if model.InfectedDuration > c.SimParams.NumGenerations {
				return fmt.Errorf("cannot create %s model if %s (%d) is greater than the number of generations (%d)",
					c.SimParams.EpidemicModel,
					"infected_duration",
					model.InfectedDuration,
					c.SimParams.NumGenerations,
				)
			}
			// Assign default value
			if model.RemovedDuration == 0 {
				model.RemovedDuration = c.SimParams.NumGenerations + 1
			}
		case "sirs":
			if model.InfectedDuration > c.SimParams.NumGenerations {
				return fmt.Errorf("cannot create %s model if %s (%d) is greater than the number of generations (%d)",
					c.SimParams.EpidemicModel,
					"infected_duration",
					model.InfectedDuration,
					c.SimParams.NumGenerations,
				)
			}
			if model.RemovedDuration > c.SimParams.NumGenerations {
				return fmt.Errorf("cannot create %s model if %s (%d) is greater than the number of generations (%d)",
					c.SimParams.EpidemicModel,
					"removed_duration",
					model.RemovedDuration,
					c.SimParams.NumGenerations,
				)
			}
		case "sei":
			if model.InfectiveDuration != 0 && model.InfectiveDuration < c.SimParams.NumGenerations {
				return fmt.Errorf("cannot create %s model if %s (%d) is less than the number of generations (%d)",
					c.SimParams.EpidemicModel,
					"infective_duration",
					model.InfectiveDuration,
					c.SimParams.NumGenerations,
				)
			}
			if model.ExposedDuration > c.SimParams.NumGenerations {
				return fmt.Errorf("cannot create %s model if %s (%d) is greater than the number of generations (%d)",
					c.SimParams.EpidemicModel,
					"exposed_duration",
					model.ExposedDuration,
					c.SimParams.NumGenerations,
				)
			}
			// Assign default value
			if model.InfectiveDuration == 0 {
				model.InfectiveDuration = c.SimParams.NumGenerations + 1
			}
		case "seir":
			if model.RemovedDuration != 0 && model.RemovedDuration < c.SimParams.NumGenerations {
				return fmt.Errorf("cannot create %s model if %s (%d) is less than the number of generations (%d)",
					c.SimParams.EpidemicModel,
					"removed_duration",
					model.RemovedDuration,
					c.SimParams.NumGenerations,
				)
			}
			if model.ExposedDuration > c.SimParams.NumGenerations {
				return fmt.Errorf("cannot create %s model if %s (%d) is greater than the number of generations (%d)",
					c.SimParams.EpidemicModel,
					"exposed_duration",
					model.ExposedDuration,
					c.SimParams.NumGenerations,
				)
			}
			if model.InfectiveDuration > c.SimParams.NumGenerations {
				return fmt.Errorf("cannot create %s model if %s (%d) is greater than the number of generations (%d)",
					c.SimParams.EpidemicModel,
					"exposed_duration",
					model.InfectiveDuration,
					c.SimParams.NumGenerations,
				)
			}
			// Assign default value
			if model.RemovedDuration == 0 {
				model.RemovedDuration = c.SimParams.NumGenerations + 1
			}
		case "seirs":
			if model.ExposedDuration > c.SimParams.NumGenerations {
				return fmt.Errorf("cannot create %s model if %s (%d) is greater than the number of generations (%d)",
					c.SimParams.EpidemicModel,
					"exposed_duration",
					model.ExposedDuration,
					c.SimParams.NumGenerations,
				)
			}
			if model.InfectiveDuration > c.SimParams.NumGenerations {
				return fmt.Errorf("cannot create %s model if %s (%d) is greater than the number of generations (%d)",
					c.SimParams.EpidemicModel,
					"exposed_duration",
					model.InfectiveDuration,
					c.SimParams.NumGenerations,
				)
			}
			if model.RemovedDuration > c.SimParams.NumGenerations {
				return fmt.Errorf("cannot create %s model if %s (%d) is greater than the number of generations (%d)",
					c.SimParams.EpidemicModel,
					"removed_duration",
					model.RemovedDuration,
					c.SimParams.NumGenerations,
				)
			}
		case "endtrans":
			// Modified version of SI or SIR
			// Endtrans adopts the SI model when
			// infected duration > number of generations

			// infected duration must be set
			if model.InfectedDuration < 1 {
				return fmt.Errorf("cannot create %s model if %s (%d) is less than 1",
					c.SimParams.EpidemicModel,
					"infected_duration",
					model.InfectedDuration,
				)
			}
			if model.RemovedDuration != 0 && model.RemovedDuration < c.SimParams.NumGenerations {
				return fmt.Errorf("cannot create %s model if %s (%d) is less than the number of generations (%d)",
					c.SimParams.EpidemicModel,
					"removed_duration",
					model.RemovedDuration,
					c.SimParams.NumGenerations,
				)
			}
			// Assign default value
			if model.RemovedDuration == 0 {
				model.RemovedDuration = c.SimParams.NumGenerations + 1
			}
		}
		//
		for _, i := range model.HostIDs {
			if _, exists := hostIDSet[i]; exists {
				return fmt.Errorf("host id "+IntKeyExists, i)
			}
			hostIDSet[i] = true
		}
	}
	// Check if all hosts have been assigned a model
	for i := 0; i < c.SimParams.HostPopSize; i++ {
		if !hostIDSet[i] {
			return fmt.Errorf("host %d was not assigned a intrahost model", i)
		}
	}

	// Validate each fitness model
	// Check if host_ids are unique
	hostIDSet = make(map[int]bool)
	for _, model := range c.FitnessModels {
		err := model.Validate()
		if err != nil {
			return err
		}
		for _, i := range model.HostIDs {
			if _, exists := hostIDSet[i]; exists {
				return fmt.Errorf("host id "+IntKeyExists, i)
			}
			hostIDSet[i] = true
		}
	}
	// Validate each transmission model
	// Check if host_ids are unique
	hostIDSet = make(map[int]bool)
	for _, model := range c.TransmissionModels {
		err := model.Validate()
		if err != nil {
			return err
		}
		for _, i := range model.HostIDs {
			if _, exists := hostIDSet[i]; exists {
				return fmt.Errorf("host id "+IntKeyExists, i)
			}
			hostIDSet[i] = true
		}
	}
	// Check if all hosts have been assigned a model
	for i := 0; i < c.SimParams.HostPopSize; i++ {
		if !hostIDSet[i] {
			return fmt.Errorf("host %d was not assigned a fitness model", i)
		}
	}

	// TODO: validate file paths
	c.validated = true
	return nil
}

// NewSimulation creates a new SingleHostSimulation simulation.
func (c *EvoEpiConfig) NewSimulation() (Epidemic, error) {
	sim := new(evoEpiSimulation)
	// Initialize maps
	sim.hosts = make(map[int]Host)
	sim.statuses = make(map[int]int)
	sim.timers = make(map[int]int)
	sim.intrahostModels = make(map[int]IntrahostModel)
	sim.fitnessModels = make(map[int]FitnessModel)
	sim.transModels = make(map[int]TransmissionModel)
	sim.hostNeighborhoods = make(map[int][]Host)
	// Create empty hosts
	for i := 0; i < c.SimParams.HostPopSize; i++ {
		sim.hosts[i] = EmptySequenceHost(i)
	}

	// Create IntrahostModels
	for i, conf := range c.IntrahostModels {
		model, err := conf.CreateModel(i)
		if err != nil {
			return nil, err
		}
		model.SetModelID(i)
		sim.intrahostModels[i] = model
		// assign to hosts
		for _, id := range conf.HostIDs {
			err := sim.hosts[id].SetIntrahostModel(model)
			if err != nil {
				return nil, err
			}
		}
	}
	// Create FitnessModels
	for i, conf := range c.FitnessModels {
		model, err := conf.CreateModel(i)
		if err != nil {
			return nil, err
		}
		model.SetModelID(i)
		sim.fitnessModels[i] = model
		// assign to hosts
		for _, id := range conf.HostIDs {
			err := sim.hosts[id].SetFitnessModel(model)
			if err != nil {
				return nil, err
			}
		}
	}
	// Create TransmissionModels
	for i, conf := range c.TransmissionModels {
		model, err := conf.CreateModel(i)
		if err != nil {
			return nil, err
		}
		model.SetModelID(i)
		sim.transModels[i] = model
		// assign to hosts
		for _, id := range conf.HostIDs {
			err := sim.hosts[id].SetTransmissionModel(model)
			if err != nil {
				return nil, err
			}
		}
	}
	// Load host connections
	var err error
	sim.hostNetwork, err = LoadAdjacencyMatrix(c.SimParams.HostNetworkPath)
	if err != nil {
		return nil, err
	}
	// Construct neighborhoods
	for id := range sim.hosts {
		neighborIDs := sim.hostNetwork.GetNeighbors(id)
		sim.hostNeighborhoods[id] = make([]Host, len(neighborIDs))
		for i, neighborID := range neighborIDs {
			sim.hostNeighborhoods[id][i] = sim.hosts[neighborID]
		}
	}
	// Initialize empty GenotypeTree
	sim.tree = EmptyGenotypeTree()
	// Load pathogens
	hostPathogenMap, err := LoadSequences(c.SimParams.PathogenSequencePath)
	if err != nil {
		return nil, err
	}
	// Seed pathogens into host/s
	for id, sequences := range hostPathogenMap {
		for _, s := range sequences {
			// Seeded pathogens are all roots
			genotype := sim.tree.NewNode(s, 0)
			sim.hosts[id].AddPathogen(genotype)
		}
	}
	// Add config to simulation
	sim.config = c

	// Add infectable status
	sim.infectableStatuses = []int{SusceptibleStatusCode}
	if c.SimParams.Coinfection {
		switch strings.ToLower(c.SimParams.EpidemicModel) {
		case "si":
			sim.infectableStatuses = append(sim.infectableStatuses, []int{
				InfectedStatusCode,
			}...)
		case "sis":
			sim.infectableStatuses = append(sim.infectableStatuses, []int{
				InfectedStatusCode,
			}...)
		case "sir":
			sim.infectableStatuses = append(sim.infectableStatuses, []int{
				InfectedStatusCode,
			}...)
		case "sirs":
			sim.infectableStatuses = append(sim.infectableStatuses, []int{
				InfectedStatusCode,
			}...)
		case "sei":
			sim.infectableStatuses = append(sim.infectableStatuses, []int{
				ExposedStatusCode,
				InfectiveStatusCode,
			}...)
		case "seir":
			sim.infectableStatuses = append(sim.infectableStatuses, []int{
				ExposedStatusCode,
				InfectiveStatusCode,
			}...)
		case "seirs":
			sim.infectableStatuses = append(sim.infectableStatuses, []int{
				ExposedStatusCode,
				InfectiveStatusCode,
			}...)
		case "endtrans":
			sim.infectableStatuses = append(sim.infectableStatuses, []int{
				InfectedStatusCode,
			}...)
		}
	}
	// Initialize host statuses to 1
	for i := range sim.hosts {
		sim.statuses[i] = 1
	}
	// Initialize host timers to -1
	for i := range sim.hosts {
		sim.timers[i] = -1
	}

	return sim, nil
}

// NumInstances returns the number of independent realizations to run.
func (c *EvoEpiConfig) NumInstances() int { return c.SimParams.NumIntances }

// NumGenerations returns the number of pathogen generation in
// a single simulation run.
func (c *EvoEpiConfig) NumGenerations() int { return c.SimParams.NumGenerations }

// LogFreq returns the number of pathogen generations in the simulation until
// data is recorded.
func (c *EvoEpiConfig) LogFreq() int { return int(c.LogParams.LogFreq) }

// LogPath returns the path where to write results.
// This can a path to a folder, or directory_path/prefix format
func (c *EvoEpiConfig) LogPath() string { return c.LogParams.LogPath }

type epidemicSimConfig struct {
	NumGenerations int    `toml:"num_generations"`
	NumIntances    int    `toml:"num_instances"`
	HostPopSize    int    `toml:"host_popsize"`
	EpidemicModel  string `toml:"epidemic_model"` // si, sir, sirs, sei, seis, seirs, endtrans
	Coinfection    bool   `toml:"coinfection"`

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

	switch strings.ToLower(c.EpidemicModel) {
	case "si":
	case "sis":
	case "sir":
	case "sirs":
	case "sei":
	case "seir":
	case "seirs":
	case "endtrans":
	default:
		return fmt.Errorf(UnrecognizedKeywordError, c.EpidemicModel, "epidemic_model")
	}
	c.validated = true
	return nil
}

type logConfig struct {
	LogFreq   uint   `toml:"log_freq"`
	LogPath   string `toml:"log_path"`
	validated bool
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
	HostIDs           []int       `toml:"host_ids"`
	MutationRate      float64     `toml:"mutation_rate"`
	TransitionMatrix  [][]float64 `toml:"transition_matrix"`
	RecombinationRate float64     `toml:"recombination_rate"`
	ReplicationModel  string      `toml:"replication_model"` // constant, bht, fitness
	ConstantPopSize   int         `toml:"constant_pop_size"` // only for constant
	MaxPopSize        int         `toml:"max_pop_size"`      // only for bht and fitness
	GrowthRate        float64     `toml:"growth_rate"`       // only for bht

	// If no duration is assigned, the default value is 0.
	// Aften calling the first process step, if the current duration is 0,
	// it turns to -1 and will never triggers the update step.
	// If running under the fitness mode, do not assign durations.
	ExposedDuration    int `toml:"exposed_duration"`
	InfectedDuration   int `toml:"infected_duration"`
	InfectiveDuration  int `toml:"infective_duration"`
	RemovedDuration    int `toml:"removed_duration"`
	RecoveredDuration  int `toml:"recovered_duration"`
	DeadDuration       int `toml:"dead_duration"`
	VaccinatedDuration int `toml:"vaccinated_duration"`

	validated bool
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

	// Check durations
	if c.ExposedDuration < 0 {
		return fmt.Errorf(InvalidIntParameterError, "exposed_duration", c.ExposedDuration, "cannot be negative")
	}
	if c.InfectedDuration < 0 {
		return fmt.Errorf(InvalidIntParameterError, "infected_duration", c.InfectedDuration, "cannot be negative")
	}
	if c.InfectiveDuration < 0 {
		return fmt.Errorf(InvalidIntParameterError, "infective_duration", c.InfectiveDuration, "cannot be negative")
	}
	if c.RemovedDuration < 0 {
		return fmt.Errorf(InvalidIntParameterError, "removed_duration", c.RemovedDuration, "cannot be negative")
	}
	if c.RecoveredDuration < 0 {
		return fmt.Errorf(InvalidIntParameterError, "recovered_duration", c.RecoveredDuration, "cannot be negative")
	}
	if c.DeadDuration < 0 {
		return fmt.Errorf(InvalidIntParameterError, "dead_duration", c.DeadDuration, "cannot be negative")
	}
	if c.VaccinatedDuration < 0 {
		return fmt.Errorf(InvalidIntParameterError, "vaccinated_duration", c.VaccinatedDuration, "cannot be negative")
	}
	// TODO: make sure host_ids are unique across models

	c.validated = true
	return nil
}

// CreateModel creates an IntrahostModel based on the configuration.
func (c *intrahostModelConfig) CreateModel(id int) (IntrahostModel, error) {
	if !c.validated {
		return nil, fmt.Errorf("validate model parameters first")
	}
	statusDuration := make(map[int]int)
	statusDuration[ExposedStatusCode] = c.ExposedDuration
	statusDuration[InfectedStatusCode] = c.InfectedDuration
	statusDuration[InfectiveStatusCode] = c.InfectiveDuration
	statusDuration[RemovedStatusCode] = c.RemovedDuration
	statusDuration[RecoveredStatusCode] = c.RecoveredDuration
	statusDuration[DeadStatusCode] = c.DeadDuration
	statusDuration[VaccinatedStatusCode] = c.VaccinatedDuration

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
		model.statusDuration = statusDuration
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
		model.statusDuration = statusDuration
		return model, nil
	case "fitness":
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
		model.statusDuration = statusDuration
		return model, nil
	}
	return nil, fmt.Errorf(UnrecognizedKeywordError, c.ReplicationModel, "replication_model")
}

// FitnessModelConfig contains parameters to create an FitnessModel.
type fitnessModelConfig struct {
	ModelName        string `toml:"model_name"`
	HostIDs          []int  `toml:"host_ids"`
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

	// TODO: make sure host_ids are unique across models
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

type transModelConfig struct {
	ModelName        string  `toml:"model_name"`
	HostIDs          []int   `toml:"host_ids"`
	Mode             string  `toml:"mode"` // poisson, constant
	TransmissionProb float64 `toml:"transmission_prob"`
	TransmissionSize float64 `toml:"transmission_size"`
	validated        bool
}

// Validate checks the validity of the transModelConfig configuration.
func (c *transModelConfig) Validate() error {
	// check keywords
	// fitness_model
	switch strings.ToLower(c.Mode) {
	case "poisson":
	case "constant":
	default:
		return fmt.Errorf(UnrecognizedKeywordError, c.Mode, "mode")
	}
	c.validated = true
	return nil
}

func (c *transModelConfig) CreateModel(id int) (TransmissionModel, error) {
	if !c.validated {
		return nil, fmt.Errorf("validate model parameters first")
	}
	switch c.Mode {
	case "poisson":
		model := new(poissonTransmitter)
		model.name = c.ModelName
		model.prob = c.TransmissionProb
		model.size = c.TransmissionSize
		return model, nil
	case "constant":
		model := new(constantTransmitter)
		model.name = c.ModelName
		model.prob = c.TransmissionProb
		model.size = int(c.TransmissionSize)
		return model, nil
	}
	return nil, fmt.Errorf(UnrecognizedKeywordError, c.Mode, "mode")
}
