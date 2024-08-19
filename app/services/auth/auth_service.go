package services

import (
	"errors"
	"os"
	requests "template-auth/app/http/requests/auth"
	"template-auth/app/models"
	repositories "template-auth/app/repositories/auth"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(request *requests.RegisterRequest) error
	Login(reqest *requests.LoginRequest) (string, error)
	DecryptToken(tokenString string) (*models.Auth, error)
}

type authService struct {
	AuthRepository repositories.AuthRepository
	JWTSecret      string
}

func NewAuthService(repo repositories.AuthRepository) AuthService {
	return &authService{
		AuthRepository: repo,
		JWTSecret:      os.Getenv("JWT_SECRET"),
	}
}

func (s *authService) Register(request *requests.RegisterRequest) error {
	if request.Password != request.ConfirmPassword {
		return errors.New("passwords not match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	auth := &models.Auth{
		AuthID:       uuid.New(),
		AuthEmail:    request.Email,
		AuthPassword: string(hashedPassword),
	}

	err = s.AuthRepository.Register(auth)
	if err != nil {
		return errors.New("duplicate")
	}

	return nil
}

func (s *authService) Login(request *requests.LoginRequest) (string, error) {
	auth, err := s.AuthRepository.FindByEmail(request.Email)
	if err != nil {
		return "", err
	}
	if auth == nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(auth.AuthPassword), []byte(request.Password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.GenerateToken(auth)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) GenerateToken(auth *models.Auth) (string, error) {
	jwtExpiration := time.Hour * 24
	claims := jwt.MapClaims{
		"sub": auth.AuthID.String(),
		"exp": time.Now().Add(jwtExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) DecryptToken(tokenString string) (*models.Auth, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(s.JWTSecret), nil
	})
	if err != nil {
		return nil, errors.New("error parsing token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	authIDStr, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	authID, err := uuid.Parse(authIDStr)
	if err != nil {
		return nil, errors.New("error parsing token")
	}

	auth, err := s.AuthRepository.FindByID(authID)
	if err != nil {
		return nil, errors.New("authentication failed")
	}

	return auth, nil
}
