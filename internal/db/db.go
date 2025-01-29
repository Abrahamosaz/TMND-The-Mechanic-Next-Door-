package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


type DBConfig struct {
	Host 		string
	User 		string
	Password 	string
	DBName 		string
	Port 		string
	SSLMode 	string
}


func ConnectDB(config DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return db, err
	}
	
	return db, nil
}


func EnableUUIDExtension(db *gorm.DB) {
	// Enable uuid-ossp extension
	err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		log.Fatalf("Failed to enable uuid-ossp extension: %v", err)
	}

	log.Println("UUID extension enabled successfully.")
}
