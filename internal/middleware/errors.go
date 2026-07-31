package middleware

import (
	"errors"
	"go-place/internal/errs"

	"go-place/internal/model/response"

	"github.com/gofiber/fiber/v3"
)

var errorToStatusCode = map[error]int{
	errs.ErrUserNotFound:        fiber.StatusNotFound,
	errs.ErrPixelNotFound:       fiber.StatusNotFound,
	errs.ErrSettingNotFound:     fiber.StatusNotFound,
	errs.ErrNotFound:            fiber.StatusNotFound,
	errs.ErrConflict:            fiber.StatusConflict,
	errs.ErrNotAuthorized:       fiber.StatusUnauthorized,
	errs.ErrInvalidCredentials:  fiber.StatusUnauthorized,
	errs.ErrNoJWTSecretFound:    fiber.StatusUnauthorized,
	errs.ErrNotEnoughCharges:    fiber.StatusBadRequest,
	errs.ErrUpgradeRequired:     fiber.StatusUpgradeRequired,
	errs.ErrFailedBCrypt:        fiber.StatusInternalServerError,
	errs.ErrFailedJWT:           fiber.StatusInternalServerError,
	errs.ErrInternalServerError: fiber.StatusInternalServerError,
	errs.ErrPixelOutOfBounds:    fiber.StatusBadRequest,
	errs.ErrJWTMalformed:        fiber.StatusBadRequest,
}

func ErrorHandler(ctx fiber.Ctx, err error) error {
	var appError *errs.AppError
	if errors.As(err, &appError) {
		return ctx.Status(appError.StatusCode).JSON(response.ErrorResponse{
			Error: appError.Message,
		})
	}

	statusCode := fiber.StatusInternalServerError
	message := "internal server error"

	for targetErr, code := range errorToStatusCode {
		if errors.Is(err, targetErr) {
			statusCode = code
			if code != fiber.StatusInternalServerError {
				message = err.Error()
			}
			break
		}
	}

	return ctx.Status(statusCode).JSON(response.ErrorResponse{
		Error: message,
	})
}
