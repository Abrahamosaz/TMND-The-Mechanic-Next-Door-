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


	type VerifyOtpResponse struct {
		User 	models.User 	`json:"user"`
		Token 	string			`json:"token"`
	}