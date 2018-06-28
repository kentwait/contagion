package contagiongo

import (
	"sync"
)

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
