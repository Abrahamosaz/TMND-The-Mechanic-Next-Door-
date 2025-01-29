package store

import (
	"gorm.io/gorm"
)


type Storage struct {
	User interface {
		Create() error
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