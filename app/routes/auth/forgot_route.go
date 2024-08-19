package auth

import (
	"template-auth/app/middlewares"
	providers "template-auth/app/providers/auth"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ForgotRoutes(app fiber.Router) {
	forgotController := providers.GetForgotProvider()

	requestOTPLimiter := middlewares.RequestOTPLimiter(1, 5*time.Minute)
	validateOTPLimiter := middlewares.ValidateOTPLimiter(3, 5*time.Minute)

	app.Post("/request", requestOTPLimiter, forgotController.RequestOTP)
	app.Post("/validate", validateOTPLimiter, forgotController.ValidateOTP)
	app.Post("/reset", forgotController.ResetPassword)
}
