package contagiongo

type SingleHostConfig struct {
	NumGenerations uint `toml:"num_generations"`
	NumReplicates  uint `toml:"num_replicates"`

	MutationRate      float64     `toml:"mutation_rate"`
	TransitionMatrix  [][]float64 `toml:"transition_matrix"`
	RecombinationRate float64     `toml:"recombination_rate"`
	ReplicationModel  string      `toml:"replication_model"` // constant, bht, fitness
	FitnessModel      string      `toml:"fitness_model"`     // multiplicative, additive, additive_motif

	PathogenSequencePath string `toml:"fasta_path"` // fasta file for seeding infections

	LogFreq         uint   `toml:"log_freq"`
	PathogenLogPath string `toml:"pathogen_log_path"`
}
