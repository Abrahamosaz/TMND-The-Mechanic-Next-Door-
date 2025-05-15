package services

import (
	"github.com/Abrahamosaz/TMND/internal/models"
	"gorm.io/gorm"
)





func (app *Application) CreateNewTransaction(tx *gorm.DB, txData *models.Transaction) error {
	return app.Store.Transaction.Create(tx, txData)
}


func (app *Application) UpdateTransaction(tx *gorm.DB, txData *models.Transaction) error {
	return app.Store.Transaction.Update(tx, txData)
}


func (app *Application) GetTransactionByTrxRef(trxRef string) (*models.Transaction, error) {
	return app.Store.Transaction.GetTransactionByTrxRef(trxRef)
}
