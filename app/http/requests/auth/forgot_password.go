package requests

type RequestOTP struct {
	Email string `json:"email" validate:"required,email"`
}

type ValidateOTP struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,len=4"`
}

type ResetPassword struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=16"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8,max=16"`
}
