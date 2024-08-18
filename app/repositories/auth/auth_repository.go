package repositories

import (
	"errors"
	"template-auth/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Register(auth *models.Auth) error
	FindByEmail(email string) (*models.Auth, error)
	FindByID(authID uuid.UUID) (*models.Auth, error)
}

type authRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{
		DB: db,
	}
}

func (r *authRepository) Register(auth *models.Auth) error {
	result := r.DB.Create(auth)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return gorm.ErrDuplicatedKey
		}

		return result.Error
	}

	return nil
}

func (r *authRepository) FindByEmail(email string) (*models.Auth, error) {
	var auth models.Auth

	result := r.DB.Where("auth_email = ?", email).First(&auth)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, result.Error
	}

	return &auth, nil
}

func (r *authRepository) FindByID(authID uuid.UUID) (*models.Auth, error) {
	var auth models.Auth

	result := r.DB.Where("auth_id = ?", authID).First(&auth)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, result.Error
	}

	return &auth, nil
}
