package postgres

import (
	"github.com/Abrahamosaz/TMND/internal/models"
	"gorm.io/gorm"
)



type VehicleRepository struct {
	DB *gorm.DB
}

func (vehicleRepo *VehicleRepository) Create(tx *gorm.DB, vehicle models.Vehicle) (models.Vehicle, error) {
	if tx.Error != nil {
		return models.Vehicle{}, tx.Error
	}
	
	if err := tx.Create(&vehicle).Error; err != nil {
		return models.Vehicle{}, err
	}
	return vehicle, nil
}