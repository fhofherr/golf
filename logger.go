package golf

import "fmt"

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

// Logger is the fundamental interface for all log operations.
//
// Log creates a log entry from kvs, a variadic sequence of alternating keys
// and values. If kvs has a length of one implementations MUST convert the
// single value to a string and treat it as a simple log message.
//
// Implementations MUST be safe for concurrent use by multiple
// goroutines. In particular, any implementation of Logger that appends to kvs
// or modifies, or retains any of its elements MUST make a copy first.
//
// In addition to the above basic contract an implementation of Logger MAY
// extend it with the following additional rules:
//
// Implementations MAY look for the keywords "level" or "lvl" and assume that
// the value of those keywords defines a logging level. Implementations MUST
// always support both keywords, "level" and "lvl". Additionally they MUST
// expect callers to use both keywords interchangeably. If kvs contains "level"
// and "lvl" at the same time they MUST give preference to "level". Should kvs
// contain neither "level" nor "lvl" they MUST assume an appropriate default
// level. Implementations MAY skip creating an log entry based on the level or
// the assumed default level.
//
// Furthermore implementations MAY look for the keywords "message" or "msg" and
// assume that the value of those keywords defines a log message which may be
// treated specially. Implementations MUST always support both keywords,
// "message" and "msg". Additionally they MUST expect callers to use both
// keywords interchangeably. If kvs contains "message" and "msg" at the same
// time they MUST give preference to "message". Should kvs contain neither
// "message" nor "msg" they MUST assume an appropriate default message.
type Logger interface {
	Log(kvs ...interface{})
}

// Log calls logger.Log with the passed kvs, if logger is not nil.
//
// Should logger.Log return an error, Log discards it.
func Log(logger Logger, kvs ...interface{}) {
	if logger == nil {
		return
	}
	logger.Log(kvs...)
}

// Logf formats kvs using fmt.Sprintf and passes it to Log as the only value
// of kvs.
func Logf(logger Logger, format string, kvs ...interface{}) {
	Log(logger, fmt.Sprintf(format, kvs...))
}
