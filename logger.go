package contagiongo

import (
	"github.com/segmentio/ksuid"
)

// DataLogger is the general definition of a logger that records
// simulation data to file whether it writes a text file or
// writes to a database.
type DataLogger interface {
	SetBasePath(path string, i int)
	// Init initializes the logger. For example, if the logger writes a
	// CSV file, Init can create a file and write header information first.
	// Or if the logger writes to a database, Init can be used to
	// create a new table.
	Init() error
	// WriteGenotypes records a new genotype's ID and sequence to file.
	WriteGenotypes(c <-chan Genotype)
	// WriteGenotypeNodes records new genotype node's ID and
	// associated genotype ID to file
	WriteGenotypeNodes(c <-chan GenotypeNode)
	// WriteGenotypeFreq records the count of unique genotype nodes
	// present within the host in a given time in the simulation.
	WriteGenotypeFreq(c <-chan GenotypeFreqPackage)
	// WriteMutations records every time a new genotype node is created.
	// It records the time and in what host this new mutation arose.
	WriteMutations(c <-chan MutationPackage)
	// WriteStatus records the status of each host every generation.
	WriteStatus(c <-chan StatusPackage)
	// WriteTransmission records the ID's of genotype node that
	// are transmitted between hosts.
	WriteTransmission(c <-chan TransmissionPackage)
}

// GenotypeFreqPackage encapsulates the data to be written everytime
// the frequency of genotypes have to be recorded.
type GenotypeFreqPackage struct {
	instanceID int
	genID      int
	hostID     int
	genotypeID ksuid.KSUID
	freq       int
}

// StatusPackage encapsulates the data to be written everytime
// the status of a host has to be recorded.
type StatusPackage struct {
	instanceID int
	genID      int
	hostID     int
	status     int
}

// MutationPackage encapsulates information to be written
// to track when and where mutations occur in the simulation.
type MutationPackage struct {
	instanceID   int
	genID        int
	hostID       int
	nodeID       ksuid.KSUID
	parentNodeID ksuid.KSUID
}

// TransmissionPackage encapsulates information to be written
// to track the movement of genotype nodes across the host
// population.
type TransmissionPackage struct {
	instanceID int
	genID      int
	fromHostID int
	toHostID   int
	nodeID     ksuid.KSUID
}
