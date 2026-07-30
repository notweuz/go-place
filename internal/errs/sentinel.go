package errs

import (
	"errors"
)

var (
	ErrNoJWTSecretFound    = errors.New("JWT secret is missing in .env")
	ErrNotAuthorized       = errors.New("not authorized")
	ErrNotFound            = errors.New("not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrPixelNotFound       = errors.New("pixel not found")
	ErrSettingNotFound     = errors.New("setting not found")
	ErrInternalServerError = errors.New("internal server error")
	ErrConflict            = errors.New("duplicated key")
	ErrFailedBCrypt        = errors.New("failed bcrypt hash")
	ErrFailedJWT           = errors.New("jwt creation failed")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUpgradeRequired     = errors.New("upgrade required")
	ErrNotEnoughCharges    = errors.New("not enough charges")
)
