package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"alhikmah-attendance-api/config"
)

func ConnectDB(cfg config.Config) (*sql.DB, error) {
	// Support DATABASE_URL (e.g., from Neon) as a single connection string.
	// Falls back to individual env vars if DATABASE_URL is not set.
	var dsn string
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		dsn = dbURL
		slog.Info("Connecting to database via DATABASE_URL")
	} else {
		sslmode := cfg.DBSSLMode
		if sslmode == "" {
			if cfg.AppEnv == "production" {
				sslmode = "require"
			} else {
				sslmode = "disable"
			}
		}
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, sslmode)
		slog.Info("Connecting to database via individual env vars", "host", cfg.DBHost, "sslmode", sslmode)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not ping the database: %w", err)
	}

	// Optimasi khusus untuk Vercel Serverless
	// Vercel membekukan (freeze) container saat tidak ada request.
	// Jika kita menyimpan koneksi menganggur (idle), koneksi tersebut akan terputus paksa oleh Neon.
	// Ini menyebabkan error 500 pada request pertama, lalu sukses pada request kedua (karena membuka koneksi baru).
	if os.Getenv("VERCEL") == "1" {
		db.SetMaxIdleConns(0)
	}

	slog.Info("Connected to PostgreSQL successfully")
	return db, nil
}
