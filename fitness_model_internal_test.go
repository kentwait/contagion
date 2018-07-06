package contagiongo

import "testing"

var fitnessTests = []struct {
	matrix   map[int]map[uint8]float64
	sequence []uint8
	output   float64
}{
	{
		matrix: map[int]map[uint8]float64{
			0: map[uint8]float64{0: 0, 1: 0},
			1: map[uint8]float64{0: 0, 1: 0},
			2: map[uint8]float64{0: 0, 1: 0},
			3: map[uint8]float64{0: 0, 1: 0},
			4: map[uint8]float64{0: 0, 1: 0},
			5: map[uint8]float64{0: 0, 1: 0},
			6: map[uint8]float64{0: 0, 1: 0},
			7: map[uint8]float64{0: 0, 1: 0},
			8: map[uint8]float64{0: 0, 1: 0},
			9: map[uint8]float64{0: 0, 1: 0},
		},
		sequence: []uint8{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		output:   0,
	},
	{
		matrix: map[int]map[uint8]float64{
			0: map[uint8]float64{0: 0, 1: 0},
			1: map[uint8]float64{0: 0, 1: 0},
			2: map[uint8]float64{0: 0, 1: 0},
			3: map[uint8]float64{0: 0, 1: 0},
			4: map[uint8]float64{0: 0, 1: 0},
			5: map[uint8]float64{0: 0, 1: 0},
			6: map[uint8]float64{0: 0, 1: 0},
			7: map[uint8]float64{0: 0, 1: 0},
			8: map[uint8]float64{0: 0, 1: 0},
			9: map[uint8]float64{0: 0, 1: 0},
		},
		sequence: []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		output:   0,
	},
	{
		matrix: map[int]map[uint8]float64{
			0: map[uint8]float64{0: 0, 1: 0},
			1: map[uint8]float64{0: 0, 1: 0},
			2: map[uint8]float64{0: 0, 1: 0},
			3: map[uint8]float64{0: 0, 1: 0},
			4: map[uint8]float64{0: 0, 1: 0},
			5: map[uint8]float64{0: 0, 1: 0},
			6: map[uint8]float64{0: 0, 1: 0},
			7: map[uint8]float64{0: 0, 1: 0},
			8: map[uint8]float64{0: 0, 1: 0},
			9: map[uint8]float64{0: 0, 1: 0},
		},
		sequence: []uint8{0, 1, 0, 1, 0, 1, 0, 1, 0, 1},
		output:   0,
	},
	{
		matrix: map[int]map[uint8]float64{
			0: map[uint8]float64{0: 0, 1: 0},
			1: map[uint8]float64{0: 0, 1: 0},
			2: map[uint8]float64{0: 0, 1: 0},
			3: map[uint8]float64{0: 0, 1: 0},
			4: map[uint8]float64{0: 0, 1: 0},
			5: map[uint8]float64{0: 0, 1: 0},
			6: map[uint8]float64{0: 0, 1: 0},
			7: map[uint8]float64{0: 0, 1: 0},
			8: map[uint8]float64{0: 0, 1: 0},
			9: map[uint8]float64{0: 0, 1: 0},
		},
		sequence: []uint8{0, 0, 0, 0, 0, 1, 1, 1, 1, 1},
		output:   0,
	},
}

var newFitnessMatrixErrors = []struct {
	desc     string
	matrix   map[int]map[uint8]float64
	sequence []uint8
	output   float64
}{
	{
		desc:     "empty matrix",
		matrix:   map[int]map[uint8]float64{},
		sequence: []uint8{0, 1, 0, 1, 0, 1, 0, 1, 0, 1},
		output:   0,
	},
	{
		desc: "empty row",
		matrix: map[int]map[uint8]float64{
			0: map[uint8]float64{0: 0, 1: 0},
			1: map[uint8]float64{},
			2: map[uint8]float64{0: 0, 1: 0},
			3: map[uint8]float64{0: 0, 1: 0},
			4: map[uint8]float64{0: 0, 1: 0},
			5: map[uint8]float64{0: 0, 1: 0},
			6: map[uint8]float64{0: 0, 1: 0},
			7: map[uint8]float64{0: 0, 1: 0},
			8: map[uint8]float64{0: 0, 1: 0},
			9: map[uint8]float64{0: 0, 1: 0},
		},
		sequence: []uint8{0, 1, 0, 1, 0, 1, 0, 1, 0, 1},
		output:   0,
	},
}

var fitnessMatrixErrors = []struct {
	desc     string
	matrix   map[int]map[uint8]float64
	sequence []uint8
	output   float64
}{
	{
		desc: "zero-length sequence",
		matrix: map[int]map[uint8]float64{
			0: map[uint8]float64{0: 0, 1: 0},
			1: map[uint8]float64{0: 0, 1: 0},
			2: map[uint8]float64{0: 0, 1: 0},
			3: map[uint8]float64{0: 0, 1: 0},
			4: map[uint8]float64{0: 0, 1: 0},
			5: map[uint8]float64{0: 0, 1: 0},
			6: map[uint8]float64{0: 0, 1: 0},
			7: map[uint8]float64{0: 0, 1: 0},
			8: map[uint8]float64{0: 0, 1: 0},
			9: map[uint8]float64{0: 0, 1: 0},
		},
		sequence: []uint8{},
		output:   0,
	},
	{
		desc: "sequence length greater than matrix size",
		matrix: map[int]map[uint8]float64{
			0: map[uint8]float64{0: 0, 1: 0},
			1: map[uint8]float64{0: 0, 1: 0},
			2: map[uint8]float64{0: 0, 1: 0},
			3: map[uint8]float64{0: 0, 1: 0},
			4: map[uint8]float64{0: 0, 1: 0},
			5: map[uint8]float64{0: 0, 1: 0},
			6: map[uint8]float64{0: 0, 1: 0},
			7: map[uint8]float64{0: 0, 1: 0},
			8: map[uint8]float64{0: 0, 1: 0},
			9: map[uint8]float64{0: 0, 1: 0},
		},
		sequence: []uint8{0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0},
		output:   0,
	},
	{
		desc: "invalid character",
		matrix: map[int]map[uint8]float64{
			0: map[uint8]float64{0: 0, 1: 0},
			1: map[uint8]float64{0: 0, 1: 0},
			2: map[uint8]float64{0: 0, 1: 0},
			3: map[uint8]float64{0: 0, 1: 0},
			4: map[uint8]float64{0: 0, 1: 0},
			5: map[uint8]float64{0: 0, 1: 0},
			6: map[uint8]float64{0: 0, 1: 0},
			7: map[uint8]float64{0: 0, 1: 0},
			8: map[uint8]float64{0: 0, 1: 0},
			9: map[uint8]float64{0: 0, 1: 0},
		},
		sequence: []uint8{0, 2, 0, 1, 0, 1, 0, 1, 0, 1},
		output:   0,
	},
}

func TestMultiplicativeFM(t *testing.T) {
	for _, tt := range fitnessTests {
		t.Run("", func(t *testing.T) {
			fm, err := NewMultiplicativeFM(0, "neutral multiplicative", tt.matrix)
			if err != nil {
				t.Errorf("error creating multiplicative fitness matrix: %v", err)
			}
			fitness, err := fm.ComputeFitness(tt.sequence...)
			if err != nil {
				t.Errorf("error computing fitness: %v", err)
			}
			if fitness != tt.output {
				t.Errorf("expected %f, got %f instead", tt.output, fitness)
			}
		})
	}
	// new fitness errors
	for _, tt := range newFitnessMatrixErrors {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := NewMultiplicativeFM(0, "neutral multiplicative", tt.matrix)
			if err == nil {
				t.Errorf("expected error: %s", tt.desc)
			}
		})
	}
	// fitness calculation errors
	for _, tt := range fitnessMatrixErrors {
		t.Run(tt.desc, func(t *testing.T) {
			fm, err := NewMultiplicativeFM(0, "neutral multiplicative", tt.matrix)
			_, err = fm.ComputeFitness(tt.sequence...)
			if err == nil {
				t.Errorf("expected error: %s", tt.desc)
			}
		})
	}
}
