package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/innovelabs/microtools-go/internal/config"
)

// GenerateJWT generates a JWT token for a user
func GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	cfg := config.LoadConfig()
	return token.SignedString([]byte(cfg.JWTSecret))
}

// ValidateJWT validates a JWT token
func ValidateJWT(tokenString string) (string, error) {
	cfg := config.LoadConfig()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["email"].(string), nil
	}
	return "", err
}
