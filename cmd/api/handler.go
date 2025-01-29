package main

import (
	"fmt"
	"net/http"

	gomail "gopkg.in/mail.v2"
)


func (app *application) healthCheckHandler (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ping from server"))
}


func  (app *application) testSendMail(w http.ResponseWriter, r *http.Request) {
	message := gomail.NewMessage()

    // Set email headers
    message.SetHeader("From", app.config.smtp.from)
    message.SetHeader("To", "abrahamosazee2@gmail.com")
    message.SetHeader("Subject", "Hello from the Mailtrap team")

    // Set email body
    message.SetBody("text/plain", "This is the Test Body")


	if err := app.smtp.DialAndSend(message); err != nil {
        fmt.Println("Error sending email:", err)
        panic(err)
    } else {
        fmt.Println("Email sent successfully!")
    }
	w.Write([]byte("done"))
}