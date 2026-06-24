package setting

import (
	"errors"
	"go-place/internal/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repository interface {
	Create(setting *model.Setting) error
	Get() (*model.Setting, error)
	Update(*model.Setting) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(setting *model.Setting) error {
	log.Debug().Msgf("Creating setting")
	err := r.db.Create(setting).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to create setting")
		return err
	}
	return nil
}

func (r *repository) Get() (*model.Setting, error) {
	log.Debug().Msg("Trying to get settings")
	var entity *model.Setting
	err := r.db.First(&entity).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to get settings")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			setting := model.Setting{}
			err2 := r.Create(&setting)
			if err2 != nil {
				log.Error().Err(err2).Msg("Failed to create setting")
				return nil, err2
			}
			return &setting, nil
		}
		return nil, err
	}
	return entity, nil
}

func (r *repository) Update(setting *model.Setting) error {
	log.Debug().Msg("Updating setting")
	err := r.db.Save(setting).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to update setting")
		return err
	}
	return nil
}
