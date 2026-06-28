package errs

import (
	"errors"
)

var (
	ErrNoJWTSecretFound    = errors.New("JWT secret is missing in .env")
	ErrNotAuthorized       = errors.New("not authorized")
	ErrNotFound            = errors.New("not found")
	ErrInternalServerError = errors.New("internal server error")
	ErrConflict            = errors.New("duplicated key")
	ErrFailedBCrypt        = errors.New("failed bcrypt hash")
	ErrFailedJWT           = errors.New("jwt creation failed")
	ErrInvalidCredentials  = errors.New("invalid credentials")
)
