package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Abrahamosaz/TMND/internal/models"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// booking cronJobs
type BookingStatusChecker struct {
	App 			*Application
    Booking   		*models.Booking
    Scheduler  		*cron.Cron
    NextCheckTime  	time.Time
    // maxRetries int // Optional: limit the number of retries
}


func CheckBookingStatusJob(app *Application, booking *models.Booking) {
	checker := &BookingStatusChecker{
		App:			app,
        Booking:    	booking,
        NextCheckTime:  time.Now().Add(time.Hour),
    }
    
    // Start the first check after 1 hour
    go checker.scheduleBookingNextCheck()
}



func (b *BookingStatusChecker) scheduleBookingNextCheck() {
    // Create a new cron scheduler
    b.Scheduler = cron.New()
    
    // Create cron expression for the next check
    cronExpression := fmt.Sprintf("%d %d %d %d %d ?",
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

    b.Scheduler.Start()
}



func (b *BookingStatusChecker) checkBookingStatus() {
	// // update the booking execution time
	// b.Booking.NextExecutionTime = &b.NextCheckTime
	// go b.App.Store.Booking.UpdateBooking(b.Booking)

	err := b.App.Store.Booking.GetBooking(b.Booking)

	if err != nil {
		fmt.Printf("Error fetching booking %v: %v\n", b.Booking.ID, err)
        // Even on error, we'll try again in an hour
        b.NextCheckTime = time.Now().Add(time.Hour)
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
	mechanic, err := b.App.Store.Mechanic.GetAvailableMechanic(b.Booking.BlacklistedMechanics)

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            fmt.Printf("no available mechanics found that aren't blacklisted\n")
			b.stopScheduler()
            return
        }
		fmt.Printf("error fetching mechanic: %v\n", err.Error())
		b.stopScheduler()
        return
    }

	b.Booking.AssignedMechanicID = mechanic.ID 
	go b.App.Store.Booking.UpdateBooking(b.Booking)
	
    // Booking still not assigned, schedule next check
    b.NextCheckTime = time.Now().Add(time.Hour)
    b.scheduleBookingNextCheck()
}



func (b *BookingStatusChecker) stopScheduler() {
    if b.Scheduler != nil {
        b.Scheduler.Stop()
        b.Scheduler = nil
    }
}



func StartAllCronJobs() {
}

//end of booking cronJobs