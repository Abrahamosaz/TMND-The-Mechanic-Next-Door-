package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Abrahamosaz/TMND/internal/models"
	// "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)



func seedDB(db *gorm.DB) {
	// seed bookings
	// bookings := []models.Booking{
	// 	{UserID: uuid.MustParse("d7831ca4-769a-446e-aa8d-4fb22d064f3b"), ServiceID: uuid.MustParse("ffb0c368-9520-4bcd-b778-250dee496cea"), BookingDate: time.Now().AddDate(0, 0, 1), Status: models.BookingPending},
	// 	{UserID: uuid.MustParse("d7831ca4-769a-446e-aa8d-4fb22d064f3b"), ServiceID: uuid.MustParse("ffb0c368-9520-4bcd-b778-250dee496cea"), BookingDate: time.Now().AddDate(0, 0, 1), Status: models.BookingPending},
	// 	{UserID: uuid.MustParse("d7831ca4-769a-446e-aa8d-4fb22d064f3b"), ServiceID: uuid.MustParse("ffb0c368-9520-4bcd-b778-250dee496cea"), BookingDate: time.Now().AddDate(0, 0, 1), Status: models.BookingPending},
	// }
	
	// if err := db.Create(&bookings).Error; err != nil{
	// 	log.Printf("log seeder: failed to create bookings %v", err)
	// }

	//set the booking fee
	var existingBookingFee models.BookingFee
	if result := db.First(&existingBookingFee); result.RowsAffected < 1 {
		bookingFee := models.BookingFee{Price: 1000}
		if err := db.Create(&bookingFee).Error; err != nil {
			log.Println("failed to create booking fee")
		}
	}

	//seed mechanics
	hashedPassword1, err := bcrypt.GenerateFromPassword([]byte("ibadin"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("failed to generate password hash for password1")
	}

	hashedPassword2, err := bcrypt.GenerateFromPassword([]byte("abraham"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("failed to generate password hash for password2")
	}

	mechanics := []models.Mechanic{
		{FullName: "Ibadin Meshach", Email: "ibmeshach@gmail.com", PhoneNumber: "08148274833", Password: string(hashedPassword1), IsEmailVerified: true, Specialty: "Repair Engines", Experience: 3, IsAvailable: true},
		{FullName: "Abraham Omorisiagbon", Email: "abrahamosazee3@gmail.com", PhoneNumber: "08061909748", Password: string(hashedPassword2), IsEmailVerified: true, Specialty: "Repair Engines", Experience: 3, IsAvailable: true},
	}

	for i, mechanic := range mechanics {
		var existingMechanic models.Mechanic
		result := db.Where("email = ?", mechanic.Email).First(&existingMechanic)
		if result.Error != nil {
			// Insert the service if it does not exist
			result := db.Create(&mechanics[i])
			if result.Error != nil {
				log.Printf("log seeder: failed to create %s mechanic", mechanic.FullName)
			}
		}
}
	//seed service category
	serviceCategories := []models.ServiceCategory{
		{Name: "Routine Maintenance", Description: "General maintenance services to keep your vehicle in top condition."},
		{Name: "Engine Services", Description: "Repair and maintenance services for your vehicle's engine."},
		{Name: "Transmission Services", Description: "Maintenance and repairs for vehicle transmissions."},
		{Name: "Electrical & Battery", Description: "Services related to vehicle electrical systems and battery."},
		{Name: "Suspension & Steering", Description: "Installing a new battery for proper power supply."},
		{Name: "Exhaust System", Description: "Services for vehicle exhaust system maintenance and repair."},
		{Name: "Air Conditioning (AC)", Description: "Services for vehicle air conditioning system."},
		{Name: "Diagnostics", Description: "Diagnostic services for vehicle issues."},
		{Name: "Bodywork & Miscellaneous", Description: "Various bodywork and miscellaneous services."},
	}

	for i, category := range serviceCategories {
		// Check if the category already exists
		var existingCategory models.ServiceCategory
		result := db.Where("name = ?", category.Name).First(&existingCategory)

		if result.Error != nil {
			// Insert the category if it does not exist
			result := db.Create(&serviceCategories[i])
			if result.Error != nil {
				log.Printf("log seeder: failed to create %s category service", category.Name)
			}

			existingCategory = serviceCategories[i]
		}
		// Assign the correct ID to avoid foreign key constraint issues
		serviceCategories[i].ID = existingCategory.ID
	}

	//seed services base on their categories
	for _, category := range serviceCategories {
		switch (category.Name) {
			case "Routine Maintenance":
				routineMaintenanceServices := []*models.Service{
					{ServiceCategoryID: category.ID, Name: "Oil Change Labor", Description: "Changing engine oil to maintain engine health.", BasePrice: 2500, Duration: 1800 * time.Second, Difficulty: models.Easy, IsAvailable: false},
					{ServiceCategoryID: category.ID, Name: "Tire Rotation & Balancing Labor", Description: "Rotating and balancing tires to ensure even wear.", BasePrice: 3000, Duration: 2700 * time.Second, Difficulty: models.Medium, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Brake Pad Replacement Labor (per axle)", Description: "Replacing brake pads to ensure safety and performance.", BasePrice: 8000, Duration: 5400 * time.Second, Difficulty: models.Easy, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Air Filter Replacement Labor", Description: "Replacing the air filter to improve engine efficiency.", BasePrice: 1500, Duration: 1200 * time.Second, Difficulty: models.Easy, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Coolant Flush Labor", Description: "Flushing and replacing coolant for optimal engine temperature regulation.", BasePrice: 5000, Duration: 3600 * time.Second, Difficulty: models.Medium, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Wheel Alignment Labor", Description: "Adjusting wheel alignment for improved handling and tire longevity.", BasePrice: 4000, Duration: 2700 * time.Second, Difficulty: models.Medium, IsAvailable: true},
				}

				for _, service := range routineMaintenanceServices {
					var existingService models.Service
					result := db.Where("name = ? AND service_category_id = ?", service.Name, service.ServiceCategoryID).First(&existingService)
					if result.Error != nil {
						// Insert the service if it does not exist
						result := db.Create(&service)
						if result.Error != nil {
							log.Printf("log seeder: failed to create %s service", service.Name)
						}
					}
				}

			case "Engine Services":
				engineServices := []*models.Service{
					{ServiceCategoryID: category.ID, Name: "Timing Belt Replacement Labor", Description: "Replacing the timing belt to prevent engine damage.", BasePrice: 15000, Duration: 14400 * time.Second, Difficulty: models.Hard, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Engine Tune-Up Labor", Description: "Performing a full engine tune-up for better performance.", BasePrice: 10000, Duration: 7200 * time.Second, Difficulty: models.Medium, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Spark Plug Replacement Labor", Description: "Replacing spark plugs for efficient engine performance.", BasePrice: 3000, Duration: 1800 * time.Second, Difficulty: models.Easy, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Engine Overhaul Labor", Description: "RComplete engine overhaul for restoring performance.", BasePrice: 70000, Duration: 86400 * time.Second, Difficulty: models.Hard, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Fuel Injector Cleaning Labor", Description: "Cleaning fuel injectors to improve fuel efficiency.", BasePrice: 6000, Duration: 3600 * time.Second, Difficulty: models.Medium, IsAvailable: true},
				}

				for _, service := range engineServices {
					var existingService models.Service
					result := db.Where("name = ? AND service_category_id = ?", service.Name, service.ServiceCategoryID).First(&existingService)
					if result.Error != nil {
						// Insert the service if it does not exist
						result := db.Create(&service)
						if result.Error != nil {
							log.Printf("log seeder: failed to create %s service", service.Name)
						}
					}
				}
			case "Transmission Services":
				transmissionServices := []*models.Service{
					{ServiceCategoryID: category.ID, Name: "Transmission Fluid Change Labor", Description: "Changing transmission fluid for smoother gear shifts.", BasePrice: 10000, Duration: 3600 * time.Second, Difficulty: models.Medium, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Clutch Replacement Labor", Description: "Replacing the clutch for optimal transmission function.", BasePrice: 20000, Duration: 14400 * time.Second, Difficulty: models.Hard, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Transmission Repair Labor", Description: "Repairing transmission components for improved operation.", BasePrice: 50000, Duration: 28800 * time.Second, Difficulty: models.Hard, IsAvailable: true},

				}

				for _, service := range transmissionServices {
					var existingService models.Service
					result := db.Where("name = ? AND service_category_id = ?", service.Name, service.ServiceCategoryID).First(&existingService)
					if result.Error != nil {
						// Insert the service if it does not exist
						result := db.Create(&service)
						if result.Error != nil {
							log.Printf("log seeder: failed to create %s service", service.Name)
						}
					}
				}
			case "Electrical & Battery":
				electricalAndBatteryServices := []*models.Service{
					{ServiceCategoryID: category.ID, Name: "Battery Installation Labor", Description: "Installing a new battery for proper power supply.", BasePrice: 2000, Duration: 1800 * time.Second, Difficulty: models.Easy, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Alternator Repair Labor", Description: "Repairing the alternator to restore charging function.", BasePrice: 15000, Duration: 10800 * time.Second, Difficulty: models.Medium, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Starter Motor Replacement Labor", Description: "Replacing the starter motor for reliable engine start.", BasePrice: 10000, Duration: 7200 * time.Second, Difficulty: models.Medium, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Wiring Repairs Labor", Description: "Fixing electrical wiring issues for optimal functionality.", BasePrice: 3000, Duration: 3600 * time.Second, Difficulty: models.Easy, IsAvailable: true},
				}

				for _, service := range electricalAndBatteryServices {
					var existingService models.Service
					result := db.Where("name = ? AND service_category_id = ?", service.Name, service.ServiceCategoryID).First(&existingService)
					if result.Error != nil {
						// Insert the service if it does not exist
						result := db.Create(&service)
						if result.Error != nil {
							log.Printf("log seeder: failed to create %s service", service.Name)
						}
					}
				}
			case "Suspension & Steering":
				suspensionAndSteeringServices := []*models.Service{
					{ServiceCategoryID: category.ID, Name: "Shock Absorber Replacement Labor (per unit)", Description: "Replacing shock absorbers for improved ride comfort.", BasePrice: 8000, Duration: 7200 * time.Second, Difficulty: models.Medium, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Ball Joint Replacement Labor", Description: "Replacing ball joints for better steering control.", BasePrice: 6000, Duration: 5400 * time.Second, Difficulty: models.Medium, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Power Steering Flush Labor", Description: "Flushing power steering fluid for smooth operation.", BasePrice: 5000, Duration: 3600 * time.Second, Difficulty: models.Medium, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Tie Rod Replacement Labor", Description: "Replacing tie rods for precise steering.", BasePrice: 7000, Duration: 5400 * time.Second, Difficulty: models.Medium, IsAvailable: true},
				}

				for _, service := range suspensionAndSteeringServices {
					var existingService models.Service
					result := db.Where("name = ? AND service_category_id = ?", service.Name, service.ServiceCategoryID).First(&existingService)
					if result.Error != nil {
						// Insert the service if it does not exist
						result := db.Create(&service)
						if result.Error != nil {
							log.Printf("log seeder: failed to create %s service", service.Name)
						}
					}
				}
			case "Exhaust System":
				exhaustSystemServices := []*models.Service{
					{ServiceCategoryID: category.ID, Name: "Muffler Replacement Labor", Description: "Replacing muffler for proper noise reduction.", BasePrice: 5000, Duration: 3600 * time.Second, Difficulty: models.Medium, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Catalytic Converter Replacement Labor", Description: "Replacing catalytic converter for emissions control.", BasePrice: 25000, Duration: 7200 * time.Second, Difficulty: models.Hard, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Exhaust Pipe Repair Labor", Description: "Repairing exhaust pipe for proper gas flow.", BasePrice: 6000, Duration: 4800 * time.Second, Difficulty: models.Medium, IsAvailable: true},
				}

				for _, service := range exhaustSystemServices {
					var existingService models.Service
					result := db.Where("name = ? AND service_category_id = ?", service.Name, service.ServiceCategoryID).First(&existingService)
					if result.Error != nil {
						// Insert the service if it does not exist
						result := db.Create(&service)
						if result.Error != nil {
							log.Printf("log seeder: failed to create %s service", service.Name)
						}
					}
				}
			case "Air Conditioning (AC)":
				airConditioningServices := []*models.Service{
					{ServiceCategoryID: category.ID, Name: "AC Gas Recharge Labor", Description: "Recharging AC gas for optimal cooling.", BasePrice: 5000, Duration: 2700 * time.Second, Difficulty: models.Easy, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Compressor Replacement Labor", Description: "Replacing AC compressor for proper system function.", BasePrice: 25000, Duration: 10800 * time.Second, Difficulty: models.Hard, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "AC Leak Repair Labor", Description: "Repairing AC system leaks.", BasePrice: 10000, Duration: 7200 * time.Second, Difficulty: models.Medium, IsAvailable: true},
				}

				for _, service := range airConditioningServices {
					var existingService models.Service
					result := db.Where("name = ? AND service_category_id = ?", service.Name, service.ServiceCategoryID).First(&existingService)
					if result.Error != nil {
						// Insert the service if it does not exist
						result := db.Create(&service)
						if result.Error != nil {
							log.Printf("log seeder: failed to create %s service", service.Name)
						}
					}
				}
			case "Diagnostics":
				diagnosticsServices := []*models.Service{
					{ServiceCategoryID: category.ID, Name: "Computer Diagnostics Labor", Description: "Running computer diagnostics to identify issues.", BasePrice: 3000, Duration: 1800 * time.Second, Difficulty: models.Easy, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Engine Light Scan Labor", Description: "Scanning engine light codes for problem identification.", BasePrice: 1500, Duration: 900 * time.Second, Difficulty: models.Easy, IsAvailable: true},
				}

				for _, service := range diagnosticsServices {
					var existingService models.Service
					result := db.Where("name = ? AND service_category_id = ?", service.Name, service.ServiceCategoryID).First(&existingService)
					if result.Error != nil {
						// Insert the service if it does not exist
						result := db.Create(&service)
						if result.Error != nil {
							log.Printf("log seeder: failed to create %s service", service.Name)
						}
					}
				}
			case "Bodywork & Miscellaneous":
				bodyWorkServices := []*models.Service{
					{ServiceCategoryID: category.ID, Name: "Windshield Replacement Labor", Description: "Replacing damaged windshield.", BasePrice: 10000, Duration: 5400 * time.Second, Difficulty: models.Medium, IsAvailable: true},
					{ServiceCategoryID: category.ID, Name: "Headlight/Taillight Installation Labor", Description: "Installing new headlights or taillights.", BasePrice: 3000, Duration: 2700 * time.Second, Difficulty: models.Easy, IsAvailable: true},
			
				}

				for _, service := range bodyWorkServices {
					var existingService models.Service
					result := db.Where("name = ? AND service_category_id = ?", service.Name, service.ServiceCategoryID).First(&existingService)
					if result.Error != nil {
						// Insert the service if it does not exist
						result := db.Create(&service)
						if result.Error != nil {
							log.Printf("log seeder: failed to create %s service", service.Name)
						}
					}
				}
			default:
				fmt.Println()
		}
	}
}
