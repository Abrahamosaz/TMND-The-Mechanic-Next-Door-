package services

import (
	"errors"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Abrahamosaz/TMND/internal/models"
	"github.com/Abrahamosaz/TMND/internal/utils"
)

// create new user
func CreateUser(app *Application, payload Signup) (models.User, int, error) {

	user, tx, err := app.Store.User.Create(models.User{
		FullName: payload.FullName,
		Email: payload.Email,
		PhoneNumber: &payload.PhoneNumber,
		Password: payload.Password,
		RegisterWithGoogle: payload.RegisterWithGoogle,
	})

	if err != nil {
		tx.Rollback()
		statusCode := http.StatusInternalServerError
		if (err.Error() == "user with email already exists") {
			statusCode = http.StatusNotFound
		}
		return models.User{}, statusCode, err
	}
	
	// generate otpCode
	otpCode, err := utils.GenerateOtpCode(4)

	if err != nil  {
		tx.Rollback()
		return models.User{}, http.StatusInternalServerError, errors.New("failed to generate otp code")
	}

	//genarate jwt token
	expireTime := 15 * time.Minute
	token, err :=  utils.GenerateJWT(user.ID.String(), expireTime, utils.PayloadClaims{
		OtpCode: &otpCode,
	})

	if err != nil {
		tx.Rollback()
		return models.User{}, http.StatusInternalServerError, errors.New("failed to generate jwt token")
	}

	user.OtpToken = &token
	// send email message
	go func() {
		if err := app.sendOtpEmail("Welcome! Please Verify Your Email to Get Started", user.Email, OtpEmailPayload{
			fullName: user.FullName,
			otpCode: otpCode,
		}, "signup"); err != nil {
			log.Println("Error sending email:", err)
			panic(err)
		} else {
			log.Println("Email sent successfully!")
		}
	}()

	result := tx.Save(&user)

	if result.Error != nil {
		tx.Rollback()
		return models.User{}, http.StatusInternalServerError, result.Error
	}

	if err := tx.Commit().Error; err != nil {
		return models.User{}, http.StatusInternalServerError, err
	}

	return user, http.StatusOK, nil
}



// login user
func LoginUser(app *Application, payload Login) (int, error) {
	email := payload.Email
	password := payload.Password
	user, err := app.Store.User.FindByEmail(email)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user with email not found" {
			statusCode = http.StatusNotFound
		}
		return statusCode,  err
	}


	// check if the user registered with google auth
	if user.RegisterWithGoogle {
		return http.StatusBadRequest, errors.New("this account was registered using Google. Please log in with Google instead")
	}

	//check if the password is correct
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return http.StatusBadRequest, errors.New("incorrect email or password")
	}

		// generate otpCode
	otpCode, err := utils.GenerateOtpCode(4)

	if err != nil  {
		return http.StatusInternalServerError, errors.New("failed to generate otp code")
	}

	//genarate jwt token
	expireTime := 15 * time.Minute
	token, err :=  utils.GenerateJWT(user.ID.String(), expireTime, utils.PayloadClaims{
		OtpCode: &otpCode,
	})
	

	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to generate jwt token")
	}

	//update the user token
	user.OtpToken = &token
	// check if the user email is verified
	// send verification email to the user if the email is not verified
	if !user.IsEmailVerified {
		go func() {
			if err := app.sendOtpEmail("Your Account is Almost Ready! Verify Your Email Now.", user.Email, OtpEmailPayload{
				fullName: user.FullName,
				otpCode: otpCode,
			}, "login"); err != nil {
				log.Println("Error sending email:", err)
				panic(err)
			} else {
				log.Println("Email sent successfully!")
			}
		}()
	}

	err = app.Store.User.Save(&user)

	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to save user")
	}

	return http.StatusOK,  nil
}



