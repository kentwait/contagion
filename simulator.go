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

// StatusUpdate is a struct for sending and receiving
// host status updates.
type StatusUpdate struct {
	hostID int
	status int
	timer  int
}

// InfectiveParams is a struct for sending and receiving
// infective host information.
type InfectiveParams struct {
	source  Host
	popSize int
}

// Unpack returns the property values of InfectiveParams.
func (p *InfectiveParams) Unpack() (Host, int) {
	return p.source, p.popSize
}

// TransParams is a struct for sending and receiving
// transmission event information.
type TransParams struct {
	destination Host
	pathogen    interface{}
}

// Unpack returns the property values of TransParams.
func (p *TransParams) Unpack() (Host, interface{}) {
	return p.destination, p.pathogen
}

// SISSimulator create and runs an SIS epidemiological simulation. Within this
// simulation, hosts may or may not run independent genetic evolution
// simulations.
type SISSimulator struct {
	Simulation
	numGenerations      int
	infectableStatuses  []int
	pathogenLoggingFreq int
	statusLoggingFreq   int
}

// Run instantiates, runs, and records the a new simulation.
func (s *SISSimulator) Run(DataRecorder interface{}) {
	s.update()
	s.record(DataRecorder)
	for t := 0; t < s.numGenerations; t++ {
		s.process()
		s.transmit(DataRecorder)
		s.update()
		s.record(DataRecorder)
	}
}

func (s *SISSimulator) statusFunc(host EpidemicHost, status, timer int, c chan<- statusUpdate, wg *sync.WaitGroup) {
	// Status is updated only if the internal host time is less than 1.
	hostID := host.HostID()
	// Add cases depending on the compartmental model being used
	// In this case, SIS only uses susceptible and infected statuses
	if timer < 1 {
		switch status {
		case SusceptibleStatusCode:
			c <- statusUpdate{hostID, InfectedStatusCode, host.TimeInterval(InfectedStatusCode)}
		case InfectedStatusCode:
			host.ClearPathogens()
			c <- statusUpdate{hostID, SusceptibleStatusCode, host.TimeInterval(SusceptibleStatusCode)}
		}
	}
	wg.Done()
}

func (s *SISSimulator) update() {
	// Update status first
	updates := make(chan statusUpdate)
	hosts := s.HostMap()
	var wg sync.WaitGroup
	// Read all hosts and process hosts concurrently
	// These succeeding steps connects the simulation's record of each host's
	// status and timer with the host's internal state
	wg.Add(len(hosts))
	for hostID, host := range hosts {
		status := s.HostStatus(hostID)
		timer := s.HostTimer(hostID)
		go s.statusFunc(host, status, timer, updates, &wg)
	}
	go func() {
		wg.Wait()
		close(updates)
	}()
	// If host status changed, update the simulation's record of the particular
	// host's status
	for u := range updates {
		s.SetHostStatus(u.hostID, u.status)
		s.SetHostTimer(u.hostID, u.timer)
	}
}

func (s *SISSimulator) record(DataRecorder interface{}) {
	// TODO: Record data
}

func (s *SISSimulator) process() {
	hosts := s.HostMap()
	var wg sync.WaitGroup
	// Read all hosts and process based on the current status of the host
	wg.Add(len(hosts))
	for hostID, host := range hosts {
		switch s.HostStatus(hostID) {
		case SusceptibleStatusCode:
			go s.SusceptibleProcess(host, &wg)
		case InfectedStatusCode:
			go s.InfectedProcess(host, &wg)
			timer := s.HostTimer(hostID)
			// Update simulator's record of infection times
			s.SetHostTimer(hostID, timer-1)
		}
	}
	wg.Wait()
}

func (s *SISSimulator) transmit(DataRecorder interface{}) {
	// Process infective hosts and transmit
	transmissions := make(chan TransParams)
	var wg sync.WaitGroup
	for hostID, src := range s.HostMap() {
		if s.HostStatus(hostID) == InfectedStatusCode {
			popSize := src.PathogenPopSize()
			for _, neighbor := range s.HostNeighbors(hostID) {
				// Spawn a new goroutine for every neighbor of the host
				neighborStatus := s.HostStatus(neighbor.HostID())
				wg.Add(1)
				go PathogenTransmitter(src, neighbor, popSize, neighborStatus, transmissions, &wg, s.infectableStatuses...)
			}
		}
	}
	go func() {
		wg.Wait()
		close(transmissions)
	}()
	// Add the new pathogen to the destination host
	for params := range transmissions {
		dst, pathogen := params.Unpack()
		dst.AddPathogen(pathogen)
	}
}
