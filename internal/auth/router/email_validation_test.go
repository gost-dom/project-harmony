package router_test

import (
	"errors"
	"testing"

	"harmony/internal/auth"
	"harmony/internal/auth/domain"
	"harmony/internal/auth/router"
	"harmony/internal/testing/domaintest"
	"harmony/internal/testing/mocks/auth/router_mock"
	"harmony/internal/testing/servertest"

	"github.com/gost-dom/browser/html"
	matchers "github.com/gost-dom/browser/testing/gomega-matchers"
	"github.com/gost-dom/shaman"
	"github.com/gost-dom/shaman/ariarole"
	. "github.com/gost-dom/shaman/predicates"
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
	validatorMock := router_mock.NewMockEmailValidator(s.T())
	validatorMock.EXPECT().
		Validate(mock.Anything, mock.Anything).
		Return(domain.AuthenticatedAccount{}, auth.ErrBadChallengeResponse)

	s.Graph = surgeon.Replace[router.EmailValidator](s.Graph, validatorMock)
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
	validatorMock := router_mock.NewMockEmailValidator(s.T())
	validatorMock.EXPECT().
		Validate(mock.Anything, mock.Anything).
		Return(domaintest.InitAuthenticatedAccount(), nil)

	s.Graph = surgeon.Replace[router.EmailValidator](s.Graph, validatorMock)
	win := s.OpenWindow("https://example.com/auth/validate-email")
	form := NewValidateEmailForm(s.T(), win)
	s.Expect(form.Alert()).To(gomega.BeNil())

	form.Email().Write("j.smith@example.com")
	form.Code().Write("123456")
	form.SubmitButton().Click()

	s.Expect(win.Location().Pathname()).To(gomega.Equal("/host"))
	s.Expect(s.Get(ByH1)).To(matchers.HaveTextContent("Host"))
}

func (s *ValidateEmailTestSuite) TestUnexpectedError() {
	validatorMock := router_mock.NewMockEmailValidator(s.T())
	validatorMock.EXPECT().
		Validate(mock.Anything, mock.Anything).
		Return(domaintest.InitAuthenticatedAccount(), errors.New("Unexpected error"))

	s.Graph = surgeon.Replace[router.EmailValidator](s.Graph, validatorMock)
	win := s.OpenWindow("https://example.com/auth/validate-email")
	form := NewValidateEmailForm(s.T(), win)
	s.Expect(form.Alert()).To(gomega.BeNil())

	form.Email().Write("j.smith@example.com")
	form.Code().Write("123456")
	form.SubmitButton().Click()

	s.Expect(form.Alert()).
		To(matchers.HaveTextContent(gomega.ContainSubstring("Unexpected error")), "Expected alert")
	s.Expect(win.Location().Pathname()).To(gomega.Equal("/auth/validate-email"))

	form = NewValidateEmailForm(s.T(), win)
	s.Expect(form.Email().Value()).To(gomega.Equal("j.smith@example.com"))
	s.Expect(form.Code().Value()).To(gomega.Equal("123456"))
}

/* -------- ValidateEmailForm -------- */

type ValidateEmailForm struct {
	shaman.Scope
}

func NewValidateEmailForm(t testing.TB, win html.Window) ValidateEmailForm {
	scope := shaman.WindowScope(t, win).
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
