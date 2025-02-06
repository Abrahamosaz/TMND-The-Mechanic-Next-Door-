package utils

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	mathRand "math/rand"
	"strings"
	"time"
)


func GenerateOtpCode(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)

	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[num.Int64()]
	}

	return string(otp), nil
}



func PtrToString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func StringToPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}


func GetUploadErrorMessage(err error) string {
	var message string;
	
	if errors.Is(err, context.DeadlineExceeded) {
        log.Printf("Cloudinary upload timeout: %s", err.Error())
        message = "The upload took too long. Please try again with a smaller file or check your internet connection."
    } else if strings.Contains(err.Error(), "format not allowed") {
        log.Printf("Cloudinary upload error (format not allowed): %s", err.Error())
        message = "Unsupported file format. Please upload a valid image (JPG, PNG, etc.)."
    } else if strings.Contains(err.Error(), "too large") {
        log.Printf("Cloudinary upload error (file too large): %s", err.Error())
        message = "The file is too large. Please upload a smaller image."
    } else {
        log.Printf("Cloudinary upload error: %s", err.Error())
        message = "There was an issue uploading your image. Please try again later."
    }

	return message
}


func GetNextNumDays(days int) []string {
	var dates []string
	today := time.Now()
	tomorrow := today.AddDate(0, 0, 1)

	for i := 0; i <= days; i++ {
		date := tomorrow.AddDate(0, 0, i) // Add 'i' days to today
		dates = append(dates, date.Format("2006-01-02")) 
	}

	return dates
}


func GenerateUniquePaymentRef() string {
	r := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	// Get the current Unix timestamp
	currentTime := time.Now().Unix()
	randomNumber := r.Intn(900000) + 100000
	transactionRef := fmt.Sprintf("REF%d%d", currentTime, randomNumber)
	return transactionRef
}
