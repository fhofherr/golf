package log

import (
	"fmt"
	stdliblog "log"
)

type stdlibAdapter struct {
	logger *stdliblog.Logger
	format Formatter
}

// NewStdlib instantiates golf's standard library logger adapter using the passed
// logger. Log entries are formatted using the passed formatter. If f is nil
// log.PlainTextFormatter is used.
func NewStdlib(l *stdliblog.Logger, f Formatter) Logger {
	if f == nil {
		f = PlainTextFormatter
	}
	return stdlibAdapter{
		logger: l,
		format: f,
	}
}

// Log writes kvs using the standard library logger.
func (l stdlibAdapter) Log(kvs ...interface{}) error {
	msg, err := l.format(kvs)
	if err != nil {
		return fmt.Errorf("formatter error: %v", err)
	}
	// According to the documentation of the Go standard library the call depth
	// passed to Output is 2 for all pre-defined paths. Since we just add another
	// predefined path, we set it to two as well.
	err = l.logger.Output(2, string(msg))
	if err != nil {
		return fmt.Errorf("logger error: %v", err)
	}
	return nil
}
