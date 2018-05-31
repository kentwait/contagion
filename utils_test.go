package contagiongo

import (
	"fmt"
	"testing"
)

func TestLoadSequences(t *testing.T) {
	hostID := 0
	infectedHosts := 1
	pathogenPopSize := 100

	path := "examples/test1.sir.pathogens.fa"
	pathogenHostMap, err := LoadSequences(path)
	if err != nil {
		t.Error(err)
	}

	if l := len(pathogenHostMap); l != infectedHosts {
		t.Errorf(UnequalIntParameterError, "number of infected hosts", infectedHosts, l)
	}
	if _, ok := pathogenHostMap[hostID]; !ok {
		t.Errorf(IntKeyNotFoundError, hostID)
	}
	if pathogens := pathogenHostMap[hostID]; len(pathogens) != pathogenPopSize {
		t.Errorf(UnequalIntParameterError, "pathogen population size in host 0", pathogenPopSize, len(pathogens))
	}
}

func TestLoadFitnessMatrix(t *testing.T) {
	sites := 100
	alleles := 2
	path := "examples/test1.sir.fm.txt"
	matrix, err := LoadFitnessMatrix(path)
	if err != nil {
		t.Error(err)
	}

	if l := len(matrix); l != sites {
		t.Errorf(UnequalIntParameterError, "number of sites", sites, l)
	}
	for i, row := range matrix {
		fmt.Println(i, row)
		if l := len(row); l != alleles {
			t.Errorf(UnequalIntParameterError, "number of alleles", alleles, l)
		}
	}
}

func TestLoadAdjacencyMatrix(t *testing.T) {
	numConnections := 44
	connectedHosts := 10
	path := "examples/test1.sir.network.txt"
	net, err := LoadAdjacencyMatrix(path)
	if err != nil {
		t.Error(err)
	}

	am := net.(adjacencyMatrix)
	counter := 0
	for _, row := range am {
		for range row {
			counter++
		}
	}
	if counter != numConnections {
		t.Errorf(UnequalIntParameterError, "number of connections", numConnections, counter)
	}
	if c := net.ConnectedPopSize(); c != connectedHosts {
		t.Errorf(UnequalIntParameterError, "number of connected hosts", connectedHosts, c)
	}
}
