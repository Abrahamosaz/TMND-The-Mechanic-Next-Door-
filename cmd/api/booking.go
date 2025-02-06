package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Abrahamosaz/TMND/internal/services"
)



func (app *application) getDisabledDateHanlder(w http.ResponseWriter, r *http.Request) {

	user, ok := app.GetUserFromContext(r)

	if !ok {
		app.responseJSON(http.StatusUnauthorized, w,  "Unauthorized: No user found", nil)
		return
	}

	serviceApp := app.createNewServiceApp()

	disabledDates, statusCode, err := services.GetDisabledDatesForUser(&serviceApp, user)

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

	var createBookingDto services.CreateBooking
	err := json.NewDecoder(r.Body).Decode(&createBookingDto)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = validate.Struct(createBookingDto)
	if err != nil {
		ValidateRequestBody(err, w)
		return
	}

	serviceApp := app.createNewServiceApp()
	newBooking, statusCode, err := services.CreateUserBooking(&serviceApp, createBookingDto, user)

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
	fee, statusCode, err := services.GetBookingFee(&serviceApp)

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
	fee, statusCode, err := services.GetBookingServices(&serviceApp)

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