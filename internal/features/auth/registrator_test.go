package auth_test

import (
	"context"
	"testing"

	"harmony/internal/features/auth"
	. "harmony/internal/features/auth"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authdomain/password"
	"harmony/internal/testing/htest"
	"harmony/internal/testing/repotest"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
)

type RegisterTestSuite struct {
	htest.GomegaSuite
	Registrator
	ctx        context.Context
	repo       *AccountRepositoryStub
	validInput RegistratorInput
}

func (s *RegisterTestSuite) SetupTest() {
	s.repo = NewAccountRepoStub(s.T())

	s.Registrator = Registrator{Repository: s.repo}
	s.ctx = context.Background()
	s.validInput = RegistratorInput{
		Email:       "jd@example.com",
		Password:    password.Parse("valid_password"),
		Name:        "John Smith",
		DisplayName: "John",
	}
}

func TestRegister(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}

func (s *RegisterTestSuite) TestValidRegistrationInput() {
	s.Register(s.ctx, s.validInput)

	entity := s.repo.Single()

	s.Assert().NotZero(entity.ID)
	s.Assert().Equal("jd@example.com", entity.Email.String())
	s.Assert().Equal("John Smith", entity.Name)
	s.Assert().Equal("John", entity.DisplayName)

	s.Expect(s.repo.Events).To(gomega.ContainElement(
		authdomain.AccountRegistered{AccountID: entity.ID}),
	)
}

func (s *RegisterTestSuite) TestActivation() {
	s.Register(s.ctx, s.validInput)

	entity := s.repo.Single()
	validationRequest := repotest.AssertOneEventOfType[authdomain.EmailValidationRequest](s.repo)
	code := validationRequest.Code

	s.Assert().False(entity.Email.Validated, "Email validated - before validation")
	s.Assert().ErrorIs(entity.ValidateEmail(
		authdomain.NewValidationCode()),
		authdomain.ErrBadEmailValidationCode, "Validating wrong code")

	s.Assert().NoError(entity.ValidateEmail(code), "Validating right code")
	s.Assert().True(entity.Email.Validated, "Email validated - after validation")
}

func (s *RegisterTestSuite) TestUnvalidatedAccountLogin() {
	s.Register(s.ctx, s.validInput)
	newAccount := s.repo.Single()

	pw := password.Parse("valid_password")
	a := auth.Authenticator{Repository: s.repo}
	_, err := a.Authenticate(s.ctx, "jd@example.com", pw)
	s.Assert().Error(err)
	s.Assert().ErrorIs(err, ErrAccountEmailNotValidated)

	newAccount.ValidateEmail(newAccount.Email.Challenge.Code)

	_, err = a.Authenticate(s.ctx, "jd@example.com", password.Parse("wrong_pw"))
	s.Assert().ErrorIs(err, ErrBadCredentials)

	actual, err := a.Authenticate(s.ctx, "jd@example.com", pw)
	s.Assert().NoError(err)
	s.Assert().Equal(newAccount.AccountID, actual.ID)
	s.Assert().Equal("jd@example.com", actual.Email.String())
}
