package contagiongo

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
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
	FitnessModel      string      `toml:"fitness_model"`     // multiplicative, additive, additive_motif

	PathogenSequencePath string `toml:"pathogen_sequence_path"` // fasta file for seeding infections
	FitnessModelPath     string `toml:"fitness_model_path"`

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
func (c *SingleHostConfig) NewSimulation() (*SingleHostSimulation, error) {
	sim := new(SingleHostSimulation)
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
		fm := NewMultiplicativeFM(0, "multiplicative fitness matrix", matrix)
		sim.fitnessModel = fm
	case "additive":
		matrix, err := LoadFitnessMatrix(c.FitnessModelPath)
		if err != nil {
			return nil, err
		}
		fm := NewAdditiveFM(0, "multiplicative fitness matrix", matrix)
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
		sim.hostStatus = InfectedStatusCode
	}
	return sim, nil
}

// LoadSequences parses a specially-formatted FASTA file to
// get the sequences, encode sequences into integers, and distribute
// to assigned hosts.
// Returns a map where the key is the host ID and the values are
// the pathogen sequences for the particular host.
func LoadSequences(path string) (map[int][][]int, error) {
	pathogenHostMap := make(map[int][][]int)

	// TODO: add pathogens to hosts using host ID list
	// Read from file
	// assign sequences based on ID
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// Regular expression to encode characters into integers
	reTranslate := regexp.MustCompile(`([A-Za-z0-9]+)\:(\d+)`)
	// Regular expression to find target host
	reHostID := regexp.MustCompile(`h\:(\d+)`)
	// seqRe := regexp.MustCompile(`[A-Za-z0-9]`)
	scanner := bufio.NewScanner(f)
	lineNum := 0
	currentHostID := -1
	translationMap := make(map[string]int)
	var currentSeq []int
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") && lineNum == 0 {
			// Check if first line starts with #
			for _, match := range reTranslate.FindAllStringSubmatch(line, -1) {
				if len(match[1]) > 0 && len(match[2]) > 0 {
					translationMap[match[1]], err = strconv.Atoi(match[2])
					if err != nil {
						return nil, fmt.Errorf(FileParsingError, lineNum, err)
					}
				}
			}
		} else if strings.HasPrefix(line, ">") {
			// Check if line starts with >
			if len(currentSeq) > 0 {
				pathogenHostMap[currentHostID] = append(pathogenHostMap[currentHostID], currentSeq)
				currentSeq = []int{}
			}
			res := reHostID.FindStringSubmatch(line)
			hostID, err := strconv.Atoi(res[1])
			if err != nil {
				return nil, fmt.Errorf(FileParsingError, lineNum, err)
			}
			currentHostID = hostID
		} else {
			// Translate string char into int
			for _, char := range line {
				if i, ok := translationMap[string(char)]; ok {
					currentSeq = append(currentSeq, int(i))
				}
			}
		}
		lineNum++
	}
	if len(currentSeq) > 0 {
		pathogenHostMap[currentHostID] = append(pathogenHostMap[currentHostID], currentSeq)
		currentSeq = []int{}
	}
	return pathogenHostMap, nil
}

// LoadFitnessMatrix parses and loads the fitness matrix encoded in the
// text file at the given path.
func LoadFitnessMatrix(path string) (map[int]map[int]float64, error) {
	/*
		Format:

		# This is a comment
		# Any line starting with a # is skipped
		default->1.0, 1.0, 1.0 1.0
		0: 1.0, 1.0, 1.0, 0.5
		1: 1.0, 1.0, 1.0, 0.5
		2: 1.0, 1.0, 1.0, 0.5
		...
		1000: 1.0, 1.0, 1.0, 0.9

	*/
	// Open file and create scanner on top of it
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	// reNum := regexp.MustCompile(`^\s*\d*`)
	rePos := regexp.MustCompile(`^\d+`)
	reValues := regexp.MustCompile(`\d*\.?\d+`)

	i := 1
	fitnessMap := make(map[int][]float64)
	var defaultValues []float64
	lastPos := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// ignore comment lines
			continue
		} else if strings.HasPrefix(line, "default") {
			// Load default values
			// These values will be used if the site is not included in the
			// file.
			splittedLine := strings.Split(line, "->")
			for _, vStr := range reValues.FindAllString(splittedLine[1], -1) {
				v, _ := strconv.ParseFloat(vStr, 64)
				defaultValues = append(defaultValues, v)
			}
		} else {
			// Load site fitness values
			splittedLine := strings.Split(line, ":")
			if len(splittedLine) != 2 {
				return nil, fmt.Errorf("missing colon delimiter in line %d", i)
			}
			prefix, valuesStr := splittedLine[0], splittedLine[1]
			// get position
			posMatch := rePos.FindString(prefix)
			pos, _ := strconv.Atoi(posMatch)
			// get fitness values
			var values []float64
			for _, vStr := range reValues.FindAllString(valuesStr, -1) {
				v, _ := strconv.ParseFloat(vStr, 64)
				values = append(values, v)
			}
			if _, ok := fitnessMap[pos]; ok {
				return nil, fmt.Errorf("duplicate position index (%d)", pos)
			}
			fitnessMap[pos] = values
			if lastPos < pos {
				lastPos = pos
			}
		}
		i++
	}
	// Transform map into FitnessMatrix
	m := make(map[int]map[int]float64)
	for i := 0; i <= lastPos; i++ {
		m[i] = make(map[int]float64)
		if v, ok := fitnessMap[i]; ok {
			for j := 0; j < len(v); j++ {
				m[i][j] = fitnessMap[i][j]
			}
		} else {
			for j := 0; j < len(v); j++ {
				m[i][j] = defaultValues[j]
			}
		}
	}
	return m, nil
}
