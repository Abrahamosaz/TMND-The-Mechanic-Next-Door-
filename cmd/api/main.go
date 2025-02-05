package main

import (
	"fmt"
	"log"
	"os"

	"strconv"

	"github.com/Abrahamosaz/TMND/internal/db"
	"github.com/Abrahamosaz/TMND/internal/models"
	"github.com/Abrahamosaz/TMND/internal/store"

	// "github.com/joho/godotenv"
	gomail "gopkg.in/mail.v2"
)


func main() {
	// for  development purpose
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	
	port := os.Getenv("PORT")

	cfg := config{
			addr: fmt.Sprintf(":%v", port),
			smtp: smtpConfig{
				user: os.Getenv("SMTP_USER"),
				from: os.Getenv("SMTP_FROM"),
				password: os.Getenv("SMTP_PASS"),
				host: os.Getenv("SMTP_HOST"),
				port: os.Getenv("SMTP_PORT"),
			},
		}

	dbConfig := db.DBConfig{
		Host: os.Getenv("DB_HOST"),
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
		Port: os.Getenv("DB_PORT"),
		SSLMode: os.Getenv("DB_MODE"),
	}



	dbCon, err := db.ConnectDB(dbConfig)
	
	if err != nil {
		log.Fatal("Failed to connect to database ", err)
	}

	// enable uuid extension
	db.EnableUUIDExtension(dbCon)

	// migrate models
	err = dbCon.AutoMigrate(
		&models.User{},
		&models.Mechanic{},
		&models.ServiceCategory{},
		&models.Service{},
		&models.Booking{},
		&models.Vehicle{},
		&models.BookingFee{},
	)

	//sendDB
	seedDB(dbCon)

	if err != nil {
		log.Fatal("Failed migrating database models")
	}
	
	log.Printf("Connected to database successfully")

	pgStore := store.PostgresStorage(dbCon)

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))

	if err != nil {
		log.Fatal("Invalid SMTP_PORT:", err)
	}
	
	dialer := gomail.NewDialer(
		cfg.smtp.host, smtpPort,
		cfg.smtp.user, cfg.smtp.password,
	)

	app := &application{
		config: cfg,
		store: pgStore,
		dbConfig: dbConfig,
		smtp: dialer,
	}

	mux := app.mount()
	err = app.run(mux)
	if err != nil {
		log.Fatal("Error running server", err)
	}
}