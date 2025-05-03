package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingPending BookingStatus = "pending"
	BookingBooked BookingStatus = "booked"
	BookingCancelled BookingStatus = "cancelled"
	BookingCompleted BookingStatus = "completed"
)


type Booking struct {
	ID       		        uuid.UUID   		`gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	PaymentRef 				string 				`gorm:"type:varchar(255);unique" json:"paymentRef"`
	UserID					uuid.UUID			`json:"userId"`
	AssignedMechanicID		uuid.UUID			`json:"-"`
	VehicleID               uuid.UUID           `gorm:"type:uuid;not null;uniqueIndex" json:"vehicleId"`
	ErrorMessage 			string 				`json:"errorMessage"`
	MechanicID				*uuid.UUID			`json:"mechanicId"`
	ServiceType 			string				`json:"serviceType"`
	ServiceDescription		*string				`json:"serviceDescription"`		
	EstimatedPrice			float64				`gorm:"default:0.0" json:"estimatedPrice"`
	BookingFee				float64				`gorm:"default:0.0" json:"bookingFee"`
	BookingDate				time.Time			`json:"bookingDate"`
	Status					BookingStatus		`json:"status"`
	Latitude				float64				`json:"latitude"`
	Longitude				float64				`json:"longitude"`
	Address					string				`json:"address"`
	BlacklistedMechanics 	datatypes.JSON 		`gorm:"type:jsonb" json:"-"`
	VisitedMechanics	 	datatypes.JSON 		`gorm:"type:jsonb" json:"-"`
	VehicleImagesUrl		datatypes.JSON		`gorm:"type:jsonb" json:"vehicleImagesUrl"`
	VehicleImagesFilename	datatypes.JSON 		`gorm:"type:jsonb" json:"-"`
	NextExecutionTime		*time.Time			`json:"-"`
	PublicIds				*string		        `json:"-"`
	CreatedAt 		        time.Time   		`json:"createdAt"`
    UpdatedAt 		        time.Time   		`json:"updatedAt"`
	
	//relationships
	Vehicle                 *Vehicle            `gorm:"foreignKey:VehicleID;constraint:OnDelete:CASCADE;" json:"vehicle"`
	User        			*User   			`gorm:"foreignKey:UserID" json:"-"`
	Mechanic        		*Mechanic   		`gorm:"foreignKey:MechanicID" json:"mechanic"`
	Services				[]*Service			`gorm:"many2many:booking_services" json:"services"`
	
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