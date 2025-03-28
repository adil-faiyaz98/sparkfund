package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func main() {
	// Create a new token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "test-user-123"
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	claims["iat"] = time.Now().Unix()

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		return
	}

	fmt.Printf("Test JWT Token:\n%s\n", tokenString)
} 