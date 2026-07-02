package router

import (
	"go-place/internal/domain/auth"
	"go-place/internal/domain/user"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

type Router interface {
	Setup()
	setupAuthRoutes(api fiber.Router)
	setupUserRoutes(api fiber.Router)
}

type router struct {
	app         *fiber.App
	authHandler auth.Handler
	userHandler user.Handler
}

func NewRouter(app *fiber.App, authHandler auth.Handler, userHandler user.Handler) Router {
	r := &router{
		app:         app,
		authHandler: authHandler,
		userHandler: userHandler,
	}

	return r
}

func (r *router) Setup() {
	log.Info().Msg("Setting up router")
	api := r.app.Group("/api")
	r.setupAuthRoutes(api)
	r.setupUserRoutes(api)
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
	userRoute.Get("/self", r.userHandler.GetCurrentUser)
	userRoute.Get("/:id", r.userHandler.GetByID)
}
