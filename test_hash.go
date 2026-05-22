package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash := "$2a$10$7R9rW7S0O6bX6b9XkZ8uOedvW86N6lRBlfWp9rY7X.W7YmZ8Kz2OW"
	passwords := []string{"password", "admin", "admin123", "super_admin", "123456", "12345678", "qwerty", "alhikmah", "secret"}
	for _, p := range passwords {
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(p))
		if err == nil {
			fmt.Println("MATCH FOUND:", p)
			return
		}
	}
	fmt.Println("No match found")
}
