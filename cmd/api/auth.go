package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/thexovc/TMND/internal/services"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

// userSignupHandler godoc
// @Summary User signup
// @Description Sign up a new user with full name, email, phone number, and password.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param signup body services.Signup true "Signup details"
// @Success 200 {object} Response{data=models.User} "User created successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "User already exists"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/signup [post]
func (app *application) userSignupHandler(w http.ResponseWriter, r *http.Request) {

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
		if statusCode == http.StatusInternalServerError {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}

	log.Println("newly created user: ", user)
	app.responseJSON(statusCode, w, "USer created successfully", user)
}

// userLoginHandler godoc
// @Summary User login
// @Description Login a user with email and password.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param login body services.Login true "Login details"
// @Success 200 {object} Response "User login successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/login [post]
func (app *application) userLoginHandler(w http.ResponseWriter, r *http.Request) {

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
		if statusCode == http.StatusInternalServerError {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}

	app.responseJSON(statusCode, w, "User login successfully", nil)
}

func (app *application) mechanicLoginHandler(w http.ResponseWriter, r *http.Request) {

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

	statusCode, err := services.LoginMechanic(&serviceApp, loginDto)

	if err != nil {
		log.Println("error login: ", err.Error())
		message := err.Error()
		if statusCode == http.StatusInternalServerError {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}

	app.responseJSON(statusCode, w, "Mechanic login successfully", nil)
}

// forgotPasswordHandler godoc
// @Summary Forgot password
// @Description Send a password reset OTP to the user's email.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param forgotPassword body services.Email true "User email"
// @Success 200 {object} Response "Password reset OTP sent to your email."
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/forgot-password [post]
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

	statusCode, err := services.SendUserOtpCode(&serviceApp, forgotPasswordDto, "forgotPassword")

	if err != nil {
		log.Println("error login: ", err.Error())
		message := err.Error()
		if statusCode == http.StatusInternalServerError {
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
	statusCode, err := services.ResetUserPassword(&serviceApp, changePasswordDto)

	if err != nil {
		log.Println("error change password: ", err.Error())
		message := err.Error()
		if statusCode == http.StatusInternalServerError {
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

	response, statusCode, err := services.VerifyUserEmail(&serviceApp, verifyOtp)

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
	statusCode, err := services.VerifyUserOtpCode(&serviceApp, verifyOtp)

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

func (app *application) resendOtpHandler(w http.ResponseWriter, r *http.Request) {
	var resendOtpDto services.Email

	err := json.NewDecoder(r.Body).Decode(&resendOtpDto)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = validate.Struct(resendOtpDto)

	if err != nil {
		ValidateRequestBody(err, w)
		return
	}

	serviceApp := app.createNewServiceApp()

	statusCode, err := services.SendUserOtpCode(&serviceApp, resendOtpDto, "resendCode")

	if err != nil {
		log.Println("error resend otp: ", err.Error())
		message := err.Error()
		if statusCode == http.StatusInternalServerError {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}
	app.responseJSON(http.StatusOK, w, "Otp code resent successfully", nil)

}
