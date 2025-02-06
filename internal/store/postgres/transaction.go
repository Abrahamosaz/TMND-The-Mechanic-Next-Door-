package postgres

import (
	"github.com/Abrahamosaz/TMND/internal/models"
	"gorm.io/gorm"
)


type TransactionRepository struct {
	DB *gorm.DB
}


func (trxRepo *TransactionRepository) Create(tx *gorm.DB,  trx *models.Transaction) (error) {

	if  err := tx.Create(trx).Error; err != nil {
		return err
	}
	return  nil
}