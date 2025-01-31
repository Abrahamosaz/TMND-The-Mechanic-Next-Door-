package store

import (
	"github.com/Abrahamosaz/TMND/internal/models"
	"gorm.io/gorm"
)


type Storage struct {
	User interface {
		Create(models.User) (models.User, *gorm.DB, error)
		FindByEmail(string) (models.User, error)
		FindByID(string) (models.User, error)
		Save(*models.User) error
	}
	Post interface {
		Create() error
	}
}


func PostgresStorage(db *gorm.DB) Storage {
	return Storage {
		User: &UserRepository{db},
	}
}