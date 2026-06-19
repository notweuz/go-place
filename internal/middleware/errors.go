package middleware

import (
	"errors"
	"go-place/internal/errs"

	"go-place/internal/model/response"

	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(ctx fiber.Ctx, err error) error {
	var appError *errs.AppError
	if errors.As(err, &appError) {
		return ctx.Status(appError.StatusCode).JSON(response.ErrorResponse{
			Error: appError.Message,
		})
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
		Error: "internal server error",
	})
}
