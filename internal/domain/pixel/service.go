package pixel

import (
	"errors"
	"go-place/internal/errs"
	"go-place/internal/model"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
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
	log.Info().Uint64("pixel_id", pixel.ID).Uint64("x", pixel.X).Uint64("y", pixel.Y).Msg("Attempting to create pixel")
	entity, err := s.repository.Create(pixel)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create pixel")
		return nil, errs.NewAppError(errs.ErrInternalServerError, fiber.StatusInternalServerError, "Couldn't create pixel")
	}
	return entity, nil
}

func (s *service) GetByID(id uint64) (*model.Pixel, error) {
	log.Info().Uint64("pixel_id", id).Msg("Attempting to get pixel by id")
	entity, err := s.repository.GetByID(id)
	if err != nil {
		log.Error().Err(err).Uint64("pixel_id", id).Msg("Failed to get pixel by id")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewAppError(errs.ErrNotFoundInDB, fiber.StatusNotFound, "Pixel not found")
		}
		return nil, errs.NewAppError(errs.ErrInternalServerError, fiber.StatusInternalServerError, "Couldn't get pixel")
	}
	return entity, nil
}

func (s *service) GetByCoordinates(x, y uint64) (*model.Pixel, error) {
	log.Info().Uint64("x", x).Uint64("y", y).Msg("Attempting to get pixel by coordinates")
	entity, err := s.repository.GetByCoordinates(x, y)
	if err != nil {
		log.Error().Err(err).Uint64("x", x).Uint64("y", y).Msg("Failed to get pixel by coordinates")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewAppError(errs.ErrNotFoundInDB, fiber.StatusNotFound, "Pixel not found")
		}
		return nil, errs.NewAppError(errs.ErrInternalServerError, fiber.StatusInternalServerError, "Couldn't get pixel")
	}
	return entity, nil
}

func (s *service) Update(pixel *model.Pixel) (*model.Pixel, error) {
	log.Info().Uint64("pixel_id", pixel.ID).Msg("Attempting to update pixel")
	entity, err := s.repository.Update(pixel)
	if err != nil {
		log.Error().Err(err).Uint64("pixel_id", pixel.ID).Msg("Failed to update pixel")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewAppError(errs.ErrNotFoundInDB, fiber.StatusNotFound, "Pixel not found")
		}
		return nil, errs.NewAppError(errs.ErrInternalServerError, fiber.StatusInternalServerError, "Couldn't get pixel")
	}
	return entity, nil
}

func (s *service) Delete(id uint64) error {
	log.Info().Uint64("pixel_id", id).Msg("Attempting to delete pixel by id")
	err := s.repository.Delete(id)
	if err != nil {
		log.Error().Err(err).Uint64("pixel_id", id).Msg("Failed to delete pixel by id")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewAppError(errs.ErrNotFoundInDB, fiber.StatusNotFound, "Pixel not found")
		}
		return errs.NewAppError(errs.ErrInternalServerError, fiber.StatusInternalServerError, "Couldn't delete pixel")
	}
	return nil
}
