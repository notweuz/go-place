package setting

import (
	"errors"
	"go-place/internal/errs"
	"go-place/internal/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Service interface {
	Create(setting *model.Setting) error
	Get() (*model.Setting, error)
	Update(*model.Setting) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Create(setting *model.Setting) error {
	log.Info().Interface("settings", setting).Msg("Attempting to create new settings profile")
	err := s.repository.Create(setting)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create new settings profile")
		return err
	}
	return nil
}

func (s *service) Get() (*model.Setting, error) {
	log.Info().Msg("Attempting to get settings profile")
	setting, err := s.repository.Get()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get settings profile")
		return nil, err
	}
	return setting, nil
}

func (s *service) Update(setting *model.Setting) error {
	log.Info().Interface("settings", setting).Msg("Attempting to update settings profile")
	err := s.repository.Update(setting)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update settings profile")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrNotFound
		}
		return errs.ErrInternalServerError
	}
	return nil
}
