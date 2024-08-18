package models

import (
	"time"

	"github.com/google/uuid"
)

type Auth struct {
	AuthID          uuid.UUID `json:"authID" gorm:"type:uuid; primaryKey"`
	AuthEmail       string    `json:"authEmail" gorm:"type:varchar(50); unique"`
	EmailToken      string    `json:"emailToken" gorm:"type:char(4)"`
	IsEmailVerified bool      `json:"isEmailVerified" gorm:"type:boolean; default:false"`
	AuthPassword    string    `json:"authPassword" gorm:"type:varchar(250)"`
	ResetOTP        string    `json:"resetOTP" gorm:"type:char(4)"`
	OTPExpire       time.Time `json:"otpExpire" gorm:"type:timestamp"`
	CreatedAt       time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (Auth) TableName() string {
	return "auth"
}
