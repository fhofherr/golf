package log_test

import (
	"testing"

	"github.com/fhofherr/golf/log"
)

func TestLog_DoesNothingIfLoggerIsNil(t *testing.T) {
	log.Log(nil)
}

func TestLog_PassesArgumentsToLogger(t *testing.T) {
	logger := &log.TestLogger{}
	log.Log(logger, "key", "value")
	logger.AssertHasMatchingLogEntries(t, 1, func(entry log.TestLogEntry) bool {
		return entry["key"] == "value"
	})
}
