package server

import (
	"context"
	"go-place/internal/config"
	"go-place/internal/domain/auth"
	"go-place/internal/domain/pixel"
	"go-place/internal/domain/setting"
	"go-place/internal/domain/user"
	"go-place/internal/middleware"
	"go-place/internal/transport/http"
	"go-place/internal/transport/websocket"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type App interface {
	Shutdown()
}

type app struct {
	DB       *gorm.DB
	Config   *config.Config
	Fiber    *fiber.App
	wsCancel context.CancelFunc
}

func NewApp(cfg *config.Config) App {
	db, err := SetupDatabase(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start database!")
	}

	settingDB := setting.NewRepository(db)
	settingSVC := setting.NewService(settingDB)

	userDB := user.NewRepository(db)
	userSVC := user.NewService(userDB)
	userHR := user.NewHandler(userSVC)

	authSVC := auth.NewService(userSVC, settingSVC, cfg)
	authHR := auth.NewHandler(authSVC)

	ctx, cancel := context.WithCancel(context.Background())
	wsHub := websocket.NewHub()
	go wsHub.Run(ctx)

	wsHR := websocket.NewHandler(wsHub)

	pixelDB := pixel.NewRepository(db)
	pixelSVC := pixel.NewService(pixelDB, userSVC, settingSVC, wsHub)
	pixelHR := pixel.NewHandler(pixelSVC)

	application := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	application.Use(cors.New())
	application.Use(middleware.Logger)

	appRouter := http.NewRouter(application, authHR, userHR, pixelHR, wsHR, cfg)
	appRouter.Setup()

	return &app{
		DB:       db,
		Config:   cfg,
		Fiber:    application,
		wsCancel: cancel,
	}
}

func (a *app) Shutdown() {
	log.Info().Msg("Shutting down server")
	a.wsCancel()
	log.Info().Msg("WS Hub stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.Fiber.ShutdownWithContext(ctx); err != nil {
		log.Panic().Err(err).Msg("Error shutting down")
	}
	log.Info().Msg("Server gracefully stopped")
}
