package pixel

import (
	"errors"
	"go-place/internal/errs"
	"go-place/internal/model"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type Service interface {
	Create(pixel *model.Pixel) (*model.Pixel, error)
	GetByID(id uint64) (*model.Pixel, error)
	GetByCoordinates(x, y uint64) (*model.Pixel, error)
	Update(pixel *model.Pixel) (*model.Pixel, error)
	Delete(id uint64) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Create(pixel *model.Pixel) (*model.Pixel, error) {
	entity, err := s.repository.Create(pixel)
	if err != nil {
		return nil, errs.NewAppError(errs.ErrInternalServerError, fiber.StatusInternalServerError, "Couldn't create pixel")
	}
	return entity, nil
}

func (s *service) GetByID(id uint64) (*model.Pixel, error) {
	entity, err := s.repository.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewAppError(errs.ErrNotFoundInDB, fiber.StatusNotFound, "Pixel not found")
		}
		return nil, errs.NewAppError(errs.ErrInternalServerError, fiber.StatusInternalServerError, "Couldn't get pixel")
	}
	return entity, nil
}

func (s *service) GetByCoordinates(x, y uint64) (*model.Pixel, error) {
	entity, err := s.repository.GetByCoordinates(x, y)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewAppError(errs.ErrNotFoundInDB, fiber.StatusNotFound, "Pixel not found")
		}
		return nil, errs.NewAppError(errs.ErrInternalServerError, fiber.StatusInternalServerError, "Couldn't get pixel")
	}
	return entity, nil
}

func (s *service) Update(pixel *model.Pixel) (*model.Pixel, error) {
	entity, err := s.repository.Update(pixel)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewAppError(errs.ErrNotFoundInDB, fiber.StatusNotFound, "Pixel not found")
		}
		return nil, errs.NewAppError(errs.ErrInternalServerError, fiber.StatusInternalServerError, "Couldn't get pixel")
	}
	return entity, nil
}

func (s *service) Delete(id uint64) error {
	err := s.repository.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewAppError(errs.ErrNotFoundInDB, fiber.StatusNotFound, "Pixel not found")
		}
		return errs.NewAppError(errs.ErrInternalServerError, fiber.StatusInternalServerError, "Couldn't delete pixel")
	}
	return nil
}
