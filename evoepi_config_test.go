package contagiongo

import (
	"fmt"
	"testing"
)

func TestEvoEpiConfig_Validate(t *testing.T) {
	conf := sampleEvoEpiConfig()
	err := conf.Validate()
	if err != nil {
		t.Error(err)
	}
	// TODO: Add errors
}

func TestEvoEpiConfig_NewSimulation(t *testing.T) {
	conf := sampleEvoEpiConfig()
	err := conf.Validate()
	if err != nil {
		t.Error(err)
	}
	sim, err := conf.NewSimulation()
	if err != nil {
		t.Error(err)
	}

	// Expected values
	numHosts := 10

	// Tests
	if l := len(sim.HostMap()); l != numHosts {
		t.Errorf(UnequalIntParameterError, "number of hosts", numHosts, l)
	}
}

func TestIntrahostModelConfig_Validate(t *testing.T) {
	conf := sampleIntrahostModelConfig()
	err := conf.Validate()
	if err != nil {
		t.Error(err)
	}
	// TODO: Add errors
}

func TestIntrahostModelConfig_CreateModel(t *testing.T) {
	conf := sampleIntrahostModelConfig()
	err := conf.Validate()
	if err != nil {
		t.Error(err)
	}
	modelID := 0
	model, err := conf.CreateModel(modelID)
	if err != nil {
		t.Error(err)
	}

	// Tests
	if i := model.ModelID(); i != 0 {
		t.Errorf(UnequalIntParameterError, "model ID", 0, i)
	}
	if s := model.NextPathogenPopSize(1); s != conf.ConstantPopSize {
		t.Errorf(UnequalIntParameterError, "next pathogen population size", conf.ConstantPopSize, s)
	}
	if s := model.ReplicationMethod(); s != "relative" {
		t.Errorf(UnequalStringParameterError, "replication method", "relative", s)
	}
	// TODO: Add cases
}

func TestFitnessModelConfig_Validate(t *testing.T) {
	conf := sampleFitnessModelConfig()
	err := conf.Validate()
	if err != nil {
		t.Error(err)
	}

	// Add errors
	// Invalid fitness_model keyword
	conf = sampleFitnessModelConfig()
	conf.FitnessModel = "intrinsic"
	err = conf.Validate()
	if err == nil {
		t.Errorf(ExpectedErrorWhileError, "validating fitness_model keyword")
	}
	// Invalid fitness_model_path
	conf = sampleFitnessModelConfig()
	conf.FitnessModelPath = "examples/mock.txt"
	err = conf.Validate()
	if err == nil {
		t.Errorf(ExpectedErrorWhileError, "validating fitness_model_path")
	}
}

func TestFitnessModelConfig_CreateModel(t *testing.T) {
	// Multiplicative
	conf := sampleFitnessModelConfig()
	err := conf.Validate()
	if err != nil {
		t.Error(err)
	}
	modelID := 0
	model, err := conf.CreateModel(modelID)
	if err != nil {
		t.Error(err)
	}

	if i := model.ModelID(); i != 0 {
		t.Errorf(UnequalIntParameterError, "model ID", 0, i)
	}
	fitness, err := model.ComputeFitness(sampleSequence(100)...)
	if err != nil {
		t.Error(err)
	}
	if fitness != 0 {
		t.Errorf(UnequalFloatParameterError, "log fitness value", 0.0, fitness)
	}

	// Additive
	conf = sampleFitnessModelConfig()
	conf.FitnessModelPath = "examples/test1.sir.fm.additive.txt"
	conf.FitnessModel = "additive"
	err = conf.Validate()
	if err != nil {
		t.Error(err)
	}
	modelID = 0
	model, err = conf.CreateModel(modelID)
	if err != nil {
		t.Error(err)
	}

	if i := model.ModelID(); i != 0 {
		t.Errorf(UnequalIntParameterError, "model ID", 0, i)
	}
	fitness, err = model.ComputeFitness(sampleSequence(100)...)
	if err != nil {
		t.Error(err)
	}
	if fmt.Sprintf("%.6f", fitness) != fmt.Sprintf("%.6f", 2.) {
		t.Errorf(UnequalFloatParameterError, "decimal fitness value", 2.0, fitness)
	}
}
