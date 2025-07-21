package authrouter_test

import (
	"harmony/internal/auth"
	"harmony/internal/auth/authdomain/password"
	"harmony/internal/auth/authrouter"
	. "harmony/internal/testing/gomegamatchers"
	"harmony/internal/testing/mocks/auth/authrouter_mock"
	"harmony/internal/testing/servertest"
	"testing"

	"github.com/gost-dom/browser/dom"
	"github.com/gost-dom/browser/html"
	. "github.com/gost-dom/browser/testing/gomega-matchers"
	"github.com/gost-dom/shaman"
	"github.com/gost-dom/shaman/ariarole"
	. "github.com/gost-dom/shaman/predicates"
	"github.com/gost-dom/surgeon"
	"github.com/onsi/gomega/gcustom"
	"github.com/onsi/gomega/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RegisterTestSuite struct {
	servertest.BrowserSuite
	registrator *authrouter_mock.MockRegistrator
}

func TestRegister(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RegisterTestSuite))
}

func (s *RegisterTestSuite) SetupTest() {
	s.BrowserSuite.SetupTest()
	s.registrator = authrouter_mock.NewMockRegistrator(s.T())
	s.Graph = surgeon.Replace[authrouter.Registrator](s.Graph, s.registrator)
	s.OpenWindow("https://example.com/auth/register")
}

func (s *RegisterTestSuite) TestSubmitValidForm() {
	s.registrator.EXPECT().Register(mock.Anything, mock.Anything).Return(nil).Once()
	s.Expect(s.Get(ByH1)).To(HaveTextContent("Register Account"))

	form := RegisterForm{s.Subscope(ByRole(ariarole.Form))}
	form.FullName().Write("John Smith")
	form.DisplayName().Write("John")
	form.Email().Write("john.smith@example.com")
	form.Password().Write("str0ngVal!dPassword")
	form.TermsOfUse().Check()
	form.Submit().Click()

	s.Expect(s.Win.Location().Pathname()).To(Equal("/auth/validate-email"),
		"Browser should be redirected to the email validation page on successful registration")
	s.Expect(s.Get(ByH1)).To(HaveTextContent("Validate Email"))
	chalRespForm := EmailChallengeResponseForm{s.Subscope(ByRole(ariarole.Form))}
	s.Expect(chalRespForm.Email()).
		To(HaveAttribute("value", "john.smith@example.com"), "Email is filled on the challenge response page")

	actualInput := s.registrator.Calls[0].Arguments[1].(auth.RegistratorInput)
	s.Expect(actualInput.DisplayName).To(Equal("John"))
	s.Expect(actualInput.Name).To(Equal("John Smith"))
	s.Expect(actualInput.Email.Address).To(Equal("john.smith@example.com"))
	s.Expect(actualInput.Password).To(BeSamePassword("str0ngVal!dPassword"))
	s.Expect(actualInput.NewsletterSignup).To(BeFalse())
}

func (s *RegisterTestSuite) TestCSRF() {
	s.AllowErrorLogs()
	s.registrator.EXPECT().Register(mock.Anything, mock.Anything).Return(nil).Maybe()

	s.CookieJar.Clear()
	s.Expect(s.Get(ByH1)).To(HaveTextContent("Register Account"))

	form := RegisterForm{s.Subscope(ByRole(ariarole.Form))}
	form.FillWithValidValues()
	form.Submit().Click()

	s.Expect(s.Win.Location().Pathname()).To(Equal("/auth/register"))
	s.Assert().Empty(s.registrator.Calls)
}

func (s *RegisterTestSuite) TestSubmitValidFormWithNewsletterSignup() {
	s.registrator.EXPECT().Register(mock.Anything, mock.Anything).Return(nil).Once()
	s.Expect(s.Get(ByH1)).To(HaveTextContent("Register Account"))

	form := RegisterForm{s.Subscope(ByRole(ariarole.Form))}
	form.FillWithValidValues()
	form.NewsletterSignup().Check()
	form.Submit().Click()

	actualInput := s.registrator.Calls[0].Arguments[1].(auth.RegistratorInput)
	s.Expect(actualInput.NewsletterSignup).To(BeTrue())
}

func (s *RegisterTestSuite) TestMissingFullname() {
	s.registrator.EXPECT().Register(mock.Anything, mock.Anything).Return(nil).Maybe()
	s.Expect(s.Get(ByH1)).To(HaveTextContent("Register Account"))

	form := RegisterForm{s.Subscope(ByRole(ariarole.Form))}
	form.FillWithValidValues()
	form.FullName().Clear()
	form.Submit().Click()

	s.Expect(s.Win.Location().Pathname()).To(Equal("/auth/register"),
		"The browser should stay on the registration page when full name is missing")
}

