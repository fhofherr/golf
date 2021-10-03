package golf_test

import (
	"fmt"
	"testing"

	"github.com/fhofherr/golf"
	"github.com/stretchr/testify/assert"
)

func TestError_NoErrMethod(t *testing.T) {
	l := golf.NewMockLogger(t)

	l.On(
		"Log", fmt.Sprintf("%s: %T does not implement Err", golf.MsgUnsupported, l),
	).Return()

	err := golf.Error(l)
	assert.NoError(t, err)

	l.AssertExpectations(t)
}

func TestTestError(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "no error",
		},
		{
			name: "error occurred",
			err:  assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			l := golf.NewMockErrorer(t)
			l.On("Err").Return(tt.err)

			err := golf.Error(l)
			assert.Equal(t, tt.err, err)

			l.AssertExpectations(t)
		})
	}
}
