package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type User struct {
	ID        		uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
    FullName      	string    `gorm:"type:varchar(100)" json:"full_name"`
    Email     		string    `gorm:"type:varchar(100);unique" json:"email"`
    PhoneNumber     *string   `gorm:"type:varchar(20);index;unique"` 
    Password  		string    `gorm:"type:varchar(100)" json:"password"`
    CreatedAt 		time.Time `json:"created_at"`
    UpdatedAt 		time.Time `json:"updated_at"`
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
