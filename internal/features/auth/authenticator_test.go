package auth_test

import (
	"testing"

	. "harmony/internal/features/auth"
	"harmony/internal/features/auth/authdomain/password"
	"harmony/internal/testing/htest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthenticatorTestSuite struct {
	htest.GomegaSuite
	Authenticator
	Account *InsertAccount
}

func (s *AuthenticatorTestSuite) SetupTest() {
	input := CreateValidInput()
	input.Email = "jd@example.com"
	input.Password = password.Parse("valid_password")
	repo := NewAccountRepoStub(s.T())

	assert.NoError(s.T(), Registrator{repo}.Register(s.Context(), input))

	s.Authenticator = Authenticator{repo}
	s.Account = repo.Single()
}

func TestAuthenticator(t *testing.T) {
	suite.Run(t, new(AuthenticatorTestSuite))
}

func (s *AuthenticatorTestSuite) TestAuthenticateUnvalidatedAccount() {
	s.Assert().False(s.Account.Email.Validated, "Guard, test assumes an unvalidated account")

	_, err := s.Authenticate(s.Context(), "jd@example.com", password.Parse("valid_password"))
	s.Assert().Error(err, "Cannot log in until the email address has been validated")
	s.Assert().ErrorIs(err, ErrAccountEmailNotValidated)
}

func (s *AuthenticatorTestSuite) validateAccount() {
	s.Account.ValidateEmail(s.Account.Email.Challenge.Code)
	s.T().Helper()
	s.Assert().True(s.Account.Email.Validated)
}

func (s *AuthenticatorTestSuite) TestAuthenticateWrongPassword() {
	s.validateAccount()

	_, err := s.Authenticate(s.Context(), "jd@example.com", password.Parse("wrong_pw"))
	s.Assert().ErrorIs(err, ErrBadCredentials, "Validating with bad credentials")
}

func (s *AuthenticatorTestSuite) TestAuthenticateCorrectPassword() {
	s.validateAccount()

	actual, err := s.Authenticate(s.Context(), "jd@example.com", password.Parse("valid_password"))
	s.Assert().NoError(err)
	s.Assert().Equal(s.Account.AccountID, actual.ID)
	s.Assert().Equal("jd@example.com", actual.Email.String())
}
