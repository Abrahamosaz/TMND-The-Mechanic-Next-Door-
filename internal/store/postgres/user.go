package postgres

import (
	"errors"

	"github.com/Abrahamosaz/TMND/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)



type UserRepository struct {
	DB *gorm.DB
}


func (userRepo *UserRepository) Create(tx *gorm.DB, user models.User) (models.User, error) {	
	// Ensure the transaction is rolled back if there’s an error
	if tx.Error != nil {
		return models.User{}, tx.Error
	}

	// Basic FOR UPDATE lock
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("email = ?", user.Email).Find(&models.User{})

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{},  result.Error
		}
	} else if result.RowsAffected > 0 {
		return models.User{}, errors.New("user with email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, errors.New("failed to hash password")
	}

	user.Password = string(hashedPassword)

	// Create new user
	if err := tx.Create(&user).Error; err != nil {
		return models.User{},  err
	}

	return user, nil
}


func (userRepo *UserRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	result := userRepo.DB.Where("email = ?", email).First(&user) 

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user with email not found")
		}
		return models.User{}, result.Error
	}
	return user, nil
}

func (userRepo *UserRepository) Save(user *models.User) (error) {
	result := userRepo.DB.Save(&user)
	if result.Error != nil {
		return result.Error // Return the error if save failed
	}
	return nil
}

func (userRepo *UserRepository) Update(user models.User) (error) {
	result := userRepo.DB.Model(&user).Updates(user)
	// Check for errors
	if result.Error != nil {
		return result.Error
	}

	return nil
}



func (userRepo *UserRepository) TrxUpdate(tx *gorm.DB, user *models.User) (error) {
	result := tx.Model(&user).Updates(user)
	// Check for errors
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (userRepo *UserRepository) FindByID(id string) (models.User, error) {

	var user models.User
	result  := userRepo.DB.First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user with id not found")
		}
		return models.User{}, result.Error
	}
	return user, nil
}

func (userRepo *UserRepository) DeductFromBalance(tx *gorm.DB, user *models.User, amount float64) error {
	if err := tx.Model(user).Update("balance", amount).Error; err != nil {
		return err
	}
	return nil
}