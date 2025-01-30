package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Abrahamosaz/TMND/internal/db"
	"github.com/Abrahamosaz/TMND/internal/services"
	"github.com/Abrahamosaz/TMND/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	gomail "gopkg.in/mail.v2"
)


type application struct {
	config  	config
	store  		store.Storage
	dbConfig 	db.DBConfig
	smtp    	*gomail.Dialer
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


type Response struct {
	Message 	string 			`json:"message"`
	StatusCode 	int 			`json:"statusCode"`
	Data  		interface {} 	`json:"data"`
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
		// rootRouter.Get("/send-email", app.testSendMail)

		rootRouter.Route("/auth", func(authRouter chi.Router) {
			authRouter.Post("/signup", app.signupHandler)
			authRouter.Post("/login", app.loginHandler)
			authRouter.Post("/forgot-password", app.forgotPasswordHandler)
			authRouter.Post("/verify-otp", app.verifyOtpHandler)
			authRouter.Post("/change-password", app.changePasswordHandlder)
			authRouter.Post("/verify-ermail", app.verifyEmailHandler)
		})
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


func (app *application) createNewServiceApp()  services.Application {
	
	config := services.Config{
		Addr: app.config.addr,
		Smtp: services.SmtpConfig{
			User: app.config.smtp.user,
			From: app.config.smtp.from,
			Password: app.config.smtp.password,
			Host: app.config.smtp.host,
			Port: app.config.smtp.port,
		},
	}

	return services.Application{
		Config: config,
		Store: app.store,
		DbConfig: app.dbConfig,
		Smtp: app.smtp,
	}
}

func (app *application) responseJSON(statusCode int, w http.ResponseWriter, message string, data interface {}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(Response{
			Message: message,
			StatusCode: statusCode,
			Data: data,
		})
}


