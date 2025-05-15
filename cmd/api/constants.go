package main

import (
	"fmt"
	"time"
)

// vehicle constants
type vehicleConstants struct {
	Label string 	`json:"label"`
	Value any		`json:"value"`
}


var VEHICLE_TYPES = []vehicleConstants{ 
	{Label: "Car", Value: "car"},
	{Label: "Truck", Value: "truck"},
}


var VEHICLE_SIZES = []vehicleConstants{
    { Label: "Small", Value: "small" },
    { Label: "Medium", Value: "medium" },
    { Label: "Large", Value: "large" },
}

var VEHICLE_MODELS = generateVehicleModels()


var VEHICLE_BRANDS = []vehicleConstants{
	{Label: "Acura", Value: "Acura"},
	{Label: "Alfa Romeo", Value: "Alfa Romeo"},
	{Label: "Aston Martin", Value: "Aston Martin"},
	{Label: "Audi", Value: "Audi"},
	{Label: "Bentley", Value: "Bentley"},
	{Label: "BMW", Value: "BMW"},
	{Label: "Bugatti", Value: "Bugatti"},
	{Label: "Buick", Value: "Buick"},
	{Label: "BYD", Value: "BYD"},
	{Label: "Cadillac", Value: "Cadillac"},
	{Label: "Chevrolet", Value: "Chevrolet"},
	{Label: "Chrysler", Value: "Chrysler"},
	{Label: "Citroën", Value: "Citroën"},
	{Label: "Dacia", Value: "Dacia"},
	{Label: "Dodge", Value: "Dodge"},
	{Label: "Ferrari", Value: "Ferrari"},
	{Label: "Fiat", Value: "Fiat"},
	{Label: "Ford", Value: "Ford"},
	{Label: "Genesis", Value: "Genesis"},
	{Label: "GMC", Value: "GMC"},
	{Label: "Honda", Value: "Honda"},
	{Label: "Hyundai", Value: "Hyundai"},
	{Label: "Infiniti", Value: "Infiniti"},
	{Label: "Jaguar", Value: "Jaguar"},
	{Label: "Jeep", Value: "Jeep"},
	{Label: "Kia", Value: "Kia"},
	{Label: "Lamborghini", Value: "Lamborghini"},
	{Label: "Land Rover", Value: "Land Rover"},
	{Label: "Lexus", Value: "Lexus"},
	{Label: "Lincoln", Value: "Lincoln"},
	{Label: "Lotus", Value: "Lotus"},
	{Label: "Maserati", Value: "Maserati"},
	{Label: "Mazda", Value: "Mazda"},
	{Label: "McLaren", Value: "McLaren"},
	{Label: "Mercedes-Benz", Value: "Mercedes-Benz"},
	{Label: "Mini", Value: "Mini"},
	{Label: "Mitsubishi", Value: "Mitsubishi"},
	{Label: "Nissan", Value: "Nissan"},
	{Label: "Peugeot", Value: "Peugeot"},
	{Label: "Porsche", Value: "Porsche"},
	{Label: "Ram", Value: "Ram"},
	{Label: "Renault", Value: "Renault"},
	{Label: "Rolls-Royce", Value: "Rolls-Royce"},
	{Label: "Saab", Value: "Saab"},
	{Label: "Subaru", Value: "Subaru"},
	{Label: "Suzuki", Value: "Suzuki"},
	{Label: "Tesla", Value: "Tesla"},
	{Label: "Toyota", Value: "Toyota"},
	{Label: "Volkswagen", Value: "Volkswagen"},
	{Label: "Volvo", Value: "Volvo"},
}


func generateVehicleModels() []vehicleConstants {
	currentYear := time.Now().Year()
	vehicleModels := make([]vehicleConstants, 0)
	
	// Loop from 2002 to the current year
	for year := 2002; year <= currentYear+1; year++ {
		vehicleModels = append(vehicleModels, vehicleConstants{
			Label: fmt.Sprintf("%d", year), // Convert int to string for the label
			Value: year,                     // Keep value as int
		})
	}

	return vehicleModels
}