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

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Create(pixel *model.Pixel) (*model.Pixel, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) GetByID(id uint64) (*model.Pixel, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) GetByCoordinates(x, y uint64) (*model.Pixel, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) Update(pixel *model.Pixel) (*model.Pixel, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) Delete(id uint64) error {
	//TODO implement me
	panic("implement me")
}
