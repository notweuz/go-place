package errs

import (
	"errors"
)

var (
	ErrNoJWTSecretFound    = errors.New("JWT secret is missing in .env")
	ErrNotAuthorized       = errors.New("not authorized")
	ErrNotFoundInDB        = errors.New("not found")
	ErrInternalServerError = errors.New("internal server error")
	ErrDuplicatedKey       = errors.New("duplicated key")
)
