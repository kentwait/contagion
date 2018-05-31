package contagiongo

import (
	"fmt"
	"math"
	"math/rand"
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
	rand.Seed(0)
	matrix := make(map[int]map[int]float64)
	for i := 0; i < 1000; i++ {
		matrix[i] = map[int]float64{0: math.Log(0.4), 1: math.Log(0.3), 2: math.Log(0.2), 3: math.Log(0.1)}
	}
	sequence := sampleSequence(1000)
	id := 1
	name := "multiplicative"
	fm := NewMultiplicativeFM(id, name, matrix)
	if res, err := fm.ComputeFitness(sequence...); err != nil {
		t.Fatalf(UnexpectedErrorWhileError, "getting the fitness of the sequence", err)
	} else if err == nil && fmt.Sprintf("%.6f", res) != "-1051.213624" {
		t.Errorf(UnequalFloatParameterError, "fitness", -1051.213624, res)
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
	rand.Seed(0)
	matrix := make(map[int]map[int]float64)
	for i := 0; i < 1000; i++ {
		matrix[i] = map[int]float64{0: 0.4, 1: 0.3, 2: 0.2, 3: 0.1}
	}
	sequence := sampleSequence(1000)
	id := 1
	name := "additive"
	fm := NewAdditiveFM(id, name, matrix)
	if res, err := fm.ComputeFitness(sequence...); err != nil {
		t.Fatalf(UnexpectedErrorWhileError, "getting the fitness of the sequence", err)
	} else if err == nil && fmt.Sprintf("%.6f", res) != "353.100000" {
		t.Errorf(UnequalFloatParameterError, "fitness", 353.100000, res)
	}
}

func TestNeutralMultiplicativeFM(t *testing.T) {
	sites := 100
	alleles := 2
	value := 0.0
	fm := NeutralMultiplicativeFM(0, "m", sites, alleles)
	cnt := 0
	for _, row := range fm.(*multiplicativeFM).matrix {
		if len(row) != alleles {
			t.Errorf(UnequalIntParameterError, "number of alleles", alleles, len(row))
		}
		for _, v := range row {
			if v != value {
				t.Errorf(UnequalFloatParameterError, "fitness value", v, value)
			}
		}
		cnt++
	}
	if cnt != sites {
		t.Errorf(UnequalIntParameterError, "number of sites", sites, cnt)
	}
}

func TestNeutralAdditiveFM(t *testing.T) {
	sites := 100
	alleles := 2
	growthRate := 2
	value := float64(growthRate) / float64(sites)
	fm := NeutralAdditiveFM(0, "m", sites, alleles, growthRate)
	cnt := 0
	for _, row := range fm.(*additiveFM).matrix {
		if len(row) != alleles {
			t.Errorf(UnequalIntParameterError, "number of alleles", alleles, len(row))
		}
		for _, v := range row {
			if v != value {
				t.Errorf(UnequalFloatParameterError, "fitness value", v, value)
			}
		}
		cnt++
	}
	if cnt != sites {
		t.Errorf(UnequalIntParameterError, "number of sites", sites, cnt)
	}
}
