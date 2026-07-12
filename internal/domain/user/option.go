package user

import (
	"go-place/internal/database"
)

func WithPixels() database.Option {
	return database.WithPreload("Pixels")
}
