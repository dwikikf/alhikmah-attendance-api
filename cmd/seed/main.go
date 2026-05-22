package main

import (
	"database/sql"
	"fmt"
	"log"

	"alhikmah-attendance-api/config"
	"alhikmah-attendance-api/pkg/database"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// 2. Connect to Database
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("Seeding database...")

	// Create 1 Admin
	createAdmin(db)
	fmt.Println("Created Admin: Dwiki Kausar Fahmi")

	// Create 1 Guru
	createGuru(db)
	fmt.Println("Created Guru: Anida Qoriyati Nur")

	fmt.Println("Seeding completed successfully!")
}

func createAdmin(db *sql.DB) {
	password := "superadmin"
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	hash := string(bytes)

	_, err := db.Exec(`
		INSERT INTO users (username, email, password_hash, full_name, role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (username) DO UPDATE SET 
			full_name = EXCLUDED.full_name,
			password_hash = EXCLUDED.password_hash,
			role = EXCLUDED.role,
			email = EXCLUDED.email
	`, "superadmin", "superadmin@school.com", hash, "Dwiki Kausar Fahmi", "admin", true)
	
	if err != nil {
		log.Printf("Error creating admin: %v\n", err)
	}
}

func createGuru(db *sql.DB) {
	password := "anidaqn"
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	hash := string(bytes)

	_, err := db.Exec(`
		INSERT INTO users (username, email, password_hash, full_name, role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (username) DO UPDATE SET 
			full_name = EXCLUDED.full_name,
			password_hash = EXCLUDED.password_hash,
			role = EXCLUDED.role,
			email = EXCLUDED.email
	`, "anidaqn", "anidaqn@school.com", hash, "Anida Qoriyati Nur", "teacher", true)
	
	if err != nil {
		log.Printf("Error creating guru: %v\n", err)
	}
}
