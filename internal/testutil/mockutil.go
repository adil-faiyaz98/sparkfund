package testutil

import (
	"github.com/stretchr/testify/mock"
)

// MockRepository is a base mock repository that can be embedded in other mock repositories
type MockRepository struct {
	mock.Mock
}

// AssertExpectations asserts that all expected calls were made
func (m *MockRepository) AssertExpectations(t interface{}) {
	m.Mock.AssertExpectations(t)
}

// On sets up an expectation for a method call
func (m *MockRepository) On(methodName string, arguments ...interface{}) *mock.Call {
	return m.Mock.On(methodName, arguments...)
}

// Return sets up the return values for a method call
func (m *MockRepository) Return(arguments ...interface{}) *mock.Call {
	return m.Mock.Return(arguments...)
}

// Times sets up how many times a method should be called
func (m *MockRepository) Times(times int) *mock.Call {
	return m.Mock.Times(times)
}

// Once sets up that a method should be called exactly once
func (m *MockRepository) Once() *mock.Call {
	return m.Mock.Once()
}

// Maybe sets up that a method might be called
func (m *MockRepository) Maybe() *mock.Call {
	return m.Mock.Maybe()
}

// ReturnError sets up an error return value for a method call
func (m *MockRepository) ReturnError(err error) *mock.Call {
	return m.Mock.Return(err)
}

// ReturnNil sets up a nil return value for a method call
func (m *MockRepository) ReturnNil() *mock.Call {
	return m.Mock.Return(nil)
}
