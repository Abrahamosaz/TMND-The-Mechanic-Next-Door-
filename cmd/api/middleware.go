package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Abrahamosaz/TMND/internal/utils"
)


type contextKey string
const userContextKey contextKey = "user"


func (app  *application) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		// get the user from the id
		user, err := app.store.User.FindByID(claims.Subject)
		if err != nil {
			app.responseJSON(http.StatusUnauthorized, w, "Unauthorized: Invalid token", nil)
			return
		}

		// Attach user to context
		ctx := context.WithValue(r.Context(), userContextKey, &user)
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