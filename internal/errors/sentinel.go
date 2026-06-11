package errors

import (
	"errors"
)

var (
	ErrNoJWTSecretFound  = errors.New("JWT secret is missing in .env")
	ErrNoBCryptSaltFound = errors.New("BCrypt salt is missing in .env")
)
