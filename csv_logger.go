package contagiongo

import (
	"bytes"
	"fmt"
	"sync"
)

// CSVLogger is a DataLogger that writes simulation data
// as comma-delimited files.
type CSVLogger struct {
	genotypePath     string
	genotypeNodePath string
	genotypeFreqPath string
	statusPath       string
	mutationPath     string
}

func (l *CSVLogger) WriteGenotypes(c <-chan Genotype, wg *sync.WaitGroup) {
	defer wg.Done()
	// Format
	// <genotypeID>  <sequence>
	const genotypeRowTemplate = "%s,%s"

	var b bytes.Buffer
	for genotype := range c {
		row := fmt.Sprintf(genotypeRowTemplate+"\n", genotype.GenotypeUID(), genotype.StringSequence())
		// TODO: log error
		b.WriteString(row)
	}
}

func (l *CSVLogger) WriteGenotypeNodes(c <-chan GenotypeNode, wg *sync.WaitGroup) {
	// Format
	// <nodeID>  <genotypeID>
	const nodeRowTemplate = "%s,%s"

}

func (l *CSVLogger) WriteGenotypeFreq(c <-chan GenotypeFreqPackage, wg *sync.WaitGroup) {
	// Format
	// <instanceID>  <generation>  <hostID>  <genotypeID>  <count>
	const genotypeFreqRowTemplate = "%d,%d,%d,%s,%d"

}

func (l *CSVLogger) WriteMutations(c <-chan MutationPackage, wg *sync.WaitGroup) {
	// Format
	// <instanceID>  <generation>  <hostID>  <parentGenotypeID>  <currentGenotypeID>
	const mutationRowTemplate = "%d,%d,%d,%s,%s"

}

func (l *CSVLogger) WriteStatus(c <-chan StatusPackage, wg *sync.WaitGroup) {
	// Format
	// <instanceID>  <generation>  <hostID>  <status>
	const statusRowTemplate = "%d,%d,%d,%d"

}

func (l *CSVLogger) WriteTransmission(c <-chan TransmissionPackage, wg *sync.WaitGroup) {
	// Format
	// <instanceID>  <generation>  <fromHostID>  <toHostID> <genotypeID>
	const transmissionRowTemplate = "%d,%d,%d,%d,%s"

}
