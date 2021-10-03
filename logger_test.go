package golf_test

import (
	"testing"

	"github.com/fhofherr/golf"
)

func TestLog_DoesNothingIfLoggerIsNil(t *testing.T) {
	golf.Log(nil)
}

func TestLog_PassesArgumentsToLogger(t *testing.T) {
	logger := golf.NewMockLogger(t)
	logger.On("Log", "key", "value").Return()

	golf.Log(logger, "key", "value")

	logger.AssertExpectations(t)
}
