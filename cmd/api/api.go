package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Abrahamosaz/TMND/internal/db"
	"github.com/Abrahamosaz/TMND/internal/models"
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
	Data  		any				`json:"data"`
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
			authRouter.Post("/signup", app.userSignupHandler)
			authRouter.Post("/login", app.userLoginHandler)
			authRouter.Post("/forgot-password", app.forgotPasswordHandler)
			authRouter.Post("/verify-otp", app.verifyOtpHandler)
			authRouter.Post("/change-password", app.changePasswordHandlder)
			authRouter.Post("/resend-otp", app.resendOtpHandler)
			authRouter.Post("/verify-email", app.verifyEmailHandler)

			authRouter.Post("/mechanic/login", app.mechanicLoginHandler)
		})


		//user routes
		rootRouter.Route("/user", func(userRouter chi.Router) {
			userRouter.Use(app.userAuthMiddleware)
			userRouter.Get("/get-user", app.getUserHandler)
			userRouter.With(app.uploadMiddleware("profile-image", CLOUDINARY_PROFILE_IMAGE_FOLDER)).Put("/edit-profile", app.editUserProfileHandler)
			
			// booking routes
			userRouter.Route("/booking", func(bookingRouter chi.Router) {
				bookingRouter.Get("/get-bookings", app.getBookingsHandler)
				bookingRouter.Post("/cancel-booking/{id}", app.cancelBookingHandler)
				bookingRouter.Get("/get-booking-fee", app.getBookingFeeHandler)
				bookingRouter.Get("/get-vehicle-details", app.getVehicleDetailsHandlder)
				bookingRouter.Get("/get-disabled-date", app.getDisabledDateHanlder)
				bookingRouter.With(app.uploadMultipleFilesMiddleware("vehicle-images", CLOUDINARY_VEHICLE_IMAGE_FOLDER)).Post("/create-booking", app.createBookingHandler)
			})

			//service routes
			userRouter.Route("/service", func(serviceRouter chi.Router) {
				// serviceRouter.Get("", app.)
				serviceRouter.Get("/get-service-categories", app.getServicesHandler)
			})

			//transaction routes
			userRouter.Route("/transaction", func(trxRouter chi.Router) {
				trxRouter.Get("/", app.getUserTransactionHandler)
			})
		})


		//mechanic routes
		rootRouter.Route("/mechanic", func(mechanicRouter chi.Router) {
			mechanicRouter.Use(app.mechanicAuthMiddleware)

			// booking routes
			mechanicRouter.Route("/booking", func(bookingRouter chi.Router) {
				bookingRouter.Post("/reject", app.rejectBookingHandler)
				bookingRouter.Post("/accept", app.acceptBookingHandler)
			})


			//transaction routes
			mechanicRouter.Route("/transaction", func(trxRouter chi.Router) {
				// trxRouter.Get()
			})

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

func (app *application) responseJSON(statusCode int, w http.ResponseWriter, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if v, ok := data.([]string); ok && len(v) == 0 {
		data = []string{}
	}
	json.NewEncoder(w).Encode(Response{
		Message: message,
		StatusCode: statusCode,
		Data: data,
	})
}


func (app *application) GetUserFromContext(r *http.Request) (*models.User, bool) {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	return user, ok
}

func (app *application) GetFileInfoFromContext(r *http.Request) (*string, *string, bool) {
	fileInfo, ok := r.Context().Value(uploadContextKey).(*UploadResult)
	if !ok {
		return nil, nil, false
	}
	return &fileInfo.URL, &fileInfo.FileName, true
}


func (app *application) GetUploadedFilesFromContext(r *http.Request) ([]string, []string, bool) {
	// Get the upload results from context
	uploadResults, ok := r.Context().Value(uploadMultipleFilesContextKey).([]UploadResult)
	if !ok {
		return nil, nil, false
	}

	// Extract URLs and filenames into separate slices
	urls := make([]string, len(uploadResults))
	filenames := make([]string, len(uploadResults))
	
	for i, result := range uploadResults {
		urls[i] = result.URL
		filenames[i] = result.FileName
	}

	return urls, filenames, true
}