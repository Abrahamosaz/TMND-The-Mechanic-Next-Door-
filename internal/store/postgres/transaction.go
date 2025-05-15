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


func (trxRepo *TransactionRepository) Update(tx *gorm.DB, trx *models.Transaction) (error) {
	if err := tx.Model(&models.Transaction{}).Where("id = ?", trx.ID).Updates(trx).Error; err != nil {
		return err
	}
	return nil
}


func (trxRepo *TransactionRepository) GetTransactionByTrxRef(trxRef string) (*models.Transaction, error) {
	var trx models.Transaction
	if err := trxRepo.DB.Where("trx_ref = ?", trxRef).First(&trx).Error; err != nil {
		return nil, err
	}
	return &trx, nil
}


func (trxRepo *TransactionRepository) GetUserTransactions(user *models.User, pgQuery *models.PaginationQuery) (*models.PaginationResponse[models.Transaction], error) {
	var transactions []models.Transaction
	var total int64

	query := trxRepo.DB.Model(&models.Transaction{}).Where("user_id = ?", user.ID).Order("created_at DESC")
	
	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	
	if (pgQuery.Limit != nil && pgQuery != nil) {
		//handle pagination if the page and limt is null
		offset := (*pgQuery.Page - 1) * *pgQuery.Limit
		result := query.Limit(*pgQuery.Limit).Offset(offset).Find(&transactions)

		if result.Error != nil {
			return nil, result.Error
		}

		// Build the response
		totalPages := int((total + int64(*pgQuery.Limit) - 1) / int64(*pgQuery.Limit))
		response := models.PaginationResponse[models.Transaction]{
			Data:       transactions,
			Total:      &total,
			Page:       pgQuery.Page,
			Limit:      pgQuery.Limit,
			TotalPages: &totalPages, // Ceil division
		}

		return &response, nil
	}

	//get all transactions without pagination
	if err := query.Find(&transactions).Error; err != nil {
		return nil, err
	}

	return &models.PaginationResponse[models.Transaction]{Data: transactions, Total: &total}, nil
}