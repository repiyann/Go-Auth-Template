package routes

import (
	"template-auth/app/routes/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	version := api.Group("/v1")

	authGroup := version.Group("/auth")
	forgotGroup := version.Group("/forgot")

	auth.AuthRoutes(authGroup)
	auth.ForgotRoutes(forgotGroup)
}
