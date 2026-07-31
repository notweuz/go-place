package setting

import (
	"github.com/gofiber/fiber/v3"
)

type Handler interface {
	Get(ctx fiber.Ctx) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{service: service}
}

func (h *handler) Get(ctx fiber.Ctx) error {
	stg, err := h.service.Get()
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(stg)
}
