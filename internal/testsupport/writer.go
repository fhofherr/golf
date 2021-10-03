package testsupport

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockWriter is a mock implementation of io.Writer.
type MockWriter struct {
	mock.Mock
}

// NewMockWriter creates a new MockWriter and registers it with t.
func NewMockWriter(t *testing.T) *MockWriter {
	w := &MockWriter{}
	w.Test(t)
	return w
}

// Write registers a call to itself and returns what it was mocked to return.
func (m *MockWriter) Write(bs []byte) (int, error) {
	args := m.Called(bs)
	return args.Int(0), args.Error(1)
}
