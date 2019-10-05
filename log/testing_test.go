package log_test

import (
	"errors"
	"testing"

	"github.com/fhofherr/golf/log"
	"github.com/stretchr/testify/assert"
)

func TestNewTestLogEntry(t *testing.T) {
	tests := []struct {
		name  string
		kvs   []interface{}
		entry log.TestLogEntry
		err   error
	}{
		{
			name:  "Create log entry from key-value pair",
			kvs:   []interface{}{"key1", "value1", "key2", "value2"},
			entry: log.TestLogEntry{"key1": "value1", "key2": "value2"},
		},
		{
			name: "Fails for odd number of entries",
			kvs:  []interface{}{"key"},
			err:  errors.New("number of kvs not even"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			entry, err := log.NewTestLogEntry(tt.kvs...)
			assert.Equal(t, tt.entry, entry)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestTestLogger_CountMatchingLogEntries(t *testing.T) {
	tests := []struct {
		name     string
		kvss     [][]interface{}
		pred     func(log.TestLogEntry) bool
		nEntries int
	}{
		{
			name: "no entries",
			pred: func(log.TestLogEntry) bool { return true },
		},
		{
			name: "no matching entries",
			kvss: [][]interface{}{{"key", "value"}},
			pred: func(log.TestLogEntry) bool { return false },
		},
		{
			name: "one matching entry",
			kvss: [][]interface{}{{"key", "value"}},
			pred: func(e log.TestLogEntry) bool {
				v, ok := e["key"]
				return ok && v == "value"
			},
			nEntries: 1,
		},
		{
			name: "multiple matching entries",
			kvss: [][]interface{}{
				{"key1", "value1"},
				{"key2", "value2"},
			},
			pred: func(e log.TestLogEntry) bool {
				for k, v := range e {
					if (k == "key1" || k == "key2") && (v == "value1" || v == "value2") {
						return true
					}
				}
				return false
			},
			nEntries: 2,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			logger := &log.TestLogger{}
			for _, kvs := range tt.kvss {
				if err := logger.Log(kvs...); err != nil {
					t.Fatal(err)
				}
			}
			nEntries := logger.CountMatchingLogEntries(tt.pred)
			assert.Equal(t, tt.nEntries, nEntries)
		})
	}
}
