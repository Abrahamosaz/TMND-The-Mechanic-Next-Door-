package services

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/Abrahamosaz/TMND/internal/models"
	"github.com/Abrahamosaz/TMND/internal/providers"
)


func (app *Application) HandleMonnifyWebhook(eventData map[string]any) error {

	eventType := eventData["eventType"].(string)
	eventData, ok := eventData["eventData"].(map[string]any)

	if !ok {
		return errors.New("invalid event data")
	}

	switch eventType {
	case "SUCCESSFUL_TRANSACTION":
		return app.handleSuccessfulTransaction(eventData)
	default:
		fmt.Println("invalid event type", eventType)
		return errors.New("invalid event type")
	}
}



func (app *Application) handleSuccessfulTransaction(eventData map[string]any) error {
	paymentReference := eventData["paymentReference"].(string)
	settlementAmount, err := strconv.ParseFloat(eventData["settlementAmount"].(string), 64)
	if err != nil {
		return fmt.Errorf("invalid settlement amount: %w", err)
	}
	
	// confirm transaction reference
	m := providers.Monnify{
		Url: os.Getenv("MONNIFY_BASE_URL"),
		ApiKey: os.Getenv("MONNIFY_API_KEY"),
		SecretKey: os.Getenv("MONNIFY_SECRET_KEY"),
	}

	err = m.ConfirmInvoicePayment(paymentReference)

	if err != nil {
		fmt.Println("failed to confirm transaction", err)
		return fmt.Errorf("failed to confirm transaction: %w", err)
	}

	pendingTrx, err := app.GetTransactionByTrxRef(paymentReference)

	if err != nil {
		return err
	}

	if pendingTrx.Status != models.StatusPending {
		return errors.New("transaction already confirmed")
	}

	user, err := app.Store.User.FindByID(pendingTrx.UserID.String())
	if err != nil {
		return err
	}

	newBalance := user.Balance + settlementAmount
	fee := pendingTrx.Amount - settlementAmount
	

	// update transaction
	pendingTrx.Fee = &fee
	pendingTrx.Status = models.StatusSuccess
	pendingTrx.CurrentBalance = newBalance


	// update user balance
	user.Balance = newBalance

	trx := app.Store.BeginTransaction()

	err = app.UpdateTransaction(trx, pendingTrx)

	if err != nil {
		trx.Rollback()
		return err
	}
	
	err = app.UpdateUserTrx(trx, &user)

	if err != nil {
		trx.Rollback()
		return err
	}

	trx.Commit()
	return nil	
}