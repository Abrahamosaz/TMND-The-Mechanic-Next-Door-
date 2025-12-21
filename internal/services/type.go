package services

import (
	"github.com/Abrahamosaz/TMND/internal/db"
	"github.com/Abrahamosaz/TMND/internal/models"
	"github.com/Abrahamosaz/TMND/internal/store"
	"github.com/resend/resend-go/v3"
)


type Application struct {
	Config  	Config
	Store  		store.Storage
	DbConfig 	db.DBConfig
	Resend     	*resend.Client
}


type Config struct {
	Addr   string
	Resend ResendConfig
}


type ResendConfig struct {
	ApiKey string
	From   string
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
	PublicId 		string
}



type BookingFeeResponse struct {
	Fee float64 `json:"fee"`
}
