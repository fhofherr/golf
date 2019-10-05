//nolint: goconst, funlen
package log_test

import (
	"fmt"
	"testing"

	"github.com/fhofherr/golf/log"
)

func TestContextualLogger(t *testing.T) {
	tests := []struct {
		name     string
		factory  func(*log.TestLogger) log.Logger
		ctxkvs   []interface{}
		kvs      []interface{}
		nEntries int
		pred     func(log.TestLogEntry) bool
	}{
		{
			name: "nil logger",
			factory: func(*log.TestLogger) log.Logger {
				return nil
			},
			ctxkvs: []interface{}{"key1", "value1"},
			kvs:    []interface{}{"key2", "value2"},
			// All entries would match the predicate, but we expect no matching
			// entries.
			pred: func(log.TestLogEntry) bool {
				return true
			},
		},
		{
			name:     "non-nil logger",
			ctxkvs:   []interface{}{"key1", "value1"},
			kvs:      []interface{}{"key2", "value2"},
			nEntries: 1,
			pred: func(e log.TestLogEntry) bool {
				return e["key1"] == "value1" && e["key2"] == "value2"
			},
		},
		{
			name:     "missing context value",
			ctxkvs:   []interface{}{"key1"},
			kvs:      []interface{}{"key2", "value2"},
			nEntries: 1,
			pred: func(e log.TestLogEntry) bool {
				return e["key1"] == "error: missing value" && e["key2"] == "value2"
			},
		},
		{
			name:     "empty context kvs",
			ctxkvs:   []interface{}{},
			kvs:      []interface{}{"key2", "value2"},
			nEntries: 1,
			pred: func(e log.TestLogEntry) bool {
				return e["key2"] == "value2"
			},
		},
		{
			name:     "empty kvs",
			ctxkvs:   []interface{}{"key1", "value1"},
			kvs:      []interface{}{},
			nEntries: 1,
			pred: func(e log.TestLogEntry) bool {
				return e["key1"] == "value1"
			},
		},
		{
			name: "nested contextual loggers",
			factory: func(tl *log.TestLogger) log.Logger {
				return log.With(tl, "key1", "value1")
			},
			ctxkvs:   []interface{}{"key2", "value2"},
			kvs:      []interface{}{"key3", "value3"},
			nEntries: 1,
			pred: func(e log.TestLogEntry) bool {
				return e["key1"] == "value1" && e["key2"] == "value2" && e["key3"] == "value3"
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.factory == nil {
				tt.factory = func(tl *log.TestLogger) log.Logger {
					return tl
				}
			}
			tl := &log.TestLogger{}
			logger := log.With(tt.factory(tl), tt.ctxkvs...)
			log.Log(logger, tt.kvs...)
			tl.AssertHasMatchingLogEntries(t, tt.nEntries, tt.pred)
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

		{"1 key value pair - writer logger", 1, &log.TestLogger{}},
		{"2 key value pairs - writer logger", 2, &log.TestLogger{}},
		{"10 key value pairs - writer logger", 10, &log.TestLogger{}},
		{"100 key value pairs - writer logger", 100, &log.TestLogger{}},
		{"1000 key value pairs - writer logger", 1000, &log.TestLogger{}},
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
				logger.Log(kvs...) // nolint: errcheck
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
				logger.Log(kvs...) // nolint: errcheck
			}
		})
	}
}
