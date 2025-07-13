package auth

import (
	"errors"
	"harmony/internal/domain"
	"harmony/internal/features/auth/authdomain"
)

// Re-export used domain errors

var ErrAccountNotValidated = authdomain.ErrAccountNotValidated

var ErrBadCredentials = errors.New("auth: bad credentials")
var ErrNotFound = domain.ErrNotFound
var ErrBadChallengeResponse = errors.New("auth: bad challenge response")
