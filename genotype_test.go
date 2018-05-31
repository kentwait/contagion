package contagiongo

import (
	"math/rand"
	"testing"
)

func TestNewGenotype(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf(UnexpectedErrorWhileError, "calling NewGenotype constructor", err)
		}
	}()
	sequence := []int{0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1}
	NewGenotype(sequence)
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
	if l := len(tree.(*genotypeTree).genotypes); l != 1 {
		t.Errorf(UnequalIntParameterError, "size of genotype map", 1, l)
	}

	sequence[0] = 1
	tree.NewNode(sequence, p1)
	if l := tree.Set().Size(); l != 2 {
		t.Errorf(UnequalIntParameterError, "size of genotype set", 2, l)
	}
	if l := len(tree.(*genotypeTree).genotypes); l != 2 {
		t.Errorf(UnequalIntParameterError, "size of genotype map", 2, l)
	}
}
