package postgres

import (
	"errors"

	"github.com/thexovc/TMND/internal/models"
	"gorm.io/gorm"
)

type MechanicRepository struct {
	DB *gorm.DB
}

func (mechanicRepo *MechanicRepository) FindByEmail(email string) (models.Mechanic, error) {
	var mechanic models.Mechanic
	result := mechanicRepo.DB.Where("email = ?", email).First(&mechanic)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.Mechanic{}, errors.New("mechanic with email not found")
		}
		return models.Mechanic{}, result.Error
	}
	return mechanic, nil
}

func (mechanicRepo *MechanicRepository) Create(tx *gorm.DB, mechanic models.Mechanic) (models.Mechanic, error) {
	return models.Mechanic{}, nil
}

func (mechanicRepo *MechanicRepository) Save(mechanic *models.Mechanic) error {
	result := mechanicRepo.DB.Save(&mechanic)
	if result.Error != nil {
		return result.Error // Return the error if save failed
	}
	return nil
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

func (mechanicRepo *MechanicRepository) GetAvailableMechanic(blackListedIDS []string, visitedIDS []string) (*models.Mechanic, error) {
	var mechanic models.Mechanic
	var highestRating float64

	query := mechanicRepo.DB.Model(&models.Mechanic{}).
		Where("is_available = ?", true)

	if len(visitedIDS) > 0 {
		query = query.Where("id NOT IN ?", visitedIDS)
	}

	if len(blackListedIDS) > 0 {
		query = query.Where("id NOT IN ?", blackListedIDS)
	}

	for {
		// Get the highest available rating
		if err := query.Select("COALESCE(MAX(rating), 0)").Scan(&highestRating).Error; err != nil {
			return nil, err
		}

		// No available mechanics left
		if highestRating == 0 {
			if len(visitedIDS) > 0 {
				return nil, errors.New("no available mechanics found")
			}
			return nil, gorm.ErrRecordNotFound
		}

		ratingQuery := mechanicRepo.DB.
			Where("is_available = ? AND rating = ?", true, highestRating).
			Order("RANDOM()").
			Limit(1)

		// Apply blacklist filtering only if necessary
		if len(blackListedIDS) > 0 {
			ratingQuery = ratingQuery.Where("id NOT IN ?", blackListedIDS)
		}

		if len(visitedIDS) > 0 {
			ratingQuery = ratingQuery.Where("id NOT IN ?", visitedIDS)
		}

		result := ratingQuery.First(&mechanic)
		if result.Error == nil {
			// Found a mechanic, return it
			return &mechanic, nil
		}

		// Remove the current highest rating and search for the next highest
		query = query.Where("rating < ?", highestRating)
	}
}
