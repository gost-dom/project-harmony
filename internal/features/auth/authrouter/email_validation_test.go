package authrouter_test

import (
	"testing"

	"harmony/internal/features/auth"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authrouter"
	ariarole "harmony/internal/testing/aria-role"
	"harmony/internal/testing/domaintest"
	. "harmony/internal/testing/mocks/features/auth/authrouter_mock"
	"harmony/internal/testing/servertest"
	"harmony/internal/testing/shaman"
	. "harmony/internal/testing/shaman/predicates"

	"github.com/gost-dom/browser/html"
	matchers "github.com/gost-dom/browser/testing/gomega-matchers"
	"github.com/gost-dom/surgeon"
	"github.com/onsi/gomega"
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
	validatorMock.EXPECT().
		Validate(mock.Anything, mock.Anything).
		Return(authdomain.AuthenticatedAccount{}, auth.ErrBadChallengeResponse)

	s.Graph = surgeon.Replace[authrouter.EmailValidator](s.Graph, validatorMock)
	win := s.OpenWindow("https://example.com/auth/validate-email")
	form := NewValidateEmailForm(s.T(), win)
	s.Expect(form.Alert()).To(gomega.BeNil())

	form.Email().Write("j.smith@example.com")
	form.Code().Write("123456")
	form.SubmitButton().Click()

	s.Expect(form.Alert()).
		To(matchers.HaveTextContent("Wrong email or validation code"), "Expected alert")
	s.Expect(win.Location().Pathname()).To(gomega.Equal("/auth/validate-email"))

	form = NewValidateEmailForm(s.T(), win)
	s.Expect(form.Email().Value()).To(gomega.Equal("j.smith@example.com"))
}

func (s *ValidateEmailTestSuite) TestValidCodeRedirects() {
	validatorMock := NewMockEmailValidator(s.T())
	validatorMock.EXPECT().
		Validate(mock.Anything, mock.Anything).
		Return(domaintest.InitAuthenticatedAccount(), nil)

	s.Graph = surgeon.Replace[authrouter.EmailValidator](s.Graph, validatorMock)
	win := s.OpenWindow("https://example.com/auth/validate-email")
	form := NewValidateEmailForm(s.T(), win)
	s.Expect(form.Alert()).To(gomega.BeNil())

	form.Email().Write("j.smith@example.com")
	form.Code().Write("123456")
	form.SubmitButton().Click()

	s.Expect(win.Location().Pathname()).To(gomega.Equal("/host"))
	shaman.NewScope(s.T(), win.Document().DocumentElement())
	s.Expect(s.Get(ByH1)).To(matchers.HaveTextContent("Host"))
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

func (f ValidateEmailForm) Alert() html.HTMLElement {
	return f.Scope.Find(ByRole(ariarole.Alert))
}
