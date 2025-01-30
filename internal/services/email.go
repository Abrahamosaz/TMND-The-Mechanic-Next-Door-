package services

import (
	"github.com/Abrahamosaz/TMND/internal/templates"
	gomail "gopkg.in/mail.v2"
)


type OtpEmailPayload struct {
	fullName string
	otpCode string
} 


func (app *Application) sendOtpEmail(subject string, to string, context OtpEmailPayload, page string) error {

	message := gomail.NewMessage()
    // Set email headers
    message.SetHeader("From", app.Config.Smtp.From)
    message.SetHeader("To", to)
    message.SetHeader("Subject", subject)

    // Set email body
    message.SetBody("text/html", templates.OtpEmail(context.fullName, context.otpCode, page))
	return app.Smtp.DialAndSend(message);
}

