package contagiongo

import (
	"bufio"
	"fmt"
	"log"
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
func LoadSequences(path string) (map[int][][]int, error) {
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
	reTranslate := regexp.MustCompile(`([A-Za-z0-9]+)\s*\:\s*(\d+)`)
	// Regular expression to find target host
	reHostID := regexp.MustCompile(`h\s*\:\s*(\d+)`)
	// seqRe := regexp.MustCompile(`[A-Za-z0-9]`)
	scanner := bufio.NewScanner(f)
	lineNum := 0
	currentHostID := -1
	translationMap := make(map[string]int)
	var currentSeq []int
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// ignore comment lines
			continue
		} else if strings.HasPrefix(line, "%") {
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
