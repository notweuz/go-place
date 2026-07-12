package user

import (
	"errors"
	"go-place/internal/database"
	"go-place/internal/errs"
	"go-place/internal/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Service interface {
	Create(user *model.User) error
	GetByID(id uint64, opts ...database.Option) (*model.User, error)
	GetByUsername(username string, opts ...database.Option) (*model.User, error)
	GetAll(opts ...database.Option) []model.User
	Update(user *model.User) error
	Delete(id uint64) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Create(user *model.User) error {
	log.Info().Uint64("id", user.ID).Str("username", user.Username).Msg("Creating user")
	err := s.repository.Create(user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrConflict
		}
		return errs.ErrInternalServerError
	}
	return nil
}

func (s *service) GetByID(id uint64, opts ...database.Option) (*model.User, error) {
	log.Info().Uint64("id", id).Msg("Getting user by id")
	user, err := s.repository.GetByID(id, opts...)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to get user by id")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, errs.ErrInternalServerError
	}
	return user, nil
}

func (s *service) GetByUsername(username string, opts ...database.Option) (*model.User, error) {
	log.Info().Str("username", username).Msg("Getting user by username")
	user, err := s.repository.GetByUsername(username, opts...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by username")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, errs.ErrInternalServerError
	}
	return user, nil
}

func (s *service) GetAll(opts ...database.Option) []model.User {
	log.Info().Msg("Getting all users")
	return s.repository.GetAll(opts...)
}

func (s *service) Update(user *model.User) error {
	log.Info().Uint64("id", user.ID).Str("username", user.Username).Msg("Updating user")
	err := s.repository.Update(user)
	if err != nil {
		log.Error().Err(err).Uint64("id", user.ID).Msg("Failed to update user")
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return errs.ErrNotFound
		default:
			return errs.ErrInternalServerError
		}
	}
	return nil
}

func (s *service) Delete(id uint64) error {
	log.Info().Uint64("id", id).Msg("Deleting user by id")
	err := s.repository.Delete(id)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to delete user by id")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrNotFound
		}
		return errs.ErrInternalServerError
	}
	return nil
}
