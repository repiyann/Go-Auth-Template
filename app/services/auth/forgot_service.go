package services

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	requests "template-auth/app/http/requests/auth"
	repositories "template-auth/app/repositories/auth"
	"time"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type ForgotService interface {
	RequestOTP(request *requests.RequestOTP) error
	ValidateOTP(request *requests.ValidateOTP) error
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
			Subject:  "Your OTP Code",
			TextPart: "Dear passenger 1, welcome to Mailjet! May the delivery force be with you!",
			HTMLPart: "<h3>Dear passenger 1, welcome to <a href=\"https://www.mailjet.com/\">Mailjet</a>!</h3><br />May the delivery force be with you!",
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mj.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
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

	return s.ForgotRepository.ValidateOTP(request.Email, auth)
}
