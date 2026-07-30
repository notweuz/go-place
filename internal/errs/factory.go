package errs

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

func newBaseError(err error, statusCode int, baseMessage string, message ...string) *AppError {
	if len(message) > 0 {
		return NewAppError(err, statusCode, strings.Join(message, " ")) // message is supposed to be only one string, but just in case...
	}
	return NewAppError(err, statusCode, baseMessage)
}

func Conflict(err error, message ...string) *AppError {
	return newBaseError(err, fiber.StatusConflict, "same data already exists", message...)
}

func NotFound(err error, message ...string) *AppError {
	return newBaseError(err, fiber.StatusNotFound, "the requested data was not found", message...)
}

func UpdateRequired(err error, message ...string) *AppError {
	return newBaseError(err, fiber.StatusUpgradeRequired, "update required", message...)
}

func Unauthorized(err error, message ...string) *AppError {
	return newBaseError(err, fiber.StatusUnauthorized, "unauthorized", message...)
}

func Internal(err error, message ...string) *AppError {
	return newBaseError(err, fiber.StatusInternalServerError, "internal server error", message...)
}

func BadRequest(err error, message ...string) *AppError {
	return newBaseError(err, fiber.StatusBadRequest, "invalid request data", message...)
}
