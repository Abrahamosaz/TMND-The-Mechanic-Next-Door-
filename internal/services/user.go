package services

import (
	"net/http"

	"github.com/Abrahamosaz/TMND/internal/models"
	"github.com/google/uuid"
)




func EditUserProfile(app *Application, userId uuid.UUID, editInfo EditProfileInfo) (int, error) {
	err := app.Store.User.Update(models.User{
		ID: userId,
		FullName: editInfo.UserProfile.FullName,
		PhoneNumber: editInfo.UserProfile.PhoneNumber,
		Address: editInfo.UserProfile.Address,
		State: editInfo.UserProfile.State,
		Lga: editInfo.UserProfile.Lga,
		ProfileFileName: &editInfo.FileName,
		ProfileImageUrl: &editInfo.URL,
	})

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil	
}


func GetUserTransaction(app *Application, user *models.User, qs *models.PaginationQuery) (*models.PaginationResponse[models.Transaction], int, error) {

	transactions, err := app.Store.Transaction.GetUserTransactions(user, qs)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return transactions, http.StatusOK, nil
}