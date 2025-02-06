package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type TrasactionStatus string

const (
	StatusPending 	TrasactionStatus = "pending"
	StatusFailed 	TrasactionStatus = "failed"
	StatusSuccess 	TrasactionStatus = "success"
)

type Transaction struct {
	ID        		        uuid.UUID   		`gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	PaymentRef 				string 				`gorm:"type:varchar(255);unique" json:"paymentRef"`
	UserID					*uuid.UUID			`json:"userId"`
	MechanicID				*uuid.UUID			`json:"mechanicId"`
	PreviousBalance			float64				`json:"previousBalance"`
	CurrentBalance			float64				`json:"currentBalance"`
	Status					TrasactionStatus	`gorm:"default:pending" json:"status"`
	Amount					float64				`json:"amount"`
	CreatedAt 		        time.Time   		`json:"createdAt"`
    UpdatedAt 		        time.Time   		`json:"updatedAt"`
}


func (Transaction) TableName() string {
	return "transactions"
}

func (u *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return
}