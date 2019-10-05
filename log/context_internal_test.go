package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextualLogger_With_NoAdditionalKVS(t *testing.T) {
	expected := contextualLogger{ctxkvs: []interface{}{"key", "value"}}
	actual := expected.With()
	assert.Equal(t, expected, actual)
}

func TestContextualLogger_With_AdditionalKVS(t *testing.T) {
	firstLogger := contextualLogger{ctxkvs: []interface{}{"key1", "value1"}}
	secondLogger := firstLogger.With([]interface{}{"key2", "value2"})
	assert.NotEqual(t, firstLogger, secondLogger)
}

func TestContextualLogger_Log(t *testing.T) {
	tests := []struct {
		name   string
		logger Logger
	}{
		{
			name:   "nil logger",
			logger: contextualLogger{ctxkvs: []interface{}{"key1", "value1"}},
		},
		{
			name: "non-nil logger",
			logger: contextualLogger{
				ctxkvs: []interface{}{"key1", "value1"},
				logger: &TestLogger{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.logger.Log("key2", "value2")
			assert.NoError(t, err)
		})
	}
}
