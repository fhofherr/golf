package log_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/fhofherr/golf/log"
	"github.com/stretchr/testify/assert"
)

func TestContextualLogger(t *testing.T) {
	tests := []struct {
		name     string
		factory  func(io.Writer, log.Formatter) log.Logger
		ctxkvs   []interface{}
		kvs      []interface{}
		expected string
	}{
		{
			name: "nil logger",
			factory: func(io.Writer, log.Formatter) log.Logger {
				return nil
			},
			ctxkvs: []interface{}{"key1", "value1"},
			kvs:    []interface{}{"key2", "value2"},
		},
		{
			name:     "non-nil logger",
			factory:  log.NewWriterLogger,
			ctxkvs:   []interface{}{"key1", "value1"},
			kvs:      []interface{}{"key2", "value2"},
			expected: "key1=value1, key2=value2\n",
		},
		{
			name:     "missing context value",
			factory:  log.NewWriterLogger,
			ctxkvs:   []interface{}{"key1"},
			kvs:      []interface{}{"key2", "value2"},
			expected: "key1=error: missing value, key2=value2\n",
		},
		{
			name:     "empty context kvs",
			factory:  log.NewWriterLogger,
			ctxkvs:   []interface{}{},
			kvs:      []interface{}{"key2", "value2"},
			expected: "key2=value2\n",
		},
		{
			name:     "empty kvs",
			factory:  log.NewWriterLogger,
			ctxkvs:   []interface{}{"key1", "value1"},
			kvs:      []interface{}{},
			expected: "key1=value1\n",
		},
		{
			name: "nested contextual loggers",
			factory: func(w io.Writer, f log.Formatter) log.Logger {
				logger := log.NewWriterLogger(w, f)
				return log.With(logger, "key1", "value1")
			},
			ctxkvs:   []interface{}{"key2", "value2"},
			kvs:      []interface{}{"key3", "value3"},
			expected: "key1=value1, key2=value2, key3=value3\n",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			w := &strings.Builder{}
			logger := tt.factory(w, log.PlainTextFormatter)
			logger = log.With(logger, tt.ctxkvs...)
			log.Log(logger, tt.kvs...)
			assert.Equal(t, tt.expected, w.String())
		})
	}
}

// BenchmarkContextualLogger_Log benchmarks the time it takes the contextual
// Logger to log its entry. In order to concatenate the key-value pairs in the
// context and the key-value pairs passed to Log the contextual logger has to
// traverse both and make a copy. As the actual formatting of the log entry
// adds some overhead as well, the benchmark is done using the nopLogger. In
// order to be able to compare it with the plain writer logger it is repeated
// using a writer logger with a plain text formatter.
func BenchmarkContextualLogger_Log(b *testing.B) {
	benchmarks := []struct {
		name   string
		nkvs   int
		logger log.Logger
	}{
		{"1 key value pair - nopLogger", 1, log.NewNOPLogger()},
		{"2 key value pairs - nopLogger", 2, log.NewNOPLogger()},
		{"10 key value pairs - nopLogger", 10, log.NewNOPLogger()},
		{"100 key value pairs - nopLogger", 100, log.NewNOPLogger()},
		{"1000 key value pairs - nopLogger", 100, log.NewNOPLogger()},

		{"1 key value pair - writer logger", 1,
			log.NewWriterLogger(ioutil.Discard, log.PlainTextFormatter)},
		{"2 key value pairs - writer logger", 2,
			log.NewWriterLogger(ioutil.Discard, log.PlainTextFormatter)},
		{"10 key value pairs - writer logger", 10,
			log.NewWriterLogger(ioutil.Discard, log.PlainTextFormatter)},
		{"100 key value pairs - writer logger", 100,
			log.NewWriterLogger(ioutil.Discard, log.PlainTextFormatter)},
		{"1000 key value pairs - writer logger", 1000,
			log.NewWriterLogger(ioutil.Discard, log.PlainTextFormatter)},
	}
	for _, bm := range benchmarks {
		bm := bm
		b.Run(bm.name, func(b *testing.B) {
			logger := bm.logger
			if bm.nkvs > 1 {
				ctxkvs := log.GenerateKEYVALs(b, bm.nkvs-1)
				logger = log.With(logger, ctxkvs...)
			}
			kvs := log.GenerateKEYVALs(b, 1)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				logger.Log(kvs...)
			}
		})
	}

}

// BenchmarkNestedContextualLogger_Log benchmarks the overhead deep nesting
// of contextual loggers adds to writing the Log entry.
// For this it creates an additional level of nesting for each key-value pair.
// Since the contextual Logger immediately concatenates new key-value pairs
// with its existing context the overhead compared to a flat contextual logger
// should be relatively small. We only use nopLogger as actual logger, since
// we already established that the actual creation of the Log entry adds
// overhead.
func BenchmarkNestedContextualLogger_Log(b *testing.B) {
	benchmarks := []struct {
		name string
		nkvs int
	}{
		{"1 key value pair", 1},
		{"2 key value pairs", 2},
		{"10 key value pairs", 10},
		{"100 key value pairs", 100},
		{"1000 key value pairs", 1000},
	}
	for _, bm := range benchmarks {
		bm := bm
		b.Run(bm.name, func(b *testing.B) {
			var logger log.Logger
			logger = log.NewNOPLogger()
			for i := 0; i < bm.nkvs-1; i++ {
				logger = log.With(logger, fmt.Sprintf("level-%d", i), fmt.Sprintf("value-%d", i))
			}
			// Pre-allocating a slice of kvs and then passing it to Log seems to
			// be faster than passing individual values.
			kvs := []interface{}{
				fmt.Sprintf("level-%d", bm.nkvs-1),
				fmt.Sprintf("value-%d", bm.nkvs-1),
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				logger.Log(kvs...)
			}
		})
	}

}
