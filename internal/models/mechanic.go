package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)



type Mechanic struct {
	ID        		        uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	FullName      	        string      `gorm:"type:varchar(100)" json:"fullName"`
    Email     		        string      `gorm:"type:varchar(100);unique" json:"email"`
    PhoneNumber             string      `gorm:"type:varchar(20);index" json:"phoneNumber"` 
    Password  		        string      `gorm:"type:varchar(100)" json:"-"`
	OtpToken                *string     `gorm:"type:text" json:"-"`
    IsEmailVerified         bool        `gorm:"type:bool;default:false" json:"isEmailVerified"`
    ProfileFileName         *string     `json:"profileFileName"`
    ProfileImageUrl         *string     `json:"profileImageUrl"`
    Address                 *string     `json:"address"`
    State                   *string     `json:"state"`
    Lga                     *string     `json:"lga"`
	Specialty				string		`json:"specialty"`
	Rating					float64		`gorm:"default:0.0" json:"rating"`
	Experience				int 		`json:"experience"`	// years of experience
	IsAvailable 			bool		`gorm:"type:bool;default:false" json:"isAvailable"`
	Bookings				[]Booking	`gorm:"foreignKey:MechanicID" json:"bookings"`
    Balance                 float64     `gorm:"default:0.0" json:"balance"`
	CreatedAt 		        time.Time   `json:"createdAt"`
    UpdatedAt 		        time.Time   `json:"updatedAt"`
}


func (Mechanic) TableName() string {
	return "mechanics"
}

func (u *Mechanic) BeforeCreate(tx *gorm.DB) (err error) {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return
}