// Code generated by mockery v2.53.3. DO NOT EDIT.

package auth_mock

import (
	context "context"
	domain "harmony/internal/auth/domain"

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

// FindPWAuthByEmail provides a mock function with given fields: ctx, email
func (_m *MockAccountEmailFinder) FindPWAuthByEmail(ctx context.Context, email string) (domain.PasswordAuthentication, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for FindPWAuthByEmail")
	}

	var r0 domain.PasswordAuthentication
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.PasswordAuthentication, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.PasswordAuthentication); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(domain.PasswordAuthentication)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAccountEmailFinder_FindPWAuthByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindPWAuthByEmail'
type MockAccountEmailFinder_FindPWAuthByEmail_Call struct {
	*mock.Call
}

// FindPWAuthByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
func (_e *MockAccountEmailFinder_Expecter) FindPWAuthByEmail(ctx interface{}, email interface{}) *MockAccountEmailFinder_FindPWAuthByEmail_Call {
	return &MockAccountEmailFinder_FindPWAuthByEmail_Call{Call: _e.mock.On("FindPWAuthByEmail", ctx, email)}
}

func (_c *MockAccountEmailFinder_FindPWAuthByEmail_Call) Run(run func(ctx context.Context, email string)) *MockAccountEmailFinder_FindPWAuthByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockAccountEmailFinder_FindPWAuthByEmail_Call) Return(_a0 domain.PasswordAuthentication, _a1 error) *MockAccountEmailFinder_FindPWAuthByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAccountEmailFinder_FindPWAuthByEmail_Call) RunAndReturn(run func(context.Context, string) (domain.PasswordAuthentication, error)) *MockAccountEmailFinder_FindPWAuthByEmail_Call {
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
