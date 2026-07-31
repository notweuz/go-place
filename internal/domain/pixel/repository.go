package pixel

import (
	"errors"
	"go-place/internal/database"
	"go-place/internal/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repository interface {
	Create(pixel *model.Pixel) error
	Update(pixel *model.Pixel) error
	GetByID(id uint64) (*model.Pixel, error)
	GetByCoordinates(x, y uint64) (*model.Pixel, error)
	GetAll(opts ...database.Option) ([]model.Pixel, error)
	Delete(id uint64) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(pixel *model.Pixel) error {
	log.Debug().Uint64("x", pixel.X).Uint64("y", pixel.Y).Msg("creating pixel")
	err := r.db.Create(pixel).Error
	if err != nil {
		log.Error().Err(err).Msg("pixel creation failed")
	}
	return err
}

func (r *repository) Update(pixel *model.Pixel) error {
	log.Debug().Uint64("x", pixel.X).Uint64("y", pixel.Y).Msg("updating pixel")
	err := r.db.Save(pixel).Error
	if err != nil {
		log.Error().Err(err).Msg("pixel updating failed")
	}
	return err
}

func (r *repository) GetByID(id uint64) (*model.Pixel, error) {
	log.Debug().Uint64("id", id).Msg("getting pixel by id")
	var entity model.Pixel
	err := r.db.First(&entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Uint64("pixel_id", id).Msg("pixel not found")
		} else {
			log.Error().Err(err).Uint64("pixel_id", id).Msg("pixel fetch failed")
		}
	}
	return &entity, err
}

func (r *repository) GetByCoordinates(x, y uint64) (*model.Pixel, error) {
	log.Debug().Uint64("x", x).Uint64("y", y).Msg("getting pixel by x,y")
	var entity model.Pixel
	err := r.db.First(&entity, "x = ? AND y = ?", x, y).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Uint64("x", x).Uint64("y", y).Msg("pixel not found")
		} else {
			log.Error().Err(err).Uint64("x", x).Uint64("y", y).Msg("pixel fetch failed")
		}
	}
	return &entity, err
}

func (r *repository) GetAll(opts ...database.Option) ([]model.Pixel, error) {
	log.Debug().Msg("getting all pixels")
	var pixels []model.Pixel
	db := r.db
	for _, opt := range opts {
		db = opt(db)
	}
	if err := db.Find(&pixels).Error; err != nil {
		log.Error().Err(err).Msg("pixel fetch failed")
		return nil, err
	}
	log.Info().Int("amount", len(pixels)).Msg("pixel fetch success")
	return pixels, nil
}

func (r *repository) Delete(id uint64) error {
	log.Debug().Uint64("id", id).Msg("deleting pixel by id")
	err := r.db.Delete(&model.Pixel{}, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Uint64("id", id).Msg("can't delete pixel, pixel not found")
		} else {
			log.Error().Err(err).Uint64("id", id).Msg("pixel delete failed")
		}
	}
	return err
}
