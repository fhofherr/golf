package golf

import (
	"log"
	"sync"
)

type stdlibAdapter struct {
	logger *log.Logger
	format Formatter

	err error
	mu  sync.Mutex // Protects err; other fields are read-only.
}

// NewStdlib instantiates golf's standard library logger adapter using the passed
// logger. Log entries are formatted using the passed formatter. If f is nil
// log.PlainTextFormatter is used.
func NewStdlib(l *log.Logger, f Formatter) Logger {
	if f == nil {
		f = PlainTextFormatter
	}
	return &stdlibAdapter{
		logger: l,
		format: f,
	}
}

// Log writes kvs using the standard library logger.
func (l *stdlibAdapter) Log(kvs ...interface{}) {
	msg, err := l.format(kvs)
	if err != nil {
		l.handleError(err)
		return
	}
	// According to the documentation of the Go standard library the call depth
	// passed to Output is 2 for all pre-defined paths. Since we just add another
	// predefined path, we set it to two as well.
	if err = l.logger.Output(2, string(msg)); err != nil {
		l.handleError(err)
		return
	}
}

// Err returns the last error that occurred during logging.
func (l *stdlibAdapter) Err() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.err
}

func (l *stdlibAdapter) handleError(err error) {
	// Try to log the error. This may fail if the error occurred due to problems
	// writing to the log stream.
	l.logger.Printf("%s: stdlibAdapter: %v", MsgError, err)

	l.mu.Lock()
	defer l.mu.Unlock()
	l.err = err
}
