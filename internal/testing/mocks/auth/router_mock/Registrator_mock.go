// Code generated by mockery v2.53.3. DO NOT EDIT.

package router_mock

import (
	context "context"
	auth "harmony/internal/auth"

	mock "github.com/stretchr/testify/mock"
)

// MockRegistrator is an autogenerated mock type for the Registrator type
type MockRegistrator struct {
	mock.Mock
}

type MockRegistrator_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRegistrator) EXPECT() *MockRegistrator_Expecter {
	return &MockRegistrator_Expecter{mock: &_m.Mock}
}

// Register provides a mock function with given fields: ctx, input
func (_m *MockRegistrator) Register(ctx context.Context, input auth.RegistratorInput) error {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for Register")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, auth.RegistratorInput) error); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRegistrator_Register_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Register'
type MockRegistrator_Register_Call struct {
	*mock.Call
}

// Register is a helper method to define mock.On call
//   - ctx context.Context
//   - input auth.RegistratorInput
func (_e *MockRegistrator_Expecter) Register(ctx interface{}, input interface{}) *MockRegistrator_Register_Call {
	return &MockRegistrator_Register_Call{Call: _e.mock.On("Register", ctx, input)}
}

func (_c *MockRegistrator_Register_Call) Run(run func(ctx context.Context, input auth.RegistratorInput)) *MockRegistrator_Register_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(auth.RegistratorInput))
	})
	return _c
}

func (_c *MockRegistrator_Register_Call) Return(_a0 error) *MockRegistrator_Register_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRegistrator_Register_Call) RunAndReturn(run func(context.Context, auth.RegistratorInput) error) *MockRegistrator_Register_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockRegistrator creates a new instance of MockRegistrator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRegistrator(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRegistrator {
	mock := &MockRegistrator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
