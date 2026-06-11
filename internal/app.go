package internal

import (
	"go-place/internal/config"

	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Config *config.Config
}

func NewApp() {

}
