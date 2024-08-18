package controllers

import (
	requests "template-auth/app/http/requests/auth"
	services "template-auth/app/services/auth"
	"template-auth/app/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	AuthService services.AuthService
}

func NewAuthController(service services.AuthService) *AuthController {
	return &AuthController{
		AuthService: service,
	}
}

func (controller *AuthController) Register(c *fiber.Ctx) error {
	register := new(requests.RegisterRequest)

	if err := c.BodyParser(register); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": fiber.Map{
				"errors": "Invalid input",
			},
		})
	}

	validationErrors := utils.Validate(register)
	if len(validationErrors.Errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": fiber.Map{
				"errors": validationErrors.Errors,
			},
		})
	}

	err := controller.AuthService.Register(register)
	if err != nil {
		if err.Error() == "duplicate" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false,
				"message": fiber.Map{
					"errors": "Email already in use",
				},
			})
		}

		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Succesfully registered",
	})
}

func (controller *AuthController) Login(c *fiber.Ctx) error {
	login := new(requests.LoginRequest)

	if err := c.BodyParser(login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": fiber.Map{
				"errors": "Invalid input",
			},
		})
	}

	validationErrors := utils.Validate(login)
	if len(validationErrors.Errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": fiber.Map{
				"errors": validationErrors.Errors,
			},
		})
	}

	token, err := controller.AuthService.Login(login)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Succesfully logged in",
		"data":    token,
	})
}

func (controller *AuthController) DecryptToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": fiber.Map{
				"errors": "Authorization header is required",
			},
		})
	}

	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	if tokenString == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": fiber.Map{
				"errors": "Token is required",
			},
		})
	}

	auth, err := controller.AuthService.DecryptToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": fiber.Map{
				"errors": "Invalid token",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Token successfully decrypted",
		"data":    auth,
	})
}
