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
	NewGenotypeNode(sequence, set)
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

// func TestNewSequenceTree(t *testing.T) {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			t.Fatalf(UnexpectedErrorWhileError, "calling NewSequenceTree constructor", err)
// 		}
// 	}()
// 	rootSequence := []int{0, 1, 0, 1, 0, 1, 0, 1, 0, 1}
// 	tree := NewSequenceTree(rootSequence)
// 	sequenceTree := tree.(*sequenceTree)

// 	if l := len(sequenceTree.roots); l != 1 {
// 		t.Errorf(UnequalIntParameterError, "number of roots", 1, l)
// 	}
// 	for _, root := range sequenceTree.roots {
// 		if fmt.Sprintf("%v", root.sequence) != fmt.Sprintf("%v", rootSequence) {
// 			t.Errorf(UnequalStringParameterError, "integer sequence", fmt.Sprintf("%v", rootSequence), fmt.Sprintf("%v", root.sequence))
// 		}
// 	}
// }

// func sampleSequenceTree() *sequenceTree {
// 	rootSequence := []int{
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 		0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1,
// 	}
// 	tree := NewSequenceTree(rootSequence)
// 	return tree.(*sequenceTree)
// }

// func TestSequenceTree_NewSub(t *testing.T) {
// 	tree := sampleSequenceTree()
// 	var key ksuid.KSUID
// 	for k := range tree.pathogens {
// 		key = k
// 	}
// 	parent := tree.Sequence(key)
// 	if parent == nil {
// 		t.Fatalf(UnexpectedErrorWhileError, "getting parent", "nil")
// 	}
// 	position := 0
// 	state := 1
// 	child := tree.NewSub(parent, position, state)

// 	if l := len(tree.pathogens); l != 2 {
// 		t.Errorf(UnequalIntParameterError, "number of pathgoens", 2, l)
// 	}

// 	if res := child.Parents(); res[0] != parent {
// 		t.Errorf(UnequalStringParameterError, "parent sequence pointer", fmt.Sprintf("%p", res), fmt.Sprintf("%p", parent))
// 	}
// 	if res := child.Sequence()[position]; res != state {
// 		t.Errorf(UnequalIntParameterError, "state at position 0", state, res)
// 	}
// 	// TODO: Add test to check if keys duplicate
// }

// func TestSequenceTree_NewRoot(t *testing.T) {
// 	tree := sampleSequenceTree()
// 	rootSequence := []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
// 	tree.NewRoot(rootSequence)

// 	if l := len(tree.roots); l != 2 {
// 		t.Errorf(UnequalIntParameterError, "number of roots", 2, l)
// 	}

// 	var key ksuid.KSUID
// 	for k := range tree.pathogens {
// 		key = k
// 	}
// 	if fmt.Sprintf("%v", tree.roots[key].sequence) != fmt.Sprintf("%v", rootSequence) {
// 		t.Errorf(UnequalStringParameterError, "integer sequence", fmt.Sprintf("%v", rootSequence), fmt.Sprintf("%v", tree.roots[key].sequence))
// 	}
// 	// TODO: Add test to check if keys duplicate
// }
