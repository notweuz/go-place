package server

import (
	"go-place/internal/config"
	"go-place/internal/domain/auth"
	"go-place/internal/domain/pixel"
	"go-place/internal/domain/setting"
	"go-place/internal/domain/user"
	"go-place/internal/middleware"
	"go-place/internal/transport/http"
	"go-place/internal/transport/websocket"

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

	settingDB := setting.NewRepository(db)
	settingSVC := setting.NewService(settingDB)

	userDB := user.NewRepository(db)
	userSVC := user.NewService(userDB)
	userHR := user.NewHandler(userSVC)

	authSVC := auth.NewService(userSVC, cfg)
	authHR := auth.NewHandler(authSVC)

	pixelDB := pixel.NewRepository(db)
	pixelSVC := pixel.NewService(pixelDB, userSVC, settingSVC)
	pixelHR := pixel.NewHandler(pixelSVC)

	wsHub := websocket.NewHub()
	go wsHub.Run()

	wsHR := websocket.NewHandler(wsHub)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	app.Use(cors.New())
	app.Use(middleware.Logger)

	appRouter := http.NewRouter(app, authHR, userHR, pixelHR, wsHR)
	appRouter.Setup()

	return &App{
		DB:     db,
		Config: cfg,
		Fiber:  app,
	}
}
