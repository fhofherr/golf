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
	logger := log.NewWriterLogger(errorWriter{Err: err}, log.PlainTextFormatter)
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

func BenchmarkWriterLogger_Log(b *testing.B) {
	logger := log.NewWriterLogger(ioutil.Discard, log.PlainTextFormatter)
	for i := 0; i < b.N; i++ {
		// TODO generate a fixed amount of key-value pairs before executing benchmark
		logger.Log("key1", "value1", "key2", 4711)
	}
}

type errorWriter struct {
	Err error
}

func (w errorWriter) Write(p []byte) (n int, err error) {
	return 0, w.Err
}
