package authrouter_test

import (
	ariarole "harmony/internal/testing/aria-role"
	"harmony/internal/testing/servertest"
	"harmony/internal/testing/shaman"
	"harmony/internal/testing/shaman/predicates"
	. "harmony/internal/testing/shaman/predicates"
	"testing"

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
	s.OpenWindow("/auth/register")
}

func (s *RegisterTestSuite) TestSubmitValidForm() {
	s.Expect(s.Scope.Get(ByH1)).To(matchers.HaveTextContent("Register Account"))

	form := RegisterForm{s.Subscope(predicates.ByRole(ariarole.Form))}
	form.FullName().Write("John Smith")
}

type RegisterForm struct{ shaman.Scope }

func (f RegisterForm) FullName() shaman.TextboxRole {
	return f.Textbox(ByName("Full name"))
}
