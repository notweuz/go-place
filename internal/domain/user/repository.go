package user

import (
	"go-place/internal/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	GetByID(id uint64) (*model.User, error)
	GetAll() []model.User
	IncreaseAllCharges(maxCharges uint) error
	Delete(id uint64) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(user *model.User) error {
	log.Debug().Uint64("user_id", user.ID).Str("username", user.Username).Msg("Creating user")
	err := r.db.Create(user).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
	}
	return err
}

func (r *repository) Update(user *model.User) error {
	log.Debug().Uint64("user_id", user.ID).Msg("Updating user")
	err := r.db.Save(user).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to update user")
	}
	return err
}

func (r *repository) GetByID(id uint64) (*model.User, error) {
	log.Debug().Uint64("user_id", id).Msg("Getting user")
	user := model.User{}
	err := r.db.First(&user, id).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user")
	}
	return &user, err
}

func (r *repository) IncreaseAllCharges(maxCharges uint) error {
	log.Debug().Msg("Increasing all paint charges")
	err := r.db.Model(&model.User{}).Where("charges + 1 <= ?", maxCharges).Update("charges", gorm.Expr("charges + 1")).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to increase all paint charges")
	}
	return err
}

func (r *repository) GetAll() []model.User {
	log.Debug().Msg("Getting all users")
	var users []model.User
	r.db.Find(&users)
	log.Info().Int("amount", len(users)).Msg("Successfully fetched all users")
	return users
}

func (r *repository) Delete(id uint64) error {
	log.Debug().Uint64("user_id", id).Msg("Deleting user")
	err := r.db.Delete(&model.User{}, id).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete user")
	}
	return err
}
