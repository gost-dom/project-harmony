package auth_test

import (
	"context"
	"reflect"
	"testing"

	. "harmony/internal/features/auth"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authdomain/password"
	"harmony/internal/testing/htest"
	"harmony/internal/testing/mocks/features/auth_mock"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RegisterTestSuite struct {
	htest.GomegaSuite
	ctx context.Context
	Registrator
	repoMock *auth_mock.MockAccountRepository
}

func (s *RegisterTestSuite) SetupTest() {
	s.repoMock = auth_mock.NewMockAccountRepository(s.T())
	s.repoMock.EXPECT().Insert(mock.Anything, mock.Anything).Return(nil)

	s.Registrator = Registrator{Repository: s.repoMock}
	s.ctx = context.Background()
}

func TestRegister(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}

func (s *RegisterTestSuite) TestValidRegistrationInput() {
	pw := password.Parse("s3cre7")
	s.Register(s.ctx, RegistratorInput{
		Email:       "jd@example.com",
		Password:    pw,
		Name:        "John Smith",
		DisplayName: "John",
	})

	res := s.repoMock.Calls[0].Arguments.Get(1).(AccountUseCaseResult)
	entity := res.Entity
	events := res.Events

	s.Assert().NotZero(entity.ID)
	s.Assert().Equal("jd@example.com", entity.Email.String())
	s.Assert().Equal("John Smith", entity.Name)
	s.Assert().Equal("John", entity.DisplayName)

	s.Expect(events).To(gomega.ContainElement(
		authdomain.AccountRegistered{AccountID: entity.ID}),
	)
}

func AssertOneElementOfType[T any](t testing.TB, e []DomainEvent) (res T) {
	t.Helper()
	var found bool
	for _, ee := range e {
		if r, ok := ee.(T); ok {
			if found {
				t.Errorf("Found multiple instances of type %s", reflect.TypeFor[T]().Name())
			}
			res = r
			found = true
		}
	}
	if !found {
		t.Errorf("Found no instance of type %s", reflect.TypeFor[T]().Name())
	}
	return
}

func (s *RegisterTestSuite) TestActivation() {
	pw := password.Parse("s3cre7")
	s.Register(s.ctx, RegistratorInput{
		Email:       "jd@example.com",
		Password:    pw,
		Name:        "John Smith",
		DisplayName: "John",
	})

	res := s.repoMock.Calls[0].Arguments.Get(1).(AccountUseCaseResult)
	entity := res.Entity
	events := res.Events
	validationRequest := AssertOneElementOfType[authdomain.EmailValidationRequest](s.T(), events)
	code := validationRequest.Code

	s.Assert().False(entity.Email.Validated, "Email validated - before validation")
	s.Assert().ErrorIs(entity.ValidateEmail(
		authdomain.NewValidationCode()),
		authdomain.ErrBadEmailValidationCode, "Validating wrong code")

	s.Assert().NoError(entity.ValidateEmail(code), "Validating right code")
	s.Assert().True(entity.Email.Validated, "Email validated - after validation")
}
