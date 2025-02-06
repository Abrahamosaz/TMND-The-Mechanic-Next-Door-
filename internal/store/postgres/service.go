package postgres

import (
	"github.com/Abrahamosaz/TMND/internal/models"
	"gorm.io/gorm"
)



type ServiceRepository struct {
	DB *gorm.DB
}


func (serviceRepo *ServiceRepository) GetServiceCategories() (*[]models.ServiceCategory, error) {
	var services []models.ServiceCategory

	if err := serviceRepo.DB.Preload("Services").Find(&services).Error; err != nil {
		return &services, err
	}

	return &services, nil
}



func (serviceRepo *ServiceRepository) GetService(service *models.Service) error {
	if err := serviceRepo.DB.First(service).Error; err != nil {
		return err
	}
	return nil
} 