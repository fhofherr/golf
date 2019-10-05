package log_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/fhofherr/golf/log"
)

// nolint: checknoglobals
var stressTests = []struct {
	nGoRoutines int
	nMessages   int
}{
	{2, 100},
	{5, 100},
	{5, 10000},
	{10, 100},
	{10, 1000},
	{10, 10000},
}

func TestContextualLogger_Log_StressTest(t *testing.T) {
	logger := log.With(&log.TestLogger{}, "test type", "contextual logger stress test")
	runLoggerStressTests(t, func() log.Logger {
		return log.With(logger, "random value", time.Now().Unix())
	})
}

func runLoggerStressTests(t *testing.T, factory func() log.Logger) {
	for _, tt := range stressTests {
		tt := tt
		name := fmt.Sprintf("%d Go routines with %d messages each", tt.nGoRoutines, tt.nMessages)
		t.Run(name, func(t *testing.T) {
			log.StressTestLogger(t, factory, tt.nGoRoutines, tt.nMessages)
		})
	}
}
