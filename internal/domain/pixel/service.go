package pixel

import (
	"encoding/binary"
	"errors"
	"go-place/internal/database"
	"go-place/internal/domain/setting"
	"go-place/internal/domain/user"
	"go-place/internal/errs"
	"go-place/internal/model"
	"go-place/internal/model/request"
	"go-place/internal/transport/websocket"
	"go-place/internal/transport/websocket/message"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Service interface {
	Create(pixel *model.Pixel) error
	GetByID(id uint64) (*model.Pixel, error)
	GetByCoordinates(x, y uint64) (*model.Pixel, error)
	GetAll(opts ...database.Option) ([]model.Pixel, error)
	GetBinaryCanvas() ([]byte, error)
	Change(pixels []request.ChangePixel, newAuthor uint64) ([]model.Pixel, error)
	Update(pixel *model.Pixel) error
	Delete(id uint64) error
}

type service struct {
	repository     Repository
	userService    user.Service
	settingService setting.Service
	wsHub          websocket.Hub
}

func NewService(repository Repository, userService user.Service, settingService setting.Service, wsHub websocket.Hub) Service {
	return &service{repository: repository, userService: userService, settingService: settingService, wsHub: wsHub}
}

func (s *service) Create(pixel *model.Pixel) error {
	log.Info().Uint64("pixel_id", pixel.ID).Uint64("x", pixel.X).Uint64("y", pixel.Y).Msg("Attempting to create pixel")
	err := s.repository.Create(pixel)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create pixel")
		return errs.ErrInternalServerError
	}
	s.wsHub.Broadcast(*message.NewMessage(message.EventPixelsChanged, message.NewPixelsChanged([]message.PixelChanged{*message.NewPixelChanged(pixel.ID, pixel.X, pixel.Y, pixel.UserID, pixel.Color)})))
	return nil
}

func (s *service) GetByID(id uint64) (*model.Pixel, error) {
	log.Info().Uint64("pixel_id", id).Msg("Attempting to get pixel by id")
	entity, err := s.repository.GetByID(id)
	if err != nil {
		log.Error().Err(err).Uint64("pixel_id", id).Msg("Failed to get pixel by id")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrPixelNotFound
		}
		return nil, errs.ErrInternalServerError
	}
	return entity, nil
}

func (s *service) GetByCoordinates(x, y uint64) (*model.Pixel, error) {
	log.Info().Uint64("x", x).Uint64("y", y).Msg("Attempting to get pixel by coordinates")
	entity, err := s.repository.GetByCoordinates(x, y)
	if err != nil {
		log.Error().Err(err).Uint64("x", x).Uint64("y", y).Msg("Failed to get pixel by coordinates")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrPixelNotFound
		}
		return nil, errs.ErrInternalServerError
	}
	return entity, nil
}

func (s *service) GetAll(opts ...database.Option) ([]model.Pixel, error) {
	log.Info().Msg("Attempting to get all pixels")
	pixels, err := s.repository.GetAll(opts...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all pixels")
		return nil, errs.ErrInternalServerError
	}
	return pixels, nil
}

func (s *service) GetBinaryCanvas() ([]byte, error) {
	settings, err := s.settingService.Get()
	if err != nil {
		return nil, err
	}
	w, h := settings.CanvasWidth, settings.CanvasHeight
	pixels, err := s.GetAll()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 8+w*h*3)
	binary.LittleEndian.PutUint32(buf[:4], uint32(w))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(h))

	for _, p := range pixels {
		if p.X < w && p.Y < h {
			index := 8 + (p.Y*w+p.X)*3
			r, g, b := parseHexToRGB(p.Color)
			buf[index] = r
			buf[index+1] = g
			buf[index+2] = b
		}
	}

	return buf, nil
}

