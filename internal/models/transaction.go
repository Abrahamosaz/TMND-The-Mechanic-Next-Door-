package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type TrasactionStatus string
type TransactionType string

const (
	StatusPending 	TrasactionStatus = "pending"
	StatusFailed 	TrasactionStatus = "failed"
	StatusSuccess 	TrasactionStatus = "success"
)

const (
	TransactionTypeDebit TransactionType = "debit"
	TransactionTypeCredit TransactionType = "credit"
)

type TransactionStatus string

type Transaction struct {
	ID        		        uuid.UUID   		`gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	TrxRef 					string 				`gorm:"type:varchar(255);unique" json:"trxRef"`
	UserID					*uuid.UUID			`json:"userId"`
	MechanicID				*uuid.UUID			`json:"mechanicId"`
	PreviousBalance			float64				`json:"previousBalance"`
	CurrentBalance			float64				`json:"currentBalance"`
	Type					TransactionType		`gorm:"default:debit" json:"type"`	
	Status					TrasactionStatus	`gorm:"default:pending" json:"status"`
	Fee						*float64			`json:"-"`
	Amount					float64				`json:"amount"`
	CreatedAt 		        time.Time   		`json:"createdAt"`
    UpdatedAt 		        time.Time   		`json:"updatedAt"`
	Description				*string				`json:"description"`


	//relationships
	User        			*User   			`gorm:"foreignKey:UserID" json:"user"`
	Mechanic        		*Mechanic   		`gorm:"foreignKey:MechanicID" json:"mechanic"`
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