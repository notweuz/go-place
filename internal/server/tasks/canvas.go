package tasks

import (
	"go-place/internal/domain/pixel"
	"go-place/internal/domain/setting"
	"go-place/internal/model"

	"github.com/rs/zerolog/log"
)

func EnsureCanvasSize(pixelService pixel.Service, settingService setting.Service) {
	s, err := settingService.Get()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current settings, aborting EnsureCanvasSize task")
		return
	}
	w, h := s.CanvasWidth, s.CanvasHeight
	for x := uint64(0); x < w; x++ {
		for y := uint64(0); y < h; y++ {
			_, err := pixelService.GetByCoordinates(x, y)
			if err != nil {
				newPixel := &model.Pixel{
					X:      x,
					Y:      y,
					UserID: 0,
				}
				if createErr := pixelService.Create(newPixel); createErr != nil {
					log.Error().Err(createErr).Uint64("x", x).Uint64("y", y).Msg("Failed to create default pixel during canvas size enforcement")
				}
			}
		}
	}
}
