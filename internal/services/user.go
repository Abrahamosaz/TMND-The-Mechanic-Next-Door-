package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Abrahamosaz/TMND/internal/models"
	"github.com/Abrahamosaz/TMND/internal/providers"
	"github.com/Abrahamosaz/TMND/internal/utils"
	"gorm.io/gorm"
)




func (app *Application) EditUserProfile(user *models.User, editInfo EditProfileInfo) (int, error) {

	if (user.PublicId != nil && user.ProfileImageUrl != nil) {
		cloudinaryURL := os.Getenv("CLOUDINARY_URL")
		cld := &utils.Cloudinary{URL: cloudinaryURL}

		folder := fmt.Sprintf("%s/%s", utils.CLOUDINARY_PROFILE_IMAGE_FOLDER, user.ID)
		publicId := fmt.Sprintf("%s/%s", folder, *user.PublicId)
		go cld.DeleteFileFromCloudinary(publicId)
	}

	err := app.UpdateUser(&models.User{
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


func (app *Application) GetUserTransaction(user *models.User, qs *models.PaginationQuery) (*models.PaginationResponse[models.Transaction], int, error) {

	transactions, err := app.Store.Transaction.GetUserTransactions(user, qs)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return transactions, http.StatusOK, nil
}

func (app *Application) CreateInvoice(user *models.User, fundAccount FundAccount) (map[string]any, int, error) {
    m := providers.Monnify{
        Url: os.Getenv("MONNIFY_BASE_URL"),
        ApiKey: os.Getenv("MONNIFY_API_KEY"),
        SecretKey: os.Getenv("MONNIFY_SECRET_KEY"),
    }

	trxRef, err := utils.GenerateUniqueTrxRef("CREDIT")
    if err != nil {
        return nil, http.StatusInternalServerError, err
    }

	expiryDate := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")
	result, err := m.CreateInvoice(&providers.CreateInvoice{
		Amount: fundAccount.Amount,
		CurrencyCode: "NGN",
		Reference: trxRef,
		CustomerName: user.FullName,
		CustomerEmail: user.Email,
		ContractCode: os.Getenv("MONNIFY_CONTRACT_CODE"),
		Description: fundAccount.Description,
		ExpiryDate: expiryDate,
		RedirectUrl: fundAccount.RedirectUrl,
		PaymentMethod: []string{"ACCOUNT_TRANSFER", "CARD"},
	})
	
	if err != nil {
		fmt.Println("error creating invoice: ", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	jsonResult, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println("result from monnify create invoice: ", string(jsonResult))

	responseBody, ok := result["responseBody"].(map[string]any)
	
	if !ok {
		return nil, http.StatusInternalServerError, errors.New("invalid response body")
	}

	tx := app.Store.BeginTransaction()
	//create transaction
	err = app.CreateNewTransaction(tx, &models.Transaction{
        TrxRef: responseBody["invoiceReference"].(string),
        UserID: &user.ID,
        PreviousBalance: user.Balance,
        CurrentBalance: user.Balance,
        Amount: fundAccount.Amount,
		Type: models.TransactionTypeCredit,
		Status: models.StatusPending,
		Description: &fundAccount.Description,
    })


	if err != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, err
	}

	tx.Commit()
	return responseBody, http.StatusOK, nil
}


func (app *Application) ConfirmUserPayment(user *models.User, paymentReference string) (int, error) {

	m := providers.Monnify{
        Url: os.Getenv("MONNIFY_BASE_URL"),
        ApiKey: os.Getenv("MONNIFY_API_KEY"),
        SecretKey: os.Getenv("MONNIFY_SECRET_KEY"),
    }

	err := m.ConfirmInvoicePayment(paymentReference)

	if err != nil {
		fmt.Println("error confirming invoice payment: ", err.Error())
		return http.StatusInternalServerError, err
	}
	

	pendingTrx, err := app.GetTransactionByTrxRef(paymentReference)
	
	if err != nil {
		fmt.Println("error getting transaction: ", err.Error())
		return http.StatusNotFound, err
	}

	if pendingTrx.Status == models.StatusSuccess {
		return http.StatusBadRequest, errors.New("invoice already confirmed")
	}
	
	currentBalance := user.Balance + pendingTrx.Amount
	
	pendingTrx.Status = models.StatusSuccess
	pendingTrx.CurrentBalance = currentBalance
	user.Balance = currentBalance
	
	tx := app.Store.BeginTransaction()

	err = app.UpdateTransaction(tx, pendingTrx)

	if err != nil {
		tx.Rollback()
		fmt.Println("error updating transaction: ", err.Error())
		return http.StatusInternalServerError, err
	}

	err = app.UpdateUserTrx(tx, user)

	if err != nil {
		tx.Rollback()
		fmt.Println("error updating user: ", err.Error())
		return http.StatusInternalServerError, err
	}
	
	tx.Commit()
	return http.StatusOK, nil
}


func (app *Application) UpdateUserTrx(tx *gorm.DB, user *models.User) (error) {
	return app.Store.User.TrxUpdate(tx, user)
}

func (app *Application) UpdateUser(user *models.User) (error) {
	return app.Store.User.Update(*user)
}
