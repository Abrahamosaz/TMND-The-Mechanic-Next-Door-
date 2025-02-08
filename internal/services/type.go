package services

import (
	"github.com/Abrahamosaz/TMND/internal/db"
	"github.com/Abrahamosaz/TMND/internal/models"
	"github.com/Abrahamosaz/TMND/internal/store"
	gomail "gopkg.in/mail.v2"
)


type Application struct {
	Config  	Config
	Store  		store.Storage
	DbConfig 	db.DBConfig
	Smtp     	*gomail.Dialer
}


type Config struct {
	Addr string
	Smtp SmtpConfig
}


type SmtpConfig struct {
	User 		string
	From 		string
	Password 	string
	Host 		string
	Port 		string
}


type VerifyUserOtpResponse struct {
	User 	models.User 	`json:"user"`
	Token 	string			`json:"token"`
}


type VerifyMechanicOtpResponse struct {
	Mechanic 	models.Mechanic 	`json:"user"`
	Token 		string			`json:"token"`
}

type EditProfileInfo struct {
	UserProfile 	EditProfile
	URL 			string
	FileName 		string
}



type BookingFeeResponse struct {
	Fee float64 `json:"fee"`
}

