package postgres

import (
	"github.com/thexovc/TMND/internal/models"
	"gorm.io/gorm"
)

type BookingRepository struct {
	DB *gorm.DB
}

func (bookingRepo *BookingRepository) GetPendingBookings() (*[]models.Booking, error) {
	var bookings []models.Booking

	err := bookingRepo.DB.
		Where("status = ?", models.BookingPending). // Filter by pending status
		Find(&bookings).Error                       // Fetch all bookings

	if err != nil {
		return nil, err
	}

	return &bookings, nil
}

func (bookingRepo *BookingRepository) GetUserBookings(user *models.User, qs *models.FilterQuery) (*models.PaginationResponse[models.Booking], error) {
	var bookings []models.Booking
	var totalCount int64

	// Initialize the query
	query := bookingRepo.DB.Model(&models.Booking{}).Where("user_id = ?", user.ID)

	// Add search condition if search parameter is provided
	if qs.Search != nil && *qs.Search != "" {
		query = query.Where("payment_ref ILIKE ?", "%"+*qs.Search+"%")
	}

	query = query.Order("created_at DESC")

	// Add status filter if status parameter is provided
	if qs.Status != nil && *qs.Status != "" {
		query = query.Where("status = ?", *qs.Status)
	}

	// Count total records before pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply pagination if parameters are provided
	if qs.Page != nil && qs.Limit != nil {
		offset := (*qs.Page - 1) * *qs.Limit
		query = query.Offset(offset).Limit(*qs.Limit)
	}

	// Execute the final query
	if err := query.Preload("Mechanic").Preload("Vehicle").Preload("Services").Find(&bookings).Error; err != nil {
		return nil, err
	}
	// Calculate total pages
	var totalPages int
	if qs.Limit != nil && *qs.Limit > 0 {
		totalPages = int((totalCount + int64(*qs.Limit) - 1) / int64(*qs.Limit))
	}

	// Calculate pagination metadata
	response := &models.PaginationResponse[models.Booking]{
		Data:       bookings,
		Total:      &totalCount,
		Page:       qs.Page,
		Limit:      qs.Limit,
		TotalPages: &totalPages,
	}

	return response, nil
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

func (bookingRepo *BookingRepository) GetBooking(booking *models.Booking) error {
	if err := bookingRepo.DB.First(booking).Error; err != nil {
		return err
	}
	return nil
}

func (bookingRepo *BookingRepository) UpdateBooking(booking *models.Booking) error {
	if err := bookingRepo.DB.Model(booking).Updates(*booking).Error; err != nil {
		return err
	}
	return nil
}
