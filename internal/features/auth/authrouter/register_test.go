package authrouter_test

import (
	ariarole "harmony/internal/testing/aria-role"
	"harmony/internal/testing/servertest"
	"harmony/internal/testing/shaman"
	"harmony/internal/testing/shaman/predicates"
	. "harmony/internal/testing/shaman/predicates"
	"testing"

	"github.com/gost-dom/browser/html"
	matchers "github.com/gost-dom/browser/testing/gomega-matchers"
	"github.com/stretchr/testify/suite"
)

type RegisterTestSuite struct {
	servertest.BrowserSuite
}

func TestRegister(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RegisterTestSuite))
}

func (s *RegisterTestSuite) SetupTest() {
	s.BrowserSuite.SetupTest()
	s.OpenWindow("https://example.com/auth/register")
}

func (s *RegisterTestSuite) TestSubmitValidForm() {
	s.Expect(s.Scope.Get(ByH1)).To(matchers.HaveTextContent("Register Account"))

	form := RegisterForm{s.Subscope(predicates.ByRole(ariarole.Form))}
	form.FullName().Write("John Smith")
	form.Email().Write("john.smith@example.com")
	form.Submit().Click()

	s.Expect(s.Scope.Get(ByH1)).To(matchers.HaveTextContent("Validate Email"))
}

type RegisterForm struct{ shaman.Scope }

func (f RegisterForm) FullName() shaman.TextboxRole {
	return f.Textbox(ByName("Full name"))
}

func (f RegisterForm) Email() shaman.TextboxRole {
	return f.Textbox(ByName("Email"))
}

func (f RegisterForm) Submit() html.HTMLElement {
	return f.Get(shaman.ByRole(ariarole.Button))
}
