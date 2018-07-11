package contagiongo

// import "testing"

// func TestSequenceHost_Empty_Getters(t *testing.T) {
// 	id := 1
// 	typeID := 1
// 	h := EmptySequenceHost(id, typeID)
// 	if h.ID() != id {
// 		t.Errorf(UnequalIntParameterError, "host ID", id, h.ID())
// 	}
// 	if h.TypeID() != typeID {
// 		t.Errorf(UnequalIntParameterError, "host type ID", typeID, h.TypeID())
// 	}
// 	if len(h.Pathogens()) != 0 {
// 		t.Errorf(UnequalIntParameterError, "number of pathogens", 0, len(h.Pathogens()))
// 	}
// 	if h.PathogenPopSize() != 0 {
// 		t.Errorf(UnequalIntParameterError, "number of pathogens", 0, h.PathogenPopSize())
// 	}
// 	if p := h.Pathogen(0); p != nil {
// 		t.Errorf(UnequalIntParameterError, "number of pathogens", 0, 1)
// 	}
// }

// func TestSequenceHost_Infected_Getters(t *testing.T) {
// 	id := 1
// 	typeID := 0
// 	numPathogens := 10
// 	sites := 100
// 	h := sampleInfectedHost(id, numPathogens, sites)
// 	if h.ID() != id {
// 		t.Errorf(UnequalIntParameterError, "host ID", id, h.ID())
// 	}
// 	if h.TypeID() != typeID {
// 		t.Errorf(UnequalIntParameterError, "host type ID", typeID, h.TypeID())
// 	}
// 	if len(h.Pathogens()) != numPathogens {
// 		t.Errorf(UnequalIntParameterError, "number of pathogens", 0, len(h.Pathogens()))
// 	}
// 	if h.PathogenPopSize() != numPathogens {
// 		t.Errorf(UnequalIntParameterError, "number of pathogens", 0, h.PathogenPopSize())
// 	}
// 	if p := h.Pathogen(0); p == nil {
// 		t.Errorf(UnequalIntParameterError, "number of pathogens", 1, 0)
// 	}
// }

// func TestSequenceHost_Setters(t *testing.T) {
// 	id := 1
// 	typeID := 0
// 	numPathogens := 10
// 	sites := 100
// 	h := sampleInfectedHost(id, numPathogens, sites)
// 	if h.ID() != id {
// 		t.Errorf(UnequalIntParameterError, "host ID", id, h.ID())
// 	}
// 	if h.TypeID() != typeID {
// 		t.Errorf(UnequalIntParameterError, "host type ID", typeID, h.TypeID())
// 	}
// 	if len(h.Pathogens()) != numPathogens {
// 		t.Errorf(UnequalIntParameterError, "number of pathogens", numPathogens, len(h.Pathogens()))
// 	}
// 	if h.PathogenPopSize() != numPathogens {
// 		t.Errorf(UnequalIntParameterError, "number of pathogens", numPathogens, h.PathogenPopSize())
// 	}

// 	// Add pathogen
// 	h.AddPathogens(h.Pathogen(0))
// 	if h.PathogenPopSize() != numPathogens+1 {
// 		t.Errorf(UnequalIntParameterError, "number of pathogens", numPathogens+1, h.PathogenPopSize())
// 	}

// 	// Remove pathogen
// 	h.RemovePathogens(10)
// 	if h.PathogenPopSize() != numPathogens {
// 		t.Errorf(UnequalIntParameterError, "number of pathogens", numPathogens, h.PathogenPopSize())
// 	}

// 	// Remove all pathogens
// 	h.RemoveAllPathogens()
// 	if h.PathogenPopSize() != 0 {
// 		t.Errorf(UnequalIntParameterError, "number of pathogens", 0, h.PathogenPopSize())
// 	}
// }

// func TestSequenceHost_SetIntrahostModel(t *testing.T) {
// 	id := 1
// 	typeID := 1
// 	h := EmptySequenceHost(id, typeID)
// 	seqH := h.(*sequenceHost)
// 	model := sampleIntrahostModel(10e-5, 1000)

// 	if seqH.IntrahostModel != nil {
// 		t.Errorf("intrahost "+ModelExistsError, seqH.IntrahostModel.ModelName(), seqH.IntrahostModel.ModelID())
// 	}
// 	// Add model
// 	h.SetIntrahostModel(model)
// 	if seqH.IntrahostModel == nil {
// 		t.Errorf("intrahost " + EmptyModelError)
// 	}
// }
