package authrouter_test

import (
	"testing"

	"harmony/internal/features/auth"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authrouter"
	ariarole "harmony/internal/testing/aria-role"
	. "harmony/internal/testing/gomegamatchers"
	. "harmony/internal/testing/mocks/features/auth/authrouter_mock"
	"harmony/internal/testing/servertest"
	"harmony/internal/testing/shaman"
	. "harmony/internal/testing/shaman/predicates"

	"github.com/gost-dom/browser/html"
	"github.com/gost-dom/surgeon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ValidateEmailTestSuite struct {
	servertest.BrowserSuite
}

func Test(t *testing.T) {
	suite.Run(t, new(ValidateEmailTestSuite))
}

func (s *ValidateEmailTestSuite) TestEmailAddressIsPrefilledFromQuery() {
	win := s.OpenWindow("https://example.com/auth/validate-email?email=jd@example.com")

	form := NewValidateEmailForm(s.T(), win)
	assert.Equal(s.T(), "jd@example.com", form.Email().Value())
}

func (s *ValidateEmailTestSuite) TestInvalidCodeShowsError() {
	validatorMock := NewMockEmailValidator(s.T())
	validatorMock.EXPECT().Validate(mock.Anything, mock.Anything).Return(nil)

	s.Graph = surgeon.Replace[authrouter.EmailValidator](s.Graph, validatorMock)
	win := s.OpenWindow("https://example.com/auth/validate-email")
	form := NewValidateEmailForm(s.T(), win)

	form.Email().Write("j.smith@example.com")
	form.Code().Write("123456")
	form.SubmitButton().Click()

	actualInput := validatorMock.Calls[0].Arguments[1].(auth.ValidateEmailInput)

	s.Expect(actualInput.Email.Address).To(Equal("j.smith@example.com"))
	s.Expect(actualInput.Code).To(Equal(authdomain.EmailValidationCode("123456")))
}

type ValidateEmailForm struct {
	shaman.Scope
}

func NewValidateEmailForm(t testing.TB, win html.Window) ValidateEmailForm {
	scope := shaman.NewScope(t, win.Document().DocumentElement()).
		Subscope(ByRole(ariarole.Main)).
		Subscope(ByRole(ariarole.Form))

	return ValidateEmailForm{scope}
}

func (f ValidateEmailForm) Email() shaman.TextboxRole {
	return f.Textbox(ByName("Email"))
}

func (f ValidateEmailForm) Code() shaman.TextboxRole { return f.Textbox(ByName("Validation code")) }

func (f ValidateEmailForm) SubmitButton() html.HTMLElement {
	return f.Scope.Get(ByRole(ariarole.Button), ByName("Validate"))
}
