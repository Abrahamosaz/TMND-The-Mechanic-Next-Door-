package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingPending BookingStatus = "pending"
	BookingBooked BookingStatus = "booked"
	BookingCancelled BookingStatus = "cancelled"
)


type Booking struct {
	ID        		        uuid.UUID   		`gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID					uuid.UUID			`json:"userId"`
	MechanicID				*uuid.UUID			`json:"mechanicId"`
	ServiceID				uuid.UUID			`json:"serviceId"`
	Vehicle 				Vehicle				`gorm:"foreignKey:BookingID" json:"vehicle"`
	EstimatedPrice			float64				`gorm:"default:0.0" json:"estimatedPrice"`
	BookingFee				float64				`gorm:"default:0.0" json:"bookingFee"`
	BookingDate				time.Time			`json:"bookingDate"`
	Status					BookingStatus		`json:"status"`
	CreatedAt 		        time.Time   		`json:"createdAt"`
    UpdatedAt 		        time.Time   		`json:"updatedAt"`
}


func (Booking) TableName() string {
	return "bookings"
}

func (u *Booking) BeforeCreate(tx *gorm.DB) (err error) {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return
}


type BookingFee struct {
	ID        		        uuid.UUID   		`gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Price 					float64				`gorm:"default:0.0" json:"price"`
	CreatedAt 		        time.Time   		`json:"createdAt"`
    UpdatedAt 		        time.Time   		`json:"updatedAt"`

}


func (BookingFee) TableName() string {
	return "booking_fee"
}


func (u *BookingFee) BeforeCreate(tx *gorm.DB) (err error) {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return
}