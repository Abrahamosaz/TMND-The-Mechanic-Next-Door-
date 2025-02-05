package services

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Abrahamosaz/TMND/internal/models"
	"github.com/Abrahamosaz/TMND/internal/utils"
)


func CreateUserBooking(app *Application, payload CreateBooking, user *models.User) (int, error) {
	return http.StatusCreated, nil
}



func GetBookingFee(app *Application)  (BookingFeeResponse, int, error) {
	bookingFee, err := app.Store.Booking.GetBookingFee()

	if err != nil {
		return BookingFeeResponse{}, http.StatusInternalServerError, err
	}

	return BookingFeeResponse{Fee: bookingFee.Price}, http.StatusOK, nil
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