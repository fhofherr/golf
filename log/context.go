package log

// Wither wraps the With method which is used by the implementing type
// to concatenate the passed key-value pairs and its internal store of key-value
// pairs. The implementor then returns a new version of itself containing the
// concatenated key value pairs.
//
// Implementors must make sure, that the wither is safe for concurrent use
// by multiple go routines. Especially they must make sure to copy their
// internal version of key-value pairs as well as the passed key-value pairs.
type Wither interface {
	With(...interface{}) Logger
}

// With creates a contextual-logger, i.e. a logger that will always add the kvs
// passed to With to the final log entry.
//
// It is safe to pass nil for the logger. In this case nil will be returned.
func With(logger Logger, kvs ...interface{}) Logger {
	// Bail out if there is nothing to copy.
	if logger == nil || len(kvs) == 0 {
		return logger
	}
	if wither, ok := logger.(Wither); ok {
		return wither.With(kvs...)
	}
	ncpy := len(kvs)
	if ncpy%2 != 0 {
		ncpy = ncpy + 1
	}
	cpy := make([]interface{}, ncpy, ncpy)
	copyWithMissing(cpy, kvs)
	return contextualLogger{
		ctxkvs: cpy,
		logger: logger,
	}
}

type contextualLogger struct {
	ctxkvs []interface{}
	logger Logger
}

// With adds the passed key value pairs to the contextual loggers internal
// representation of key value pairs and returns a new contextualLogger.
func (c contextualLogger) With(kvs ...interface{}) Logger {
	if len(kvs) == 0 {
		return c
	}
	ctxsiz := len(c.ctxkvs)
	ncpy := ctxsiz + len(kvs)
	if ncpy%2 != 0 {
		ncpy++
	}
	cpy := make([]interface{}, ncpy, ncpy)
	// We know ctxsiz is even.
	copy(cpy, c.ctxkvs)
	copyWithMissing(cpy[ctxsiz:], kvs)
	return contextualLogger{
		ctxkvs: cpy,
		logger: c.logger,
	}
}

// Log merges the passed kvs into c's kvs and then calls log on c.logger.
func (c contextualLogger) Log(kvs ...interface{}) error {
	if c.logger == nil {
		return nil
	}
	nctx, n := len(c.ctxkvs), len(kvs)
	merged := make([]interface{}, nctx+n)
	copy(merged, c.ctxkvs)
	copy(merged[nctx:], kvs)
	return c.logger.Log(merged...)
}

func copyWithMissing(dst, src []interface{}) {
	copy(dst, src)
	if len(src) < len(dst) {
		dst[len(dst)-1] = "error: missing value"
	}
}
