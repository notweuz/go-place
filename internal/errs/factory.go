package errs

import "github.com/gofiber/fiber/v3"

func Conflict(err error, message string) *AppError {
	return NewAppError(err, fiber.StatusConflict, message)
}

func FailedJWT(err error) *AppError {
	return NewAppError(err, fiber.StatusInternalServerError, "failed to generate token")
}

func InvalidCredentials(err error) *AppError {
	return NewAppError(err, fiber.StatusUnauthorized, "invalid credentials")
}

func NotFound(err error, message string) *AppError {
	return NewAppError(err, fiber.StatusNotFound, message)
}

func UnauthorizedTokenError(err error, message string) *AppError {
	return NewAppError(err, fiber.StatusUnauthorized, message)
}

func Unauthorized(err error) *AppError {
	return NewAppError(err, fiber.StatusUnauthorized, "unauthorized")
}

func FailedBCrypt(err error) *AppError {
	return NewAppError(err, fiber.StatusInternalServerError, "failed to hash password")
}

func Internal(err error) *AppError {
	return NewAppError(err, fiber.StatusInternalServerError, "internal server error")
}

func BadRequest(err error, message string) *AppError {
	return NewAppError(err, fiber.StatusBadRequest, message)
}
