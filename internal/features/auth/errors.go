package auth

import (
	"errors"
	"harmony/internal/domain"
	"harmony/internal/features/auth/authdomain"
)

// Re-export used domain errors

var ErrAccountNotValidated = authdomain.ErrAccountNotValidated

var ErrBadCredentials = errors.New("authenticate: Bad credentials")
var ErrNotFound = domain.ErrNotFound
