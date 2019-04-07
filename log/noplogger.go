package log

type nopLogger struct{}

// NewNOPLogger creates a logger that just discards its entries.
func NewNOPLogger() Logger {
	return nopLogger{}
}

// Log discards the passed kvs.
func (nopLogger) Log(kvs ...interface{}) error {
	return nil
}