func (s *RegisterTestSuite) TestMissingEmail() {
	s.registrator.EXPECT().Register(mock.Anything, mock.Anything).Return(nil).Maybe()
	s.Expect(s.Get(ByH1)).To(HaveTextContent("Register Account"))

	form := RegisterForm{s.Subscope(ByRole(ariarole.Form))}
	form.FillWithValidValues()
	form.Email().Clear()
	form.Submit().Click()

	s.Expect(s.Win.Location().Pathname()).To(Equal("/auth/register"),
		"The browser should stay on the registration page when email is missing")

	s.Expect(form.Email()).ToNot(HaveARIADescription("Must be a valid email address"))

	s.Expect(form.Email()).To(HaveARIADescription("Must be filled out"))
}

func (s *RegisterTestSuite) TestInvalidEmail() {
	s.registrator.EXPECT().Register(mock.Anything, mock.Anything).Return(nil).Maybe()
	s.Expect(s.Get(ByH1)).To(HaveTextContent("Register Account"))

	form := RegisterForm{s.Subscope(ByRole(ariarole.Form))}
	form.FullName().Write("John Smith")
	form.DisplayName().Write("John")
	form.Email().Write("invalid.email.example.com")
	form.Password().Write("str0ngVal!dPassword")
	form.TermsOfUse().Check()
	form.Submit().Click()

	s.Expect(s.Win.Location().Pathname()).To(Equal("/auth/register"),
		"The browser should stay on the registration page when email is invalid")

	form = RegisterForm{s.Subscope(ByRole(ariarole.Form))}
	s.Expect(form.FullName()).To(HaveAttribute("value", "John Smith"))
	s.Expect(form.DisplayName()).To(HaveAttribute("value", "John"))
	s.Expect(form.Email()).To(HaveAttribute("value", "invalid.email.example.com"))
	s.Expect(form.Password()).To(HaveAttribute("value", ""))
	s.Expect(form.Email()).To(HaveARIADescription("Must be a valid email address"))
}

func (s *RegisterTestSuite) TestMissingAccept() {
	s.registrator.EXPECT().Register(mock.Anything, mock.Anything).Return(nil).Maybe()
	s.Expect(s.Get(ByH1)).To(HaveTextContent("Register Account"))

	form := RegisterForm{s.Subscope(ByRole(ariarole.Form))}
	form.FillWithValidValues()
	form.TermsOfUse().Uncheck()
	form.Submit().Click()

	s.Expect(s.Win.Location().Pathname()).To(Equal("/auth/register"),
		"The browser should stay on the registration page when terms are not accepted")

	form = RegisterForm{s.Subscope(ByRole(ariarole.Form))}
	s.Expect(form.TermsOfUse()).To(HaveARIADescription("You must accept the terms of use"))
}

type RegisterForm struct{ shaman.Scope }

func (f RegisterForm) FullName() shaman.TextboxRole    { return f.Textbox(ByName("Full name")) }
func (f RegisterForm) DisplayName() shaman.TextboxRole { return f.Textbox(ByName("Display name")) }
func (f RegisterForm) Email() shaman.TextboxRole       { return f.Textbox(ByName("Email")) }
func (f RegisterForm) Password() shaman.TextboxRole    { return f.PasswordText(ByName("Password")) }
func (f RegisterForm) TermsOfUse() shaman.CheckboxRole {
	return f.Checkbox(ByName("I agree to the terms of use"))
}
func (f RegisterForm) NewsletterSignup() shaman.CheckboxRole {
	return f.Checkbox(ByName("Sign up for the newsletter"))
}

func (f RegisterForm) Submit() html.HTMLElement { return f.Get(shaman.ByRole(ariarole.Button)) }

func (f RegisterForm) FillWithValidValues() {
	f.FullName().Write("John Smith")
	f.DisplayName().Write("John")
	f.Email().Write("john.smith@example.com")
	f.Password().Write("str0ngVal!dPassword")
	f.TermsOfUse().Check()
}

type EmailChallengeResponseForm struct{ shaman.Scope }

func (f EmailChallengeResponseForm) Email() shaman.TextboxRole { return f.Textbox(ByName("Email")) }

func HaveARIADescription(expected any) types.GomegaMatcher {
	matcher, ok := expected.(types.GomegaMatcher)
	if !ok {
		return HaveARIADescription(Equal(expected))
	}
	var data = struct {
		Matcher     types.GomegaMatcher
		Expected    any
		Description string
	}{Matcher: matcher, Expected: expected}
	return gcustom.MakeMatcher(func(e dom.Element) (bool, error) {
		data.Description = shaman.GetDescription(e)
		return data.Matcher.Match(data.Description)
	}).WithTemplate("Expected:\n{{.FormattedActual}}\n{{.To}} have ARIA Description: {{.Data.Expected}}\n{{.Data.Matcher.FailureMessage .Data.Description}}", &data)
}

func BeSamePassword(pw string) types.GomegaMatcher {
	return gcustom.MakeMatcher(func(actual password.Password) (bool, error) {
		return password.Parse(pw).Equals(actual), nil
	}).WithTemplate("Expected password {{.To}} be the: {{.Data}}", pw)
}
