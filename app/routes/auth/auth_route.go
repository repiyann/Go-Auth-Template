package auth

import (
	providers "template-auth/app/providers/auth"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app fiber.Router) {
	authController := providers.GetAuthProvider()

	app.Post("/register", authController.Register)
	app.Post("/login", authController.Login)
	app.Get("/decrypt", authController.DecryptToken)
}
