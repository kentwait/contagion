package contagiongo

import (
	"sync"

	"github.com/segmentio/ksuid"
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

// Epidemic encapsulates the set of hosts, its connections,
// the pathogen tree lineage and the host types used to create
// a simulated epidemic.
type Epidemic interface {
	// Host returns the selected host in the simulation.
	Host(id int) Host

	// HostStatus retrieves the current status of the selected host.
	HostStatus(id int) int
	// SetHostStatus sets the current status of the selected host
	// to a given status code.
	SetHostStatus(id, status int)
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

	// NewInstance creates a NewInstancete a new instance from
	// the stored configuration
	NewInstance() (Epidemic, error)

	GenotypeNodeMap() map[ksuid.KSUID]GenotypeNode
	GenotypeSet() GenotypeSet

	// The following methods perform intrahost processes associated with
	// the status. For every generation, one of the following is called for
	// each host.

	// SusceptibleProcess performs intrahost processes while the host is in
	// the susceptible status.
	SusceptibleProcess(i, t int, host Host, wg *sync.WaitGroup)
	// ExposedProcess performs intrahost processes while the host is in
	// the exposed status.
	ExposedProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup)
	// InfectedProcess performs intrahost processes while the host is in
	// the infected status.
	InfectedProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup)
	// InfectiveProcess performs intrahost processes while the host is in
	// the infective status.
	InfectiveProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup)
	// RemovedProcess performs intrahost processes while the host is in
	// the removed status.
	RemovedProcess(i, t int, host Host, wg *sync.WaitGroup)
	// RecoveredProcess performs intrahost processes while the host is in
	// the recovered status.
	RecoveredProcess(i, t int, host Host, wg *sync.WaitGroup)
	// DeadProcess performs intrahost processes while the host is in
	// the dead status.
	DeadProcess(i, t int, host Host, wg *sync.WaitGroup)
	// DeadProcess performs intrahost processes while the host is in
	// the dead status.
	VaccinatedProcess(i, t int, host Host, wg *sync.WaitGroup)
}

// EpidemicSimulation is a simulation environment that simulates
// the spread of disease between hosts in a connected host network.
type EpidemicSimulation interface {
	Epidemic
	// Run runs the whole simulation
	Init(params ...interface{})
	Run(i int)
	Update(t int)
	Process(t int)
	Transmit(t int)
}

// Infection encapsulates a single host and the pathogen tree lineage
// to trace the evolution of pathogens within one host. This is useful
// to strudy intrahost evolutionary dynamics especially in chronic diseases.
type Infection interface {
	// Host returns the selected host in the simulation.
	Host() Host

	// HostStatus retrieves the current status of the selected host.
	HostStatus() int
	// SetHostStatus sets the current status of the selected host
	// to a given status code.
	SetHostStatus(status int)
	// HostStatusDuration returns the number of generations a host
	// remains in a given status.
	HostStatusDuration(status int) int
	// HostTime returns the current number of generations remaining
	// before the host changes status.
	HostTimer() int
	// SetHostTimer sets the number of generations for the host to
	// remain in its current status.
	SetHostTimer(interval int)

	// The following methods perform intrahost processes associated with
	// the status. For every generation, one of the following is called for
	// each host.

	// SusceptibleProcess performs intrahost processes while the host is in
	// the susceptible status.
	SusceptibleProcess(i, t int, host Host, wg *sync.WaitGroup)
	// ExposedProcess performs intrahost processes while the host is in
	// the exposed status.
	ExposedProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup)
	// InfectedProcess performs intrahost processes while the host is in
	// the infected status.
	InfectedProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup)
	// InfectiveProcess performs intrahost processes while the host is in
	// the infective status.
	InfectiveProcess(i, t int, host Host, c chan<- MutationPackage, wg *sync.WaitGroup)
	// RemovedProcess performs intrahost processes while the host is in
	// the removed status.
	RemovedProcess(i, t int, host Host, wg *sync.WaitGroup)
	// RecoveredProcess performs intrahost processes while the host is in
	// the recovered status.
	RecoveredProcess(i, t int, host Host, wg *sync.WaitGroup)
	// DeadProcess performs intrahost processes while the host is in
	// the dead status.
	DeadProcess(i, t int, host Host, wg *sync.WaitGroup)
	// DeadProcess performs intrahost processes while the host is in
	// the dead status.
	VaccinatedProcess(i, t int, host Host, wg *sync.WaitGroup)
}

// InfectionSimulation is a simulation environment that simulates
// the infection within a single hosts or in an evironment where a network
// configuration is not necessary.
type InfectionSimulation interface {
	Infection
	// Run runs the whole simulation
	Init(params ...interface{})
	Run(i int)
	Update(t int)
	Process(t int)
	Transmit(t int)
}

// TransmissionEvent is a struct for sending and receiving
// transmission event information.
type TransmissionEvent struct {
	destination Host
	pathogen    GenotypeNode
}
