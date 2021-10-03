package golf

// With tries to create a logger that always adds any values in kvs to the
// arguments passed to Logger.Log.
//
// To enable this the passed logger MUST provide a With method with the
// following signature:
//
//     With(...interface{}) interface{}
//
// The value returned by the With method must implement the Logger interface.
// If this is not the case a message is logged using logger and
// the return value of the With method is ignored. In this case logger is
// used to log a message and then returned instead.
//
// It is safe to pass nil for the logger. In this case nil will be returned.
func With(logger Logger, kvs ...interface{}) Logger {
	// Bail out if there is nothing to copy.
	if logger == nil || len(kvs) == 0 {
		return logger
	}
	w, ok := logger.(wither)
	if !ok {
		Logf(logger, "%s: %T does not implement With", MsgUnsupported, logger)
		return logger
	}
	l, ok := w.With(kvs...).(Logger)
	if !ok {
		Logf(logger, "%s: (%T).With did not return a Logger", MsgError, logger)
		return logger
	}
	return l
}

type wither interface {
	With(...interface{}) interface{}
}
