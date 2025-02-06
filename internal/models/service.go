package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type ServiceCategory struct {
	ID        		        uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name					string		`json:"name"`
	Description				string		`json:"description"`
	Services				[]Service	`gorm:"foreignKey:ServiceCategoryID" json:"services"`
	CreatedAt 		        time.Time   `json:"createdAt"`
    UpdatedAt 		        time.Time   `json:"updatedAt"`
}


func (ServiceCategory) TableName() string {
	return "service_categories"
}

func (u *ServiceCategory) BeforeCreate(tx *gorm.DB) (err error) {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return
}


type ServiceDifficulty string 

const (
	Easy ServiceDifficulty = "Easy"
	Medium ServiceDifficulty = "Medium"
	Hard ServiceDifficulty = "Hard"
)


type Service struct {
	ID        		        uuid.UUID   		`gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	ServiceCategoryID		uuid.UUID			`json:"serviceCategoryId"`	
	Name					string				`json:"name"`
	Description				string				`json:"description"`
	BasePrice				float64				`json:"basePrice"`		
	Duration				time.Duration		`json:"duration"`
	Difficulty				ServiceDifficulty 	`json:"difficulty"`
	IsAvailable				bool				`json:"isAvailable"`
	Bookings 				[]*Booking			`gorm:"many2many:booking_services" json:"bookings"`
	CreatedAt 		        time.Time   		`json:"createdAt"`
    UpdatedAt 		        time.Time   		`json:"updatedAt"`
}

func (Service) TableName() string {
	return "services"
}

func (u *Service) BeforeCreate(tx *gorm.DB) (err error) {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return
}