package contagiongo

import (
	"sync"
)

// The following are status codes for different preset compartments that
// describe the current epidemiological status of a host in the simulation.
const (
	SusceptibleStatusCode = 1
	ExposedStatusCode     = 2
	InfectedStatusCode    = 3
	InfectiveStatusCode   = 4
	RemovedStatusCode     = 5
	RecoveredStatusCode   = 6
	DeadStatusCode        = 7
	VaccinatedStatusCode  = 8
)

// Simulation encapsulates the set of hosts, its connections,
// the pathogen tree lineage and the host types used to create
// a simulated epidemic.
type Simulation interface {
	// Host returns a single host based on the given host ID.
	Host(id int) EpidemicHost
	// HostMap returns the map of all hosts in the simulation.
	HostMap() map[int]EpidemicHost
	// HostStatus returns the status of a particular host that matches
	// the given host ID.
	HostStatus(id int) int
	// SetHostStatus assigns an encoded integer host status to a particular host.
	SetHostStatus(id, status int)
	// HostTimer
	HostTimer(id int) int
	// SetHostTimer
	SetHostTimer(id, interval int)
	// HostNeighbors returns the list of hosts that are directly connected to
	// the host with the given host ID.
	HostNeighbors(id int) []EpidemicHost

	// SusceptibleProcess runs the processes associated with the "susceptible"
	// host status.
	SusceptibleProcess(host EpidemicHost, wg *sync.WaitGroup)
	// ExposedProcess runs the processes associated with the "exposed"
	// host status.
	ExposedProcess(host EpidemicHost, wg *sync.WaitGroup)
	// InfectedProcess runs the processes associated with the "infected"
	// host status.
	InfectedProcess(host EpidemicHost, wg *sync.WaitGroup)
	// InfectiveProcess runs the processes associated with the "infective"
	// host status.
	InfectiveProcess(host EpidemicHost, wg *sync.WaitGroup)
	// RemovedProcess runs the processes associated with the "removed"
	// host status.
	RemovedProcess(host EpidemicHost, wg *sync.WaitGroup)
	// RecoveredProcess runs the processes associated with the "recovered"
	// host status.
	RecoveredProcess(host EpidemicHost, wg *sync.WaitGroup)
	// DeadProcess runs the processes associated with the "dead"
	// host status.
	DeadProcess(host EpidemicHost, wg *sync.WaitGroup)
	// VaccinatedProcess runs the processes associated with the "vaccinated"
	// host status.
	VaccinatedProcess(host EpidemicHost, wg *sync.WaitGroup)
}

// Simulator is a Simulation with a Run method to perform the simulation and
// record the results.
type Simulator interface {
	Simulation
	// Run
	Run(DataRecorder interface{})
}
