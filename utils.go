package main

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	config := LoadConfig()
	return token.SignedString([]byte(config.JWTSecret))
}

func ValidateJWT(tokenString string) (string, error) {
	config := LoadConfig()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["email"].(string), nil
	}
	return "", err
}
