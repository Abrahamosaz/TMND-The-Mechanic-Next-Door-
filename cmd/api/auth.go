package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Abrahamosaz/TMND/internal/services"
	"github.com/go-playground/validator/v10"
)


var validate = validator.New(validator.WithRequiredStructEnabled())

func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {

	var signUpDto services.Signup

	err := json.NewDecoder(r.Body).Decode(&signUpDto)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = validate.Struct(signUpDto)

	if err != nil {
		ValidateRequestBody(err, w)
		return
	}

	serviceApp := app.createNewServiceApp()

	user, statusCode, err := services.CreateUser(&serviceApp, signUpDto)

	if err != nil {
		log.Println("error creating new user: ", err.Error())
		message := err.Error()
		if (statusCode == http.StatusInternalServerError) {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}

	log.Println("newly created user: ", user)
	app.responseJSON(statusCode, w, "USer created successfully", user)
}


func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {

	var loginDto services.Login

	err := json.NewDecoder(r.Body).Decode(&loginDto)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = validate.Struct(loginDto)

	if err != nil {
		ValidateRequestBody(err, w)
		return
	}

	serviceApp := app.createNewServiceApp()

	statusCode, err := services.LoginUser(&serviceApp, loginDto)

	if err != nil {
		log.Println("error login: ", err.Error())
		message := err.Error()
		if (statusCode == http.StatusInternalServerError) {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}

	app.responseJSON(statusCode, w, "User login successfully", nil)
}

func (app *application) forgotPasswordHandler(w http.ResponseWriter, r *http.Request) {

	var forgotPasswordDto services.Email

	err := json.NewDecoder(r.Body).Decode(&forgotPasswordDto)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = validate.Struct(forgotPasswordDto)

	if err != nil {
		ValidateRequestBody(err, w)
		return
	}

	serviceApp := app.createNewServiceApp()

	statusCode, err := services.ForgotPassword(&serviceApp, forgotPasswordDto)

	if err != nil {
		log.Println("error login: ", err.Error())
		message := err.Error()
		if (statusCode == http.StatusInternalServerError) {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}

	app.responseJSON(statusCode, w, "Password reset OTP sent to your email.", nil)
}

func (app *application) changePasswordHandlder(w http.ResponseWriter, r *http.Request) {
	var changePasswordDto services.ChangePassword

	err := json.NewDecoder(r.Body).Decode(&changePasswordDto)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = validate.Struct(changePasswordDto)

	if err != nil {
		ValidateRequestBody(err, w)
		return
	}
	serviceApp := app.createNewServiceApp()
	statusCode, err := services.ResetPassword(&serviceApp, changePasswordDto)

	if err != nil {
		log.Println("error change password: ", err.Error())
		message := err.Error()
		if (statusCode == http.StatusInternalServerError) {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}
	app.responseJSON(http.StatusOK, w, "Password changed successfully", nil)
}

func (app *application) verifyEmailHandler(w http.ResponseWriter, r *http.Request) {

	var verifyOtp services.VerifyOtp

	err := json.NewDecoder(r.Body).Decode(&verifyOtp)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = validate.Struct(verifyOtp)

	if err != nil {
		ValidateRequestBody(err, w)
		return
	}

	serviceApp := app.createNewServiceApp()

	response, statusCode, err := services.VerifyEmail(&serviceApp, verifyOtp)

	if err != nil {
		log.Println("error verifying email: ", err.Error())
		message := err.Error()
		if statusCode == http.StatusInternalServerError {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}

	app.responseJSON(statusCode, w, "User email verified successfully", response)
}


func (app *application) verifyOtpHandler(w http.ResponseWriter, r *http.Request) {
	var verifyOtp services.VerifyOtp
	err := json.NewDecoder(r.Body).Decode(&verifyOtp)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = validate.Struct(verifyOtp)

	if err != nil {
		ValidateRequestBody(err, w)
		return
	}

	serviceApp := app.createNewServiceApp()
	statusCode, err := services.VerifyOtpCode(&serviceApp, verifyOtp)

	if err != nil {
		log.Println("error verifying otp code: ", err.Error())
		message := err.Error()
		if statusCode == http.StatusInternalServerError {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}
	app.responseJSON(statusCode, w, "Otp code verified successfuly", nil)

}