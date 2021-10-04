package golfstdlib

import (
	"log"
	"sync"

	"github.com/fhofherr/golf"
)

type options struct {
	formatter Formatter
}

// Option provides an optional argument to NewLogger.
type Option func(*options)

// WithFormatter instructs NewLogger to create an logger using f as Formatter.
func WithFormatter(f Formatter) Option {
	return func(opts *options) {
		opts.formatter = f
	}
}

// Logger adapts an instance of the log.Logger to provide the golf.Logger
// interface.
type Logger struct {
	logger *log.Logger
	format Formatter

	err error
	mu  sync.Mutex // Protects err; other fields are read-only.
}

// NewLogger instantiates golf's standard library logger adapter using the
// passed logger.
func NewLogger(l *log.Logger, opts ...Option) *Logger {
	var loggerOpts options

	for _, o := range opts {
		o(&loggerOpts)
	}
	if loggerOpts.formatter == nil {
		loggerOpts.formatter = PlainTextFormatter
	}
	return &Logger{
		logger: l,
		format: loggerOpts.formatter,
	}
}

// Log writes kvs using the standard library logger.
func (l *Logger) Log(kvs ...interface{}) {
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
func (l *Logger) Err() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.err
}

func (l *Logger) handleError(err error) {
	// Try to log the error. This may fail if the error occurred due to problems
	// writing to the log stream.
	l.logger.Printf("%s: stdlibAdapter: %v", golf.MsgError, err)

	l.mu.Lock()
	defer l.mu.Unlock()
	l.err = err
}
