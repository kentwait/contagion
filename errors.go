package contagiongo

import (
	"fmt"
	"strings"
)

const (
	// IntKeyNotFoundError is the message for "Integer key not found" errors
	IntKeyNotFoundError = "key %d not found"

	// IntKeyExists is the message printed when a given key already exists
	IntKeyExists = "key %d already exists"

	// ModelExistsError is the message printed when an intrahost model
	// has already been assigned to a host.
	ModelExistsError = "model %s (%d) already exists"
	EmptyModelError  = "model does not exist"

	// GraphPathogenTypeAssertionError is the message printed when
	// an GraphPathogen cannot be asserted for an interface
	GraphPathogenTypeAssertionError = "error asserting PathogenNode interface"

	// ZeroItemsError is the message for errors where at least one item must be
	// passed into the function.
	ZeroItemsError = "one or more items are required"

	InvalidFloatParameterError  = "invalid %s %f, %s"
	InvalidIntParameterError    = "invalid %s %d, %s"
	InvalidStringParameterError = "invalid %s %s, %s"
)

const (
	UnequalFloatParameterError = "expected %s %f, instead got %f"
	EqualFloatParameterError   = "%s should not be equal: %f, %f"
	FloatNotBetweenError       = "expected %s between %f and %f, instead got %f"

	UnequalIntParameterError = "expected %s %d, instead got %d"
	EqualIntParameterError   = "%s should not be equal: %d, %d"
	IntNotBetweenError       = "expected %s between %d and %d, instead got %d"

	UnequalStringParameterError = "expected %s %s, instead got %s"
	EqualStringParameterError   = "%s should not be idenitical: %s, %s"

	UnexpectedErrorWhileError = "encountered error while %s: %s"
	ExpectedErrorWhileError   = "expected an error while %s, instead got none"
	UnrecognizedKeywordError  = "%s is not a valid keyword for %s"
	FileParsingError          = "error in line %d: %s"
)

const (
	IdenticalPointerError    = "memory address of %s (%p) and %s (%p) are identical"
	NotIdenticalPointerError = "memory address of %s (%p) and %s (%p) are not identical"
)

// Errors related to the adjacency matrix
const (
	// ConnectionExistsError indicates that a connection between
	// the source host and the destination host exists and
	// has the following value in float64.
	ConnectionExistsError = "connection (%d,%d): %f already exists"
	// ConnectionExistsError indicates that a connection between
	// hosts a and b (int) does not exist.
	ConnectionDoesNotExistError = "connection (%d,%d) does not exist"
	// SelfLoopError indicates that the start and end host are the
	// same based on host ID, which results in a self-loop.
	SelfLoopError = "connection stats and ends at the same host (%d)"
)

// Errors related to the data logger
const (
	// FileExistsError indicates that a file exists at the
	// given path (string).
	FileExistsError = "error creating new file: %s already exists"
	FileOpenError   = "error opening file: %v"
	FileWriteError  = "error writing to file: %v"
	FileSyncError   = "error commiting contents of the file to disk: %v"
	// SQLOpenError indicates that an error was encountered while
	// open a database connection. Includes the error returned
	// by sql.Open.
	SQLOpenError = "error opening SQL connection: %v"
	// SQLExecError indicates that an error was encountered while executing
	// an SQL statement. Returns the error raised by the database
	// connection and the SQL statement that produced the error.
	SQLExecError = "error executing SQL statement: %v (SQL statement: %s)"
)

func checkKeyword(text, category string, keywords ...string) error {
	for _, kw := range keywords {
		if strings.ToLower(text) == strings.ToLower(kw) {
			return nil
		}
	}
	return fmt.Errorf(UnrecognizedKeywordError, text, category)
}
