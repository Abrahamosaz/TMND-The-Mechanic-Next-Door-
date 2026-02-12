package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/thexovc/TMND/internal/models"
	"github.com/thexovc/TMND/internal/utils"
	"gorm.io/gorm"
)

// booking cronJobs
type BookingStatusChecker struct {
	App           *Application
	Booking       *models.Booking
	Scheduler     *cron.Cron
	NextCheckTime time.Time
	// maxRetries int // Optional: limit the number of retries
}

func CheckBookingStatusJob(app *Application, booking *models.Booking, currentTime time.Time) {
	checker := &BookingStatusChecker{
		App:           app,
		Booking:       booking,
		NextCheckTime: currentTime.Add(time.Hour),
	}

	// Start the first check after 1 hour
	go checker.scheduleBookingNextCheck()
}

func (b *BookingStatusChecker) scheduleBookingNextCheck() {
	b.stopScheduler()
	// Create a new cron scheduler
	b.Scheduler = cron.New()

	// Create cron expression for the next check
	cronExpression := fmt.Sprintf("%d %d %d %d %d",
		b.NextCheckTime.Minute(),
		b.NextCheckTime.Hour(),
		b.NextCheckTime.Day(),
		b.NextCheckTime.Month(),
		b.NextCheckTime.Weekday(),
	)

	_, err := b.Scheduler.AddFunc(cronExpression, func() {
		b.checkBookingStatus()
	})

	if err != nil {
		// Log the error but don't stop the application
		fmt.Printf("Error scheduling booking status check for booking %v: %v\n",
			b.Booking.ID, err)
		return
	}

	log.Printf("CronJob started for booking ID %v - Next check scheduled for %v",
		b.Booking.ID,
		b.NextCheckTime.Format("2006-01-02 15:04:05"))
	b.Scheduler.Start()
}

func (b *BookingStatusChecker) checkBookingStatus() {
	err := b.App.Store.Booking.GetBooking(b.Booking)

	if err != nil {
		fmt.Printf("Error fetching booking %v: %v\n", b.Booking.ID, err)
		// Even on error, we'll try again
		b.NextCheckTime = time.Now().Add(2 * time.Minute)
		b.scheduleBookingNextCheck()
		return
	}

	if b.Booking.MechanicID != nil {
		// Mechanic has been assigned, stop the scheduler
		fmt.Printf("Booking %v has been assigned to mechanic, stopping scheduler\n",
			b.Booking.ID,
		)
		b.stopScheduler()
		return
	}

	// assign the booking to another mechanic
	var blackListedIDS []string
	var visitedIDS []string

	err = utils.DecodeJSONSlice(b.Booking.BlacklistedMechanics, &blackListedIDS)
	if err != nil {
		fmt.Println("Error decoding blackListedIDS JSON:", err)
		b.stopScheduler()
		return
	}

	err = utils.DecodeJSONSlice(b.Booking.VisitedMechanics, &visitedIDS)
	if err != nil {
		fmt.Println("Error decoding visitedIDS JSON:", err)
		b.stopScheduler()
		return
	}

	mechanic, err := b.App.Store.Mechanic.GetAvailableMechanic(blackListedIDS, visitedIDS)

	if err != nil {
		fmt.Printf("error fetching mechanic: %v\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Printf("no available mechanics found that aren't blacklisted\n")
			b.Booking.ErrorMessage = "No mechanic available at this time"
			b.App.Store.Booking.UpdateBooking(b.Booking)
			b.stopScheduler()
			return
		}

		if err.Error() == "no available mechanics found" {
			emptyJSON, _ := json.Marshal([]string{})
			b.Booking.VisitedMechanics = emptyJSON
			b.App.Store.Booking.UpdateBooking(b.Booking)
		}

		b.NextCheckTime = time.Now().Add(2 * time.Minute)
		b.scheduleBookingNextCheck()
		return
	}

	newVisitedIDS := append(visitedIDS, mechanic.ID.String())
	newVisitedJsonIDS, err := utils.EncodeJSONSlice(newVisitedIDS)

	if err != nil {
		fmt.Println("Error encoding newVisitedIDS JSON:", err)
		b.stopScheduler()
		return
	}

	b.Booking.VisitedMechanics = newVisitedJsonIDS
	b.Booking.AssignedMechanicID = mechanic.ID
	b.NextCheckTime = time.Now().Add(time.Hour)
	go b.App.Store.Booking.UpdateBooking(b.Booking)

	// Booking still not assigned, schedule next check
	b.scheduleBookingNextCheck()
}

func (b *BookingStatusChecker) stopScheduler() {
	if b.Scheduler != nil {
		b.Scheduler.Stop()
		b.Scheduler = nil
	}
}

func StartAllBookingCronJobs(app *Application) {
	// get all the pending bookings
	bookings, err := app.Store.Booking.GetPendingBookings()

	fmt.Println("all bookigns", len(*bookings))
	if err != nil {
		return
	}

	for _, booking := range *bookings {
		if booking.NextExecutionTime == nil {
			CheckBookingStatusJob(app, &booking, time.Now().Add(-58*time.Minute))
			continue
		}

		if time.Now().Before(*booking.NextExecutionTime) {
			executionTime := *booking.NextExecutionTime
			CheckBookingStatusJob(app, &booking, executionTime.Add(-time.Hour))
		} else {
			CheckBookingStatusJob(app, &booking, time.Now().Add(-58*time.Minute))
		}
	}
}

//end of booking cronJobs
