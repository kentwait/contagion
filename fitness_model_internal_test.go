package contagiongo

import "testing"

func TestMultiplicativeFM(t *testing.T) {
	fitnessTests := []struct {
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
			sequence: []uint8{0, 1, 0, 1, 0, 1, 0, 1, 0, 1},
			output:   0,
		},
	}
}
