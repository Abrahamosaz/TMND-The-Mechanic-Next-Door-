package services

import (
	"github.com/Abrahamosaz/TMND/internal/models"
	"gorm.io/gorm"
)





func CreateNewTransaction(app  *Application, tx *gorm.DB, txData *models.Transaction) error {
	return app.Store.Transaction.Create(tx, txData)
}