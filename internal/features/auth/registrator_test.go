package auth_test

import (
	"net/mail"
	"testing"
	"testing/synctest"
	"time"

	. "harmony/internal/features/auth"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authdomain/password"
	"harmony/internal/testing/htest"
	"harmony/internal/testing/repotest"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
)

func CreateValidInput() RegistratorInput {
	email, _ := mail.ParseAddress("jd@example.com")
	return RegistratorInput{
		Email:       *email,
		Password:    password.Parse("valid_password"),
		Name:        "John Smith",
		DisplayName: "John",
	}

}

type RegisterTestSuite struct {
	htest.GomegaSuite
	Registrator
	repo       *AccountRepositoryStub
	validInput RegistratorInput
}

func (s *RegisterTestSuite) SetupTest() {
	s.repo = NewAccountRepoStub(s.T())

	s.Registrator = Registrator{Repository: s.repo}
	s.validInput = CreateValidInput()
}

func TestRegister(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}

func (s *RegisterTestSuite) TestValidRegistrationInput() {
	s.Register(s.Context(), s.validInput)

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
	s.Register(s.Context(), s.validInput)
	entity := s.repo.Single()

	s.Assert().False(entity.Email.Validated, "Email validated - before validation")

	s.Assert().ErrorIs(entity.ValidateEmail(
		authdomain.NewValidationCode()),
		authdomain.ErrBadEmailChallengeResponse, "Validating wrong code")

	code := repotest.SingleEventOfType[authdomain.EmailValidationRequest](s.repo).Code
	s.Assert().NoError(entity.ValidateEmail(code), "Validating right code")
	s.Assert().True(entity.Email.Validated, "Email validated - after validation")
}

func (s *RegisterTestSuite) TestActivationCodeBeforeExpiry() {
	synctest.Run(func() {
		s.Register(s.Context(), s.validInput)
		entity := s.repo.Single()
		code := repotest.SingleEventOfType[authdomain.EmailValidationRequest](
			s.repo,
		).Code

		time.Sleep(14 * time.Minute)
		synctest.Wait()

		s.Assert().NoError(entity.ValidateEmail(code), "Validation error")
		s.Assert().True(entity.Email.Validated, "Email validated")
	})
}

func (s *RegisterTestSuite) TestActivationCodeExpired() {
	synctest.Run(func() {
		s.Register(s.Context(), s.validInput)
		entity := s.repo.Single()
		validationRequest := repotest.SingleEventOfType[authdomain.EmailValidationRequest](
			s.repo,
		)
		code := validationRequest.Code

		s.Assert().False(entity.Email.Validated, "Email validated - before validation")

		time.Sleep(16 * time.Minute)
		synctest.Wait()

		s.Assert().ErrorIs(entity.ValidateEmail(code), authdomain.ErrBadEmailChallengeResponse)
		s.Assert().False(entity.Email.Validated, "Email validated - after validation")
	})
}
