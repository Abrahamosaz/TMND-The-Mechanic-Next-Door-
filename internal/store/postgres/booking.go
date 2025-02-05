package postgres

import (
	"github.com/Abrahamosaz/TMND/internal/models"
	"gorm.io/gorm"
)



type BookingRepository struct {
	DB *gorm.DB
}


func (bookingRepo *BookingRepository) GetPendingBookings() (*[]models.Booking, error) {
	var bookings []models.Booking

	err := bookingRepo.DB.
		Where("status = ?", models.BookingPending). // Filter by pending status
		Find(&bookings).Error // Fetch all bookings

	if err != nil {
		return nil, err
	}

	return &bookings, nil
}


func (bookingRepo *BookingRepository) Create(tx *gorm.DB, booking models.Booking) (models.Booking, error) {
	if tx.Error != nil {
		return models.Booking{}, tx.Error
	}

	if err := tx.Create(&booking).Error; err != nil {
		return models.Booking{}, err
	}

	return booking, nil
}


func (bookingRepo *BookingRepository) GetBookingFee() (models.BookingFee, error) {

	var bookingFee models.BookingFee
	if err := bookingRepo.DB.First(&bookingFee).Error; err != nil {
		return bookingFee, err
	}
	
	return bookingFee, nil
}