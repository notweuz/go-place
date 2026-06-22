package errs

import (
	"errors"
)

var (
	ErrNoJWTSecretFound    = errors.New("JWT secret is missing in .env")
	ErrNoBCryptSaltFound   = errors.New("BCrypt salt is missing in .env")
	ErrNotFoundInDB        = errors.New("not found")
	ErrInternalServerError = errors.New("internal server error")
)
