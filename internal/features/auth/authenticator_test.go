package auth_test

import (
	"net/mail"
	"testing"

	. "harmony/internal/features/auth"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authdomain/password"
	"harmony/internal/testing/htest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthenticatorTestSuite struct {
	htest.GomegaSuite
	Authenticator
	Account *authdomain.PasswordAuthentication
}

// MustParseEmail creates a *mail.Address from an email string. The function
// panics if the address is not a valid address. This is intended for use in
// test scenarios where an example email address is hardcoded or generated in a
// way that is assumed to generate valid email addresses.
func MustParseEmail(address string) *mail.Address {
	email, err := mail.ParseAddress(address)
	if err != nil {
		panic(err)
	}
	return email
}

func (s *AuthenticatorTestSuite) SetupTest() {
	input := CreateValidInput()
	input.Email = MustParseEmail("jd@example.com")
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
	s.Assert().ErrorIs(err, authdomain.ErrAccountNotValidated)
}

func (s *AuthenticatorTestSuite) validateAccount() {
	s.Assert().NoError(s.Account.ValidateEmail(s.Account.Email.Challenge.Code))
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
	x, _ := s.Repository.FindByEmail(s.Context(), "jd@example.com")
	s.Assert().True(x.Email.Validated, "Email validated")

	actual, err := s.Authenticate(s.Context(), "jd@example.com", password.Parse("valid_password"))
	s.Assert().NoError(err)
	s.Assert().Equal(s.Account.ID, actual.ID)
	s.Assert().Equal("jd@example.com", actual.Email.String())
}
