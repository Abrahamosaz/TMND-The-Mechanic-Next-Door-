package main

import (
	"log"
	"net/http"

	"github.com/Abrahamosaz/TMND/internal/models"
	"github.com/Abrahamosaz/TMND/internal/services"
	"github.com/Abrahamosaz/TMND/internal/utils"
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

	profileInfo := app.GetProfileInfoFromContext(r)
	var profileURL = utils.StringToPtr(user.ProfileImageUrl)
	var fileName = utils.StringToPtr(user.ProfileFileName)
	
	if profileInfo != nil {
		profileURL = profileInfo.URL
		fileName = profileInfo.FileName
	}

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
		Address: utils.PtrToString(r.FormValue("address")),
		State:  utils.PtrToString(r.FormValue("state")),
		Lga: utils.PtrToString(r.FormValue("lga")),
	}

	err = validate.Struct(editProfile)

	if err != nil {
		ValidateRequestBody(err, w)
		return
	}

	serviceApp := app.createNewServiceApp()
	statusCode, err := services.EditUserProfile(&serviceApp, user.ID, services.EditProfileInfo{
		UserProfile: editProfile,
		URL: profileURL,
		FileName: fileName,
	})

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
	trxs, statusCode, err := services.GetUserTransaction(
		&serviceApp,
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