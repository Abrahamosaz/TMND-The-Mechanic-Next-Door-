package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type User struct {
	ID        		        uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
    FullName      	        string      `gorm:"type:varchar(100)" json:"fullName"`
    Email     		        string      `gorm:"type:varchar(100);unique" json:"email"`
    PhoneNumber             *string     `gorm:"type:varchar(20);index" json:"phoneNumber"` 
    Password  		        string      `gorm:"type:varchar(100)" json:"-"`
    RegisterWithGoogle      bool        `gorm:"type:bool" json:"registerWithGoogle"`
    OtpToken                *string     `gorm:"type:text" json:"-"`
    IsEmailVerified         bool        `gorm:"type:bool;default:false" json:"isEmailVerified"`
    CreatedAt 		        time.Time   `json:"createdAt"`
    UpdatedAt 		        time.Time   `json:"updatedAt"`
}


func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return
}
