package core

import (
	"bombayv/logiapp-monorepo/logi_api/internal/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims defines the structure of the JWT payload.
// It includes the standard registered claims provided by the jwt-go library,
// as well as custom claims specific to this application (UserID, Role).
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT for a given user ID and role.
// The token will have an expiration time of 24 hours.
func GenerateJWT(userID, role string) (string, error) {
	// Define the expiration time for the token (24 hours from now).
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the JWT claims, which includes the user ID, role, and expiration time.
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds.
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			// Issuer of the token.
			Issuer: "logiapp",
		},
	}

	// Create a new token object, specifying the signing method (HS256) and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key to get the complete, signed token string.
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT parses and validates a token string.
// It returns the claims if the token is valid, otherwise it returns an error.
func ValidateJWT(tokenString string) (*Claims, error) {
	// Initialize a new instance of the Claims struct.
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key function, which validates the signing algorithm
	// and provides the secret key for verification.
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// return the secret key
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	// Check if the token is valid.
	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	// Return the claims from the token.
	return claims, nil
}
