package contagiongo

import "testing"

func TestNewEmptySequenceHost(t *testing.T) {
	id := 1
	h := NewEmptySequenceHost(id)
	if h.ID() != id {
		t.Errorf(UnequalIntParameterError, "host ID", id, h.ID())
	}
	if h.TypeID() != 0 {
		t.Errorf(UnequalIntParameterError, "host type ID", 0, h.TypeID())
	}
	if len(h.Pathogens()) > 0 {
		t.Errorf(UnequalIntParameterError, "number of pathogens", 0, len(h.Pathogens()))
	}

	id = 1
	typeID := 1
	h = NewEmptySequenceHost(id, typeID)
	if h.ID() != id {
		t.Errorf(UnequalIntParameterError, "host ID", id, h.ID())
	}
	if h.TypeID() != typeID {
		t.Errorf(UnequalIntParameterError, "host type ID", typeID, h.TypeID())
	}
	if len(h.Pathogens()) > 0 {
		t.Errorf(UnequalIntParameterError, "number of pathogens", 0, len(h.Pathogens()))
	}
}

func TestSequenceHost_SetIntrahostModel(t *testing.T) {
	id := 1
	typeID := 1
	h := NewEmptySequenceHost(id, typeID)
	seqH := h.(*sequenceHost)
	model := sampleIntrahostModel()

	if seqH.IntrahostModel != nil {
		t.Errorf(IntrahostModelExistsError, seqH.IntrahostModel.ModelName(), seqH.IntrahostModel.ModelID())
	}
	// Add model
	h.SetIntrahostModel(model)
	if seqH.IntrahostModel == nil {
		t.Errorf(EmptyIntrahostModelError)
	}
}
