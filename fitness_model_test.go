package contagiongo

import (
	"fmt"
	"math"
	"testing"
)

func TestNewMultiplicativeFM(t *testing.T) {
	id := 1
	name := "multiplicative"
	matrix := make(map[int]map[int]float64)
	for i := 0; i < 1000; i++ {
		matrix[i] = map[int]float64{0: math.Log(0.4), 1: math.Log(0.3), 2: math.Log(0.2), 3: math.Log(0.1)}
	}

	fm := NewMultiplicativeFM(id, name, matrix)
	for i, row := range fm.(*multiplicativeFM).matrix {
		for j := range row {
			if res := fm.(*multiplicativeFM).matrix[i][j]; res != matrix[i][j] {
				t.Errorf(UnequalFloatParameterError, "fitness", matrix[i][j], res)
			}
		}
	}
}

func TestMultiplicativeFM_Fitness(t *testing.T) {
	matrix := make(map[int]map[int]float64)
	for i := 0; i < 1000; i++ {
		matrix[i] = map[int]float64{0: math.Log(0.4), 1: math.Log(0.3), 2: math.Log(0.2), 3: math.Log(0.1)}
	}
	sequence := sampleSequenceLong()
	id := 1
	name := "multiplicative"
	fm := NewMultiplicativeFM(id, name, matrix)
	if res, err := fm.ComputeFitness(sequence...); err != nil {
		t.Fatalf(UnexpectedErrorWhileError, "getting the fitness of the sequence", err)
	} else if err == nil && fmt.Sprintf("%.6f", res) != "-1060.131768" {
		t.Errorf(UnequalFloatParameterError, "fitness", -1060.131768, res)
	}
}

func TestNewAdditiveFM(t *testing.T) {
	matrix := make(map[int]map[int]float64)
	for i := 0; i < 1000; i++ {
		matrix[i] = map[int]float64{0: 0.4, 1: 0.3, 2: 0.2, 3: 0.1}
	}
	id := 1
	name := "additive"
	fm := NewAdditiveFM(id, name, matrix)
	for i, row := range fm.(*additiveFM).matrix {
		for j := range row {
			if res := fm.(*additiveFM).matrix[i][j]; res != matrix[i][j] {
				t.Errorf(UnequalFloatParameterError, "fitness", matrix[i][j], res)
			}
		}
	}
}

func TestAdditiveFM_Fitness(t *testing.T) {
	matrix := make(map[int]map[int]float64)
	for i := 0; i < 1000; i++ {
		matrix[i] = map[int]float64{0: 0.4, 1: 0.3, 2: 0.2, 3: 0.1}
	}
	sequence := sampleSequenceLong()
	id := 1
	name := "additive"
	fm := NewAdditiveFM(id, name, matrix)
	if res, err := fm.ComputeFitness(sequence...); err != nil {
		t.Fatalf(UnexpectedErrorWhileError, "getting the fitness of the sequence", err)
	} else if err == nil && fmt.Sprintf("%.6f", res) != "350.000000" {
		t.Errorf(UnequalFloatParameterError, "fitness", 350.000000, res)
	}
}
