package log

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

// GenerateKEYVALs generates a fixed amount of key value pairs.
func GenerateKEYVALs(tb testing.TB, n int) []interface{} {
	if n < 0 {
		tb.Fatalf("expected n >= 0; got %d", n)
	}
	var kvs []interface{}
	for i := 0; i < n; i++ {
		kvs = append(kvs, fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	return kvs
}

// StressTestLogger exercises the logger returned by factory using multiple Go routines.
//
// Each Go routine writes nMSGs before it terminates.
func StressTestLogger(tb testing.TB, factory func() Logger, nGoRoutines, nMSGs int) {
	if nGoRoutines < 1 {
		tb.Fatalf("expected nGoRoutines >= 1; got %d", nGoRoutines)
	}
	if nMSGs < 1 {
		tb.Fatalf("expected nMSGs >= 1; got %d", nMSGs)
	}
	start := make(chan struct{})
	errc := make(chan error, nGoRoutines)
	wg := &sync.WaitGroup{}
	wg.Add(nGoRoutines)
	for i := 0; i < nGoRoutines; i++ {
		go func(i int) {
			defer wg.Done()
			// Block until the start channel is closed, then write the messages
			<-start
			logger := factory()
			for j := 0; j < nMSGs; j++ {
				err := logger.Log("go routine", i+1, "msg no", j+1)
				if err != nil {
					errc <- err
					return
				}
			}

		}(i)
	}
	close(start) // Trigger the loggers
	wg.Wait()
	for {
		select {
		case err := <-errc:
			tb.Error(err)
		default:
			return
		}
	}
}

// TestLogEntry represents a single log entry stored by the TestLogger.
type TestLogEntry map[interface{}]interface{}

// NewTestLogEntry creates a new TestLogEntry from the passed kvs. The number
// of kvs must be even, otherwise an error is returned.
func NewTestLogEntry(kvs ...interface{}) (TestLogEntry, error) {
	if len(kvs)%2 != 0 {
		return nil, errors.New("number of kvs not even")
	}
	entry := make(map[interface{}]interface{}, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		k, v := kvs[i], kvs[i+1]
		entry[k] = v
	}
	return entry, nil
}

// TestLogger stores the entries internally without formatting them.
// TestLogger is save for concurrent use.
type TestLogger struct {
	mu  sync.Mutex
	log []TestLogEntry
}

// Log adds kvs as an entry to the TestLoggers internal structure of log
// entries.
func (l *TestLogger) Log(kvs ...interface{}) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	entry, err := NewTestLogEntry(kvs...)
	if err != nil {
		return err
	}
	l.log = append(l.log, entry)
	return nil
}

// CountMatchingLogEntries returns the number of log entries matching pred.
func (l *TestLogger) CountMatchingLogEntries(pred func(TestLogEntry) bool) int {
	var count int

	for _, entry := range l.log {
		if pred(entry) {
			count++
		}
	}
	return count
}

// AssertHasMatchingLogEntries searches the TestLogger's log entries for any
// that match pred. It adds a test error if expected does not match the number
// of matching entries. AssertHasMatchingLogEntries returns the number of
// matching entries.
func (l *TestLogger) AssertHasMatchingLogEntries(t *testing.T, expected int, pred func(TestLogEntry) bool) int {
	count := l.CountMatchingLogEntries(pred)
	if count != expected {
		t.Errorf("Expected %d entries; got %d", expected, count)
	}
	return count
}
