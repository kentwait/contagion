package contagiongo

import (
	"math"
	"math/rand"
	"sync"

	rv "github.com/kentwait/randomvariate"
)

// TransmissionEvent is a struct for sending and receiving
// transmission event information.
type TransmissionEvent struct {
	destination Host
	pathogen    GenotypeNode
}

// ExchangeEvent is a struct for sending and receiving
// exchange event information.
type ExchangeEvent struct {
	source      Host
	destination Host
	// pathogenIndex int
	pathogen GenotypeNode
}

// TransmitPathogens transmits the pathogen to its neighboring host/s.
// If transmission occurs, sends transmitted node over the channel to
// be added to the recepient. Also sends node information in order to
// record the event.
func TransmitPathogens(i, t int, src, dst Host, numMigrants int, transmissionProb float64, count int, c chan<- TransmissionEvent, d chan<- TransmissionPackage, wg *sync.WaitGroup) {
	defer wg.Done()
	// Check if migration size if larger than the current population size
	// If larger, make migrants equal to existing size
	if count < 1 {
		return
	}
	// numMigrants := src.GetTransmissionModel().TransmissionSize()
	if numMigrants > count {
		numMigrants = count
	}
	// Determine if tranmission occurs or not based on source's
	// transmission probability
	// transmissionProb := src.GetTransmissionModel().TransmissionProb()
	if rv.Binomial(1, transmissionProb) == 1.0 {
		// If transmission occurs, randomly pick pathogens to transmit
		for _, p := range src.PickPathogens(numMigrants) {
			if p != nil {
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

// ExchangePathogens exchanges pathogens between neighboring hosts.
func ExchangePathogens(i, t int, h1, h2 Host, h1Count, h2Count int, c chan<- ExchangeEvent, d chan<- TransmissionPackage, wg *sync.WaitGroup) {
	defer wg.Done()
	// Check if migration size if larger than the current population size
	// If larger, skip
	if h1.GetTransmissionModel().TransmissionSize() > h1Count {
		return
	} else if h2.GetTransmissionModel().TransmissionSize() > h2Count {
		return
	}
	// Assumes transmission size if equal in all hosts
	numMigrants := h1.GetTransmissionModel().TransmissionSize()

	// Determine if exchange occurs or not based on square of the source's
	// transmission probability. This assumes that transmission prob is equal
	// between any two hosts.
	transmissionProb := math.Pow(h1.GetTransmissionModel().TransmissionProb(), 2)
	if rv.Binomial(1, transmissionProb) == 1.0 {
		// If exchange occurs, randomly pick pathogens in the h1 and h2 hosts
		// h1 -> h2
		for _, p := range h1.PickPathogens(numMigrants) {
			if p != nil {
				c <- ExchangeEvent{
					source:      h1,
					destination: h2,
					pathogen:    p,
				}
				d <- TransmissionPackage{
					instanceID: i,
					genID:      t,
					fromHostID: h1.ID(),
					toHostID:   h2.ID(),
					nodeID:     p.UID(),
				}
			}
		}
		// h2 -> h1
		for _, p := range h2.PickPathogens(numMigrants) {
			if p != nil {
				c <- ExchangeEvent{
					source:      h2,
					destination: h1,
					pathogen:    p,
				}
				d <- TransmissionPackage{
					instanceID: i,
					genID:      t,
					fromHostID: h2.ID(),
					toHostID:   h1.ID(),
					nodeID:     p.UID(),
				}
			}
		}
	}
}

func pickPathogens(count, numMigrants int) []int {
	return rand.Perm(count)[:numMigrants]
}
