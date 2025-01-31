package main

import (
	"net/http"
)



func (app  *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	user, ok := app.GetUserFromContext(r)
	
	if !ok {
		app.responseJSON(http.StatusUnauthorized, w,  "Unauthorized: No user found", nil)
		return
	}
	app.responseJSON(http.StatusOK, w, "User details retrieve successfully", user)

}