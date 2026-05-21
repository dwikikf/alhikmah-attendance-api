package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"alhikmah-attendance-api/config"
)

func ConnectDB(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not ping the database: %w", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return db, nil
}
