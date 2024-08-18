package controllers

import (
	requests "template-auth/app/http/requests/auth"
	services "template-auth/app/services/auth"
	"template-auth/app/utils"

	"github.com/gofiber/fiber/v2"
)

type ForgotController struct {
	ForgotService services.ForgotService
}

func NewForgotController(service services.ForgotService) *ForgotController {
	return &ForgotController{
		ForgotService: service,
	}
}

func (controller *ForgotController) RequestOTP(c *fiber.Ctx) error {
	forgot := new(requests.RequestOTP)

	if err := c.BodyParser(forgot); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": fiber.Map{
				"errors": "Invalid input",
			},
		})
	}

	validationErrors := utils.Validate(forgot)
	if len(validationErrors.Errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": fiber.Map{
				"errors": validationErrors.Errors,
			},
		})
	}

	if err := controller.ForgotService.RequestOTP(forgot); err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": fiber.Map{
					"errors": "Email not found",
				},
			})
		}

		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Successfully generate OTP",
	})
}

func (controller *ForgotController) ValidateOTP(c *fiber.Ctx) error {
	valiOTP := new(requests.ValidateOTP)

	if err := c.BodyParser(valiOTP); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": fiber.Map{
				"errors": "Invalid input",
			},
		})
	}

	validationErrors := utils.Validate(valiOTP)
	if len(validationErrors.Errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": fiber.Map{
				"errors": validationErrors.Errors,
			},
		})
	}

	if err := controller.ForgotService.ValidateOTP(valiOTP); err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": fiber.Map{
					"errors": "Email not found",
				},
			})
		}

		if err.Error() == "invalid OTP" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": fiber.Map{
					"errors": "Invalid OTP",
				},
			})
		}

		if err.Error() == "OTP has expired" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": fiber.Map{
					"errors": "OTP has expired",
				},
			})
		}

		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Successfully validated OTP",
	})
}
