package pixel

import (
	"go-place/internal/errs"
	"go-place/internal/middleware"
	"go-place/internal/model/request"
	"go-place/internal/model/response"
	"go-place/internal/validation"

	"github.com/gofiber/fiber/v3"
)

type Handler interface {
	GetByID(ctx fiber.Ctx) error
	GetByCoordinates(ctx fiber.Ctx) error
	GetAll(ctx fiber.Ctx) error
	Change(ctx fiber.Ctx) error
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
		return err
	}

	pixelFull := response.NewPixelFull(pixel.ID, pixel.X, pixel.Y, pixel.UserID, pixel.Color, pixel.CreatedAt, pixel.UpdatedAt)

	return ctx.Status(fiber.StatusOK).JSON(pixelFull)
}

func (h *handler) GetByCoordinates(ctx fiber.Ctx) error {
	x := fiber.Query[uint64](ctx, "x")
	y := fiber.Query[uint64](ctx, "y")
	pixel, err := h.service.GetByCoordinates(x, y)
	if err != nil {
		return err
	}

	pixelFull := response.NewPixelFull(pixel.ID, pixel.X, pixel.Y, pixel.UserID, pixel.Color, pixel.CreatedAt, pixel.UpdatedAt)

	return ctx.Status(fiber.StatusOK).JSON(pixelFull)
}

func (h *handler) GetAll(ctx fiber.Ctx) error {
	pixels, err := h.service.GetAll()
	if err != nil {
		return err
	}

	pixelsFull := make([]response.PixelFull, len(pixels))
	for i, pixel := range pixels {
		pixelsFull[i] = *response.NewPixelFull(pixel.ID, pixel.X, pixel.Y, pixel.UserID, pixel.Color, pixel.CreatedAt, pixel.UpdatedAt)
	}

	return ctx.Status(fiber.StatusOK).JSON(pixelsFull)
}

func (h *handler) Change(ctx fiber.Ctx) error {
	var changePixel request.ChangePixel
	if err := ctx.Bind().Body(&changePixel); err != nil {
		return errs.BadRequest(err, "invalid request body")
	}
	if err := validation.Validate(&changePixel); err != nil {
		return err
	}
	userID, err := middleware.GetCurrentUserID(ctx)
	if err != nil {
		return err
	}

	pixel, err := h.service.Change(changePixel.X, changePixel.Y, changePixel.Color, userID)
	if err != nil {
		return err
	}

	pixelFull := response.NewPixelFull(pixel.ID, pixel.X, pixel.Y, pixel.UserID, pixel.Color, pixel.CreatedAt, pixel.UpdatedAt)

	return ctx.Status(fiber.StatusOK).JSON(pixelFull)
}
