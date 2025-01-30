package main

import (
	"net/http"
)


func (app *application) healthCheckHandler (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ping from server"))
}


// func  (app *application) testSendMail(w http.ResponseWriter, r *http.Request) {
// 	if err := app.sendVerificationEmail("Test email", "Ibmeshach@gmail.com", VerifyEmailPayload{
// 		fullName: "abraham",
// 		otpCode: "1232",
// 	}); err != nil {
//         fmt.Println("Error sending email:", err)
//         panic(err)
//     } else {
//         fmt.Println("Email sent successfully!")
//     }
// 	w.Write([]byte("done"))
// }