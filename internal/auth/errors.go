package auth

import (
	"errors"
	"harmony/internal/core"
	"harmony/internal/auth/authdomain"
)

// Re-export used core errors

var ErrAccountNotValidated = authdomain.ErrAccountNotValidated

var ErrBadCredentials = errors.New("auth: bad credentials")
var ErrNotFound = core.ErrNotFound
var ErrBadChallengeResponse = errors.New("auth: bad challenge response")
