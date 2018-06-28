package contagiongo

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/segmentio/ksuid"
	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// DataLogger is the general definition of a logger that records
// simulation data to file whether it writes a text file or
// writes to a database.
type DataLogger interface {
	// SetBasePath sets the base path of the logger.
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

// NewCSVLogger creates a new logger that writes data into CSV files.
func NewCSVLogger(basepath string, i int) *CSVLogger {
	l := new(CSVLogger)
	l.SetBasePath(basepath, i)
	return l
}

// SetBasePath sets the base path of the logger.
func (l *CSVLogger) SetBasePath(basepath string, i int) {
	if info, err := os.Stat(basepath); err == nil && info.IsDir() {
		basepath += fmt.Sprintf("log")
	}
	l.genotypePath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%03d.%s.csv", i, "g")
	l.genotypeNodePath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%03d.%s.csv", i, "n")
	l.genotypeFreqPath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%03d.%s.csv", i, "freq")
	l.statusPath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%03d.%s.csv", i, "status")
	l.transmissionPath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%03d.%s.csv", i, "trans")
	l.mutationPath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%03d.%s.csv", i, "tree")
}

// Init creates CSV files and writes header information for each file.
func (l *CSVLogger) Init() error {
	newFile := func(path, header string) error {
		var b bytes.Buffer
		_, err := b.WriteString(header)
		if err != nil {
			return err
		}
		err = NewFile(path, b.Bytes())
		if err != nil {
			return err
		}
		return nil
	}

	err := newFile(l.genotypePath, "genotypeID,sequence\n")
	if err != nil {
		return err
	}
	err = newFile(l.genotypeNodePath, "nodeID,genotypeID\n")
	if err != nil {
		return err
	}
	err = newFile(l.genotypeFreqPath, "instance,generation,hostID,genotypeID,freq\n")
	if err != nil {
		return err
	}
	err = newFile(l.mutationPath, "instance,generation,hostID,parentNodeID,nodeID\n")
	if err != nil {
		return err
	}
	err = newFile(l.statusPath, "instance,generation,hostID,status\n")
	if err != nil {
		return err
	}
	err = newFile(l.transmissionPath, "instance,generation,fromHostID,toHostID,nodeID\n")
	if err != nil {
		return err
	}
	return nil
}

// WriteGenotypes records a new genotype's ID and sequence to file.
func (l *CSVLogger) WriteGenotypes(c <-chan Genotype) {
	// Format
	// <genotypeID>  <sequence>
	const template = "%s,%s\n"
	var b bytes.Buffer
	// b.WriteString("genotypeID,sequence\n")
	for genotype := range c {
		row := fmt.Sprintf(template,
			genotype.GenotypeUID().String(),
			genotype.StringSequence(),
		)
		// TODO: log error
		b.WriteString(row)
	}
	AppendToFile(l.genotypePath, b.Bytes())
}

// WriteGenotypeNodes records new genotype node's ID and
// associated genotype ID to file
func (l *CSVLogger) WriteGenotypeNodes(c <-chan GenotypeNode) {
	// Format
	// <nodeID>  <genotypeID>
	const template = "%s,%s\n"
	var b bytes.Buffer
	// b.WriteString("nodeID,genotypeID\n")
	for node := range c {
		row := fmt.Sprintf(template,
			node.UID().String(),
			node.GenotypeUID().String(),
		)
		// TODO: log error
		b.WriteString(row)
	}
	AppendToFile(l.genotypeNodePath, b.Bytes())
}

// WriteGenotypeFreq records the count of unique genotype nodes
// present within the host in a given time in the simulation.
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

// WriteMutations records every time a new genotype node is created.
// It records the time and in what host this new mutation arose.
func (l *CSVLogger) WriteMutations(c <-chan MutationPackage) {
	// Format
	// <instanceID>  <generation>  <hostID>  <parentNodeID>  <nodeID>
	const template = "%d,%d,%d,%s,%s\n"
	var b bytes.Buffer
	// b.WriteString("instance,generation,hostID,parentNodeID,nodeID\n")
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

// WriteStatus records the status of each host every generation.
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

// WriteTransmission records the ID's of genotype node that
// are transmitted between hosts.
func (l *CSVLogger) WriteTransmission(c <-chan TransmissionPackage) {
	// Format
	// <instanceID>  <generation>  <fromHostID>  <toHostID> <genotypeNodeID>
	const template = "%d,%d,%d,%d,%s\n"
	var b bytes.Buffer
	// b.WriteString("instance,generation,fromHostID,toHostID,nodeID\n")
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

// NewFile creates a new file on the given path if it does not exist.
// Returns an error if the file exists.
func NewFile(path string, b []byte) error {
	// Create file
	if exists, _ := Exists(path); exists {
		return fmt.Errorf("%s already exists", path)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
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

// SQLiteLogger is a DataLogger that writes simulation data
// t0 SQLite databases. Each writer function writes to an independent SQLite
// database and foreign keys are added to each database at the closing
// phase after the simulation is completed.
type SQLiteLogger struct {
	genotypePath     string
	genotypeNodePath string
	genotypeFreqPath string
	statusPath       string
	transmissionPath string
	mutationPath     string
	instanceID       int
}

// NewSQLiteLogger creates a new logger that writes to a SQLite database.
func NewSQLiteLogger(basepath string, i int) *SQLiteLogger {
	l := new(SQLiteLogger)
	l.SetBasePath(basepath, i)
	return l
}

// SetBasePath sets the base path of the logger.
func (l *SQLiteLogger) SetBasePath(basepath string, i int) {
	if info, err := os.Stat(basepath); err == nil && info.IsDir() {
		basepath += fmt.Sprintf("log.%03d", i)
	}
	l.genotypePath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%s.db", "g")
	l.genotypeNodePath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%s.db", "n")
	l.genotypeFreqPath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%s.db", "freq")
	l.statusPath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%s.db", "status")
	l.transmissionPath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%s.db", "trans")
	l.mutationPath = strings.TrimSuffix(basepath, ".") + fmt.Sprintf(".%s.db", "tree")

	// set instance
	l.instanceID = i
}

// Init creates a new tables in the database.
// For example, each new realization of the simulation creates a new table
// for transmissions, frequencies, statuses, nodes and genotypes.
func (l *SQLiteLogger) Init() error {
	// General function to create a new table
	newTable := func(path, tableName, cols string) error {
		db, err := OpenSQLiteDBOptimized(path)
		if err != nil {
			return err
		}
		defer db.Close()
		// cols example:
		// (id integer not null primary key, genotypeID text, sequence text)
		_sqlStmt := `
	create table %s %s;
	delete from %s;
	`
		fullTableName := fmt.Sprintf("%s%03d", tableName, l.instanceID)
		sqlStmt := fmt.Sprintf(_sqlStmt, fullTableName, cols, fullTableName)
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return fmt.Errorf("%q: %s", err, sqlStmt)
		}
		return nil
	}

	// Create tables
	err := newTable(l.genotypePath, "Genotype", "(id integer not null primary key, genotypeID text, sequence text)")
	if err != nil {
		return err
	}
	err = newTable(l.genotypeNodePath, "Node", "(id integer not null primary key, nodeID text, genotypeID text)")
	if err != nil {
		return err
	}
	err = newTable(l.genotypeFreqPath, "GenotypeFreq", "(id integer not null primary key, generation int, hostID int, genotypeID text, freq int)")
	if err != nil {
		return err
	}
	err = newTable(l.mutationPath, "Tree", "(id integer not null primary key, generation int, hostID int, parentNodeID text, nodeID text)")
	if err != nil {
		return err
	}
	err = newTable(l.statusPath, "Status", "(id integer not null primary key, generation int, hostID int, status int)")
	if err != nil {
		return err
	}
	err = newTable(l.transmissionPath, "Transmission", "(id integer not null primary key, generation int, fromHostID int, toHostID int, nodeID text)")
	if err != nil {
		return err
	}
	return nil
}

// WriteGenotypes records a new genotype's ID and sequence to file.
func (l *SQLiteLogger) WriteGenotypes(c <-chan Genotype) {
	tableName := fmt.Sprintf("Genotype%03d", l.instanceID)
	path := l.genotypePath
	_stmt := "insert into " + tableName + "(genotypeID, sequence) values(?, ?)"
	// Database ops below
	db, err := OpenSQLiteDBOptimized(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return
	}
	stmt, err := tx.Prepare(_stmt)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	for genotype := range c {
		_, err = stmt.Exec(
			genotype.GenotypeUID().String(),
			genotype.StringSequence(),
		)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	// Commit at the end
	tx.Commit()
}

// WriteGenotypeNodes records new genotype node's ID and
// associated genotype ID to file
func (l *SQLiteLogger) WriteGenotypeNodes(c <-chan GenotypeNode) {
	tableName := fmt.Sprintf("Node%03d", l.instanceID)
	path := l.genotypeNodePath
	_stmt := "insert into " + tableName + "(nodeID, genotypeID) values(?, ?)"
	// Database ops below
	db, err := OpenSQLiteDBOptimized(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return
	}
	stmt, err := tx.Prepare(_stmt)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	for node := range c {
		_, err = stmt.Exec(
			node.UID().String(),
			node.GenotypeUID().String(),
		)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	// Commit at the end
	tx.Commit()
}

// WriteGenotypeFreq records the count of unique genotype nodes
// present within the host in a given time in the simulation.
func (l *SQLiteLogger) WriteGenotypeFreq(c <-chan GenotypeFreqPackage) {
	tableName := fmt.Sprintf("GenotypeFreq%03d", l.instanceID)
	path := l.genotypeFreqPath
	_stmt := "insert into " + tableName + "(generation, hostID, genotypeID, freq) values(?, ?, ?, ?)"
	// Database ops below
	db, err := OpenSQLiteDBOptimized(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return
	}
	stmt, err := tx.Prepare(_stmt)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	for pack := range c {
		_, err = stmt.Exec(
			pack.genID,
			pack.hostID,
			pack.genotypeID.String(),
			pack.freq,
		)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	// Commit at the end
	tx.Commit()
}

// WriteMutations records every time a new genotype node is created.
// It records the time and in what host this new mutation arose.
func (l *SQLiteLogger) WriteMutations(c <-chan MutationPackage) {
	tableName := fmt.Sprintf("Tree%03d", l.instanceID)
	path := l.mutationPath
	_stmt := "insert into " + tableName + "(generation, hostID, parentNodeID, nodeID) values(?, ?, ?, ?)"
	// Database ops below
	db, err := OpenSQLiteDBOptimized(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return
	}
	stmt, err := tx.Prepare(_stmt)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	for pack := range c {
		_, err = stmt.Exec(
			pack.genID,
			pack.hostID,
			pack.parentNodeID.String(),
			pack.nodeID.String(),
		)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	// Commit at the end
	tx.Commit()
}

// WriteStatus records the status of each host every generation.
func (l *SQLiteLogger) WriteStatus(c <-chan StatusPackage) {
	tableName := fmt.Sprintf("Status%03d", l.instanceID)
	path := l.statusPath
	_stmt := "insert into " + tableName + "(generation, hostID, status) values(?, ?, ?)"
	// Database ops below
	db, err := OpenSQLiteDBOptimized(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return
	}
	stmt, err := tx.Prepare(_stmt)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	for pack := range c {
		_, err = stmt.Exec(
			pack.genID,
			pack.hostID,
			pack.status,
		)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	// Commit at the end
	tx.Commit()
}

// WriteTransmission records the ID's of genotype node that
// are transmitted between hosts.
func (l *SQLiteLogger) WriteTransmission(c <-chan TransmissionPackage) {
	tableName := fmt.Sprintf("Transmission%03d", l.instanceID)
	path := l.transmissionPath
	_stmt := "insert into " + tableName + "(generation, fromHostID, toHostID, nodeID) values(?, ?, ?, ?)"
	// Database ops below
	db, err := OpenSQLiteDBOptimized(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return
	}
	stmt, err := tx.Prepare(_stmt)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	for pack := range c {
		_, err = stmt.Exec(
			pack.genID,
			pack.fromHostID,
			pack.toHostID,
			pack.nodeID.String(),
		)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	// Commit at the end
	tx.Commit()
}

// OpenSQLiteDBOptimized establishes a database connection using WAL
// and exclusive locking.
func OpenSQLiteDBOptimized(path string) (*sql.DB, error) {
	return OpenSQLiteDB(path, "?_journal=WAL&_locking=EXCLUSIVE&_sync=NORMAL")
}

// OpenSQLiteDB establishes a database connection using
// the given connection string.
func OpenSQLiteDB(path, connectionString string) (*sql.DB, error) {
	dsn := "file:%s%s"
	db, err := sql.Open("sqlite3", fmt.Sprintf(dsn, path, connectionString))
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Finalize should create summary tables that join the databases together
// Attach databases first then create new tables using create table as

// select a.nodeID, b.genotypeID, b.sequence
// from genotype.Node001 as a
// inner join nodes.Genotype001 as b
// on a.genotypeID = b.genotypeID
// order by b.sequence asc;
