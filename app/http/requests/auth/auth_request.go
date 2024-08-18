package requests

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=16"`
}

type RegisterRequest struct {
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required,min=8,max=16"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8,max=16"`
}
