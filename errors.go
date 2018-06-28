package contagiongo

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
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

// FileExistsError indicates that a file exists at the given path.
func FileExistsError(path string) error {
	return fmt.Errorf("%s already exists", path)
}

func FileOpenError(err error) error {
	return errors.Wrap(err, "opening file failed")
}
func FileWriteError(err error) error {
	return errors.Wrap(err, "writing to file filed")
}
func FileSyncError(err error) error {
	return errors.Wrap(err, "commiting file to disk failed")
}

// SQLOpenError indicates that an error was encountered while
// open a database connection. Includes the error returned
// by sql.Open.
func SQLOpenError(err error) error {
	return errors.Wrap(err, "opening SQL connection failed")
}

// SQLExecError indicates that an error was encountered while executing
// an SQL statement. Returns the error raised by the database
// connection and the SQL statement that produced the error.
func SQLExecError(err error, stmt string) error {
	return errors.Wrapf(err, "executing SQL statement failed (%s)", stmt)
}

// SQLBeginTransactionError indicates that an error was encountered
// while a transaction was being initialized.
func SQLBeginTransactionError(err error) error {
	return errors.Wrap(err, "creating SQL transaction failed")
}

// SQLPrepareStatementError indicates that an error was encountered
// while a template SQL statement was being initialized.
func SQLPrepareStatementError(err error, stmt string) error {
	return errors.Wrapf(err, "preparing SQL template statement failed (%s)", stmt)
}

// SQLExecStatementError indicates that an error was encountered
// while a template statement was being substituted with actual values.
func SQLExecStatementError(err error) error {
	return errors.Wrap(err, "executing SQL template statement failed")
}

// Errors related to the motif model

// MotifExistsError indicates that a motif with the same
// sequence and positions already exists in the model.
func MotifExistsError(motifID string) error {
	return fmt.Errorf("motif %s already exists", motifID)
}

// OverlappingMotifError indicates that a particular site
// cannot be used again because it is already being considered
// by another motif in the model.
func OverlappingMotifError(pos int) error {
	return fmt.Errorf("site %d is already considered by another motif", pos)
}

// General errors

// ZeroItemsError indicates that the length of a list or set is
// empty but at least one item is required.
func ZeroItemsError() error {
	return fmt.Errorf("one or more items are required")
}

func checkKeyword(text, category string, keywords ...string) error {
	for _, kw := range keywords {
		if strings.ToLower(text) == strings.ToLower(kw) {
			return nil
		}
	}
	return fmt.Errorf(UnrecognizedKeywordError, text, category)
}
