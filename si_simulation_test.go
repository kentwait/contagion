package contagiongo

import (
	"fmt"
	"os"
	"testing"
)

func CreateTestSim() (*SISimulation, *EvoEpiConfig, *CSVLogger, error) {
	conf := sampleEvoEpiConfig()
	err := conf.Validate()
	if err != nil {
		return nil, nil, nil, err
	}
	logger := NewCSVLogger(conf.LogPath())
	sim, err := NewSISimulation(conf, logger)
	if err != nil {
		return nil, nil, nil, err
	}
	return sim, conf, logger, nil
}

func TestSISimulation_Update(t *testing.T) {
	sim, _, logger, err := CreateTestSim()
	if err != nil {
		t.Error(err)
	}
	originalStatus := sim.Epidemic.HostStatus(0)
	sim.Update(0)
	newStatus := sim.Epidemic.HostStatus(0)

	// Host 0 is supposed to be infected
	if originalStatus == newStatus {
		t.Errorf(EqualIntParameterError, "host 0 status before and after", originalStatus, newStatus)
	}
	if newStatus != InfectedStatusCode {
		t.Errorf(UnequalIntParameterError, "host 0 status", InfectedStatusCode, newStatus)
	}
	// Remove written files
	os.Remove(logger.statusPath)
	os.Remove(logger.genotypeFreqPath)
}

func TestSISimulation_Process(t *testing.T) {
	sim, _, logger, err := CreateTestSim()
	if err != nil {
		t.Error(err)
	}
	sim.Update(0)
	// Get genotype UID
	originalUID := sim.Host(0).Pathogen(0).GenotypeUID()
	sim.Process(1)
	sim.Update(1)

	if sim.Epidemic.HostStatus(0) != InfectedStatusCode {
		t.Errorf(UnequalIntParameterError, "host 0 status", InfectedStatusCode, sim.Epidemic.HostStatus(0))
	}
	if popSize := sim.Epidemic.Host(0).PathogenPopSize(); popSize != 100 {
		t.Errorf(UnequalIntParameterError, "host 0 pathogen population size", 100, popSize)
	}
	count := 0
	for _, p := range sim.Host(0).Pathogens() {
		if p.GenotypeUID() == originalUID {
			count++
		}
	}
	if count == 100 {
		t.Errorf(EqualIntParameterError,
			fmt.Sprintf("frequency of genotype %s before and after", originalUID.String()),
			100,
			count,
		)
	}
	// Remove written files
	os.Remove(logger.statusPath)
	os.Remove(logger.genotypeFreqPath)
	os.Remove(logger.mutationPath)
}

func TestSISimulation_Transmit(t *testing.T) {
	sim, _, logger, err := CreateTestSim()
	if err != nil {
		t.Error(err)
	}
	sim.Update(0)
	// Check if hosts were cleared
	for i, host := range sim.HostMap() {
		popSize := host.PathogenPopSize()
		if i == 0 && popSize != 100 {
			t.Errorf(UnequalIntParameterError,
				fmt.Sprintf("host %d pathogen population size", i),
				100,
				popSize)
		} else if i != 0 && popSize > 0 {
			t.Errorf(UnequalIntParameterError,
				fmt.Sprintf("host %d pathogen population size", i),
				0,
				popSize)
		}
	}
	sim.Transmit(0)
	sim.Update(1)

	// TODO: Read output files
	// Check if hosts were cleared
	for i, host := range sim.HostMap() {
		popSize := host.PathogenPopSize()
		if i == 0 && popSize != 100 {
			t.Errorf(UnequalIntParameterError,
				fmt.Sprintf("host %d pathogen population size", i),
				100,
				popSize)
		} else if i != 0 && popSize > 1 {
			t.Errorf(UnequalIntParameterError,
				fmt.Sprintf("host %d pathogen population size", i),
				0,
				popSize)
		}
	}
	// Remove written files
	os.Remove(logger.statusPath)
	os.Remove(logger.genotypeFreqPath)
	os.Remove(logger.transmissionPath)
}

func TestSISimulation_Finalize(t *testing.T) {
	sim, _, logger, err := CreateTestSim()
	if err != nil {
		t.Error(err)
	}
	sim.Update(0)
	sim.Finalize()

	// TODO: Read output files
	// Check if hosts were cleared
	for _, host := range sim.HostMap() {
		if popSize := host.PathogenPopSize(); popSize > 0 {
			t.Errorf(UnequalIntParameterError, "host 0 pathogen population size", 0, popSize)
		}
	}
	// Remove written files
	os.Remove(logger.statusPath)
	os.Remove(logger.genotypeFreqPath)
	os.Remove(logger.mutationPath)
	os.Remove(logger.genotypePath)
	os.Remove(logger.genotypeNodePath)
}
