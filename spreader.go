package contagiongo

import (
	"fmt"
	"math/rand"
	"sync"

	rv "github.com/kentwait/randomvariate"
)

// Spreader is an interface for methods that facilitate spreading of pathogens.
type Spreader interface {
	// SpreaderDescription returns a short description about the spreader type
	SpreaderDescription() string

	// TransmissionProb returns the probability that a transmission event
	// occurs between one host and one neighbor (per capita event) occurs.
	TransmissionProb() float64

	// TransmissionSize returns the number of pathogens transmitted given
	// a transmission event occurs.
	TransmissionSize() int
}

type spreader struct {
	desc string
	prob float64
	size int
}

// NewSpreader creates a new spreader with transmission probability p
// and transmission size n.
func NewSpreader(p float64, n int, desc ...string) (Spreader, error) {
	var description string
	if len(desc) > 0 {
		description = desc[0]
	}
	if p < 0 {
		return nil, fmt.Errorf(InvalidFloatParameterError, "transmission probability", p, "p < 0")
	} else if p > 1.0 {
		return nil, fmt.Errorf(InvalidFloatParameterError, "transmission probability", p, "p > 1.0")
	}
	if n < 0 {
		return nil, fmt.Errorf(InvalidIntParameterError, "transmission size", n, "n < 0")
	}
	return &spreader{description, p, n}, nil
}

func (s *spreader) SpreaderDescription() string {
	return s.desc
}

func (s *spreader) TransmissionProb() float64 {
	return s.prob
}

func (s *spreader) TransmissionSize() int {
	return s.size
}

// PathogenTransmitter transmits the pathogen to its neighboring host/s.
func PathogenTransmitter(source, neighbor EpidemicHost, pathogenPopSize, neighborStatus int, c chan<- TransParams, wg *sync.WaitGroup, infectableStatuses ...int) {
	defer wg.Done()
	// Check if migration size if larger than the current population size
	// If larger, skip
	migrationSize := source.TransmissionSize()
	if migrationSize > pathogenPopSize {
		return
	}
	// Check if neighbor has an infectable status
	// If current status matches at least one infectable status, proceed
	var infectable bool
	for _, infectableStatus := range infectableStatuses {
		if neighborStatus == infectableStatus {
			infectable = true
			break
		}
	}
	if !infectable {
		return
	}
	// Determine if tranmission occurs or not based on source's
	// transmission probability
	transmissionProb := source.TransmissionProb()
	if rv.Binomial(1, transmissionProb) != 1.0 {
		return
	}
	// If transmission occurs, randomly pick pathogens to transmit
	for _, pathogenPos := range rand.Perm(pathogenPopSize)[:migrationSize] {
		if pathogen := source.Pathogen(pathogenPos); pathogen != nil {
			c <- TransParams{neighbor, pathogen}
		}
	}
}
