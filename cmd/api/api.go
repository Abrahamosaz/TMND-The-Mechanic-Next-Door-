package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Abrahamosaz/TMND/internal/db"
	"github.com/Abrahamosaz/TMND/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	gomail "gopkg.in/mail.v2"
)


type application struct {
	config  	config
	store  		store.Storage
	dbConfig 	db.DBConfig
	smtp    *gomail.Dialer
}

type config struct {
	addr string
	smtp smtpConfig
}


type smtpConfig struct {
	user 		string
	from 		string
	password 	string
	host 		string
	port 		string
}


func (app *application) mount() http.Handler {

	router := chi.NewRouter()

  // A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))


	// health check
	router.Route("/api/v1", func(rootRouter chi.Router) {
		rootRouter.Get("/health", app.healthCheckHandler)
		rootRouter.Get("/send-email", app.testSendMail)
	})
	
	return router
}


func (app *application) run(mux http.Handler) (error) {

	srv := &http.Server{
		Addr: app.config.addr,
		Handler: mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 10,
		IdleTimeout: time.Minute,
	}

	log.Printf("Server running on %s", app.config.addr)
	return srv.ListenAndServe()
}

