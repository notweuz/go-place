package http

import (
	"go-place/internal/config"
	"go-place/internal/domain/auth"
	"go-place/internal/domain/pixel"
	"go-place/internal/domain/setting"
	"go-place/internal/domain/user"
	"go-place/internal/middleware"
	"go-place/internal/transport/websocket"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

type Router interface {
	Setup()
	setupAuthRoutes(api fiber.Router)
	setupUserRoutes(api fiber.Router)
	setupPixelRoutes(api fiber.Router)
	setupSettingRoutes(api fiber.Router)
	setupWebsocketRoutes(api fiber.Router)
}

type router struct {
	app              *fiber.App
	authHandler      auth.Handler
	userHandler      user.Handler
	pixelHandler     pixel.Handler
	settingHandler   setting.Handler
	websocketHandler websocket.Handler
	cfg              *config.Config
}

func NewRouter(app *fiber.App, authHandler auth.Handler, userHandler user.Handler, pixelHandler pixel.Handler, settingHandler setting.Handler, websocketHandler websocket.Handler, cfg *config.Config) Router {
	return &router{
		app:              app,
		authHandler:      authHandler,
		userHandler:      userHandler,
		pixelHandler:     pixelHandler,
		settingHandler:   settingHandler,
		websocketHandler: websocketHandler,
		cfg:              cfg,
	}
}

func (r *router) Setup() {
	log.Info().Msg("Setting up router")
	api := r.app.Group("/api")
	r.setupAuthRoutes(api)
	r.setupUserRoutes(api)
	r.setupPixelRoutes(api)
	r.setupSettingRoutes(api)
	r.setupWebsocketRoutes(api)
}

func (r *router) setupAuthRoutes(api fiber.Router) {
	log.Debug().Msg("Setting up auth routes")
	authRoute := api.Group("/auth")
	authRoute.Post("/register", r.authHandler.Register)
	authRoute.Post("/login", r.authHandler.Login)
}

func (r *router) setupUserRoutes(api fiber.Router) {
	log.Debug().Msg("Setting up user routes")
	userRoute := api.Group("/user")
	userRoute.Get("/self", middleware.AuthProtected(r.cfg), r.userHandler.GetCurrentUser)
	userRoute.Get("/:id", r.userHandler.GetByID)
}

func (r *router) setupPixelRoutes(api fiber.Router) {
	log.Debug().Msg("Setting up pixel routes")
	pixelRoute := api.Group("/pixel")
	pixelRoute.Get("/canvas", r.pixelHandler.GetBinaryCanvas)
	pixelRoute.Get("/search", r.pixelHandler.GetByCoordinates)
	pixelRoute.Get("/:id", r.pixelHandler.GetByID)
	pixelRoute.Get("/", r.pixelHandler.GetAll)
	pixelRoute.Post("/", middleware.AuthProtected(r.cfg), r.pixelHandler.Change)
}

func (r *router) setupSettingRoutes(api fiber.Router) {
	log.Debug().Msg("Setting up setting routes")
	settingRoute := api.Group("/setting")
	settingRoute.Get("/", r.settingHandler.Get)
}

func (r *router) setupWebsocketRoutes(api fiber.Router) {
	log.Debug().Msg("Setting up websocket routes")
	api.Get("/ws", middleware.WebsocketUpgrade, r.websocketHandler.Upgrade())
}
