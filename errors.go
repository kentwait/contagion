package contagiongo

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
	UnrecognizedKeywordError  = "unrecognized keyword %s for %s"
	FileParsingError          = "error in line %d: %s"
)

const (
	IdenticalPointerError    = "memory address of %s (%p) and %s (%p) are identical"
	NotIdenticalPointerError = "memory address of %s (%p) and %s (%p) are not identical"
)
