package log_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/fhofherr/golf/log"
	"github.com/stretchr/testify/assert"
)

func TestWriterLogger_Log(t *testing.T) {
	tests := []struct {
		name         string
		kvs          []interface{}
		expectedJSON string
	}{
		{
			name:         "even number of kvs",
			kvs:          []interface{}{"key1", "value1", "key2", "value2"},
			expectedJSON: `{"key1": "value1", "key2": "value2"}`,
		},
		{
			name:         "odd number of kvs",
			kvs:          []interface{}{"key1", "value1", "key2"},
			expectedJSON: `{"key1": "value1", "key2": "error: missing value"}`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := log.NewWriterLogger(&buf, log.JSONFormatter)
			err := logger.Log(tt.kvs...)
			assert.JSONEq(t, tt.expectedJSON, buf.String())
			assert.NoError(t, err)
		})
	}
}

func TestWriterLogger_Log_WriterError(t *testing.T) {
	err := errors.New("some error")
	logger := log.NewWriterLogger(log.ErrorWriter{Err: err}, log.PlainTextFormatter)
	actual := logger.Log("key", "value")
	assert.EqualError(t, actual, fmt.Sprintf("write log entry: %v", err))
}

func TestWriterLogger_Log_FormatterError(t *testing.T) {
	err := errors.New("some error")
	logger := log.NewWriterLogger(ioutil.Discard, func([]interface{}) ([]byte, error) {
		return nil, err
	})
	actual := logger.Log("key", "value")
	assert.EqualError(t, actual, fmt.Sprintf("format log entry: %v", err))
}

// BenchmarkWriterLogger_Log benchmarks the time it takes for preparing
// a log entry and writing it to the Writer.
//
// In order to prepare a log entry Log has to iterate over all key-value pairs,
// therefore the time it takes to log should increase with the number of
// key-value pairs. Another important source of overhead is the used formatter.
// BenchmarkWriterLogger_Log thus uses different formatters with the same number
// of key-value pairs.
func BenchmarkWriterLogger_Log(b *testing.B) {
	benchmarks := []struct {
		name      string
		nkvs      int
		formatter log.Formatter
	}{
		{"1 key value pair - plain text", 1, log.PlainTextFormatter},
		{"10 key value pairs - plain text", 10, log.PlainTextFormatter},
		{"100 key value pairs - plain text", 100, log.PlainTextFormatter},
		{"1000 key value pairs - plain text", 1000, log.PlainTextFormatter},
		{"1 key value pair - JSON", 1, log.JSONFormatter},
		{"10 key value pairs - JSON", 10, log.JSONFormatter},
		{"100 key value pairs - JSON", 100, log.JSONFormatter},
		{"1000 key value pairs - JSON", 1000, log.JSONFormatter},
	}
	for _, bm := range benchmarks {
		bm := bm
		b.Run(bm.name, func(b *testing.B) {
			logger := log.NewWriterLogger(ioutil.Discard, bm.formatter)
			kvs := log.GenerateKEYVALs(b, bm.nkvs)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				logger.Log(kvs...)
			}
		})
	}

}
