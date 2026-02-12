package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/resend/resend-go/v3"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	apiKey := os.Getenv("RESEND_API_KEY")
	fromEmail := os.Getenv("RESEND_FROM_EMAIL")

	if apiKey == "" || fromEmail == "" {
		log.Fatal("RESEND_API_KEY or RESEND_FROM_EMAIL not set in .env")
	}

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run scripts/test_send_email.go <recipient-email>")
	}

	toEmail := os.Args[1]

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    fromEmail,
		To:      []string{toEmail},
		Subject: "TMND Test Email",
		Html:    "<strong>Hello from TMND!</strong><p>This is a test email sent via Resend.</p>",
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	fmt.Printf("Email sent successfully! ID: %s\n", sent.Id)
}
