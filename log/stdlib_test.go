package log_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	stdliblog "log"
	"strings"
	"testing"

	"github.com/fhofherr/golf/log"
	"github.com/stretchr/testify/assert"
)

func TestStdlibAdapter_Log(t *testing.T) {
	tests := []struct {
		name      string
		formatter log.Formatter
		kvs       []interface{}
		expected  string
	}{
		{
			name:     "default formatter",
			kvs:      []interface{}{"key", "value"},
			expected: "key=value\n",
		},
		{
			name:      "JSON formatter",
			formatter: log.JSONFormatter,
			kvs:       []interface{}{"key", "value"},
			expected:  `{"key":"value"}` + "\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			w := &strings.Builder{}
			logger := log.NewStdlib(stdliblog.New(w, "", 0), tt.formatter)
			logger.Log(tt.kvs...)
			assert.Equal(t, tt.expected, w.String())
		})
	}
}

func TestStdlibAdapter_Log_FormatterError(t *testing.T) {
	err := errors.New("some error")
	errorFormatter := func(kvs []interface{}) ([]byte, error) {
		return nil, err
	}
	logger := log.NewStdlib(stdliblog.New(ioutil.Discard, "", 0), errorFormatter)
	actual := logger.Log("key", "value")
	assert.EqualError(t, actual, fmt.Sprintf("formatter error: %v", err))
}

func TestStdlibAdapter_Log_LoggerError(t *testing.T) {
	w := log.ErrorWriter{Err: errors.New("some error")}
	logger := log.NewStdlib(stdliblog.New(w, "", 0), nil)
	actual := logger.Log("key", "value")
	assert.EqualError(t, actual, fmt.Sprintf("logger error: %v", w.Err))
}
