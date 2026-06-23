package server

import (
	"go-place/internal/config"
	"go-place/internal/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Config *config.Config
	Fiber  *fiber.App
}

func NewApp(cfg *config.Config) *App {
	db, err := SetupDatabase(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start database!")
	}

	//userRepo := user.NewRepository(db)
	//pixelRepo := pixel.NewRepository(db)

	//pixelService := pixel.NewService(pixelRepo)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	app.Use(cors.New())

	return &App{
		DB:     db,
		Config: cfg,
		Fiber:  app,
	}
}
