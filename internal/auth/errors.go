package auth

import (
	"errors"
	"harmony/internal/auth/authdomain"
	"harmony/internal/core"
)

// ErrAccountNotValidated is re-exported from authdom so callers need a single
// import path
var ErrAccountNotValidated = authdomain.ErrAccountNotValidated

// ErrNotFound is re-exported from core so callers need a single import path
var ErrNotFound = core.ErrNotFound

// ErrBadCredentials indicates that the user has supplied the wrong username or
// password.
var ErrBadCredentials = errors.New("auth: bad credentials")

// ErrBadChallengeResponse indicates that the email validation challenge failed
// with a bad code.
var ErrBadChallengeResponse = errors.New("auth: bad challenge response")
