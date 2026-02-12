package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/resend/resend-go/v3"
	"github.com/thexovc/TMND/internal/db"
	"github.com/thexovc/TMND/internal/models"
	"github.com/thexovc/TMND/internal/services"
	"github.com/thexovc/TMND/internal/store"
)

// @title TMND API
// @version 1.0
// @description This is the backend API for The Mechanic Next Door (TMND).
// @termsOfService http://swagger.io/terms/

// @contact.name Daniel Osariemen
// @contact.url http://github.com/thexovc
// @contact.email osazeepeter79@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("PORT not set, using default: 8080")
	}

	cfg := config{
		addr: fmt.Sprintf(":%v", port),
		resend: resendConfig{
			apiKey: os.Getenv("RESEND_API_KEY"),
			from:   os.Getenv("RESEND_FROM_EMAIL"),
		},
	}

	dbConfig := db.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_MODE"),
	}

	// Validate required database configuration
	if dbConfig.Host == "" || dbConfig.User == "" || dbConfig.DBName == "" {
		log.Fatal("Missing required database configuration. Please set DB_HOST, DB_USER, and DB_NAME environment variables")
	}

	if dbConfig.Port == "" {
		dbConfig.Port = "5432" // Default PostgreSQL port
		log.Println("DB_PORT not set, using default: 5432")
	}

	if dbConfig.SSLMode == "" {
		dbConfig.SSLMode = "require" // Default to require for cloud databases
		log.Println("DB_MODE not set, using default: require")
	}

	log.Printf("Connecting to database: host=%s port=%s dbname=%s user=%s", dbConfig.Host, dbConfig.Port, dbConfig.DBName, dbConfig.User)

	dbCon, err := db.ConnectDB(dbConfig)

	if err != nil {
		log.Fatal("Failed to connect to database ", err)
	}

	// enable uuid extension
	db.EnableUUIDExtension(dbCon)

	// migrate models only for development purpose
	err = dbCon.AutoMigrate(
		&models.User{},
		&models.Mechanic{},
		&models.ServiceCategory{},
		&models.Service{},
		&models.Booking{},
		&models.Vehicle{},
		&models.BookingFee{},
		&models.Transaction{},
	)

	//sendDB
	seedDB(dbCon)

	if err != nil {
		log.Fatal("Failed migrating database models")
	}

	log.Printf("Connected to database successfully")

	pgStore := store.PostgresStorage(dbCon)

	// Initialize Resend client
	if cfg.resend.apiKey == "" {
		log.Println("Warning: RESEND_API_KEY not set. Email functionality will not work.")
	}
	if cfg.resend.from == "" {
		log.Println("Warning: RESEND_FROM_EMAIL not set. Email functionality will not work.")
	}

	resendClient := resend.NewClient(cfg.resend.apiKey)

	app := &application{
		config:   cfg,
		store:    pgStore,
		dbConfig: dbConfig,
		resend:   resendClient,
	}

	//Run cronJob
	serviceApp := app.createNewServiceApp()
	services.StartAllBookingCronJobs(&serviceApp)

	mux := app.mount()
	err = app.run(mux)
	if err != nil {
		log.Fatal("Error running server", err)
	}
}
