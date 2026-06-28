package auth

import (
	"go-place/internal/auth"
	"go-place/internal/config"
	"go-place/internal/domain/user"
	"go-place/internal/errs"
	"go-place/internal/model"
	"go-place/internal/model/request"

	"github.com/rs/zerolog/log"
)

type Service interface {
	Register(credentials *request.AuthCredentials) (*string, error)
	Login(credentials *request.AuthCredentials) (*string, error)
}

type service struct {
	userService user.Service
	config      *config.Config
}

func NewService(userService user.Service, cfg *config.Config) Service {
	return &service{userService: userService, config: cfg}
}

func (s *service) Register(credentials *request.AuthCredentials) (*string, error) {
	log.Info().Str("username", credentials.Login).Msg("Registering user")
	passwordHash, err := auth.GeneratePasswordHash(credentials.Password)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate password hash")
		return nil, errs.ErrFailedBCrypt
	}

	newUser := model.NewUser(credentials.Login, passwordHash)
	err = s.userService.Create(newUser)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		return nil, err
	}

	token, err := auth.GenerateToken(newUser.ID, s.config.JwtSecret)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate token")
		return nil, errs.ErrFailedJWT
	}

	log.Info().Str("username", newUser.Username).Msg("Successfully registered user")
	return &token, nil
}

func (s *service) Login(credentials *request.AuthCredentials) (*string, error) {
	log.Info().Str("username", credentials.Login).Msg("User login")
	userInDB, err := s.userService.GetByUsername(credentials.Login)
	if err != nil {
		return nil, err
	}
	if !auth.VerifyPasswordHash(credentials.Password, userInDB.Password) {
		log.Error().Err(err).Msg("User login failed")
		return nil, errs.ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(userInDB.ID, s.config.JwtSecret)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate token")
		return nil, errs.ErrFailedJWT
	}

	log.Info().Str("username", userInDB.Username).Msg("Successfully logged in")
	return &token, nil
}
