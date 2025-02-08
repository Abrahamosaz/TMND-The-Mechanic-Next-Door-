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

// ----- USER SECTION ---------
// create new user
func CreateUser(app *Application, payload Signup) (models.User, int, error) {
	// Start a transaction manually
	tx := app.Store.BeginTransaction()
	
	user, err := app.Store.User.Create(tx, models.User{
		FullName: payload.FullName,
		Email: payload.Email,
		PhoneNumber: payload.PhoneNumber,
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

	sendEmail := func(subject string) {
		if err := app.sendOtpEmail(subject, user.Email, OtpEmailPayload{
			fullName: user.FullName,
			otpCode: otpCode,
		}, "login"); err != nil {
			log.Println("Error sending email:", err)
			panic(err)
		} else {
			log.Println("Email sent successfully!")
		}
	}
	//update the user token
	user.OtpToken = &token
	// check if the user email is verified
	// send verification email to the user if the email is not verified
	if !user.IsEmailVerified {
		go sendEmail("Your Account is Almost Ready! Verify Your Email Now.")
	} else {
		go sendEmail("Your One-Time Password (OTP) for Login")	
	}

	err = app.Store.User.Save(&user)

	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to save user")
	}

	return http.StatusOK, nil
}

func SendUserOtpCode(app  *Application, payload Email, route string) (int, error) {
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
	sendEmail := func(subject string) {
		if err := app.sendOtpEmail(subject, user.Email, OtpEmailPayload{
			fullName: user.FullName,
			otpCode: otpCode,
		}, "forgotPassword"); err != nil {
			log.Println("Error sending email:", err)
			panic(err)
		} else {
			log.Println("Email sent successfully!")
		}
	}

	if route == "forgotPassword" {
		go sendEmail("Reset Your Password")
	} else {
		message := "Your New OTP Code for Secure Login"
		if !user.IsEmailVerified {
			message = "New OTP Code: Verify Your Email to Complete Registration"
		}
		go sendEmail(message)
	}

	err = app.Store.User.Save(&user)

	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to save user")
	}
	return http.StatusOK,  nil

}


func ResetUserPassword(app *Application, payload ChangePassword) (int, error) {
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


func VerifyUserEmail(app *Application, payload VerifyOtp) (VerifyUserOtpResponse, int, error) {
	email := payload.Email
	otpCode := payload.OtpCode

	user, err := app.Store.User.FindByEmail(email)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user with email not found" {
			statusCode = http.StatusNotFound
		}
		return VerifyUserOtpResponse{}, statusCode, err
	}

	token := user.OtpToken

	if token == nil {
		return  VerifyUserOtpResponse{}, http.StatusNotAcceptable, errors.New("please login and request another otpcode")
	}

	//decode otp token
	claims, err := utils.DecodeJWT(*token)

	if err != nil {
		return VerifyUserOtpResponse{}, http.StatusBadRequest, errors.New("invalid otp code")
	}

	if otpCode != *claims.OtpCode {
		return VerifyUserOtpResponse{}, http.StatusBadRequest, errors.New("invalid otp code")
	}

	type TokenPayload struct {
		token *string
		errorMessage error
	}

	tokenChannel := make(chan TokenPayload)

	go func() {
		// generate otpCode
		otpCode, err := utils.GenerateOtpCode(4)

		if err != nil  {
			tokenChannel <- TokenPayload{nil, errors.New("failed to generate otp code")}
			return
		}

		//genarate jwt token
		expireTime := 30 * 24 * time.Hour
		token, err :=  utils.GenerateJWT(user.ID.String(), expireTime, utils.PayloadClaims{
			OtpCode: &otpCode,
			Role: "USER",
		})
		
		if err != nil {
			tokenChannel <- TokenPayload{nil, errors.New("failed to generate jwt token")}
			return
		}

		tokenChannel <- TokenPayload{&token, nil}
	}()

	//set the otpToken to nill
	user.OtpToken = nil
	user.IsEmailVerified = true
	err = app.Store.User.Save(&user)

	if err != nil {
		return VerifyUserOtpResponse{}, http.StatusInternalServerError, err
	}

	tokenResponse := <-tokenChannel

	if tokenResponse.errorMessage != nil {
		return VerifyUserOtpResponse{}, http.StatusInternalServerError, tokenResponse.errorMessage
	}
	
	return VerifyUserOtpResponse{User: user, Token: *tokenResponse.token}, http.StatusOK, nil
}


