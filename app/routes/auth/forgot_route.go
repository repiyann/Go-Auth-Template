package auth

import (
	providers "template-auth/app/providers/auth"

	"github.com/gofiber/fiber/v2"
)

func ForgotRoutes(app fiber.Router) {
	forgotController := providers.GetForgotProvider()

	app.Post("/request", forgotController.RequestOTP)
	app.Post("/validate", forgotController.ValidateOTP)
}
