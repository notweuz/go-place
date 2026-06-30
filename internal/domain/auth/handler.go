package auth

import (
	"errors"
	"go-place/internal/errs"
	"go-place/internal/model/request"
	"go-place/internal/model/response"
	"go-place/internal/validation"

	"github.com/gofiber/fiber/v3"
)

type Handler interface {
	Register(ctx fiber.Ctx) error
	Login(ctx fiber.Ctx) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{service: service}
}

func (h *handler) Register(ctx fiber.Ctx) error {
	var req request.AuthCredentials
	if err := ctx.Bind().Body(&req); err != nil {
		return errs.BadRequest(err, "invalid request body")
	}
	if err := validation.Validate(&req); err != nil {
		return err
	}
	token, err := h.service.Register(&req)

	if err != nil {
		switch {
		case errors.Is(err, errs.ErrConflict):
			return errs.Conflict(err, "user with same username already exist")
		case errors.Is(err, errs.ErrFailedJWT):
			return errs.FailedJWT(err)
		case errors.Is(err, errs.ErrFailedBCrypt):
			return errs.FailedBCrypt(err)
		}
		return errs.Internal(err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.AuthToken{Token: *token})
}

func (h *handler) Login(ctx fiber.Ctx) error {
	var req request.AuthCredentials
	if err := ctx.Bind().Body(&req); err != nil {
		return errs.BadRequest(err, "invalid request body")
	}
	if err := validation.Validate(&req); err != nil {
		return err
	}

	token, err := h.service.Login(&req)

	if err != nil {
		switch {
		case errors.Is(err, errs.ErrFailedJWT):
			return errs.FailedJWT(err)
		case errors.Is(err, errs.ErrFailedBCrypt):
			return errs.FailedBCrypt(err)
		}
		return errs.Internal(err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.AuthToken{Token: *token})
}
