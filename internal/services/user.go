package services

import (
	"net/http"

	"github.com/Abrahamosaz/TMND/internal/models"
	"github.com/google/uuid"
)




func EditUserProfile(app *Application, userId uuid.UUID, editInfo EditProfileInfo) (int, error) {
	err := app.Store.User.Update(models.User{
		ID: userId,
		FullName: editInfo.UserProfile.Name,
		PhoneNumber: editInfo.UserProfile.PhoneNumber,
		Address: editInfo.UserProfile.Address,
		Location: editInfo.UserProfile.Location,
		ProfileFileName: &editInfo.FileName,
		ProfileImageUrl: &editInfo.URL,
	})

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil	
}