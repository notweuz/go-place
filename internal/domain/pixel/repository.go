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

func NewPixelRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(pixel *model.Pixel) (*model.Pixel, error) {
	err := r.db.Create(pixel).Error
	return pixel, err
}

func (r *repository) Update(pixel *model.Pixel) (*model.Pixel, error) {
	err := r.db.Save(pixel).Error
	return pixel, err
}

func (r *repository) GetByID(id uint64) (*model.Pixel, error) {
	var entity model.Pixel
	err := r.db.Where("id = ?", id).First(&entity).Error
	return &entity, err
}

func (r *repository) GetByCoordinates(x, y uint64) (*model.Pixel, error) {
	var entity model.Pixel
	err := r.db.Where("x = ? AND y = ?", x, y).First(&entity).Error
	return &entity, err
}

func (r *repository) GetAll() []model.Pixel {
	var entities []model.Pixel
	r.db.Find(&entities)
	return entities
}

func (r *repository) Delete(id uint64) error {
	return r.db.Delete(&model.Pixel{}, id).Error
}
