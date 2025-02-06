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

## Login flow

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

## Login page

These tests verify the login page specifically, including required fields,
validation errors, etc.

These tests are expressed in terms of accessibility properties. This is
recommended as it doesn't resist refactoring the UI, and it promotes an
accessible design.

## Shaman

This package contains helpers for querying the DOM. The shaman helps drive
ghosts away (couldn't call it excorsist, too complicated to type in Go code).

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
