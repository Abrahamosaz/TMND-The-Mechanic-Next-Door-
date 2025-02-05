package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)



func ValidateRequestBody(err error, w http.ResponseWriter) {
	validationErrors := make(map[string]string)

	// Loop through validation errors and generate custom messages
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
			case "required":
				validationErrors[err.Field()] = fmt.Sprintf("%s is required", err.Field())
			case "email":
				validationErrors[err.Field()] = fmt.Sprintf("%s must be a valid email", err.Field())
			case "gte":
				validationErrors[err.Field()] = fmt.Sprintf("%s must be greater than or equal to %s", err.Field(), err.Param())
			case "lte":
				validationErrors[err.Field()] = fmt.Sprintf("%s must be less than or equal to %s", err.Field(), err.Param())
			case "eqfield":
				validationErrors[err.Field()] = fmt.Sprintf("%s must match %s", err.Field(), err.Param())
			case "min":
				validationErrors[err.Field()] = fmt.Sprintf("%s must contain items greater than %s", err.Field(), err.Param())
			}
	}

	// Custom response format
	response := map[string]interface{}{
		"message":    "Bad Request",
		"statusCode": http.StatusBadRequest, // HTTP status code
		"errors":     validationErrors,
	}

	// Return validation errors as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}