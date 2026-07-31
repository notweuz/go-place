package pixel

import (
	"errors"
	"go-place/internal/database"
	"go-place/internal/domain/setting"
	"go-place/internal/domain/user"
	"go-place/internal/errs"
	"go-place/internal/model"
	"go-place/internal/model/request"
	"go-place/internal/transport/websocket"
	"go-place/internal/transport/websocket/message"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Service interface {
	Create(pixel *model.Pixel) error
	GetByID(id uint64) (*model.Pixel, error)
	GetByCoordinates(x, y uint64) (*model.Pixel, error)
	GetAll(opts ...database.Option) ([]model.Pixel, error)
	Change(pixels []request.ChangePixel, newAuthor uint64) ([]model.Pixel, error)
	Update(pixel *model.Pixel) error
	Delete(id uint64) error
}

type service struct {
	repository     Repository
	userService    user.Service
	settingService setting.Service
	wsHub          websocket.Hub
}

func NewService(repository Repository, userService user.Service, settingService setting.Service, wsHub websocket.Hub) Service {
	return &service{repository: repository, userService: userService, settingService: settingService, wsHub: wsHub}
}

func (s *service) Create(pixel *model.Pixel) error {
	log.Info().Uint64("pixel_id", pixel.ID).Uint64("x", pixel.X).Uint64("y", pixel.Y).Msg("Attempting to create pixel")
	err := s.repository.Create(pixel)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create pixel")
		return errs.ErrInternalServerError
	}
	s.wsHub.Broadcast(*message.NewMessage(message.EventPixelChanged, *message.NewPixelChanged(pixel.ID, pixel.X, pixel.Y, pixel.UserID, pixel.Color)))
	return nil
}

func (s *service) GetByID(id uint64) (*model.Pixel, error) {
	log.Info().Uint64("pixel_id", id).Msg("Attempting to get pixel by id")
	entity, err := s.repository.GetByID(id)
	if err != nil {
		log.Error().Err(err).Uint64("pixel_id", id).Msg("Failed to get pixel by id")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrPixelNotFound
		}
		return nil, errs.ErrInternalServerError
	}
	return entity, nil
}

func (s *service) GetByCoordinates(x, y uint64) (*model.Pixel, error) {
	log.Info().Uint64("x", x).Uint64("y", y).Msg("Attempting to get pixel by coordinates")
	entity, err := s.repository.GetByCoordinates(x, y)
	if err != nil {
		log.Error().Err(err).Uint64("x", x).Uint64("y", y).Msg("Failed to get pixel by coordinates")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrPixelNotFound
		}
		return nil, errs.ErrInternalServerError
	}
	return entity, nil
}

func (s *service) GetAll(opts ...database.Option) ([]model.Pixel, error) {
	log.Info().Msg("Attempting to get all pixels")
	pixels, err := s.repository.GetAll(opts...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all pixels")
		return nil, errs.ErrInternalServerError
	}
	return pixels, nil
}

func (s *service) Change(pixels []request.ChangePixel, newAuthor uint64) ([]model.Pixel, error) {
	log.Info().Int("count", len(pixels)).Uint64("author", newAuthor).Msg("Attempting to change pixels")

	settings, err := s.settingService.Get()
	if err != nil {
		return nil, err
	}

	for _, p := range pixels {
		if p.X >= settings.CanvasWidth || p.Y >= settings.CanvasHeight {
			return nil, errs.ErrPixelOutOfBounds
		}
	}

	err = s.userService.SpendCharge(newAuthor, settings.MaxCharges, uint(len(pixels)), time.Duration(settings.CooldownSeconds)*time.Second)
	if err != nil {
		return nil, err
	}

	models := make([]model.Pixel, len(pixels))
	for i, p := range pixels {
		models[i] = model.Pixel{X: p.X, Y: p.Y, Color: p.Color, UserID: newAuthor}
	}

	result, err := s.repository.Upsert(models)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upsert pixels")
		return nil, errs.ErrInternalServerError
	}

	changed := make([]message.PixelChanged, len(result))
	for i, px := range result {
		changed[i] = *message.NewPixelChanged(px.ID, px.X, px.Y, px.UserID, px.Color)
	}
	s.wsHub.Broadcast(*message.NewMessage(message.EventPixelsChanged, message.NewPixelsChanged(changed)))

	return result, nil
}

func (s *service) Update(pixel *model.Pixel) error {
	log.Info().Uint64("pixel_id", pixel.ID).Msg("Attempting to update pixel")
	err := s.repository.Update(pixel)
	if err != nil {
		log.Error().Err(err).Uint64("pixel_id", pixel.ID).Msg("Failed to update pixel")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrPixelNotFound
		}
		return errs.ErrInternalServerError
	}
	s.wsHub.Broadcast(*message.NewMessage(message.EventPixelChanged, *message.NewPixelChanged(pixel.ID, pixel.X, pixel.Y, pixel.UserID, pixel.Color)))
	return nil
}

func (s *service) Delete(id uint64) error {
	log.Info().Uint64("pixel_id", id).Msg("Attempting to delete pixel by id")
	err := s.repository.Delete(id)
	if err != nil {
		log.Error().Err(err).Uint64("pixel_id", id).Msg("Failed to delete pixel by id")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrPixelNotFound
		}
		return errs.ErrInternalServerError
	}
	return nil
}
