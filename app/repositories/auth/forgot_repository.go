package repositories

import (
	"errors"
	"template-auth/app/models"
	"time"

	"gorm.io/gorm"
)

type ForgotRepository interface {
	GenerateOTP(email string, auth *models.Auth) error
	ValidateOTP(email string, auth *models.Auth) error
	ResetPassword(email string, auth *models.Auth) error
}

type forgotRepository struct {
	DB *gorm.DB
}

func NewForgotRepository(db *gorm.DB) ForgotRepository {
	return &forgotRepository{
		DB: db,
	}
}

func (r *forgotRepository) GenerateOTP(email string, auth *models.Auth) error {
	result := r.DB.Model(&models.Auth{}).Where("auth_email = ?", email).Updates(map[string]interface{}{
		"reset_otp":  auth.ResetOTP,
		"otp_expire": auth.OTPExpire,
	})

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}

		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (r *forgotRepository) ValidateOTP(email string, auth *models.Auth) error {
	result := r.DB.Model(&models.Auth{}).Where("auth_email = ?", email).Updates(map[string]interface{}{
		"reset_otp":  "",
		"otp_expire": time.Time{},
	})

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (r *forgotRepository) ResetPassword(email string, auth *models.Auth) error {
	result := r.DB.Model(&models.Auth{}).Where("auth_email = ?", email).Updates(map[string]interface{}{
		"auth_password": auth.AuthPassword,
	})

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}
