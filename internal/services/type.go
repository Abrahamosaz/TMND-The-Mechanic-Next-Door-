package services

import (
	"github.com/resend/resend-go/v3"
	"github.com/thexovc/TMND/internal/db"
	"github.com/thexovc/TMND/internal/models"
	"github.com/thexovc/TMND/internal/store"
)

type Application struct {
	Config   Config
	Store    store.Storage
	DbConfig db.DBConfig
	Resend   *resend.Client
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
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

type VerifyMechanicOtpResponse struct {
	Mechanic models.Mechanic `json:"user"`
	Token    string          `json:"token"`
}

type EditProfileInfo struct {
	UserProfile EditProfile
	URL         string
	FileName    string
	PublicId    string
}

type BookingFeeResponse struct {
	Fee float64 `json:"fee"`
}
