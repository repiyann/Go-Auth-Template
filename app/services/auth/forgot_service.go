package services

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	requests "template-auth/app/http/requests/auth"
	"template-auth/app/models"
	repositories "template-auth/app/repositories/auth"
	"time"

	"github.com/mailjet/mailjet-apiv3-go/v4"
	"golang.org/x/crypto/bcrypt"
)

type ForgotService interface {
	RequestOTP(request *requests.RequestOTP) error
	ValidateOTP(request *requests.ValidateOTP) error
	ResetPassword(request *requests.ResetPassword) error
}

type forgotService struct {
	ForgotRepository repositories.ForgotRepository
	AuthRepository   repositories.AuthRepository
	publicKey        string
	secretKey        string
}

func NewForgotService(repo repositories.ForgotRepository, authRepo repositories.AuthRepository) ForgotService {
	return &forgotService{
		ForgotRepository: repo,
		AuthRepository:   authRepo,
		publicKey:        os.Getenv("MJ_APIKEY_PUBLIC"),
		secretKey:        os.Getenv("MJ_APIKEY_PRIVATE"),
	}
}

func (s *forgotService) RequestOTP(request *requests.RequestOTP) error {
	auth, err := s.AuthRepository.FindByEmail(request.Email)
	if err != nil {
		return err
	}
	if auth == nil {
		return errors.New("record not found")
	}

	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return err
	}

	otp := rand.Intn(10000)
	otpStr := fmt.Sprintf("%04d", otp)
	expiration := time.Now().In(location).Add(5 * time.Minute).UTC()

	mj := mailjet.NewMailjetClient(s.publicKey, s.secretKey)
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "akundavee02@gmail.com",
				Name:  "TravelinAja",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: request.Email,
				},
			},
			Subject: "Your OTP Code",
			HTMLPart: fmt.Sprintf(`
			<html>
			<body>
				<h2>Hi there!</h2>
				<p>We received a request to verify your email address. To proceed, please use the OTP code below:</p>
				<h1 style="font-size: 36px; color: #007BFF;">%s</h1>
				<p>This code is valid for the next 5 minutes. If you did not request this, please ignore this email.</p>
				<p>Best regards,<br>The TravelinAja Team</p>
				<footer style="margin-top: 20px; font-size: 12px; color: #888;">
					<p>TravelinAja</p>
					<p>1234 Travel St, Adventure City, AC 12345</p>
					<p><a href="https://www.travelinaja.com">Visit our website</a></p>
				</footer>
			</body>
			</html>
			`, otpStr),
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mj.SendMailV31(&messages)
	if err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}
	fmt.Printf("Data: %+v\n", res)

	auth.ResetOTP = otpStr
	auth.OTPExpire = expiration

	return s.ForgotRepository.GenerateOTP(request.Email, auth)
}

func (s *forgotService) ValidateOTP(request *requests.ValidateOTP) error {
	auth, err := s.AuthRepository.FindByEmail(request.Email)
	if err != nil {
		return err
	}
	if auth == nil {
		return errors.New("record not found")
	}

	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return err
	}

	currentTime := time.Now().In(location).UTC()
	OTPExpiry := auth.OTPExpire.UTC()

	if currentTime.After(OTPExpiry) {
		return errors.New("OTP has expired")
	}
	if auth.ResetOTP != request.OTP {
		return errors.New("invalid OTP")
	}

	return s.ForgotRepository.ValidateOTP(request.Email)
}

func (s *forgotService) ResetPassword(request *requests.ResetPassword) error {
	if request.Password != request.ConfirmPassword {
		return errors.New("passwords not match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	auth := &models.Auth{
		AuthEmail:    request.Email,
		AuthPassword: string(hashedPassword),
	}

	return s.ForgotRepository.ResetPassword(auth)
}
