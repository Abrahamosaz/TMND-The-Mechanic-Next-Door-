package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"os"

	"github.com/Abrahamosaz/TMND/internal/utils"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)


type contextKey string
const userContextKey contextKey = "user"

const mechanicContextKey contextKey = "mechanic"

type uploadContext string
const uploadContextKey uploadContext = "uploadedFileURL"


type UploadResult struct {
	URL      string
	FileName string
}



func (app *application) userAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.authenticateJWT(r, w, next, "USER")
	})
}



func (app *application) mechanicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.authenticateJWT(r, w, next, "MECHANIC")
	})
}


func (app *application) authenticateJWT(r *http.Request, w http.ResponseWriter, next http.Handler, role string) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			app.responseJSON(http.StatusUnauthorized, w, "Unauthorized: Missing Authorization header", nil)
			return
		}

		// Token format should be "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.responseJSON(http.StatusUnauthorized, w, "Unauthorized: Invalid token format", nil)
			return
		}

		tokenString := parts[1]

		// Parse and validate the token
		claims, err := validateToken(tokenString)
		if err != nil {

			message := fmt.Sprintf("Unauthorized: %s", err.Error())
			app.responseJSON(http.StatusUnauthorized, w, message, nil)
			return
		}

		payloadRole := claims.Role

		if payloadRole != role {
			app.responseJSON(http.StatusUnauthorized, w, "Unauthorized: Invalid token", nil)
			return
		}
		
		// get the user from the id
		user, err := app.store.User.FindByID(claims.Subject)
		if err != nil {
			app.responseJSON(http.StatusUnauthorized, w, "Unauthorized: Invalid token", nil)
			return
		}

		// Attach user to context
		ctx := context.WithValue(r.Context(), userContextKey, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
}




func (app  *application) uploadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get the uploaded file
		file, fileHeader, err := r.FormFile("profile-image")

		if err != nil {
			// If no file is uploaded, continue without modifying the request context
			if err == http.ErrMissingFile {
				next.ServeHTTP(w, r)
				return
			}
			log.Printf("File upload error: %s", err.Error())
			app.responseJSON(http.StatusBadRequest, w, utils.GetUploadErrorMessage(err), nil)
			return
		}
		defer file.Close()

		// Initialize Cloudinary

		cloudinaryUrl := os.Getenv("CLOUDINARY_URL")
		cld, err := cloudinary.NewFromURL(cloudinaryUrl) // Replace with your Cloudinary URL
		if err != nil {
			log.Printf("Cloudinary setup error: %s", err.Error())
			app.responseJSON(http.StatusInternalServerError, w, "internal server error", nil)
			return
		}

		// Upload file to Cloudinary
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()
		uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
			PublicID: fileHeader.Filename, // Use original filename
			Folder:   CLOUDINARY_PROFILE_IMAGE_FOLDER,           // Upload to "uploads" folder
			AllowedFormats: CLOUDINARY_ALLOWED_FORMATS,
			Transformation: CLOUDINARY_PROFILE_PICTURE_TRANSFORMATION,
		})

		if err != nil {
			log.Printf("Cloudinary setup error: %s", err.Error())
			app.responseJSON(http.StatusInternalServerError, w, "internal server error", nil)
			return
		}

		// Add uploaded file URL to request context
		// Define uploadResult struct with exported fields (Uppercase)
		ctx = context.WithValue(r.Context(), uploadContextKey, &UploadResult{
			URL:      uploadResult.SecureURL, // Correct field name and reference
			FileName: uploadResult.OriginalFilename,    // Assign original filename
		})

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


// validateToken parses and verifies the JWT token
func validateToken(tokenString string) (*utils.MyCustomClaims, error) {
	claims, err := utils.DecodeJWT(tokenString)

	if err != nil {
		return claims, err
	}
	return claims, nil
}