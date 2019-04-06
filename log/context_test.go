package log_test

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/fhofherr/golf/log"
	"github.com/fhofherr/golf/log/logtest"
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

func BenchmarkContextualLogger_Log(b *testing.B) {
	benchmarks := []struct {
		name    string
		nctxkvs int
		nkvs    int
	}{
		{"1 key value pair", 0, 1},
		{"2 key value pairs", 1, 1},
		{"10 key value pairs", 5, 5},
		{"100 key value pairs", 50, 50},
		{"1000 key value pairs", 500, 500},
	}
	for _, bm := range benchmarks {
		bm := bm
		b.Run(bm.name, func(b *testing.B) {
			ctxkvs := logtest.GenerateKEYVALs(b, bm.nctxkvs)
			kvs := logtest.GenerateKEYVALs(b, bm.nkvs)
			logger := log.NewWriterLogger(ioutil.Discard, log.PlainTextFormatter)
			logger = log.With(logger, ctxkvs...)
			for i := 0; i < b.N; i++ {
				logger.Log(kvs...)
			}
		})
	}

}

func BenchmarkNestedContextualLogger_Log(b *testing.B) {
	benchmarks := []struct {
		name     string
		nctx1kvs int
		nctx2kvs int
		nkvs     int
	}{
		{"1 key value pair", 0, 0, 1},
		{"2 key value pairs", 0, 1, 1},
		{"10 key value pairs", 3, 2, 5},
		{"100 key value pairs", 25, 25, 50},
		{"1000 key value pairs", 250, 250, 500},
	}
	for _, bm := range benchmarks {
		bm := bm
		b.Run(bm.name, func(b *testing.B) {
			ctx1kvs := logtest.GenerateKEYVALs(b, bm.nctx1kvs)
			ctx2kvs := logtest.GenerateKEYVALs(b, bm.nctx2kvs)
			kvs := logtest.GenerateKEYVALs(b, bm.nkvs)
			logger := log.NewWriterLogger(ioutil.Discard, log.PlainTextFormatter)
			logger = log.With(logger, ctx1kvs...)
			logger = log.With(logger, ctx2kvs...)
			for i := 0; i < b.N; i++ {
				logger.Log(kvs...)
			}
		})
	}

}
