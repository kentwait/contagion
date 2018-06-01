package contagiongo

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

// CSVLogger is a DataLogger that writes simulation data
// as comma-delimited files.
type CSVLogger struct {
	genotypePath     string
	genotypeNodePath string
	genotypeFreqPath string
	statusPath       string
	transmissionPath string
	mutationPath     string
}

func NewCSVLogger(basepath string, i int) *CSVLogger {
	l := new(CSVLogger)
	l.SetBasePath(basepath, i)
	return l
}

func (l *CSVLogger) SetBasePath(basepath string, i int) {
	if info, err := os.Stat(basepath); err == nil && info.IsDir() {
		basepath += fmt.Sprintf("log.%03d", i)
	}
	l.genotypePath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%03d.%s.csv", i, "g")
	l.genotypeNodePath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%03d.%s.csv", i, "n")
	l.genotypeFreqPath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%03d.%s.csv", i, "freq")
	l.statusPath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%03d.%s.csv", i, "status")
	l.transmissionPath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%03d.%s.csv", i, "trans")
	l.mutationPath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%03d.%s.csv", i, "tree")
}

func (l *CSVLogger) WriteGenotypes(c <-chan Genotype) {
	// Format
	// <genotypeID>  <sequence>
	const template = "%s,%s\n"
	var b bytes.Buffer
	// b.WriteString("genotypeID,sequence\n")
	for genotype := range c {
		row := fmt.Sprintf(template,
			genotype.GenotypeUID(),
			genotype.StringSequence(),
		)
		// TODO: log error
		b.WriteString(row)
	}
	AppendToFile(l.genotypePath, b.Bytes())
}

func (l *CSVLogger) WriteGenotypeNodes(c <-chan GenotypeNode) {
	// Format
	// <nodeID>  <genotypeID>
	const template = "%s,%s\n"
	var b bytes.Buffer
	// b.WriteString("nodeID,genotypeID\n")
	for node := range c {
		row := fmt.Sprintf(template,
			node.UID(),
			node.GenotypeUID(),
		)
		// TODO: log error
		b.WriteString(row)
	}
	AppendToFile(l.genotypeNodePath, b.Bytes())
}

func (l *CSVLogger) WriteGenotypeFreq(c <-chan GenotypeFreqPackage) {
	// Format
	// <instanceID>  <generation>  <hostID>  <genotypeID>  <freq>
	const template = "%d,%d,%d,%s,%d\n"
	var b bytes.Buffer
	// b.WriteString("instance,generation,hostID,genotypeID,freq\n")
	for pack := range c {
		row := fmt.Sprintf(template,
			pack.instanceID,
			pack.genID,
			pack.hostID,
			pack.genotypeID.String(),
			pack.freq,
		)
		// TODO: log error
		b.WriteString(row)
	}
	AppendToFile(l.genotypeFreqPath, b.Bytes())
}

func (l *CSVLogger) WriteMutations(c <-chan MutationPackage) {
	// Format
	// <instanceID>  <generation>  <hostID>  <parentGenotypeID>  <currentGenotypeID>
	const template = "%d,%d,%d,%s,%s\n"
	var b bytes.Buffer
	// b.WriteString("instance,generation,hostID,parentGenotypeID,currentGenotypeID\n")
	for pack := range c {
		row := fmt.Sprintf(template,
			pack.instanceID,
			pack.genID,
			pack.hostID,
			pack.parentNodeID.String(),
			pack.nodeID.String(),
		)
		// TODO: log error
		b.WriteString(row)
	}
	AppendToFile(l.mutationPath, b.Bytes())
}

func (l *CSVLogger) WriteStatus(c <-chan StatusPackage) {
	// Format
	// <instanceID>  <generation>  <hostID>  <status>
	const template = "%d,%d,%d,%d\n"
	var b bytes.Buffer
	// b.WriteString("instance,generation,hostID,status\n")
	for pack := range c {
		row := fmt.Sprintf(template,
			pack.instanceID,
			pack.genID,
			pack.hostID,
			pack.status,
		)
		// TODO: log error
		b.WriteString(row)
	}
	AppendToFile(l.statusPath, b.Bytes())
}

func (l *CSVLogger) WriteTransmission(c <-chan TransmissionPackage) {
	// Format
	// <instanceID>  <generation>  <fromHostID>  <toHostID> <genotypeID>
	const template = "%d,%d,%d,%d,%s\n"
	var b bytes.Buffer
	// b.WriteString("instance,generation,fromHostID,toHostID,genotypeID\n")
	for pack := range c {
		row := fmt.Sprintf(template,
			pack.instanceID,
			pack.genID,
			pack.fromHostID,
			pack.toHostID,
			pack.nodeID.String(),
		)
		// TODO: log error
		b.WriteString(row)
	}
	AppendToFile(l.transmissionPath, b.Bytes())
}

// AppendToFile creates a new file on the given path if it does not exist, or
// appends to the end of the existing file if the file exists.
func AppendToFile(path string, b []byte) error {
	// Create file
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}
	return nil
}
