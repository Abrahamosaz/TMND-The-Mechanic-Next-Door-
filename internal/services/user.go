package services

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Abrahamosaz/TMND/internal/models"
	"github.com/Abrahamosaz/TMND/internal/utils"
)




func EditUserProfile(app *Application, user *models.User, editInfo EditProfileInfo) (int, error) {

	if (user.PublicId != nil && user.ProfileImageUrl != nil) {
		cloudinaryURL := os.Getenv("CLOUDINARY_URL")
		cld := &utils.Cloudinary{URL: cloudinaryURL}

		folder := fmt.Sprintf("%s/%s", utils.CLOUDINARY_PROFILE_IMAGE_FOLDER, user.ID)
		publicId := fmt.Sprintf("%s/%s", folder, *user.PublicId)
		go cld.DeleteFileFromCloudinary(publicId)
	}

	err := app.Store.User.Update(models.User{
		ID: user.ID,
		FullName: editInfo.UserProfile.FullName,
		PhoneNumber: editInfo.UserProfile.PhoneNumber,
		Address: editInfo.UserProfile.Address,
		State: editInfo.UserProfile.State,
		Lga: editInfo.UserProfile.Lga,
		ProfileFileName: &editInfo.FileName,
		ProfileImageUrl: &editInfo.URL,
		PublicId: &editInfo.PublicId,
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