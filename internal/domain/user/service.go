package user

import "go-place/internal/model"

type Service interface {
	Create(user *model.User) (*model.User, error)
	GetByID(id uint64) (*model.User, error)
	GetAll() ([]model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(id uint64) error
}

type service struct {
	repository Repository
}

func (s service) Create(user *model.User) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetByID(id uint64) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetAll() ([]model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) Update(user *model.User) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) Delete(id uint64) error {
	//TODO implement me
	panic("implement me")
}
