package router

import (
	"go-place/internal/domain/auth"

	"github.com/gofiber/fiber/v3"
)

type Router interface {
	Setup()
	setupAuthRoutes(api fiber.Router)
}

type router struct {
	app         *fiber.App
	authHandler auth.Handler
}

func NewRouter(app *fiber.App, authHandler auth.Handler) Router {
	r := &router{
		app:         app,
		authHandler: authHandler,
	}

	return r
}

func (r *router) Setup() {
	api := r.app.Group("/api")
	r.setupAuthRoutes(api)
}

func (r *router) setupAuthRoutes(api fiber.Router) {
	authRoute := api.Group("/auth")
	authRoute.Post("/register", r.authHandler.Register)
	authRoute.Post("/login", r.authHandler.Login)
}
