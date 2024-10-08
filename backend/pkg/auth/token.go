package auth

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// ValidateToken parses and validates a JWT token and returns the user ID
func ValidateToken(tokenString string) (int, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		log.Printf("Token validation failed: %v", err)
		if err == jwt.ErrSignatureInvalid {
			return 0, errors.New("invalid token signature")
		}
		return 0, fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		log.Println("Invalid token")
		return 0, errors.New("invalid token")
	}

	// Kontrollige, kas token ei ole aegunud
	if claims.ExpiresAt < time.Now().Unix() {
		log.Println("Token has expired")
		return 0, errors.New("token has expired")
	}

	return claims.UserID, nil
}
