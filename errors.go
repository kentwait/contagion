package contagiongo

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
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
)

const (
	IdenticalPointerError    = "memory address of %s (%p) and %s (%p) are identical"
	NotIdenticalPointerError = "memory address of %s (%p) and %s (%p) are not identical"
)

// Errors related to creating a simulation from the configuration file

// DurationTooShortError indicates that the duration of a particular interval
// is shorter than expected.
func DurationTooShortError(interval string, intervalDuration int, condition string, conditionValue int) error {
	return fmt.Errorf("%s (%d) is less than the %s (%d)", interval, intervalDuration, condition, conditionValue)
}

// ExposedDurationTooShortError indicates that the duration in the
// exposed state is too short.
func ExposedDurationTooShortError(intervalDuration int, conditionValue int) error {
	return DurationTooShortError("exposed_duration", intervalDuration, "number of generations", conditionValue)
}

// InfectedDurationTooShortError indicates that the duration in the
// infected state is too short.
func InfectedDurationTooShortError(intervalDuration int, conditionValue int) error {
	return DurationTooShortError("infected_duration", intervalDuration, "number of generations", conditionValue)
}

// InfectiveDurationTooShortError indicates that the duration in the
// infective state is too short.
func InfectiveDurationTooShortError(intervalDuration int, conditionValue int) error {
	return DurationTooShortError("infective_duration", intervalDuration, "number of generations", conditionValue)
}

// RemovedDurationTooShortError indicates that the duration in the
// removed state is too short.
func RemovedDurationTooShortError(intervalDuration int, conditionValue int) error {
	return DurationTooShortError("removed_duration", intervalDuration, "number of generations", conditionValue)
}

// DurationTooLongError indicates that the duration of a particular interval
// is longer than expected.
func DurationTooLongError(interval string, intervalDuration int, condition string, conditionValue int) error {
	return fmt.Errorf("%s (%d) is less than the %s (%d)", interval, intervalDuration, condition, conditionValue)
}

// ExposedDurationTooLongError indicates that the duration in the exposed
// state is too long.
func ExposedDurationTooLongError(intervalDuration int, conditionValue int) error {
	return DurationTooLongError("exposed_duration", intervalDuration, "number of generations", conditionValue)
}

// InfectedDurationTooLongError indicates that the duration in the
// infected state is too long.
func InfectedDurationTooLongError(intervalDuration int, conditionValue int) error {
	return DurationTooLongError("infected_duration", intervalDuration, "number of generations", conditionValue)
}

// InfectiveDurationTooLongError indicates that the duration in the
// infective state is too long.
func InfectiveDurationTooLongError(intervalDuration int, conditionValue int) error {
	return DurationTooLongError("infective_duration", intervalDuration, "number of generations", conditionValue)
}

// RemovedDurationTooLongError indicates that the duration in the removed
// state is too long.
func RemovedDurationTooLongError(intervalDuration int, conditionValue int) error {
	return DurationTooLongError("removed_duration", intervalDuration, "number of generations", conditionValue)
}

// Errors related to model assignment

// ModelExistsError indicates that an existing model already exists
// a new model cannot be assigned to replace it.
func ModelExistsError(modelName string, modelID int) error {
	return fmt.Errorf("model %s (%d) already exists", modelName, modelID)
}

// SetIntrahostModelExistsError indicates that an intrahost model
// has already been assigned to a host.
func SetIntrahostModelExistsError(modelName string, modelID int) error {
	err := ModelExistsError(modelName, modelID)
	return errors.Wrap(err, "setting intrahost model failed")
}

// SetFitnessModelExistsError indicates that a fitness model
// has already been assigned to a host.
func SetFitnessModelExistsError(modelName string, modelID int) error {
	err := ModelExistsError(modelName, modelID)
	return errors.Wrap(err, "setting fitness model failed")
}

