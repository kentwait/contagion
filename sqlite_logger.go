package contagiongo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

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
		db, err := OpenSQLiteDB(path)
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
	db, err := OpenSQLiteDB(path)
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
	db, err := OpenSQLiteDB(path)
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
	db, err := OpenSQLiteDB(path)
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
	db, err := OpenSQLiteDB(path)
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
	db, err := OpenSQLiteDB(path)
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
	db, err := OpenSQLiteDB(path)
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

func OpenSQLiteDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}
