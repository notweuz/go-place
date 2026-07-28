package pixel

import "go-place/internal/database"

func WithUser() database.Option {
	return database.WithPreload("User")
}