// SetTransmissionModelExistsError indicates that a transmission model
// has already been assigned to a host.
func SetTransmissionModelExistsError(modelName string, modelID int) error {
	err := ModelExistsError(modelName, modelID)
	return errors.Wrap(err, "setting transmission model failed")
}

// EmptyModelError indicates that a model should exist but instead
// is nil.
func EmptyModelError() error {
	return fmt.Errorf("model does not exist")
}

// InvalidStateCharError indicates that a character encountered is
// not in the set of expected characters.
func InvalidStateCharError(char string, pos int) error {
	return fmt.Errorf("char %s at position %d is not in the set of expected characters")
}

// Errors related to reading and parsing files

// FileParsingError indicates a parsing error was encountered
// at a particular line in the file. Most likely the file
// was not properly formatted.
func FileParsingError(err error, lineNum int) error {
	return errors.Wrapf(err, "error in line %d", lineNum)
}

// DuplicateSitePositionError indicates that the site in the
// file is not unique and has been included more than once.
func DuplicateSitePositionError(pos int, lineNum int) error {
	err := fmt.Errorf("duplicate position index (%d)", pos)
	return FileParsingError(err, lineNum)
}

// UnequalNumStatesError indicates that the number of states
// specified in a site does not match the number of states
// in another site.
func UnequalNumStatesError(numStates, prevNumStates int, site int, lineNum int) error {
	err := fmt.Errorf("site %d has %d states instead of %d", site, numStates, prevNumStates)
	return FileParsingError(err, lineNum)
}

// InvalidConnectionWeightError indicates that the given connection
// weight is less than 0.
func InvalidConnectionWeightError(wt float64, lineNum int) error {
	err := fmt.Errorf("weight must be non-negative value: %f", wt)
	return FileParsingError(err, lineNum)
}

// Errors related to the adjacency matrix

// ConnectionExistsError indicates that a connection between
// the source host and the destination host exists and
// has the following value in float64.
func ConnectionExistsError(a, b int, value float64) error {
	return fmt.Errorf("connection (%d,%d): %f already exists", a, b, value)
}

// ConnectionDoesNotExistError indicates that a connection between
// hosts a and b (int) does not exist but is expected to exist.
func ConnectionDoesNotExistError(a, b int) error {
	return fmt.Errorf("connection (%d,%d) does not exist", a, b)
}

// SelfLoopError indicates that the start and end host are the
// same based on host ID, which results in a self-loop.
func SelfLoopError(hostID int) error {
	return fmt.Errorf("connection stats and ends at the same host (%d)", hostID)
}

// Errors related to the data logger

// FileExistsError indicates that a file exists at the given path.
func FileExistsError(path string) error {
	return fmt.Errorf("%s already exists", path)
}

// FileExistsCheckError indicates an error was encountered while
// checking if the files exists. This is not the same with an error
// because a file exists.
func FileExistsCheckError(err error, path string) error {
	return errors.Wrapf(err, "check to see if file exists at %s failed", path)
}

// FileDoesNotExistError indicates that a file does not exist
// at the given path when it is expected to.
func FileDoesNotExistError(path string) error {
	return fmt.Errorf("%s does not exist", path)
}

// FileOpenError indicates that an error was encountered while
// opening a file.
func FileOpenError(err error) error {
	return errors.Wrap(err, "opening file failed")
}

// FileWriteError indicates that an error was encountered while
// writing to the file in memory.
func FileWriteError(err error) error {
	return errors.Wrap(err, "writing to file filed")
}

// FileSyncError indicates that an error was encountered while
// the file was being flushed from memory and being written
// to disk.
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

// IntKeyExists indicates that the given integer key already exists.
func IntKeyExists(key int) error {
	return fmt.Errorf("key %d already exists", key)
}

// IntKeyNotFoundError indicates that the given integer key does not exist.
func IntKeyNotFoundError(key int) error {
	return fmt.Errorf("key %d not found", key)
}
