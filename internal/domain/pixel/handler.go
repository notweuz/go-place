package pixel

import (
	"errors"
	"go-place/internal/errs"
	"go-place/internal/model/response"

	"github.com/gofiber/fiber/v3"
)

type Handler interface {
	GetByID(ctx fiber.Ctx) error
	GetByCoordinates(ctx fiber.Ctx) error
	GetAll(ctx fiber.Ctx) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{service: service}
}

func (h *handler) GetByID(ctx fiber.Ctx) error {
	id := fiber.Params[uint64](ctx, "id")
	pixel, err := h.service.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			return errs.NotFound(err, "pixel with the given ID doesnt exist")
		}
		return errs.Internal(err)
	}

	pixelFull := response.NewPixelFull(pixel.ID, pixel.X, pixel.Y, pixel.UserID, pixel.Color, pixel.CreatedAt, pixel.UpdatedAt)

	return ctx.Status(fiber.StatusOK).JSON(pixelFull)
}

func (h *handler) GetByCoordinates(ctx fiber.Ctx) error {
	x := fiber.Query[uint64](ctx, "x")
	y := fiber.Query[uint64](ctx, "y")
	pixel, err := h.service.GetByCoordinates(x, y)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			return errs.NotFound(err, "pixel with the given coordinates doesnt exist")
		}
		return errs.Internal(err)
	}

	pixelFull := response.NewPixelFull(pixel.ID, pixel.X, pixel.Y, pixel.UserID, pixel.Color, pixel.CreatedAt, pixel.UpdatedAt)

	return ctx.Status(fiber.StatusOK).JSON(pixelFull)
}

func (h *handler) GetAll(ctx fiber.Ctx) error {
	pixels, err := h.service.GetAll()
	if err != nil {
		return errs.Internal(err)
	}

	pixelsFull := make([]response.PixelFull, len(pixels))
	for i, pixel := range pixels {
		pixelsFull[i] = *response.NewPixelFull(pixel.ID, pixel.X, pixel.Y, pixel.UserID, pixel.Color, pixel.CreatedAt, pixel.UpdatedAt)
	}

	return ctx.Status(fiber.StatusOK).JSON(pixelsFull)
}
