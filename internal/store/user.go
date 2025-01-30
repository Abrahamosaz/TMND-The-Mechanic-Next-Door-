package store

import (
	"errors"

	"github.com/Abrahamosaz/TMND/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)



type UserRepository struct {
	db *gorm.DB
}


func (userRepo *UserRepository) Create(user models.User) (models.User, *gorm.DB, error) {
	// Start a transaction manually
	tx := userRepo.db.Begin()

	// Ensure the transaction is rolled back if there’s an error
	if tx.Error != nil {
		return models.User{}, tx, tx.Error
	}

	// Basic FOR UPDATE lock
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("email = ?", user.Email).Find(&models.User{})

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, tx, result.Error
		}
	} else if result.RowsAffected > 0 {
		return models.User{}, tx, errors.New("user with email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, tx, errors.New("failed to hash password")
	}

	user.Password = string(hashedPassword)

	// Create new user
	if err := tx.Create(&user).Error; err != nil {
		return models.User{}, tx, err
	}

	return user, tx, nil
}


func (userRepo *UserRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	result := userRepo.db.Where("email = ?", email).First(&user) 

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user with email not found")
		}
		return models.User{}, result.Error
	}
	return user, nil
}

func (userRepo *UserRepository) Save(user *models.User) (error) {
	result := userRepo.db.Save(&user)
	if result.Error != nil {
		return result.Error // Return the error if save failed
	}
	return nil
}