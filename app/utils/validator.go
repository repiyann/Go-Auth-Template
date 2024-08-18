package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var customMessages = map[string]string{
	"required": "This field is required.",
	"email":    "Please provide a valid email address.",
	"min":      "This field must be at least %s characters long.",
	"max":      "This field must be no more than %s characters long.",
	"len":      "This field must be exactly %s characters long.",
}

type ValidationErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

func customValidationErrors(err error) ValidationErrorResponse {
	response := ValidationErrorResponse{Errors: make(map[string]string)}
	if _, ok := err.(*validator.InvalidValidationError); ok {
		response.Errors["general"] = err.Error()
		return response
	}
	if err == nil {
		return response
	}

	validationErrors := err.(validator.ValidationErrors)
	for _, err := range validationErrors {
		field := err.Field()
		tag := err.Tag()
		param := err.Param()

		if msg, exists := customMessages[tag]; exists {
			if tag == "min" || tag == "max" || tag == "len" {
				msg = fmt.Sprintf(msg, param)
			}
			response.Errors[field] = msg
		} else {
			response.Errors[field] = err.Error()
		}
	}

	return response
}

func Validate(data interface{}) ValidationErrorResponse {
	validate := validator.New()

	err := validate.Struct(data)
	if err != nil {
		return customValidationErrors(err)
	}

	return ValidationErrorResponse{Errors: make(map[string]string)}
}
