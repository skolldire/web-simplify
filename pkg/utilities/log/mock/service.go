// Code generated by mockery v2.52.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// Debug provides a mock function with given fields: ctx, fields
func (_m *Service) Debug(ctx context.Context, fields map[string]interface{}) {
	_m.Called(ctx, fields)
}

// Error provides a mock function with given fields: ctx, err, msg, fields
func (_m *Service) Error(ctx context.Context, err error, msg string, fields map[string]interface{}) {
	_m.Called(ctx, err, msg, fields)
}

// FatalError provides a mock function with given fields: ctx, err, fields
func (_m *Service) FatalError(ctx context.Context, err error, fields map[string]interface{}) {
	_m.Called(ctx, err, fields)
}

// Info provides a mock function with given fields: ctx, msg, fields
func (_m *Service) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	_m.Called(ctx, msg, fields)
}

// Warn provides a mock function with given fields: ctx, fields
func (_m *Service) Warn(ctx context.Context, fields map[string]interface{}) {
	_m.Called(ctx, fields)
}

// WrapError provides a mock function with given fields: err, msg
func (_m *Service) WrapError(err error, msg string) error {
	ret := _m.Called(err, msg)

	if len(ret) == 0 {
		panic("no return value specified for WrapError")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(error, string) error); ok {
		r0 = rf(err, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
