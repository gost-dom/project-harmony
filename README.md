# Harmony - Gost example app

This web application serves an an example of usage patters for Gost. The
application uses
  - HTMX - which Go is specifically written to support
  - Templ - seems to be the best template engine atm
  - Tailwind - to quickly make stuff not like like c*** (need to update to v4)

> [!NOTE]
>
> This is early version, with very hacked backend, e.g., a global boolean if you
> are logged in. Focus right now is on expressive test cases.

## The tests

Note: Code will likely be restructured, so if the files don't exist; they have
been moved.

### Login flow

This is in `server_test.go`

This verifies the overall login flow, that requesting a private resource goes
through a login page, but the user eventually ends up on the requested page
after successful login

- User opens the index page, i.e., `/`.
- User clicks "hosts" that links `/hosts`, which requires an authenticated user.
- The system redirects the user to the login page
  - Verify: There is a new history entry.
  - Verify: The current location is `/auth/login`.
- The user types in username and password, and clicks the submit button
  - Verify: A new history entry
  - Verify: The current location is `/hosts`.

The tests do not care about _how_ authentication works. The browser has a cookie
jar, and any session cookies are to be considered an implementation detail.

> [!Note]
>
> This purpose of this test is to verify the behaviour, seen from the user's
> perspective.
>
> That doesn't rule out that you have other tests more tightly coupled to
> implementation details to verify security properties of the cookies, e.g.,
> `http-only`, `secure`, no leaking of information, etc.

### Login page

This is in `login_test.go`

These tests verify the login page specifically, including required fields,
validation errors, etc.

These tests are expressed in terms of accessibility properties. This is
recommended as it doesn't resist refactoring the UI, and it promotes an
accessible design.

## Shaman

This package contains helpers for querying the DOM. The shaman helps drive
ghosts away (couldn't call it excorsist, too complicated to type in Go code).

This would eventually be extracted to a separate project; this just serves as a
place to experiment with patterns.

All of these are incomplete, and just works on the element types that are
actually used in this codebase

### EventSync

Helps synchronize to events, e.g., don't start clicking elements until HTMX has
installed event handlers

### QueryHelper

Helps func multiple, or single elements mathching a set of predicates.
Predicates are specified in the predicates packages. This was it's own package,
primarily as you might want to dot-import these for readability; so not pollute
global scope in one file.

### aria-role

Contains definitions of aria roles, and a function to determine the role of a
specific element.

### Various helpers

`GetName`, `GetDescription`. Gets the accessibility name, as well as the
description for various elements. Can be used for querying, or veriying.

The examples use `GetName` for querying elements, and `GetDescription` to verify
error messages are attached to input fields.

## Running

There's a makefile with a "live" target as the default target, which starts the
server with live-reload capabilities on port 7331.

## Testing frameworks

The structure use [testify](https://github.com/stretchr/testify) suites.

For assertions, the suite uses mostly testify, but a few use gomega.
[gomega](https://onsi.github.io/gomega). Gomega allows creation of custom
matches, that can significantly increase the expressiveness of tests, hiding
irrelevant details. 

E.g., checking the the value of an attribute of a variable of type `Node`
requires 3 check, one for type cast, one of the presence of the attribute, and
one for the value of the attribute.

Use whatever you want; this is meant as a source of test patterns, and as such
try to cover multiple ways to achieve the result.

## Dependency injection

I looked through a lot of IoC containers to automate dependency injection,
focusing on two properties:

    1. Easy dependency replacement in a larger hierarchy
    2. Simple configuration with sensible defaults

I didn't find one satisfying both cases, but 1 is more important than 2, so I
opted for [samber/do](https://pkg.go.dev/github.com/samber/do), as this supports
cloning and replacement in the dependency tree.

I think I will build my own that supports both premises.

### Example

To provide an example of what Do is used for, here is an example of verifying
a failed login attempt:

```go
func (s *LoginPageSuite) SetupTest() {
	s.BrowserSuite.SetupTest()
	s.authMock = mocks.NewAuthenticator(s.T())
	do.OverrideValue[server.Authenticator](s.injector, s.authMock)
	s.OpenWindow("/auth/login")
	s.WaitFor("htmx:load")
	s.loginForm = NewLoginForm(s.Scope)
}

func (s *LoginPageSuite) TestInvalidCredentials() {
	s.authMock.EXPECT().
		Authenticate(mock.Anything, "bad-user@example.com", "s3cret").
		Return(server.Account{}, server.ErrBadCredentials).Once()
	s.loginForm.Email().SetAttribute("value", "bad-user@example.com")
	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()
	s.WaitFor("htmx:afterSettle")

    // Verify the document is still on the login page.
	s.Equal("/auth/login", s.win.Location().Href())

    // Verify that an alert was shown.
    // The expectation is written in a high level of abstraction, expressing how
    // the user interacts with the page. This promotes an accessible design.
	alert := s.Get(ByRole(ariarole.Alert))
	s.Assert().Equal("Email or password did not match", alert.TextContent())
}
```

The general setup replaces the `Authenticator` component with a mocked instance
supplied by the test.

Each test case sets up the specific expectations on the mock, and programmed
result; in this case an `ErrBadCredentials` error result.

## Test data

- John Doe 
  - email: jd@example.com
  - pw: 1234
