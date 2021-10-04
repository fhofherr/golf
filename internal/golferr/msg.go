package golferr

// The following constants represent prefixes to messages logged by golf.
//
// They are exported for documentation purposes, but should usually not be used
// in client code.
const (
	// MsgUnsupported indicates that golf logged a message because an operation
	// was not supported.
	//
	// Finding a message containing this value in your log messages indicates
	// that you use a feature of golf which is not supported by your adapter.
	MsgUnsupported = "[GOLF UNSUPPORTED]"

	// MsgError indicates that golf could not perform an operation because
	// the underlying adapter's implementation does not meet golf's
	// expectations.
	//
	// Finding a message containing this value in your log messages usually
	// indicates a bug or incompatibility in the adapter's implementation.
	MsgError = "[GOLF ERROR]"
)
