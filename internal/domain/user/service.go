package user

import (
	"errors"
	"go-place/internal/errs"
	"go-place/internal/model"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Service interface {
	Create(user *model.User) (*model.User, error)
	GetByID(id uint64) (*model.User, error)
	GetAll() []model.User
	Update(user *model.User) (*model.User, error)
	Delete(id uint64) error
}

type service struct {
	repository Repository
}

func (s *service) Create(user *model.User) (*model.User, error) {
	log.Info().Uint64("id", user.ID).Str("username", user.Username).Msg("Creating user")
	err := s.repository.Create(user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errs.NewAppError(errs.ErrDuplicatedKey, fiber.StatusConflict, "User with same username already exists")
		}
		return nil, errs.NewAppError(errs.ErrInternalServerError, fiber.StatusInternalServerError, "Failed to create user")
	}
	return user, nil
}

func (s *service) GetByID(id uint64) (*model.User, error) {
	log.Info().Uint64("id", id).Msg("Getting user by id")
	user, err := s.repository.GetByID(id)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to get user by id")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewAppError(errs.ErrInternalServerError, fiber.StatusNotFound, "User not found")
		}
		return nil, errs.NewAppError(errs.ErrInternalServerError, fiber.StatusInternalServerError, "Failed to get user by id")
	}
	return user, nil
}

func (s *service) GetAll() []model.User {
	log.Info().Msg("Getting all users")
	return s.repository.GetAll()
}

func (s *service) Update(user *model.User) (*model.User, error) {
	log.Info().Uint64("id", user.ID).Str("username", user.Username).Msg("Updating user")
	err := s.repository.Update(user)
	if err != nil {
		log.Error().Err(err).Uint64("id", user.ID).Msg("Failed to update user")
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, errs.NewAppError(errs.ErrInternalServerError, fiber.StatusNotFound, "User not found")
		case errors.Is(err, errs.NewAppError(errs.ErrInternalServerError, fiber.StatusInternalServerError, "Failed to update user")):
			return nil, errs.NewAppError(errs.ErrInternalServerError, fiber.StatusInternalServerError, "Failed to update user")
		default:
			return nil, errs.NewAppError(errs.ErrInternalServerError, fiber.StatusInternalServerError, "Failed to update user")
		}
	}
	return user, nil
}

func (s *service) Delete(id uint64) error {
	log.Info().Uint64("id", id).Msg("Deleting user by id")
	err := s.repository.Delete(id)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to delete user by id")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewAppError(errs.ErrInternalServerError, fiber.StatusNotFound, "User not found")
		}
		return errs.NewAppError(errs.ErrInternalServerError, fiber.StatusInternalServerError, "Failed to delete user")
	}
	return nil
}
