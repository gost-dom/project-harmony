// Code generated by mockery v2.53.3. DO NOT EDIT.

package auth_mock

import (
	context "context"
	authdomain "harmony/internal/features/auth/authdomain"

	mock "github.com/stretchr/testify/mock"
)

// MockAccountEmailFinder is an autogenerated mock type for the AccountEmailFinder type
type MockAccountEmailFinder struct {
	mock.Mock
}

type MockAccountEmailFinder_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAccountEmailFinder) EXPECT() *MockAccountEmailFinder_Expecter {
	return &MockAccountEmailFinder_Expecter{mock: &_m.Mock}
}

// FindByEmail provides a mock function with given fields: ctx, id
func (_m *MockAccountEmailFinder) FindByEmail(ctx context.Context, id string) (authdomain.PasswordAuthentication, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for FindByEmail")
	}

	var r0 authdomain.PasswordAuthentication
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (authdomain.PasswordAuthentication, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) authdomain.PasswordAuthentication); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(authdomain.PasswordAuthentication)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAccountEmailFinder_FindByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByEmail'
type MockAccountEmailFinder_FindByEmail_Call struct {
	*mock.Call
}

// FindByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *MockAccountEmailFinder_Expecter) FindByEmail(ctx interface{}, id interface{}) *MockAccountEmailFinder_FindByEmail_Call {
	return &MockAccountEmailFinder_FindByEmail_Call{Call: _e.mock.On("FindByEmail", ctx, id)}
}

func (_c *MockAccountEmailFinder_FindByEmail_Call) Run(run func(ctx context.Context, id string)) *MockAccountEmailFinder_FindByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockAccountEmailFinder_FindByEmail_Call) Return(_a0 authdomain.PasswordAuthentication, _a1 error) *MockAccountEmailFinder_FindByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAccountEmailFinder_FindByEmail_Call) RunAndReturn(run func(context.Context, string) (authdomain.PasswordAuthentication, error)) *MockAccountEmailFinder_FindByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAccountEmailFinder creates a new instance of MockAccountEmailFinder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAccountEmailFinder(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAccountEmailFinder {
	mock := &MockAccountEmailFinder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
