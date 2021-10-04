package golf_test

import (
	"fmt"
	"testing"

	"github.com/fhofherr/golf"
	"github.com/fhofherr/golf/internal/golferr"
	"github.com/stretchr/testify/assert"
)

func TestWith_NilLogger(t *testing.T) {
	golf.With(nil, "key", "value")
}

func TestWith_NoWithMethod(t *testing.T) {
	logger := golf.NewMockLogger(t)

	logger.On(
		"Log",
		fmt.Sprintf("%s: %T does not implement With", golferr.MsgUnsupported, logger),
	).Return()

	actual := golf.With(logger, "key", "value")

	assert.Same(t, logger, actual)
	logger.AssertExpectations(t)
}

func TestWith_CallWithMethod(t *testing.T) {
	type testCase struct {
		name    string
		prepare func(t *testing.T, tt *testCase)
		assert  func(t *testing.T, tt *testCase)
		args    []interface{}

		// Set during run
		wither *golf.MockWither
		actual golf.Logger
	}

	tests := []testCase{
		{
			name: "With returns invalid value",
			args: []interface{}{"key", "value"},
			prepare: func(t *testing.T, tt *testCase) {
				tt.wither.On("With", tt.args...).Return(struct{}{})
				tt.wither.On(
					"Log", fmt.Sprintf("%s: (%T).With did not return a Logger", golferr.MsgError, tt.wither),
				).Return(struct{}{})
			},
			assert: func(t *testing.T, tt *testCase) {
				assert.Same(t, tt.wither, tt.actual)
			},
		},
		{
			name: "With returns correct value",
			args: []interface{}{"another key", "another value"},
			prepare: func(t *testing.T, tt *testCase) {
				logger := golf.NewMockLogger(t)

				tt.wither.On("With", tt.args...).Return(logger)
				tt.assert = func(t *testing.T, tt *testCase) {
					assert.Same(t, tt.actual, logger)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.wither = golf.NewMockWither(t)
			tt.prepare(t, &tt)

			tt.actual = golf.With(tt.wither, tt.args...)
			tt.assert(t, &tt)

			tt.wither.AssertExpectations(t)
		})
	}
}
