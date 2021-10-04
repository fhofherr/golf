package golfstdlib_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/fhofherr/golf"
	"github.com/fhofherr/golf/golfstdlib"
	"github.com/fhofherr/golf/internal/golferr"
	"github.com/fhofherr/golf/internal/testsupport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStdlibAdapter_Log(t *testing.T) {
	tests := []struct {
		name      string
		formatter golfstdlib.Formatter
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
			formatter: golfstdlib.JSONFormatter,
			kvs:       []interface{}{"key", "value"},
			expected:  `{"key":"value"}` + "\n",
		},
		{
			name:     "single value default formatter",
			kvs:      []interface{}{"just some message"},
			expected: "just some message\n",
		},
		{
			name: "Formatter error",
			formatter: func(_ []interface{}) ([]byte, error) {
				return nil, assert.AnError
			},
			expected: fmt.Sprintf("%s: stdlibAdapter: %v\n", golferr.MsgError, assert.AnError),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			w := &strings.Builder{}
			logger := golfstdlib.NewLogger(log.New(w, "", 0), golfstdlib.WithFormatter(tt.formatter))
			logger.Log(tt.kvs...)
			assert.Equal(t, tt.expected, w.String())
		})
	}
}

func TestStdlibAdapter_Log_IOError(t *testing.T) {
	w := testsupport.NewMockWriter(t)
	w.On("Write", mock.Anything).Return(0, assert.AnError)

	logger := golfstdlib.NewLogger(log.New(w, "", 0))
	logger.Log("key", "value")

	err := golf.Error(logger)
	assert.Same(t, assert.AnError, err)
}
