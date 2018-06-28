package contagiongo

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

// LoadSingleHostConfig parses a TOML config file and
// creates a SingleHostConfig configuration.
func LoadSingleHostConfig(path string) (*SingleHostConfig, error) {
	spec := new(SingleHostConfig)
	_, err := toml.DecodeFile(path, spec)
	if err != nil {
		return nil, err
	}
	return spec, nil
}

// LoadSequences parses a specially-formatted FASTA file to
// get the sequences, encode sequences into integers, and distribute
// to assigned hosts.
// Returns a map where the key is the host ID and the values are
// the pathogen sequences for the particular host.
func LoadSequences(path string) (map[int][][]uint8, error) {
	/*
		Format:

		# This is a comment
		% U:0 P:1
		>h:0
		UUPUUPUPUUPPUUPUUPUPPPUU
		PUPUUPUPUPUPPUUUUPPPPPPP

		First line indicates how characters are translated into integer
		encoding. The desciption line is similar to FASTA format except that
		the pattern h:\d+ must be present somewhere in the line.

		This indicates which host the pathogen will be inoculated in, where
		\d+ is the host ID of the target host.
	*/

	pathogenHostMap := make(map[int][][]uint8)

	// TODO: add pathogens to hosts using host ID list
	// Read from file
	// assign sequences based on ID
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// Regular expression to encode characters into integers
	reTranslate := regexp.MustCompile(`([A-Za-z0-9]+)\s*\:\s*(\d+)`)
	// Regular expression to find target host
	reHostID := regexp.MustCompile(`h\s*\:\s*(\d+)`)
	// seqRe := regexp.MustCompile(`[A-Za-z0-9]`)
	scanner := bufio.NewScanner(f)
	lineNum := 0
	currentHostID := -1
	translationMap := make(map[string]uint8)
	var currentSeq []uint8
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// ignore comment lines
			continue
		} else if strings.HasPrefix(line, "%") {
			// Check if first line starts with #
			for _, match := range reTranslate.FindAllStringSubmatch(line, -1) {
				if len(match[1]) > 0 && len(match[2]) > 0 {
					encoded, err := strconv.Atoi(match[2])
					if err != nil {
						return nil, fmt.Errorf(FileParsingError, lineNum, err)
					}
					translationMap[match[1]] = uint8(encoded)
				}
			}
		} else if strings.HasPrefix(line, ">") {
			// Check if line starts with >
			if len(currentSeq) > 0 {
				pathogenHostMap[currentHostID] = append(pathogenHostMap[currentHostID], currentSeq)
				currentSeq = []uint8{}
			}
			res := reHostID.FindStringSubmatch(line)
			hostID, err := strconv.Atoi(res[1])
			if err != nil {
				return nil, fmt.Errorf(FileParsingError, lineNum, err)
			}
			currentHostID = hostID
		} else {
			// Translate string char into uint8
			for _, char := range line {
				if i, ok := translationMap[string(char)]; ok {
					currentSeq = append(currentSeq, uint8(i))
				}
			}
		}
		lineNum++
	}
	if len(currentSeq) > 0 {
		pathogenHostMap[currentHostID] = append(pathogenHostMap[currentHostID], currentSeq)
		currentSeq = []uint8{}
	}
	return pathogenHostMap, nil
}

