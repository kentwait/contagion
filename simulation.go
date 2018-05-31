package contagiongo

import "sync"

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
	// Host returns the selected host in the simulation.
	Host(id int) Host

	// HostStatus retrieves the current status of the selected host.
	HostStatus(id int) int
	// SetHostStatus sets the current status of the selected host
	// to a given status code.
	SetHostStatus(id, status int)
	// HostStatusDuration returns the number of generations a host
	// remains in a given status.
	HostStatusDuration(id, status int) int
	// HostTime returns the current number of generations remaining
	// before the host changes status.
	HostTimer(id int) int
	// SetHostTimer sets the number of generations for the host to
	// remain in its current status.
	SetHostTimer(id, interval int)
	// InfectableStatuses returns the list of statuses that infected
	// hosts can transmit to.
	InfectableStatuses() []int

	// HostMap returns the hosts in the simulation in the form of a map.
	// The key is the host's ID and the value is the pointer to the host.
	HostMap() map[int]Host
	// HostNeighbors retrieves the directly connected hosts to the current
	// host based on the supplied adjacency matrix.
	HostNeighbors(id int) []Host

	// The following methods perform intrahost processes associated with
	// the status. For every generation, one of the following is called for
	// each host.

	// SusceptibleProcess performs intrahost processes while the host is in
	// the susceptible status.
	SusceptibleProcess(host Host, wg *sync.WaitGroup)
	// ExposedProcess performs intrahost processes while the host is in
	// the exposed status.
	ExposedProcess(host Host, wg *sync.WaitGroup)
	// InfectedProcess performs intrahost processes while the host is in
	// the infected status.
	InfectedProcess(host Host, wg *sync.WaitGroup)
	// InfectiveProcess performs intrahost processes while the host is in
	// the infective status.
	InfectiveProcess(host Host, wg *sync.WaitGroup)
	// RemovedProcess performs intrahost processes while the host is in
	// the removed status.
	RemovedProcess(host Host, wg *sync.WaitGroup)
	// RecoveredProcess performs intrahost processes while the host is in
	// the recovered status.
	RecoveredProcess(host Host, wg *sync.WaitGroup)
	// DeadProcess performs intrahost processes while the host is in
	// the dead status.
	DeadProcess(host Host, wg *sync.WaitGroup)
	// DeadProcess performs intrahost processes while the host is in
	// the dead status.
	VaccinatedProcess(host Host, wg *sync.WaitGroup)
	// Run runs the whole simulation
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
