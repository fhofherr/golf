package log_test

import (
	"testing"

	"github.com/fhofherr/golf/log"
	"github.com/stretchr/testify/assert"
)

func TestNOPLogger_Log(t *testing.T) {
	logger := log.NewNOPLogger()
	assert.NoError(t, logger.Log("key", "value"))
}
