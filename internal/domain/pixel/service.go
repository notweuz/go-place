package pixel

import "go-place/internal/model"

type Service interface {
	Create(pixel *model.Pixel) (*model.Pixel, error)
	GetByID(id uint64) (*model.Pixel, error)
	GetByCoordinates(x, y uint64) (*model.Pixel, error)
	Update(pixel *model.Pixel) (*model.Pixel, error)
	Delete(id uint64) error
}

type service struct {
	repository Repository
}
