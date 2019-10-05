package log

// Log calls logger.Log with the passed kvs, if logger is not nil.
//
// Should logger.Log return an error, Log discards it.
func Log(logger Logger, kvs ...interface{}) {
	if logger == nil {
		return
	}
	logger.Log(kvs...)
}
