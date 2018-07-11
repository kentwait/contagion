package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Flags that the program accepts
	var ind bool
	flag.BoolVar(&ind, "independent", false, "Assumes multiple independent starts with its own config and pathogen sequence file")
	var endCommit bool
	flag.BoolVar(&endCommit, "commit_once", false, "Commits only once")
	var outPath string
	flag.StringVar(&outPath, "out", "", "Location where to create sqlite3 file (Required)")
	// Flags to skip tables
	var skipGenotypeFreq bool
	flag.BoolVar(&skipGenotypeFreq, "skip_freq", false, "Skip GenotypeFreq tables")
	var skipGenotype bool
	flag.BoolVar(&skipGenotype, "skip_genotype", false, "Skip Genotype tables")
	var skipNode bool
	flag.BoolVar(&skipNode, "skip_node", false, "Skip Node tables")
	var skipStatus bool
	flag.BoolVar(&skipStatus, "skip_status", false, "Skip Status tables")
	var skipTrans bool
	flag.BoolVar(&skipTrans, "skip_trans", false, "Skip Transmission tables")
	var skipTree bool
	flag.BoolVar(&skipTree, "skip_tree", false, "Skip Tree tables")
	// Create virtual tables
	var makeGenotypeFreqView bool
	flag.BoolVar(&makeGenotypeFreqView, "genotype_freq_view", false, "Create GenotypeFreq - Genotype virtual table")
	flag.Parse()

	// Check number of arguments
	if flag.NArg() < 1 {
		fmt.Println("CSV basepath was not specified!")
		flag.Usage()
	}
	if flag.NArg() > 1 && ind {
		fmt.Println("Only one CSV basepath should be provided if using the -independent flag")
		flag.Usage()
	}
	// Check if out Path is set
	if outPath == "" {
		fmt.Println("-out was not specified")
	}
	// make basepath each subdir if ind is specified
	var csvDirPaths []string
	if ind {
		if flag.NArg() == 1 {
			baseDirPath := filepath.Clean(flag.Arg(0))
			_basepaths, err := filepath.Glob(baseDirPath)
			if err != nil {
				panic(err)
			}
			for _, path := range _basepaths {
				fi, err := os.Stat(path)
				if err != nil {
					panic(err)
				}
				if fi.IsDir() {
					csvDirPaths = append(csvDirPaths, path)
				}
			}
		} else {
			panic("only one directory can be set when the -independent flag is used")
		}
	} else {
		for c := 0; c < flag.NArg(); c++ {
			path := filepath.Clean(flag.Arg(c))
			csvDirPaths = append(csvDirPaths, path)
		}
	}

	// Open database connection
	db, err := openSQLiteDBOptimized(outPath)
	if err != nil {
		panic(err)
	}
	// Table names based on content type of CSV
	tableNameMap := make(map[string]string)
	tableNameMap["freq"] = "GenotypeFreq"
	tableNameMap["g"] = "Genotype"
	tableNameMap["n"] = "Node"
	tableNameMap["status"] = "Status"
	tableNameMap["trans"] = "Transmission"
	tableNameMap["tree"] = "Tree"
	// Columns based on content type of CSV
	columnNameMap := make(map[string]string)
	columnNameMap["freq"] = "(id integer not null primary key, instance int, generation int, hostID int, genotypeID text, freq int)"
	columnNameMap["g"] = "(id integer not null primary key, genotypeID text, sequence text)"
	columnNameMap["n"] = "(id integer not null primary key, nodeID text, genotypeID text)"
	columnNameMap["status"] = "(id integer not null primary key, instance int, generation int, hostID int, status int)"
	columnNameMap["trans"] = "(id integer not null primary key, instance int, generation int, fromHostID int, toHostID int, nodeID text)"
	columnNameMap["tree"] = "(id integer not null primary key, instance int, generation int, hostID int, parentNodeID text, nodeID text)"
	// Insert statement based on content type of CSV
	insertStmtMap := make(map[string]string)
	insertStmtMap["freq"] = "insert into %s (instance, generation, hostID, genotypeID, freq) values(?, ?, ?, ?, ?)"
	insertStmtMap["g"] = "insert into %s (genotypeID, sequence) values(?, ?)"
	insertStmtMap["n"] = "insert into %s (nodeID, genotypeID) values(?, ?)"
	insertStmtMap["status"] = "insert into %s (instance, generation, hostID, status) values(?, ?, ?, ?)"
	insertStmtMap["trans"] = "insert into %s (instance, generation, fromHostID, toHostID, nodeID) values(?, ?, ?, ?, ?)"
	insertStmtMap["tree"] = "insert into %s (instance, generation, hostID, parentNodeID, nodeID) values(?, ?, ?, ?, ?)"

	// Path to folder with CSV files to process
	// Accepts one or more args, each representing a folder path
	// If -independent flag is on, the path should be the root path of
	// multiple independent runs each with its own config and pathogens file
	fileCounter := 0
	startTime := time.Now()
	for c, csvDirPath := range csvDirPaths {

		// Find all CSVs in the folder
		globString := filepath.Join(csvDirPath, "*.csv")
		csvPaths, err := filepath.Glob(globString)
		if err != nil {
			panic(err)
		}
		// Check if anything matches
		if len(csvPaths) < 1 {
			log.Fatalf("%s did not return any matches", globString)
		}

		var tx *sql.Tx
		if endCommit {
			tx, err = db.Begin()
			if err != nil {
				panic(err)
			}
		}
		for _, csvPath := range csvPaths {
			// Open CSV
			f, err := os.Open(csvPath)
			if err != nil {
				f.Close()
				panic(err)
			}
			// defer f.Close()

			// Get dirpath, filename
			// Then split filename to get content type and iteration number
			_, csvFilename := filepath.Split(csvPath)
			splittedCsvFilename := strings.Split(csvFilename, ".")
			contentType := splittedCsvFilename[len(splittedCsvFilename)-2]
			// iter, _ := strconv.Atoi(splittedCsvFilename[len(splittedCsvFilename)-3])

			// Set table name and column titles based on contentType
			tableName := tableNameMap[contentType]
			columnNames := columnNameMap[contentType]
			insertStmt := fmt.Sprintf(insertStmtMap[contentType], tableName)

			// Skip table?
			switch {
			case tableName == "GenotypeFreq" && skipGenotypeFreq:
				continue
			case tableName == "Genotype" && skipGenotype:
				continue
			case tableName == "Node" && skipNode:
				continue
			case tableName == "Status" && skipStatus:
				continue
			case tableName == "Transmission" && skipTrans:
				continue
			case tableName == "Tree" && skipTree:
				continue
			}

			// Read using buffered reader
			scanner := bufio.NewScanner(f)

			// Adjust the capacity to your need (max characters in line)
			// TODO: get size of second line and adjust buffer size by lengh of
			// sequence
			// const maxCapacity = 512*1024
			// buf := make([]byte, maxCapacity)
			// scanner.Buffer(buf, maxCapacity)
			splitter := regexp.MustCompile(`\s*\,\s*`)

			// Each file begins a transaction
			if !endCommit {
				tx, err = db.Begin()
				if err != nil {
					panic(err)
				}
			}
			i := 0
			for scanner.Scan() {
				scanner.Text()
				// Create a new table is a table doesnt exist
				createStmt := "create table if not exists %s %s;"
				createStmt = fmt.Sprintf(createStmt, tableName, columnNames)
				_, err = tx.Exec(createStmt)
				if err != nil {
					log.Fatalf("%q: %s", err, createStmt)
				}
				// tx.Commit()
				break
			}

			// Prepare the statement
			stmt, err := tx.Prepare(insertStmt)
			if err != nil {
				panic(err)
			}
			defer stmt.Close()
			for scanner.Scan() {
				line := scanner.Text()
				// Convert to []interface{}
				stringValues := splitter.Split(line, -1)
				if ind {
					stringValues[0] = strconv.Itoa(c)
				}
				values := make([]interface{}, len(stringValues))
				for i, v := range stringValues {
					values[i] = v
				}
				// Execute statement using
				_, err = stmt.Exec(values...)
				if err != nil {
					panic(fmt.Sprintln(err, stringValues))
				}

				i++
			}
			fmt.Print(csvFilename)
			if !endCommit {
				// Commit at the end of the file
				tx.Commit()
				fmt.Print(", committed.")
			}
			fmt.Print("\n")
			f.Close()
			fileCounter++
		}
		if endCommit {
			// Commit at the end of the file
			tx.Commit()
		}
	}
	elapsed := time.Since(startTime)

	if makeGenotypeFreqView {
		db.Exec(`create view if not exists GenotypeFreqView as
				 select 
					 GenotypeFreq.instance, 
					 GenotypeFreq.generation, 
					 GenotypeFreq.hostID, 
					 GenotypeFreq.genotypeID, 
					 GenotypeFreq.freq, 
					 Genotype.sequence 
				 from GenotypeFreq 
				 left join 
				 	 Genotype on GenotypeFreq.genotypeID == Genotype.genotypeID;`)
	}
	db.Close()
	fmt.Println("Finished.")
	fmt.Printf("Completed in %v\n", elapsed)
}

func newTableIfNot(path, tableName, cols string) error {
	db, err := openSQLiteDBOptimized(path)
	if err != nil {
		return err
	}
	defer db.Close()
	// cols example:
	// (id integer not null primary key, genotypeID text, sequence text)
	_sqlStmt := `
create table if not exists %s %s;
delete from %s;
`
	sqlStmt := fmt.Sprintf(_sqlStmt, tableName, cols, tableName)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("%q: %s", err, sqlStmt)
	}
	return nil
}

func openSQLiteDBOptimized(path string) (*sql.DB, error) {
	return openSQLiteDB(path, "?_journal=WAL&_locking=EXCLUSIVE&_sync=NORMAL")
}
func openSQLiteDB(path, connectionString string) (*sql.DB, error) {
	dsn := "file:%s%s"
	db, err := sql.Open("sqlite3", fmt.Sprintf(dsn, path, connectionString))
	if err != nil {
		return nil, err
	}
	return db, nil
}
