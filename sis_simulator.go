package contagiongo

import "sync"

// SISSimulator creates and runs an SIS epidemiological simulation. Within this
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
