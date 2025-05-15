package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Abrahamosaz/TMND/internal/models"
	"github.com/Abrahamosaz/TMND/internal/services"
	"github.com/Abrahamosaz/TMND/internal/utils"
	"github.com/go-chi/chi/v5"
)



func (app  *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	user, ok := app.GetUserFromContext(r)

	if !ok {
		app.responseJSON(http.StatusUnauthorized, w,  "Unauthorized: No user found", nil)
		return
	}
	app.responseJSON(http.StatusOK, w, "User details retrieve successfully", user)

}


func (app *application) editUserProfileHandler(w http.ResponseWriter, r *http.Request) {

	user, ok := app.GetUserFromContext(r)

	if !ok {
		app.responseJSON(http.StatusUnauthorized, w,  "Unauthorized: No user found", nil)
		return
	}

	profileURL, fileName, publicId, _ := app.GetFileInfoFromContext(r)

	// var editProfile services.EditProfile
	err := r.ParseMultipartForm(10 << 20) // 10MB max size

	if err != nil {
		log.Printf("Error parsing form")
		app.responseJSON(http.StatusInternalServerError, w,  "internal server error", nil)
		return
	}

	editProfile := services.EditProfile{
		FullName: r.FormValue("fullName"),
		PhoneNumber: r.FormValue("phoneNumber"),
		Address: utils.StringToPtr(r.FormValue("address")),
		State:  utils.StringToPtr(r.FormValue("state")),
		Lga: utils.StringToPtr(r.FormValue("lga")),
	}

	err = validate.Struct(editProfile)

	if err != nil {
		ValidateRequestBody(err, w)
		return
	}

	serviceApp := app.createNewServiceApp()

	var editProfileInfo services.EditProfileInfo
	editProfileInfo.UserProfile = editProfile

	if profileURL != nil && fileName != nil && publicId != nil {
		editProfileInfo.URL = *profileURL
		editProfileInfo.FileName = *fileName
		editProfileInfo.PublicId = *publicId
	}

	statusCode, err := serviceApp.EditUserProfile(user, editProfileInfo)

	if err != nil {
		log.Println("error editing user profile: ", err.Error())
		message := err.Error()
		if statusCode == http.StatusInternalServerError {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}
	app.responseJSON(statusCode, w, "User profile updated successfully", nil)
}



func (app *application) getUserTransactionHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := app.GetUserFromContext(r)

	if !ok {
		app.responseJSON(http.StatusUnauthorized, w,  "Unauthorized: No user found", nil)
		return
	}

	//get query strings
	qs := r.URL.Query()

	page := qs.Get("page")
	if page == "" {
		page = "1"
	}
	
	limit := qs.Get("limit")
	if limit == "" {
		limit = "10"
	}

	serviceApp := app.createNewServiceApp()
	trxs, statusCode, err := serviceApp.GetUserTransaction(
		user,
		&models.PaginationQuery{Page: utils.ConvertStrToPtrInt(page), Limit: utils.ConvertStrToPtrInt(limit)},
	)

	if err != nil {
		log.Println("error getting user transactions: ", err.Error())
		message := err.Error()
		if statusCode == http.StatusInternalServerError {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}
	app.responseJSON(statusCode, w, "User transactions retrieved successfully", trxs)

}


func (app *application) createInvoiceHandler(w http.ResponseWriter, r *http.Request) {

	var fundAccount services.FundAccount

	err := json.NewDecoder(r.Body).Decode(&fundAccount)

	if err != nil {
		app.responseJSON(http.StatusBadRequest, w, "Invalid request body", nil)
		return
	}

	err = validate.Struct(fundAccount)

	if err != nil {
		ValidateRequestBody(err, w)
		return
	}

	user, ok := app.GetUserFromContext(r)

	if !ok {
		app.responseJSON(http.StatusUnauthorized, w,  "Unauthorized: No user found", nil)
		return
	}

	serviceApp := app.createNewServiceApp()

	result, statusCode, err := serviceApp.CreateInvoice(user, fundAccount)

	if err != nil {
		log.Println("error creating invoice: ", err.Error())
		message := err.Error()
		if statusCode == http.StatusInternalServerError {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}
	app.responseJSON(statusCode, w, "Invoice created successfully", result)
	
}


func (app *application) confirmPaymentHandler(w http.ResponseWriter, r *http.Request) {
	paymentReference := chi.URLParam(r, "paymentReference")

	user, ok := app.GetUserFromContext(r)

	if !ok {
		app.responseJSON(http.StatusUnauthorized, w,  "Unauthorized: No user found", nil)
		return
	}

	serviceApp := app.createNewServiceApp()

	statusCode, err := serviceApp.ConfirmUserPayment(user, paymentReference)

	if err != nil {
		log.Println("error confirming payment: ", err.Error())
		app.responseJSON(statusCode, w, err.Error(), nil)
		return
	}

	app.responseJSON(statusCode, w, "Payment confirmed successfully", nil)
}