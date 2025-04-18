package auth

import (
	"errors"
	"harmony/internal/domain"
	"harmony/internal/features/auth/authdomain"
)

// Re-export used domain errors

var ErrAccountEmailNotValidated = authdomain.ErrAccountEmailNotValidated

var ErrBadCredentials = errors.New("authenticate: Bad credentials")
var ErrNotFound = domain.ErrNotFound
