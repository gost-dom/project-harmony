package authrouter_test

import (
	"harmony/internal/features/auth/authrouter"
	ariarole "harmony/internal/testing/aria-role"
	. "harmony/internal/testing/gomegamatchers"
	"harmony/internal/testing/mocks/features/auth/authrouter_mock"
	"harmony/internal/testing/servertest"
	"harmony/internal/testing/shaman"
	. "harmony/internal/testing/shaman/predicates"
	"testing"

	"github.com/gost-dom/browser/dom"
	"github.com/gost-dom/browser/html"
	. "github.com/gost-dom/browser/testing/gomega-matchers"
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
	s.registrator.EXPECT().Register(mock.Anything, mock.Anything).Return(nil).Maybe()
	s.Expect(s.Get(ByH1)).To(HaveTextContent("Register Account"))

	form := RegisterForm{s.Subscope(ByRole(ariarole.Form))}
	form.FullName().Write("John Smith")
	form.DisplayName().Write("John")
	form.Email().Write("john.smith@example.com")
	form.Password().Write("str0ngVal!dPassword")
	form.Submit().Click()

	// Verify that the valid form directs to the email validation page with the
	// email field filled out.
	s.Expect(s.Win.Location().Pathname()).To(Equal("/auth/validate-email"))
	s.Expect(s.Get(ByH1)).To(HaveTextContent("Validate Email"))
	chalRespForm := EmailChallengeResponseForm{s.Subscope(ByRole(ariarole.Form))}
	s.Expect(chalRespForm.Email()).To(HaveAttribute("value", "john.smith@example.com"))
}

type RegisterForm struct{ shaman.Scope }

func (f RegisterForm) FullName() shaman.TextboxRole    { return f.Textbox(ByName("Full name")) }
func (f RegisterForm) DisplayName() shaman.TextboxRole { return f.Textbox(ByName("Display name")) }
func (f RegisterForm) Email() shaman.TextboxRole       { return f.Textbox(ByName("Email")) }
func (f RegisterForm) Password() shaman.TextboxRole    { return f.PasswordText(ByName("Password")) }

func (f RegisterForm) Submit() html.HTMLElement { return f.Get(shaman.ByRole(ariarole.Button)) }

type EmailChallengeResponseForm struct{ shaman.Scope }

func (f EmailChallengeResponseForm) Email() shaman.TextboxRole { return f.Textbox(ByName("Email")) }

func HaveARIADescription(expected string) types.GomegaMatcher {
	var data = struct {
		Matcher     types.GomegaMatcher
		Expected    string
		Description string
	}{Matcher: Equal(expected), Expected: expected}
	return gcustom.MakeMatcher(func(e dom.Element) (bool, error) {
		data.Description = shaman.GetDescription(e)
		return data.Matcher.Match(data.Description)
	}).WithTemplate("Expected:\n{{.FormattedActual}}\n{{.To}} have ARIA Description: {{.Data.Expected}}\n{{.Data.Matcher.FailureMessage .Data.Description}}", &data)
}

// func HaveTag(expected string) GomegaMatcher {
// 	matcher := gomega.Equal(expected)
// 	return gcustom.MakeMatcher(func(e dom.Element) (bool, error) {
// 		return matcher.Match(e.TagName())
// 	}).WithTemplate("Expected:\n{{.FormattedActual}}\n{{.To}} have tag {{.Data.FailureMessage .Actual.TagName}}", matcher)
// }
