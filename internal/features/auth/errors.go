package auth

import (
	"errors"
	"harmony/internal/features/auth/authdomain"
)

// "Reexport" used domain errors
var ErrAccountEmailNotValidated = authdomain.ErrAccountEmailNotValidated
var ErrBadCredentials = errors.New("authenticate: Bad credentials")
var ErrNotFound = errors.New("Not found")
