package middleware

import (
	"errors"
	"go-place/internal/errs"

	"go-place/internal/model/response"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
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
	errs.ErrRegistrationClosed:  fiber.StatusForbidden,
}

func ErrorHandler(ctx fiber.Ctx, err error) error {
	var appError *errs.AppError
	statusCode := fiber.StatusInternalServerError
	message := "internal server error"

	if errors.As(err, &appError) {
		statusCode = appError.StatusCode
		message = appError.Message
	} else {
		for targetErr, code := range errorToStatusCode {
			if errors.Is(err, targetErr) {
				statusCode = code
				if code != fiber.StatusInternalServerError {
					message = err.Error()
				}
				break
			}
		}
	}

	log.Error().Err(err).Int("status", statusCode).Msgf("HTTP %s %s failed", ctx.Method(), ctx.Path())

	return ctx.Status(statusCode).JSON(response.ErrorResponse{
		Error: message,
	})
}
