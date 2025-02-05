package store

import (
	"github.com/Abrahamosaz/TMND/internal/models"
	"github.com/Abrahamosaz/TMND/internal/store/postgres"
	"gorm.io/gorm"
)


type Storage struct {
	BeginTransaction func() *gorm.DB
	User interface {
		Create(*gorm.DB, models.User) (models.User, error)
		FindByEmail(string) (models.User, error)
		FindByID(string) (models.User, error)
		Save(*models.User) error
		Update(models.User) error
	}
	Mechanic interface {
		Create(*gorm.DB, models.Mechanic) (models.Mechanic, error)
		GetAllAvailableMechanics() (*[]models.Mechanic, error)
	}
	Booking interface {
		GetPendingBookings() (*[]models.Booking, error)
		Create(*gorm.DB, models.Booking) (models.Booking, error)
		GetBookingFee() (models.BookingFee, error)
	}
}


func PostgresStorage(db *gorm.DB) Storage {
	return Storage {
		BeginTransaction: func() *gorm.DB {
			return db.Begin()
		},
		User: &postgres.UserRepository{DB: db},
		Mechanic: &postgres.MechanicRepository{DB: db},
		Booking: &postgres.BookingRepository{DB: db},
	}
}