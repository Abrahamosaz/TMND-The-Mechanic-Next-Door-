package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"os"

	"github.com/thexovc/TMND/internal/utils"
)

type contextKey string

const userContextKey contextKey = "user"

const mechanicContextKey contextKey = "mechanic"

type uploadContext string

const uploadContextKey uploadContext = "uploadedFileURL"
const uploadMultipleFilesContextKey uploadContext = "uploadedFilesURL"

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
	key := userContextKey
	if role == "MECHANIC" {
		key = mechanicContextKey
	}

	ctx := context.WithValue(r.Context(), key, &user)
	next.ServeHTTP(w, r.WithContext(ctx))
}

// Create a middleware generator that takes a file key parameter
func (app *application) uploadMiddleware(fileKey string, folder string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to get the uploaded file using the provided key
			file, fileHeader, err := r.FormFile(fileKey)

			user, ok := app.GetUserFromContext(r)

			if !ok {
				app.responseJSON(http.StatusUnauthorized, w, "Unauthorized: No user found", nil)
				return
			}

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

			// Upload file to Cloudinary
			ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
			defer cancel()

			uploadResult, err := handleUploadToCloudinary(&utils.CloudinaryUploadParams{
				Ctx:        ctx,
				File:       file,
				FileHeader: fileHeader,
				Folder:     folder,
				User:       user,
			})

			if err != nil {
				log.Printf("%s", err.Error())
				app.responseJSON(http.StatusInternalServerError, w, "internal server error", nil)
				return
			}

			// Add uploaded file URL to request context
			ctx = context.WithValue(r.Context(), uploadContextKey, uploadResult)
			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Create a middleware generator that takes a file key parameter for multiple files
func (app *application) uploadMultipleFilesMiddleware(fileKey string, folder string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Parse the multipart form with a reasonable size limit (e.g., 32MB)
			if err := r.ParseMultipartForm(32 << 20); err != nil {
				log.Printf("Form parsing error: %s", err.Error())
				app.responseJSON(http.StatusBadRequest, w, "Unable to process the uploaded files. Please ensure each file is less than 32MB and try again.", nil)
				return
			}

			user, ok := app.GetUserFromContext(r)

			if !ok {
				app.responseJSON(http.StatusUnauthorized, w, "Unauthorized: No user found", nil)
				return
			}

			// Get all files for the given key
			files := r.MultipartForm.File[fileKey]
			if len(files) == 0 {
				// If no files are uploaded, continue without modifying the request context
				next.ServeHTTP(w, r)
				return
			}

			// Create a slice to store all upload results
			uploadResults := make([]utils.UploadResult, 0, len(files))
			// Create a context with timeout for all uploads
			ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
			defer cancel()

			// Process each file
			for _, fileHeader := range files {
				// Open the file
				file, err := fileHeader.Open()
				if err != nil {
					log.Printf("File open error: %s", err.Error())
					continue
				}

				// Upload file to Cloudinary
				uploadResult, err := handleUploadToCloudinary(&utils.CloudinaryUploadParams{
					Ctx:        ctx,
					File:       file,
					FileHeader: fileHeader,
					Folder:     folder,
					User:       user,
				})
				file.Close()

				if err != nil {
					log.Printf("%s", err.Error())
					continue
				}

				// Add to results
				uploadResults = append(uploadResults, *uploadResult)
			}

			// jsonBytes, _ := json.MarshalIndent(uploadResults, "", "\t")
			// log.Printf("uploadResults: %v all files", string(jsonBytes))

			// If no files were successfully uploaded, return an error
			if len(uploadResults) == 0 {
				app.responseJSON(http.StatusInternalServerError, w, "Failed to upload any files", nil)
				return
			}

			// Add uploaded files URLs to request context
			ctx = context.WithValue(r.Context(), uploadMultipleFilesContextKey, uploadResults)

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func handleUploadToCloudinary(params *utils.CloudinaryUploadParams) (*utils.UploadResult, error) {
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	cld := &utils.Cloudinary{URL: cloudinaryURL}

	// Upload file to Cloudinary using the utility
	uploadResult, err := cld.UploadFileToCloudinary(params)
	return uploadResult, err
}

// validateToken parses and verifies the JWT token
func validateToken(tokenString string) (*utils.MyCustomClaims, error) {
	claims, err := utils.DecodeJWT(tokenString)

	if err != nil {
		return claims, err
	}
	return claims, nil
}
