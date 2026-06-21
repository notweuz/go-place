package pixel

import (
	"go-place/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	Create(pixel *model.Pixel) (*model.Pixel, error)
	Update(pixel *model.Pixel) (*model.Pixel, error)
	GetByID(id uint64) (*model.Pixel, error)
	GetByCoordinates(x, y uint64) (*model.Pixel, error)
	GetAll() []model.Pixel
	Delete(id uint64) error
}

type repository struct {
	db *gorm.DB
}
