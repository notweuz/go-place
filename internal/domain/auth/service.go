package auth

import "go-place/internal/model/request"

type Service interface {
	Register(credentials *request.AuthCredentials) (string, error)
	Login(credentials *request.AuthCredentials) (string, error)
}
