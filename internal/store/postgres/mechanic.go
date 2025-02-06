package postgres

import (
	"github.com/Abrahamosaz/TMND/internal/models"
	"gorm.io/gorm"
)



type MechanicRepository struct {
	DB *gorm.DB
}



func (mechanicRepo *MechanicRepository) Create(tx *gorm.DB, mechanic models.Mechanic) (models.Mechanic, error) {
	return models.Mechanic{}, nil
}


func (mechanicRepo *MechanicRepository) GetAllAvailableMechanics() (*[]models.Mechanic, error) {
	var mechanics []models.Mechanic

	err := mechanicRepo.DB.
		Where("is_available = ?", true).
		Find(&mechanics).Error

	if err != nil {
		return nil, err
	}

	return &mechanics, nil
}


func (mechanicRepo *MechanicRepository) GetAvailableMechanic(blackListedIDS []string) (*models.Mechanic, error) {
	var mechanic models.Mechanic

	query := mechanicRepo.DB.Where("is_available = ?", true)
    
    if len(blackListedIDS) > 0 {
        query = query.Where("id NOT IN (?)", blackListedIDS)
    }
    
	result := query.Order("rating DESC").
        Limit(1).
        First(&mechanic)

	
    if err := result.Error; err != nil {
        return &mechanic, err
    }

	return &mechanic, nil
}