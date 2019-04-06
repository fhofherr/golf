package log

import (
	"fmt"
	"io"
	"sync"
)

type writerLogger struct {
	mu *sync.Mutex
	f  Formatter
	w  io.Writer
}

// NewWriterLogger creates a very basic logger.
//
// The WriterLogger formats every log entry as JSON and writes it to the passed
// writer. It uses a sync.Mutex to protect w from concurrent access. Callers
// must not access w concurrently.
func NewWriterLogger(w io.Writer, f Formatter) Logger {
	return &writerLogger{
		w:  w,
		f:  f,
		mu: &sync.Mutex{},
	}
}

// Log creates a log entry from kvs and writes it to the basic logger's writer.
func (l *writerLogger) Log(kvs ...interface{}) error {
	bs, err := l.f(kvs)
	l.mu.Lock()
	defer l.mu.Unlock()
	_, err = l.w.Write(bs)
	if err != nil {
		return fmt.Errorf("write log entry: %v", err)
	}
	return nil
}
