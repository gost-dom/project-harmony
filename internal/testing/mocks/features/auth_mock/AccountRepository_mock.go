// Code generated by mockery v2.53.2. DO NOT EDIT.

package auth_mock

import (
	context "context"
	auth "harmony/internal/features/auth"

	mock "github.com/stretchr/testify/mock"
)

// MockAccountRepository is an autogenerated mock type for the AccountRepository type
type MockAccountRepository struct {
	mock.Mock
}

type MockAccountRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAccountRepository) EXPECT() *MockAccountRepository_Expecter {
	return &MockAccountRepository_Expecter{mock: &_m.Mock}
}

// Insert provides a mock function with given fields: _a0, _a1
func (_m *MockAccountRepository) Insert(_a0 context.Context, _a1 auth.UseCaseResult[auth.Account, auth.AccountID]) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Insert")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, auth.UseCaseResult[auth.Account, auth.AccountID]) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAccountRepository_Insert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Insert'
type MockAccountRepository_Insert_Call struct {
	*mock.Call
}

// Insert is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 auth.UseCaseResult[auth.Account,auth.AccountID]
func (_e *MockAccountRepository_Expecter) Insert(_a0 interface{}, _a1 interface{}) *MockAccountRepository_Insert_Call {
	return &MockAccountRepository_Insert_Call{Call: _e.mock.On("Insert", _a0, _a1)}
}

func (_c *MockAccountRepository_Insert_Call) Run(run func(_a0 context.Context, _a1 auth.UseCaseResult[auth.Account, auth.AccountID])) *MockAccountRepository_Insert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(auth.UseCaseResult[auth.Account, auth.AccountID]))
	})
	return _c
}

func (_c *MockAccountRepository_Insert_Call) Return(_a0 error) *MockAccountRepository_Insert_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAccountRepository_Insert_Call) RunAndReturn(run func(context.Context, auth.UseCaseResult[auth.Account, auth.AccountID]) error) *MockAccountRepository_Insert_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAccountRepository creates a new instance of MockAccountRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAccountRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAccountRepository {
	mock := &MockAccountRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
