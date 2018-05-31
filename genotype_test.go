package contagiongo

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestNewGenotype(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf(UnexpectedErrorWhileError, "calling NewGenotype constructor", err)
		}
	}()
	sequence := []uint8{0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1}
	NewGenotype(sequence)
}

func TestGenotype_Fitness(t *testing.T) {
	sites := 100
	fm := NeutralMultiplicativeFM(0, "m", sites, 2)
	genotype := NewGenotype(sampleSequence(sites))
	logFitness := genotype.Fitness(fm)

	if logFitness != 0.0 {
		t.Errorf(UnequalFloatParameterError, "log fitness", 0.0, logFitness)
	}
	// Fitness has been assigned during the first call
	// If fitness model has the same ID, then just recalls previous value.
	logFitness = genotype.Fitness(fm)
	if logFitness != 0.0 {
		t.Errorf(UnequalFloatParameterError, "log fitness", 0.0, logFitness)
	}
}

func TestEmptyGenotypeSet(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf(UnexpectedErrorWhileError, "calling EmptyGenotypeSet constructor", err)
		}
	}()
	EmptyGenotypeSet()
}

func TestGenotypeSet_Add(t *testing.T) {
	g := sampleGenotype()
	set := EmptyGenotypeSet()
	set.Add(g)
	if l := set.Size(); l != 1 {
		t.Errorf(UnequalIntParameterError, "size of genotype set", 1, l)
	}
	set.Add(g)
	if l := set.Size(); l != 1 {
		t.Errorf(UnequalIntParameterError, "size of genotype set", 1, l)
	}
}

func TestGenotypeSet_AddSequence(t *testing.T) {
	sequence := sampleSequence(1000)
	set := EmptyGenotypeSet()
	set.AddSequence(sequence)
	if l := set.Size(); l != 1 {
		t.Errorf(UnequalIntParameterError, "size of genotype set", 1, l)
	}
	set.AddSequence(sequence)
	if l := set.Size(); l != 1 {
		t.Errorf(UnequalIntParameterError, "size of genotype set", 1, l)
	}
}

func TestGenotypeSet_Remove(t *testing.T) {
	sequence := sampleSequence(1000)
	set := EmptyGenotypeSet()
	set.AddSequence(sequence)
	set.Remove(sequence)
	if l := set.Size(); l != 0 {
		t.Errorf(UnequalIntParameterError, "size of genotype set", 0, l)
	}
}

func TestNewGenotypeNode(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf(UnexpectedErrorWhileError, "calling NewGenotype constructor", err)
		}
	}()
	sequence := sampleSequence(1000)
	set := EmptyGenotypeSet()
	newGenotypeNode(sequence, set)
}

func TestNewGenotypeNode_Getters(t *testing.T) {
	set := EmptyGenotypeSet()
	p1 := newGenotypeNode(sampleSequence(1000), set)
	p2 := newGenotypeNode(sampleSequence(1000), set, p1)
	p3 := newGenotypeNode(sampleSequence(1000), set, p2)

	if l := len(p2.Parents()); l != 1 {
		t.Errorf(UnequalIntParameterError, "number of parents", 1, l)
	}
	if n := p2.Parents()[0]; n.UID() != p1.UID() {
		t.Errorf(UnequalStringParameterError, "parent UID", fmt.Sprint(p1.UID()), fmt.Sprint(n.UID()))
	}
	if l := len(p2.Children()); l != 1 {
		t.Errorf(UnequalIntParameterError, "number of childen", 1, l)
	}
	if n := p2.Children()[0]; n.UID() != p3.UID() {
		t.Errorf(UnequalStringParameterError, "child UID", fmt.Sprint(p3.UID()), fmt.Sprint(n.UID()))
	}
}

func TestEmptyGenotypeTree(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf(UnexpectedErrorWhileError, "calling EmptyGenotypeSet constructor", err)
		}
	}()
	EmptyGenotypeTree()
}

func TestGenotypeTree_NewNode(t *testing.T) {
	rand.Seed(0)
	tree := EmptyGenotypeTree()
	sequence := sampleSequence(1000)
	p1 := tree.NewNode(sequence)
	if l := tree.Set().Size(); l != 1 {
		t.Errorf(UnequalIntParameterError, "size of genotype set", 1, l)
	}
	if l := len(tree.Nodes()); l != 1 {
		t.Errorf(UnequalIntParameterError, "size of genotype map", 1, l)
	}

	sequence[0] = 1
	tree.NewNode(sequence, p1)
	if l := tree.Set().Size(); l != 2 {
		t.Errorf(UnequalIntParameterError, "size of genotype set", 2, l)
	}
	if l := len(tree.Nodes()); l != 2 {
		t.Errorf(UnequalIntParameterError, "size of genotype map", 2, l)
	}
}
