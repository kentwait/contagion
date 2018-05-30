package contagiongo

import (
	"fmt"
	"testing"
)

func TestConstantPopModel(t *testing.T) {
	model := new(ConstantPopModel)
	model.mutationRate = 1e-5
	model.transitionMatrix = [][]float64{
		[]float64{0, 1},
		[]float64{1, 0},
	}
	model.recombinationRate = 0
	model.popSize = 1000

	if exp := model.mutationRate; exp != model.MutationRate() {
		t.Errorf(UnequalFloatParameterError, "mutation rate", exp, model.MutationRate())
	}
	for i, row := range model.transitionMatrix {
		for j := range row {
			if exp := model.transitionMatrix[i][j]; exp != model.TransitionMatrix()[i][j] {
				t.Errorf(UnequalFloatParameterError, fmt.Sprintf("transition rate %d -> %d", i, j), exp, model.TransitionMatrix()[i][j])
			}
		}
	}

	if exp := model.recombinationRate; exp != model.RecombinationRate() {
		t.Errorf(UnequalFloatParameterError, "recombination rate", exp, model.RecombinationRate())
	}
	if exp := model.popSize; exp != model.MaxPathogenPopSize() {
		t.Errorf(UnequalIntParameterError, "maximum population size", exp, model.MaxPathogenPopSize())
	}
	n := 1
	if exp := model.popSize; exp != model.NextPathogenPopSize(n) {
		t.Errorf(UnequalIntParameterError, "next population size", exp, model.NextPathogenPopSize(n))
	}
}

func TestBevertonHoltThresholdPopModel(t *testing.T) {
	model := new(BevertonHoltThresholdPopModel)
	model.mutationRate = 1e-5
	model.transitionMatrix = [][]float64{
		[]float64{0, 1},
		[]float64{1, 0},
	}
	model.recombinationRate = 0
	model.maxPopSize = 1000
	model.growthRate = 2

	if exp := model.mutationRate; exp != model.MutationRate() {
		t.Errorf(UnequalFloatParameterError, "mutation rate", exp, model.MutationRate())
	}
	for i, row := range model.transitionMatrix {
		for j := range row {
			if exp := model.transitionMatrix[i][j]; exp != model.TransitionMatrix()[i][j] {
				t.Errorf(UnequalFloatParameterError, fmt.Sprintf("transition rate %d -> %d", i, j), exp, model.TransitionMatrix()[i][j])
			}
		}
	}

	if exp := model.recombinationRate; exp != model.RecombinationRate() {
		t.Errorf(UnequalFloatParameterError, "recombination rate", exp, model.RecombinationRate())
	}
	if exp := model.maxPopSize; exp != model.MaxPathogenPopSize() {
		t.Errorf(UnequalIntParameterError, "maximum population size", exp, model.MaxPathogenPopSize())
	}
	// Different cases
	ns := []int{0, 1, 2, 10, 999, 1000, 2000}
	exps := []int{0, 2, 4, 20, 1000, 1000, 1000}
	for i := range ns {
		n := ns[i]
		if exp := exps[i]; exp != model.NextPathogenPopSize(n) {
			t.Errorf(UnequalIntParameterError, "next population size", exp, model.NextPathogenPopSize(n))
		}
	}
}

func TestFitnessDependentPopModel(t *testing.T) {
	model := new(FitnessDependentPopModel)
	model.mutationRate = 1e-5
	model.transitionMatrix = [][]float64{
		[]float64{0, 1},
		[]float64{1, 0},
	}
	model.recombinationRate = 0
	model.maxPopSize = 1000

	if exp := model.mutationRate; exp != model.MutationRate() {
		t.Errorf(UnequalFloatParameterError, "mutation rate", exp, model.MutationRate())
	}
	for i, row := range model.transitionMatrix {
		for j := range row {
			if exp := model.transitionMatrix[i][j]; exp != model.TransitionMatrix()[i][j] {
				t.Errorf(UnequalFloatParameterError, fmt.Sprintf("transition rate %d -> %d", i, j), exp, model.TransitionMatrix()[i][j])
			}
		}
	}

	if exp := model.recombinationRate; exp != model.RecombinationRate() {
		t.Errorf(UnequalFloatParameterError, "recombination rate", exp, model.RecombinationRate())
	}
	if exp := model.maxPopSize; exp != model.MaxPathogenPopSize() {
		t.Errorf(UnequalIntParameterError, "maximum population size", exp, model.MaxPathogenPopSize())
	}
	n := 1
	if exp := -1; exp != model.NextPathogenPopSize(n) {
		t.Errorf(UnequalIntParameterError, "next population size", exp, model.NextPathogenPopSize(n))
	}
}