func ForgotPassword(app  *Application, payload Email) (int, error) {
	email := payload.Email
	user, err := app.Store.User.FindByEmail(email)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user with email not found" {
			statusCode = http.StatusNotFound
		}
		return statusCode, err
	}

	// generate otpCode
	otpCode, err := utils.GenerateOtpCode(4)

	if err != nil  {
		return http.StatusInternalServerError, errors.New("failed to generate otp code")
	}

	//genarate jwt token
	expireTime := 15 * time.Minute
	token, err :=  utils.GenerateJWT(user.ID.String(), expireTime, utils.PayloadClaims{
		OtpCode: &otpCode,
	})
	

	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to generate jwt token")
	}

	//update the user token
	user.OtpToken = &token

	go func() {
		if err := app.sendOtpEmail("Reset Your Password", user.Email, OtpEmailPayload{
			fullName: user.FullName,
			otpCode: otpCode,
		}, "forgotPassword"); err != nil {
			log.Println("Error sending email:", err)
			panic(err)
		} else {
			log.Println("Email sent successfully!")
		}
	}()

	err = app.Store.User.Save(&user)

	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to save user")
	}
	return http.StatusOK,  nil

}


func ResetPassword(app *Application, payload ChangePassword) (int, error) {
	email := payload.Email
	password := payload.Password
	otpCode := payload.OtpCode

	user, err := app.Store.User.FindByEmail(email)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user with email not found" {
			statusCode = http.StatusNotFound
		}
		return statusCode, err
	}

	token := user.OtpToken

	if token == nil {
		return  http.StatusNotAcceptable, errors.New("please request otpcode")
	}

	//decode otp token
	claims, err := utils.DecodeJWT(*token)

	if err != nil {
		return  http.StatusBadRequest, errors.New("invalid otp code")
	}

	if otpCode != *claims.OtpCode {
		return http.StatusBadRequest, errors.New("invalid otp code")
	}

	//set the otpToken to nill
	user.OtpToken = nil
	user.IsEmailVerified = true

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to hash password")
	}

	user.Password = string(hashedPassword)

	err = app.Store.User.Save(&user)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}


func VerifyEmail(app *Application, payload VerifyOtp) (VerifyOtpResponse, int, error) {
	email := payload.Email
	otpCode := payload.OtpCode

	user, err := app.Store.User.FindByEmail(email)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user with email not found" {
			statusCode = http.StatusNotFound
		}
		return VerifyOtpResponse{}, statusCode, err
	}

	token := user.OtpToken

	if token == nil {
		return  VerifyOtpResponse{}, http.StatusNotAcceptable, errors.New("please login and request another otpcode")
	}

	//decode otp token
	claims, err := utils.DecodeJWT(*token)

	if err != nil {
		return VerifyOtpResponse{}, http.StatusBadRequest, errors.New("invalid otp code")
	}

	if otpCode != *claims.OtpCode {
		return VerifyOtpResponse{}, http.StatusBadRequest, errors.New("invalid otp code")
	}

	//set the otpToken to nill
	user.OtpToken = nil
	user.IsEmailVerified = true
	err = app.Store.User.Save(&user)
	if err != nil {
		return VerifyOtpResponse{}, http.StatusInternalServerError, err
	}
	
	return VerifyOtpResponse{User: user, Token: *token}, http.StatusOK, nil
}


func VerifyOtpCode(app *Application, payload VerifyOtp) (int, error) {
	email := payload.Email
	otpCode := payload.OtpCode

	user, err := app.Store.User.FindByEmail(email)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user with email not found" {
			statusCode = http.StatusNotFound
		}
		return statusCode, err
	}

	token := user.OtpToken

	if token == nil {
		return http.StatusNotAcceptable, errors.New("please request for otpcode")
	}

	//decode otp token
	claims, err := utils.DecodeJWT(*token)

	if err != nil {
		return http.StatusBadRequest, errors.New("invalid otp code")
	}

	if otpCode != *claims.OtpCode {
		return http.StatusBadRequest, errors.New("invalid otp code")
	}
	
	return http.StatusOK, nil
}


func GoolgeAuth(app *Application) error {
	return nil
}




