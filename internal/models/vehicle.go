package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type Vehicle struct {
	ID        		        uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Vtype 					string		`gorm:"type:varchar(10)" json:"type"`
	Brand					string		`gorm:"type:varchar(20)" json:"brand"`
	Size					string		`gorm:"type:varchar(10)" json:"size"`
	Model					int			`gorm:"type:varchar(4)" json:"model"`
	BookingID 				uuid.UUID	`json:"bookingId"`
	CreatedAt 		        time.Time   `json:"createdAt"`
    UpdatedAt 		        time.Time   `json:"updatedAt"`

}

func (Vehicle) TableName() string {
	return "vehicles"
}

func (u *Vehicle) BeforeCreate(tx *gorm.DB) (err error) {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return
}

