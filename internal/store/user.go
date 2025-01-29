package store

import (
	"gorm.io/gorm"
)



type UserRepository struct {
	db *gorm.DB
}


func (userRepo *UserRepository) Create() error {
	return nil

}