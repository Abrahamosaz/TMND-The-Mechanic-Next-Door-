package services

import (
	"log"

	"github.com/Abrahamosaz/TMND/internal/templates"
	"github.com/resend/resend-go/v3"
)

type OtpEmailPayload struct {
	fullName string
	otpCode  string
}

func (app *Application) sendOtpEmail(subject string, to string, context OtpEmailPayload, page string) error {
	// Prepare email parameters according to Resend API
	params := &resend.SendEmailRequest{
		From:    app.Config.Resend.From,
		To:      []string{to},
		Subject: subject,
		Html:    templates.OtpEmail(context.fullName, context.otpCode, page),
	}

	// Send email via Resend API
	sent, err := app.Resend.Emails.Send(params)
	if err != nil {
		log.Printf("Error sending email via Resend: %v", err)
		return err
	}

	log.Printf("Email sent successfully via Resend. Email ID: %s", sent.Id)
	return nil
}
