package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type PayloadClaims struct {
	OtpCode *string `json:"otpCode"`
}


type MyCustomClaims struct {
	PayloadClaims
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT token with user ID and expiration time
func GenerateJWT(userID string, expireTime time.Duration, payloadClaims PayloadClaims) (string, error) {

	claims := MyCustomClaims{
		payloadClaims,
		jwt.RegisteredClaims{
			// Also fixed dates can be used for the NumericDate
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	// Create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	mySigningKey := []byte(os.Getenv("SECRET_KEY"))
	signedToken, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}



func DecodeJWT(tokenString string) (*MyCustomClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	// Extract claims
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {

		if claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, errors.New("token has expired")
		}

		return claims, nil
	}
	return nil, errors.New("invalid token")
}