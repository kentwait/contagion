package contagiongo

const (
	// IntKeyNotFoundError is the message for "Integer key not found" errors
	IntKeyNotFoundError = "key %d not found"

	// IntKeyExists is the message printed when a given key already exists
	IntKeyExists = "key %d already exists"

	// IntrahostModelExistsError is the message printed when an intrahost model
	// has already been assigned to a host.
	IntrahostModelExistsError = "intrahost model %s (%d) already exists"
	EmptyIntrahostModelError  = "intrahost model does not exist"

	// FitnessModelExistsError is the message printed when a fitness model
	// has already been assigned to a host.
	FitnessModelExistsError = "fitness model %s (%d) already exists"
	EmptyFitnessModelError  = "fitness model does not exist"

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
	UnequalFloatParameterError  = "expected %s %f, instead got %f"
	UnequalIntParameterError    = "expected %s %d, instead got %d"
	FloatNotBetweenError        = "expected %s between %f and %f, instead got %f"
	IntNotBetweenError          = "expected %s between %d and %d, instead got %d"
	UnequalStringParameterError = "expected %s %s, instead got %s"
	UnexpectedErrorWhileError   = "encountered error while %s: %s"
	ExpectedErrorWhileError     = "expected an error while %s, instead got none"
	UnrecognizedKeywordError    = "unrecognized keyword %s for %s"
	FileParsingError            = "error in line %d: %s"
)
