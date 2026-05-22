package main

import (
	"fmt"
	"time"
	"alhikmah-attendance-api/pkg/jwt"
)

func main() {
	secret := "supersecretjwttoken"
	token, err := jwt.GenerateToken("123", "admin", "admin", secret, 24*time.Hour)
	if err != nil {
		fmt.Println("Error generating:", err)
		return
	}
	fmt.Println("Generated:", token)

	claims, err := jwt.ValidateToken(token, secret)
	if err != nil {
		fmt.Println("Error validating:", err)
		return
	}
	fmt.Println("Validated:", claims.Username)
}