func (s *service) Change(pixels []request.ChangePixel, newAuthor uint64) ([]model.Pixel, error) {
	log.Info().Int("count", len(pixels)).Uint64("author", newAuthor).Msg("Attempting to change pixels")

	settings, err := s.settingService.Get()
	if err != nil {
		return nil, err
	}

	dedupedMap := make(map[struct{ X, Y uint64 }]request.ChangePixel, len(pixels))
	for _, p := range pixels {
		if p.X >= settings.CanvasWidth || p.Y >= settings.CanvasHeight {
			return nil, errs.ErrPixelOutOfBounds
		}
		dedupedMap[struct{ X, Y uint64 }{X: p.X, Y: p.Y}] = p
	}

	dedupedPixels := make([]request.ChangePixel, 0, len(dedupedMap))
	for _, p := range dedupedMap {
		dedupedPixels = append(dedupedPixels, p)
	}

	err = s.userService.SpendCharge(newAuthor, uint(len(dedupedPixels)))
	if err != nil {
		return nil, err
	}

	models := make([]model.Pixel, len(dedupedPixels))
	for i, p := range dedupedPixels {
		models[i] = model.Pixel{X: p.X, Y: p.Y, Color: p.Color, UserID: newAuthor}
	}

	result, err := s.repository.Upsert(models)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upsert pixels")
		return nil, errs.ErrInternalServerError
	}

	changed := make([]message.PixelChanged, len(result))
	for i, px := range result {
		changed[i] = *message.NewPixelChanged(px.ID, px.X, px.Y, px.UserID, px.Color)
	}
	s.wsHub.Broadcast(*message.NewMessage(message.EventPixelsChanged, message.NewPixelsChanged(changed)))

	return result, nil
}

func (s *service) Update(pixel *model.Pixel) error {
	log.Info().Uint64("pixel_id", pixel.ID).Msg("Attempting to update pixel")
	err := s.repository.Update(pixel)
	if err != nil {
		log.Error().Err(err).Uint64("pixel_id", pixel.ID).Msg("Failed to update pixel")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrPixelNotFound
		}
		return errs.ErrInternalServerError
	}
	s.wsHub.Broadcast(*message.NewMessage(message.EventPixelsChanged, message.NewPixelsChanged([]message.PixelChanged{*message.NewPixelChanged(pixel.ID, pixel.X, pixel.Y, pixel.UserID, pixel.Color)})))
	return nil
}

func (s *service) Delete(id uint64) error {
	log.Info().Uint64("pixel_id", id).Msg("Attempting to delete pixel by id")
	err := s.repository.Delete(id)
	if err != nil {
		log.Error().Err(err).Uint64("pixel_id", id).Msg("Failed to delete pixel by id")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrPixelNotFound
		}
		return errs.ErrInternalServerError
	}
	return nil
}

func parseHexToRGB(color string) (r, g, b byte) {
	cleanColor := strings.ToUpper(strings.TrimSpace(strings.TrimPrefix(color, "#")))

	if len(cleanColor) == 6 {
		rVal, err1 := strconv.ParseUint(cleanColor[0:2], 16, 8)
		gVal, err2 := strconv.ParseUint(cleanColor[2:4], 16, 8)
		bVal, err3 := strconv.ParseUint(cleanColor[4:6], 16, 8)
		if err1 == nil && err2 == nil && err3 == nil {
			return byte(rVal), byte(gVal), byte(bVal)
		}
	} else if len(cleanColor) == 3 {
		rVal, err1 := strconv.ParseUint(cleanColor[0:1]+cleanColor[0:1], 16, 8)
		gVal, err2 := strconv.ParseUint(cleanColor[1:2]+cleanColor[1:2], 16, 8)
		bVal, err3 := strconv.ParseUint(cleanColor[2:3]+cleanColor[2:3], 16, 8)
		if err1 == nil && err2 == nil && err3 == nil {
			return byte(rVal), byte(gVal), byte(bVal)
		}
	} else if id, err := strconv.Atoi(color); err == nil && id >= 0 && id <= 255 {
		bVal := byte(id)
		return bVal, bVal, bVal
	}

	return 0, 0, 0
}
