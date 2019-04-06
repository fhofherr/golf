package log

// Logger is the fundamental interface for all log operations.
//
// Logger is the same interface as GoKit's Logger interface
// (https://godoc.org/github.com/go-kit/kit/log#Logger) and follows the same
// contract: Log creates a log event from kvs, a variadic sequence of alternating
// keys and values. Implementations must be safe for concurrent use by multiple
// goroutines. In particular, any implementation of Logger that appends to
// kvs or modifies or retains any of its elements must make a copy first.
//
// In addition to the basic contract specified by GoKit implementations of
// Logger MAY extend it with the following additional rules.
//
// Implementations MAY look for the keywords "level" or "lvl" and assume that
// the value of those keywords defines a logging level. Implementations must
// always support both keywords, "level" and "lvl". Additionally they must
// expect callers to use both keywords interchangeably. If kvs contains "level"
// and "lvl" at the same time they must give preference to "level". Should kvs
// contain neither "level" nor "lvl" they must assume an appropriate default
// level. Implementations MAY skip creating an log entry based on the level or
// the assumed default level.
//
// Furthermore implementations MAY look for the keywords "message" or "msg" and
// assume that the value of those keywords defines a log message which may be
// treated specially. Implementations must  always support both keywords,
// "message" and "msg". Additionally they must expect callers to use both
// keywords interchangeably. If kvs contains "message" and "msg" at the same
// time they must give preference to "message". Should kvs contain neither
// "message" nor "msg" they must assume an appropriate default message.
type Logger interface {
	Log(kvs ...interface{}) error
}
