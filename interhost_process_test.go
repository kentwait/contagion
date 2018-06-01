package contagiongo

import (
	"sync"
	"testing"
)

func TestTransmitPathogens(t *testing.T) {
	count := 10
	sites := 100
	h1 := sampleInfectedHost(1, count, sites)
	model := new(constantTransmitter)
	model.prob = 1.0
	model.size = 1
	h1.SetTransmissionModel(model)
	h2 := EmptySequenceHost(2)
	eventC := make(chan TransmissionEvent)
	packC := make(chan TransmissionPackage)
	i := 0
	time := 0

	var wg sync.WaitGroup
	wg.Add(1)
	go TransmitPathogens(i, time, h1, h2, count, eventC, packC, &wg)
	go func() {
		wg.Wait()
		close(eventC)
		close(packC)
	}()
	var wg2 sync.WaitGroup
	wg2.Add(2)
	go func() {
		defer wg2.Done()
		cnt := 0
		for evt := range eventC {
			if evt.destination.ID() != 2 {
				t.Errorf(UnequalIntParameterError, "destination host ID", 2, evt.destination.ID())
			}
			cnt++
		}
		if cnt != model.size {
			t.Errorf(UnequalIntParameterError, "number of transmitted pathogens", model.size, cnt)
		}
	}()
	go func() {
		defer wg2.Done()
		cnt := 0
		for pack := range packC {
			if pack.fromHostID != 1 {
				t.Errorf(UnequalIntParameterError, "source host ID", 1, pack.fromHostID)
			}
			if pack.toHostID != 2 {
				t.Errorf(UnequalIntParameterError, "destination host ID", 2, pack.toHostID)
			}
			cnt++
		}
		if cnt != model.size {
			t.Errorf(UnequalIntParameterError, "number of transmitted pathogens", model.size, cnt)
		}
	}()
	wg2.Wait()
}
