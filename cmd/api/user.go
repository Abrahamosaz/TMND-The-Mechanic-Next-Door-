package main

import (
	"log"
	"net/http"

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

	profileInfo, ok := app.GetProfileInfoFromContext(r)

	if !ok {
		log.Printf("Failed to get profile information")
		app.responseJSON(http.StatusInternalServerError, w,  "internal server error", nil)
		return
	}

	// var editProfile services.EditProfile
	
	err := r.ParseMultipartForm(10 << 20) // 10MB max size

	if err != nil {
		log.Printf("Error parsing form")
		app.responseJSON(http.StatusInternalServerError, w,  "internal server error", nil)
		return
	}

	editProfile := services.EditProfile{
		Name: r.FormValue("fullName"),
		Address: utils.StringPtr(r.FormValue("address")),
		PhoneNumber: utils.StringPtr(r.FormValue("phoneNumber")),
		Location:  utils.StringPtr(r.FormValue("location")),
	}

	err = validate.Struct(editProfile)

	if err != nil {
		ValidateRequestBody(err, w)
		return
	}

	serviceApp := app.createNewServiceApp()
	statusCode, err := services.EditUserProfile(&serviceApp, user.ID, services.EditProfileInfo{
		UserProfile: editProfile,
		URL: profileInfo.URL,
		FileName: profileInfo.FileName,
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