package user

import (
	"errors"
	"go-place/internal/database"
	"go-place/internal/domain/setting"
	"go-place/internal/errs"
	"go-place/internal/model"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Service interface {
	Create(user *model.User) error
	GetByID(id uint64, opts ...database.Option) (*model.User, error)
	GetByUsername(username string, opts ...database.Option) (*model.User, error)
	GetAll(opts ...database.Option) ([]model.User, error)
	GetSelf(id uint64) (*model.User, error)
	Update(user *model.User, opts ...database.Option) error
	Delete(id uint64) error
	SyncAndGet(id uint64) (*model.User, error)
	SpendCharge(id uint64, pixelsToSpend uint, opts ...database.Option) error
}

type service struct {
	repository     Repository
	settingService setting.Service
}

func NewService(repository Repository, settingService setting.Service) Service {
	return &service{
		repository:     repository,
		settingService: settingService,
	}
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
			return nil, errs.ErrUserNotFound
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
			return nil, errs.ErrUserNotFound
		}
		return nil, errs.ErrInternalServerError
	}
	return user, nil
}

func (s *service) GetAll(opts ...database.Option) ([]model.User, error) {
	log.Info().Msg("Getting all users")
	users, err := s.repository.GetAll(opts...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all users")
		return nil, errs.ErrInternalServerError
	}
	return users, nil
}

func (s *service) Update(user *model.User, opts ...database.Option) error {
	log.Info().Uint64("id", user.ID).Str("username", user.Username).Msg("Updating user")
	err := s.repository.Update(user, opts...)
	if err != nil {
		log.Error().Err(err).Uint64("id", user.ID).Msg("Failed to update user")
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return errs.ErrUserNotFound
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
			return errs.ErrUserNotFound
		}
		return errs.ErrInternalServerError
	}
	return nil
}

func (s *service) SyncAndGet(id uint64) (*model.User, error) {
	stg, err := s.settingService.Get()
	if err != nil {
		return nil, err
	}

	fetchOpts := []database.Option{database.WithLock()}
	user, err := s.GetByID(id, fetchOpts...)
	if err != nil {
		return nil, err
	}

	user.SyncCharges(stg.MaxCharges, time.Duration(stg.CooldownSeconds)*time.Second)

	err = s.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetSelf(id uint64) (*model.User, error) {
	log.Info().Uint64("id", id).Msg("Getting self profile with in-memory charge sync")
	user, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	stg, err := s.settingService.Get()
	if err != nil {
		return nil, err
	}

	// I'm not sure, but I think it's possible to use SyncAndGet function, but it's gonna block user entry, so ig I'll leave it as it is
	user.SyncCharges(stg.MaxCharges, time.Duration(stg.CooldownSeconds)*time.Second)
	return user, nil
}

func (s *service) SpendCharge(id uint64, pixelsToSpend uint, opts ...database.Option) error {
	stg, err := s.settingService.Get()
	if err != nil {
		return err
	}

	fetchOpts := append([]database.Option{database.WithLock()}, opts...)
	user, err := s.GetByID(id, fetchOpts...)
	if err != nil {
		return err
	}

	user.SyncCharges(stg.MaxCharges, time.Duration(stg.CooldownSeconds)*time.Second)

	if user.Charges < pixelsToSpend {
		return errs.ErrNotEnoughCharges
	}

	user.Charges -= pixelsToSpend

	return s.Update(user, opts...)
}
