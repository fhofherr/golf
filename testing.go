package golf

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockLogger is a mock logger implementation useful for testing.
type MockLogger struct {
	mock.Mock
}

// NewMockLogger creates a new MockLogger and registers it with t.
func NewMockLogger(t *testing.T) *MockLogger {
	m := &MockLogger{}
	m.Test(t)
	return m
}

// Log registers a call to itself.
func (m *MockLogger) Log(kvs ...interface{}) {
	m.Called(kvs...)
}

// MockWither is a MockLogger that implements With.
type MockWither struct {
	MockLogger
}

// NewMockWither creates a new MockWither and registers it with t.
func NewMockWither(t *testing.T) *MockWither {
	m := &MockWither{}
	m.Test(t)
	return m
}

// With registers a call to itself and returns any argument it was mocked for.
func (m *MockWither) With(kvs ...interface{}) interface{} {
	args := m.Called(kvs...)
	return args.Get(0)
}

// MockErrorer is a MockLogger that implements Err.
type MockErrorer struct {
	MockLogger
}

// NewMockErrorer creates a new MockErrorer and registers it with t.
func NewMockErrorer(t *testing.T) *MockErrorer {
	m := &MockErrorer{}
	m.Test(t)
	return m
}

// Err registers a call to itself and returns any argument it was mocked for.
func (m *MockErrorer) Err() error {
	args := m.Called()
	return args.Error(0)
}
