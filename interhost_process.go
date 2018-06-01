package contagiongo

import (
	"math/rand"
	"sync"

	rv "github.com/kentwait/randomvariate"
)

// TransmitPathogens transmits the pathogen to its neighboring host/s.
// If transmission occurs, sends transmitted node over the channel to
// be added to the recepient. Also sends node information in order to
// record the event.
func TransmitPathogens(i, t int, src, dst Host, count, status int, c chan<- TransmissionEvent, d chan<- TransmissionPackage, wg *sync.WaitGroup) {
	defer wg.Done()
	// Check if migration size if larger than the current population size
	// If larger, skip
	numMigrants := src.GetTransmissionModel().TransmissionSize()
	if numMigrants > count {
		return
	}
	// Determine if tranmission occurs or not based on source's
	// transmission probability
	transmissionProb := src.GetTransmissionModel().TransmissionProb()
	if rv.Binomial(1, transmissionProb) == 1.0 {
		// If transmission occurs, randomly pick pathogens to transmit
		ids := pickPathogens(count, numMigrants)
		for _, id := range ids {
			if p := src.Pathogen(id); p != nil {
				c <- TransmissionEvent{dst, p}
				d <- TransmissionPackage{
					instanceID: i,
					genID:      t,
					fromHostID: src.ID(),
					toHostID:   dst.ID(),
					nodeID:     p.UID(),
				}
			}
		}
	}
}

func pickPathogens(count, numMigrants int) []int {
	return rand.Perm(count)[:numMigrants]
}
