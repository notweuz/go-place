package user

import (
	"go-place/internal/middleware"
	"go-place/internal/model/response"

	"github.com/gofiber/fiber/v3"
)

type Handler interface {
	GetByID(ctx fiber.Ctx) error
	GetCurrentUser(ctx fiber.Ctx) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{service: service}
}

func (h *handler) GetByID(ctx fiber.Ctx) error {
	id := fiber.Params[uint64](ctx, "id")
	user, err := h.service.GetByID(id)
	if err != nil {
		return err
	}

	userResponse := response.NewUserPublic(user.ID, user.Username, user.CreatedAt)

	return ctx.Status(fiber.StatusOK).JSON(userResponse)
}

func (h *handler) GetCurrentUser(ctx fiber.Ctx) error {
	id, err := middleware.GetCurrentUserID(ctx)
	if err != nil {
		return err
	}
	user, err := h.service.GetByID(id)
	if err != nil {
		return err
	}

	userResponse := response.NewUserPrivateDetailed(user.ID, user.Username, user.Charges, user.LastChargeAt, user.CreatedAt, user.UpdatedAt)

	return ctx.Status(fiber.StatusOK).JSON(userResponse)
}
