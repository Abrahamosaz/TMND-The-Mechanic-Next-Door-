package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Abrahamosaz/TMND/internal/models"
	"github.com/Abrahamosaz/TMND/internal/services"
	"github.com/Abrahamosaz/TMND/internal/utils"
	"github.com/go-chi/chi/v5"
)



func (app *application) getDisabledDateHanlder(w http.ResponseWriter, r *http.Request) {

	user, ok := app.GetUserFromContext(r)

	if !ok {
		app.responseJSON(http.StatusUnauthorized, w,  "Unauthorized: No user found", nil)
		return
	}

	serviceApp := app.createNewServiceApp()

	disabledDates, statusCode, err := serviceApp.GetDisabledDatesForUser(user)

	if err != nil {
		log.Println("error getting disabled dates: ", err.Error())
		message := err.Error()
		if statusCode == http.StatusInternalServerError {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}

	app.responseJSON(statusCode, w, "User disabled dates retrieve successfully", disabledDates)
}


func (app *application) createBookingHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := app.GetUserFromContext(r)

	if !ok {
		app.responseJSON(http.StatusUnauthorized, w,  "Unauthorized: No user found", nil)
		return
	}

	filesUrl, filenames, publicIds, _ := app.GetUploadedFilesFromContext(r)

    // Create booking DTO from form values
    var createBookingDto services.CreateBooking

	if filesUrl != nil && filenames != nil && publicIds != nil {
		createBookingDto.VehicleImagesUrl = filesUrl
		createBookingDto.VehicleImagesFilename = filenames
		createBookingDto.PublicIds = publicIds
	}

    // Get the form data as JSON string
    jsonData := r.FormValue("data")
    if jsonData == "" {
        app.responseJSON(http.StatusBadRequest, w, "Missing form data", nil)
        return
    }
	
	// fmt.Println("jsonData: ", jsonData)
    // Decode the JSON string into the createBookingDto
    err := json.Unmarshal([]byte(jsonData), &createBookingDto)
    if err != nil {
        app.responseJSON(http.StatusBadRequest, w, "Invalid JSON in form data", nil)
        return
    }

    err = validate.Struct(createBookingDto)
    if err != nil {
        ValidateRequestBody(err, w)
        return
    }

	serviceApp := app.createNewServiceApp()
	newBooking, statusCode, err := serviceApp.CreateUserBooking(createBookingDto, user)

	if err != nil {
		log.Println("error creating new booking: ", err.Error())
		message := err.Error()
		if (statusCode == http.StatusInternalServerError) {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}

	app.responseJSON(http.StatusOK, w, "Booking created successfully", newBooking)
}


func (app *application) getBookingFeeHandler(w http.ResponseWriter, r *http.Request) {
	serviceApp := app.createNewServiceApp()
	fee, statusCode, err := serviceApp.GetBookingFee()

	if err != nil {
		log.Println("error getting booking fee: ", err.Error())
		message := err.Error()
		if (statusCode == http.StatusInternalServerError) {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}

	app.responseJSON(statusCode, w, "Booking fee retrieve successfully", fee)

}


func (app *application) cancelBookingHandler(w http.ResponseWriter, r *http.Request) {
	bookingID := chi.URLParam(r, "id")
	serviceApp := app.createNewServiceApp()
	statusCode, err := serviceApp.CancelBooking(bookingID)

	if err != nil {
		log.Println("error canceling booking: ", err.Error())
		message := err.Error()
		if (statusCode == http.StatusInternalServerError) {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}

	app.responseJSON(statusCode, w, "Booking cancelled successfully", nil)
}


func (app *application) getBookingsHandler(w http.ResponseWriter, r *http.Request) {

	user, ok := app.GetUserFromContext(r)

	if !ok {
		app.responseJSON(http.StatusUnauthorized, w,  "Unauthorized: No user found", nil)
		return
	}

	//get query strings
	qs := r.URL.Query()

	page := qs.Get("page")
	if page == "" {
		page = "1"
	}
	limit := qs.Get("limit")
	if limit == "" {
		limit = "10"
	}

	search := qs.Get("search")
	status := qs.Get("status")

	serviceApp := app.createNewServiceApp()
	bookings, statusCode, err := serviceApp.GetUserBookings(
		user, 
		&models.FilterQuery{
			Search: &search,
			Status: &status,
			PaginationQuery: &models.PaginationQuery{
				Page: utils.ConvertStrToPtrInt(page),
				Limit: utils.ConvertStrToPtrInt(limit),
			},
		},
	)

	if err != nil {
		log.Println("error getting bookings: ", err.Error())
		message := err.Error()
		if (statusCode == http.StatusInternalServerError) {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}

	app.responseJSON(statusCode, w, "Bookings retrieved successfully", bookings)
}


func (app *application) getVehicleDetailsHandlder(w http.ResponseWriter, r *http.Request) {

	type VehicleDetailsResponse struct {
		Types	[]vehicleConstants	`json:"types"`
		Sizes 	[]vehicleConstants	`json:"sizes"`
		Models 	[]vehicleConstants	`json:"models"`
		Brands 	[]vehicleConstants 	`json:"brands"`
	}

	app.responseJSON(http.StatusOK, w, "Vehicle details retrieve successfully", VehicleDetailsResponse{
		Types: VEHICLE_TYPES,
		Sizes: VEHICLE_SIZES,
		Models: VEHICLE_MODELS,
		Brands: VEHICLE_BRANDS,
	})
}


func (app *application) getServicesHandler(w http.ResponseWriter, r *http.Request) {

	serviceApp := app.createNewServiceApp()
	fee, statusCode, err := serviceApp.GetBookingServices()

		if err != nil {
		log.Println("error getting booking services: ", err.Error())
		message := err.Error()
		if (statusCode == http.StatusInternalServerError) {
			message = "internal server error"
		}
		app.responseJSON(statusCode, w, message, nil)
		return
	}

	app.responseJSON(statusCode, w, "Services retrieve successfully", fee)
}