// LoadFitnessMatrix parses and loads the fitness matrix encoded in the
// text file at the given path.
func LoadFitnessMatrix(path string, valueType string) (map[int]map[uint8]float64, error) {
	/*
		Format:

		# This is a comment
		# Any line starting with a # is skipped
		log
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
	reValues := regexp.MustCompile(`-?\d*\.?\d+`)

	// Assumes that the default input type is base 10
	inputValueType := "dec" // or log

	i := 1
	fitnessMap := make(map[int][]float64)
	var defaultValues []float64
	alleles := 0
	lastPos := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		i++

		if strings.HasPrefix(line, "#") {
			// ignore comment lines
			continue
		} else if strings.HasPrefix(line, "dec") || strings.HasPrefix(line, "base10") {
			inputValueType = "dec"
		} else if strings.HasPrefix(line, "log") || strings.HasPrefix(line, "ln") {
			inputValueType = "log"
		} else if strings.HasPrefix(line, "default") {
			// Load default values
			// These values will be used if the site is not included in the
			// file.
			splittedLine := strings.Split(line, "->")
			for _, vStr := range reValues.FindAllString(splittedLine[1], -1) {
				v, _ := strconv.ParseFloat(vStr, 64)
				defaultValues = append(defaultValues, v)

			}
			alleles = len(defaultValues)
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
				// Convert values if there is a mismatch
				switch {
				case valueType == "log" && inputValueType == "base10":
					fallthrough
				case valueType == "log" && inputValueType == "dec":
					v = math.Log(v)
				case valueType == "base10" && inputValueType == "log":
					fallthrough
				case valueType == "dec" && inputValueType == "log":
					v = math.Exp(v)
				}
				values = append(values, v)
			}
			if _, ok := fitnessMap[pos]; ok {
				return nil, fmt.Errorf("duplicate position index (%d)", pos)
			}
			fitnessMap[pos] = values
			if lastPos < pos {
				lastPos = pos
			}
			if alleles == 0 {
				alleles = len(values)
			} else if alleles != len(values) {
				return nil, fmt.Errorf("number of alleles in site %d (%d) is not equal to the previous count (%d)", pos, len(values), alleles)
			}
		}
	}
	// Transform map into FitnessMatrix
	m := make(map[int]map[uint8]float64)
	for i := 0; i <= lastPos; i++ {
		m[i] = make(map[uint8]float64)
		for j := 0; j < alleles; j++ {
			if _, ok := fitnessMap[i]; ok {
				m[i][uint8(j)] = fitnessMap[i][j]
			} else {
				m[i][uint8(j)] = defaultValues[j]
			}
		}
	}
	return m, nil
}

// LoadAdjacencyMatrix creates a new 2D mapping based on a text file.
func LoadAdjacencyMatrix(path string) (HostNetwork, error) {
	m := make(adjacencyMatrix)

	/*
		Parses text file for connection information between hosts.
		Ignores lines that starts with #.
		The text file should be formatted as follows for every line:

			from_uid<int>    to_uid<int>    weight<float64>

		For every line in the file, the AddWeightedConnection method is called
		to create a directed edge between the source and recepient host
		specified by the given UIDs.

		If the edge is undirected, two declarations are expected per source-
		recepient pair: one in the forward direction and the other in the
		reverse direction.

	*/
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	re := regexp.MustCompile(`(\d+)\s+(\d+)\s+(\d*\.?\d+)`)
	scanner := bufio.NewScanner(f)
	i := 0
	for scanner.Scan() {
		i++
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// ignore comment lines
			continue
		}
		res := re.FindStringSubmatch(line)
		if len(res) == 0 {
			continue
		}
		if len(res) < 3 {
			return nil, fmt.Errorf("invalid format in line %d", i)
		}
		a, err := strconv.Atoi(res[1])
		if err != nil {
			return nil, fmt.Errorf("%s in line %d", err, i)
		}
		b, err := strconv.Atoi(res[2])
		if err != nil {
			return nil, fmt.Errorf("%s in line %d", err, i)
		}
		wt, err := strconv.ParseFloat(res[3], 64)
		if err != nil {
			return nil, fmt.Errorf("%s in line %d", err, i)
		}
		m.AddWeightedConnection(a, b, wt)
	}
	return m, nil
}

// LoadEvoEpiConfig creates an EvoEpiConfig struct from a TOML file.
func LoadEvoEpiConfig(path string) (*EvoEpiConfig, error) {
	var conf EvoEpiConfig
	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		return &EvoEpiConfig{}, err
	}
	return &conf, nil
}
