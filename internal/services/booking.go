package services

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Abrahamosaz/TMND/internal/models"
	"github.com/Abrahamosaz/TMND/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


func CreateUserBooking(app *Application, payload CreateBooking, user *models.User) (models.Booking, int, error) {
    // begin transaction
    tx := app.Store.BeginTransaction()
    bookingFee, err := app.Store.Booking.GetBookingFee()

    if err != nil {
        return models.Booking{}, http.StatusInternalServerError, err
    }

    bookingDate, err := time.Parse("2006-01-02", payload.Date)

    if err != nil {
        return models.Booking{}, http.StatusInternalServerError, err
    }

    // deduct booking fee from user balance
    if user.Balance < bookingFee.Price {
        return models.Booking{}, http.StatusNotAcceptable, errors.New("insufficient funds")
    }

    amount := user.Balance - bookingFee.Price
    err = app.Store.User.DeductFromBalance(tx, user, amount)

    if err != nil {
        tx.Rollback()
        return models.Booking{}, http.StatusInternalServerError, err
    }

    tempMechanicChan := make(chan models.Mechanic)
    errChan := make(chan error)
    // get temporary  assigned mechanic
    go getMechanicForBooking(app, tempMechanicChan, errChan)

    // create all services
    var services []*models.Service 
    for _, serviceID := range payload.ServiceDetails.Services {
        service := models.Service{ID: uuid.MustParse(serviceID)}
        err := app.Store.Service.GetService(&service)
        if err != nil {
            services = append(services, &service)
        }
    }

    var temporaryMechanic models.Mechanic
    select {
        case tempMechanic := <-tempMechanicChan:
            temporaryMechanic = tempMechanic
        case err := <-errChan:
            log.Println("error getting temporary mechanics: ", err.Error())
            return models.Booking{}, http.StatusInternalServerError, err
    }

    // found mechanic
    booking := models.Booking{
		UserID:     user.ID,
		Services:   services,
        Latitude: payload.Location.Lat,
        Longitude: payload.Location.Lng,
        Address: payload.Location.Address,
        AssignedMechanicID: temporaryMechanic.ID,
		Vehicle: models.Vehicle{
            Vtype: payload.VehicleDetails.VehicleType,
            Brand: payload.VehicleDetails.Brand,
            Size: payload.VehicleDetails.Size,
            Model: payload.VehicleDetails.Model,
		},
		BookingFee:     bookingFee.Price,
		BookingDate:    bookingDate,
		Status:         models.BookingPending,
	}

    //create booking
    newBooking, err := app.Store.Booking.Create(tx, booking)

    if err != nil {
        tx.Rollback()
        return models.Booking{}, http.StatusInternalServerError, err
    }

    tx.Commit()
    //set a cronjob to run after one hour to check if the booking has been assigned to a given mechanic
	return newBooking, http.StatusCreated, nil
}


func GetBookingFee(app *Application)  (BookingFeeResponse, int, error) {
	bookingFee, err := app.Store.Booking.GetBookingFee()

	if err != nil {
		return BookingFeeResponse{}, http.StatusInternalServerError, err
	}

	return BookingFeeResponse{Fee: bookingFee.Price}, http.StatusOK, nil
}


func GetBookingServices(app *Application) (*[]models.ServiceCategory, int, error) {
	serviceCategories, err := app.Store.Service.GetServiceCategories()
	
	if err != nil {
		return serviceCategories, http.StatusInternalServerError, err
	}
	return serviceCategories, http.StatusOK, nil
}


func GetDisabledDatesForUser(app *Application, user *models.User) ([]string, int, error) {
	dates := utils.GetNextNumDays(30)

	fmt.Println("dates", dates)
	pendingChan := make(chan []models.Booking)
	availableChan := make(chan []models.Mechanic)
	errorChan := make(chan error, 2) 

	go fetchPendingBookings(app, pendingChan, errorChan)
    go fetchAvailableMechanics(app, availableChan, errorChan)

	// collect result
	pendingBookings, availableMechanics, err := collectResults(pendingChan, availableChan, errorChan)
    if err != nil {
        return nil, http.StatusInternalServerError, err
    }

	// Create booking count map for each date
    bookingCountMap := make(map[string]int)
    for _, booking := range pendingBookings {
        date := booking.BookingDate.Truncate(24 * time.Hour).Format("2006-01-02")
        bookingCountMap[date]++
    }

	resultChan := make(chan string, len(dates))
    var wg sync.WaitGroup

	mechanicCount := len(availableMechanics)

    for _, dateStr := range dates {
        wg.Add(1)
        go func(dateStr string) {
            defer wg.Done()
            
            if isDateDisabledByCount(dateStr, bookingCountMap, mechanicCount) {
                resultChan <- dateStr
            }
        }(dateStr)
    }


	go func() {
        wg.Wait()
        close(resultChan)
    }()

    var disabledDates []string
    for date := range resultChan {
        disabledDates = append(disabledDates, date)
    }

    return disabledDates, http.StatusOK, nil

}


func getMechanicForBooking(app *Application, ch chan<- models.Mechanic, errCh chan<- error) {
    mechanic, err := app.Store.Mechanic.GetAvailableMechanic([]string{})

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            errCh <- fmt.Errorf("no available mechanics found that aren't blacklisted")
            return
        }
        errCh <- fmt.Errorf("error fetching mechanic: %v", err.Error())
        return
    }

    ch <- *mechanic
}


func fetchPendingBookings(app *Application, ch chan<- []models.Booking, errCh chan<- error) {
    bookings, err := app.Store.Booking.GetPendingBookings()
    if err != nil {
        errCh <- err
        return
    }
    ch <- *bookings
}

func fetchAvailableMechanics(app *Application, ch chan<- []models.Mechanic, errCh chan<- error) {
    mechanics, err := app.Store.Mechanic.GetAllAvailableMechanics()
    if err != nil {
        errCh <- err
        return
    }
    ch <- *mechanics
}

func collectResults(pendingChan <-chan []models.Booking, availableChan <-chan []models.Mechanic, errorChan <-chan error) ([]models.Booking, []models.Mechanic, error) {
    var pendingBookings []models.Booking
    var availableMechanics []models.Mechanic

    for i := 0; i < 2; i++ {
        select {
			case booked := <-pendingChan:
				pendingBookings = booked
			case available := <-availableChan:
				availableMechanics = available
			case err := <-errorChan:
				return nil, nil, err
        }
    }

    return pendingBookings, availableMechanics, nil
}

func isDateDisabledByCount(date string, bookingCountMap map[string]int, availableMechanicCount int) bool {
	if bookingCount, exists := bookingCountMap[date]; exists {
        return bookingCount > availableMechanicCount
    }
    return false
}