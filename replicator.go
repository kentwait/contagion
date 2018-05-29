package contagiongo

import (
	"fmt"
	"math"
)

// Replicator is an interface for methods related to the intrahost population.
type Replicator interface {
	// ReplicatorDescription returns a short description about the Replicator type
	ReplicatorDescription() string

	// MaxPathogenPopSize returns the maximum number of pathogens allowed within
	// a single host of this particular host type.
	MaxPathogenPopSize() int

	// NextPathogenPopSize returns the pathogen population size for the next
	// generation of pathogens given the current population size.
	NextPathogenPopSize(n int) int

	// TODO: transfer to host
	// TimeInterval returns the number of generations the host remains
	// in a given state.
	// TimeInterval(status int) int
}

type bhtReplicator struct {
	desc string
	r    float64 // growth rate
	k    int     // carrying capacity
}

// NewBhtReplicator creates a new bhtReplicator with growth rate r and
// maximum population size k.
func NewBhtReplicator(r float64, k int, desc ...string) (Replicator, error) {
	var description string
	if len(desc) > 0 {
		description = desc[0]
	}
	if r < 0 {
		return nil, fmt.Errorf(InvalidFloatParameterError, "growth rate", r, "r < 0")
	}
	if k < 0 {
		return nil, fmt.Errorf(InvalidIntParameterError, "maximum population size", k, "k < 0")
	}
	return &bhtReplicator{description, r, k}, nil
}

func (c *bhtReplicator) ReplicatorDescription() string {
	return c.desc
}

func (c *bhtReplicator) GrowthRate() float64 {
	return c.r
}

func (c *bhtReplicator) MaxPathogenPopSize() int {
	return c.k
}

func (c *bhtReplicator) NextPathogenPopSize(n int) int {
	n64 := float64(n)
	k64 := float64(c.k)
	res := (c.r * n64 * k64) / (k64 + ((c.r - 1.0) * n64))
	roundedRes := int(math.Ceil(res))
	if c.k > roundedRes {
		return roundedRes
	}
	return c.k
}

type constReplicator struct {
	desc string
	k    int
}

// NewConstReplicator creates a new constReplicator with maximum population size k.
func NewConstReplicator(k int, desc ...string) (Replicator, error) {
	var description string
	if len(desc) > 0 {
		description = desc[0]
	}
	if k < 0 {
		return nil, fmt.Errorf(InvalidIntParameterError, "maximum population size", k, "k < 0")
	}
	return &constReplicator{description, k}, nil
}

func (c *constReplicator) ReplicatorDescription() string {
	return c.desc
}

func (c *constReplicator) MaxPathogenPopSize() int {
	return c.k
}

func (c *constReplicator) NextPathogenPopSize(n int) int {
	if n == 0 {
		return 0
	}
	return c.k
}

// TODO: Add Rikers, Logistic map
