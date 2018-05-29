package contagiongo

import (
	"fmt"
	"math"
)

// Colonizer is an interface for methods related to the intrahost population.
type Colonizer interface {
	// ColonizerDescription returns a short description about the colonizer type
	ColonizerDescription() string

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

type bhtColonizer struct {
	desc string
	r    float64 // growth rate
	k    int     // carrying capacity
}

// NewBhtColonizer creates a new bhtColonizer with growth rate r and
// maximum population size k.
func NewBhtColonizer(r float64, k int, desc ...string) (Colonizer, error) {
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
	return &bhtColonizer{description, r, k}, nil
}

func (c *bhtColonizer) ColonizerDescription() string {
	return c.desc
}

func (c *bhtColonizer) GrowthRate() float64 {
	return c.r
}

func (c *bhtColonizer) MaxPathogenPopSize() int {
	return c.k
}

func (c *bhtColonizer) NextPathogenPopSize(n int) int {
	n64 := float64(n)
	k64 := float64(c.k)
	res := (c.r * n64 * k64) / (k64 + ((c.r - 1.0) * n64))
	roundedRes := int(math.Ceil(res))
	if c.k > roundedRes {
		return roundedRes
	}
	return c.k
}

type constColonizer struct {
	desc string
	k    int
}

// NewConstColonizer creates a new constColonizer with maximum population size k.
func NewConstColonizer(k int, desc ...string) (Colonizer, error) {
	var description string
	if len(desc) > 0 {
		description = desc[0]
	}
	if k < 0 {
		return nil, fmt.Errorf(InvalidIntParameterError, "maximum population size", k, "k < 0")
	}
	return &constColonizer{description, k}, nil
}

func (c *constColonizer) ColonizerDescription() string {
	return c.desc
}

func (c *constColonizer) MaxPathogenPopSize() int {
	return c.k
}

func (c *constColonizer) NextPathogenPopSize(n int) int {
	if n == 0 {
		return 0
	}
	return c.k
}

// TODO: Add Rikers, Logistic map
