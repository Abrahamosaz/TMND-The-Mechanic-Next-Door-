package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func (app *application) handleMonnifyWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Unable to read request body", err)
		app.responseJSON(http.StatusBadRequest, w, "Unable to read request body", nil)
		return
	}
	defer r.Body.Close()

	// Get the signature from request header
	signature := r.Header.Get("monnify-signature")
	if signature == "" {
		fmt.Println("Missing monnify-signature header")
		app.responseJSON(http.StatusBadRequest, w, "Missing monnify-signature header", nil)
		return
	}

	// Compute SHA-512 HMAC hash
	clientSecret := os.Getenv("MONNIFY_SECRET_KEY")
	h := hmac.New(sha512.New, []byte(clientSecret))
	h.Write(body)
	computedHash := hex.EncodeToString(h.Sum(nil))

    fmt.Println("computedHash", computedHash)
    fmt.Println("signature", signature)

	// Compare the computed hash with the received signature
	if computedHash != signature {
		fmt.Println("Invalid signature")
		app.responseJSON(http.StatusUnauthorized, w, "Invalid signature", nil)
		return
	}

    var eventData map[string]any
    err = json.Unmarshal(body, &eventData)
    if err != nil {
        fmt.Println("Unable to unmarshal event data", err)
        app.responseJSON(http.StatusBadRequest, w, "Unable to unmarshal event data", nil)
        return
    }

    fmt.Println("event data body",  eventData)

    serviceApp := app.createNewServiceApp()
    err = serviceApp.HandleMonnifyWebhook(eventData)

    if err != nil {
        fmt.Println("error handling monnify webhook", err)
        app.responseJSON(http.StatusInternalServerError, w, "error handling monnify webhook", nil)
        return
    }

    app.responseJSON(http.StatusOK, w, "monnify webhook handled successfully", nil)
}

