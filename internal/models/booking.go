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
	ID       		        uuid.UUID   		`gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	PaymentRef 				string 				`gorm:"type:varchar(255);unique" json:"paymentRef"`
	UserID					uuid.UUID			`json:"userId"`
	AssignedMechanicID		uuid.UUID			`json:"-"`
	ErrorMessage 			string 				`json:"errorMessage"`
	MechanicID				*uuid.UUID			`json:"mechanicId"`
	Services				[]*Service			`gorm:"many2many:booking_services" json:"services"`
	ServiceType 			string				`json:"serviceType"`
	ServiceDescription		*string				`json:"serviceDescription"`		
	Vehicle 				Vehicle				`gorm:"foreignKey:BookingID" json:"vehicle"`
	EstimatedPrice			float64				`gorm:"default:0.0" json:"estimatedPrice"`
	BookingFee				float64				`gorm:"default:0.0" json:"bookingFee"`
	BookingDate				time.Time			`json:"bookingDate"`
	Status					BookingStatus		`json:"status"`
	Latitude				float64				`json:"latitude"`
	Longitude				float64				`json:"longitude"`
	Address					string				`json:"address"`
	BlacklistedMechanics 	[]string 			`gorm:"type:jsonb" json:"-"`
	NextExecutionTime		*time.Time			`json:"-"`
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