func VerifyUserOtpCode(app *Application, payload VerifyOtp) (int, error) {
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



// -------- MECHANIC SECTION ----------
// login mechanic
func LoginMechanic(app *Application, payload Login) (int, error) {
	email := payload.Email
	password := payload.Password
	mechanic, err := app.Store.Mechanic.FindByEmail(email)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "mechanic with email not found" {
			statusCode = http.StatusNotFound
		}
		return statusCode,  err
	}

	//check if the password is correct
	if err = bcrypt.CompareHashAndPassword([]byte(mechanic.Password), []byte(password)); err != nil {
		return http.StatusBadRequest, errors.New("incorrect email or password")
	}

	// generate otpCode
	otpCode, err := utils.GenerateOtpCode(4)

	if err != nil  {
		return http.StatusInternalServerError, errors.New("failed to generate otp code")
	}

	//genarate jwt token
	expireTime := 15 * time.Minute
	token, err :=  utils.GenerateJWT(mechanic.ID.String(), expireTime, utils.PayloadClaims{
		OtpCode: &otpCode,
	})
	

	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to generate jwt token")
	}

	sendEmail := func(subject string) {
		if err := app.sendOtpEmail(subject, mechanic.Email, OtpEmailPayload{
			fullName: mechanic.FullName,
			otpCode: otpCode,
		}, "login"); err != nil {
			log.Println("Error sending email:", err)
			panic(err)
		} else {
			log.Println("Email sent successfully!")
		}
	}
	//update the mechanic token
	mechanic.OtpToken = &token
	// check if the mechainc email is verified
	// send verification email to the mechanic if the email is not verified
	if !mechanic.IsEmailVerified {
		go sendEmail("Your Account is Almost Ready! Verify Your Email Now.")
	} else {
		go sendEmail("Your One-Time Password (OTP) for Login")	
	}

	err = app.Store.Mechanic.Save(&mechanic)

	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to save mechanic")
	}

	return http.StatusOK, nil
}

// verify email
func VerifyMechanicEmail(app *Application, payload VerifyOtp) (VerifyMechanicOtpResponse, int, error) {
	email := payload.Email
	otpCode := payload.OtpCode

	mechanic, err := app.Store.Mechanic.FindByEmail(email)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "mechanic with email not found" {
			statusCode = http.StatusNotFound
		}
		return VerifyMechanicOtpResponse{}, statusCode, err
	}

	token := mechanic.OtpToken

	if token == nil {
		return  VerifyMechanicOtpResponse{}, http.StatusNotAcceptable, errors.New("please login and request another otpcode")
	}

	//decode otp token
	claims, err := utils.DecodeJWT(*token)

	if err != nil {
		return VerifyMechanicOtpResponse{}, http.StatusBadRequest, errors.New("invalid otp code")
	}

	if otpCode != *claims.OtpCode {
		return VerifyMechanicOtpResponse{}, http.StatusBadRequest, errors.New("invalid otp code")
	}

	type TokenPayload struct {
		token *string
		errorMessage error
	}

	tokenChannel := make(chan TokenPayload)

	go func() {
		// generate otpCode
		otpCode, err := utils.GenerateOtpCode(4)

		if err != nil  {
			tokenChannel <- TokenPayload{nil, errors.New("failed to generate otp code")}
			return
		}

		//genarate jwt token
		expireTime := 30 * 24 * time.Hour
		token, err :=  utils.GenerateJWT(mechanic.ID.String(), expireTime, utils.PayloadClaims{
			OtpCode: &otpCode,
			Role: "MECHANIC",
		})
		
		if err != nil {
			tokenChannel <- TokenPayload{nil, errors.New("failed to generate jwt token")}
			return
		}

		tokenChannel <- TokenPayload{&token, nil}
	}()

	//set the otpToken to nill
	mechanic.OtpToken = nil
	mechanic.IsEmailVerified = true
	err = app.Store.Mechanic.Save(&mechanic)

	if err != nil {
		return VerifyMechanicOtpResponse{}, http.StatusInternalServerError, err
	}

	tokenResponse := <-tokenChannel

	if tokenResponse.errorMessage != nil {
		return VerifyMechanicOtpResponse{}, http.StatusInternalServerError, tokenResponse.errorMessage
	}
	
	return VerifyMechanicOtpResponse{Mechanic: mechanic, Token: *tokenResponse.token}, http.StatusOK, nil
}

