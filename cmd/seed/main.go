package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"alhikmah-attendance-api/config"
	"alhikmah-attendance-api/pkg/database"

	"github.com/google/uuid"
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

	// Create 1 Guru
	guruID := createGuru(db)
	fmt.Printf("Created Guru with ID: %s\n", guruID)

	// Create 6 Classes
	classIDs := createClasses(db, guruID)
	fmt.Printf("Created 6 Classes\n")

	// Create 20 Students per Class
	var allStudents []studentData
	for i, classID := range classIDs {
		students := createStudents(db, classID, i+1, 20)
		allStudents = append(allStudents, students...)
		fmt.Printf("Created 20 students for Class %d\n", i+1)
	}

	// Create 1 month of attendance for May 2026
	createAttendances(db, allStudents, guruID)
	fmt.Println("Created attendance records for May 2026")

	fmt.Println("Seeding completed successfully!")
}

func createGuru(db *sql.DB) string {
	password := "password123"
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	hash := string(bytes)

	var id string
	err := db.QueryRow(`
		INSERT INTO users (username, email, password_hash, full_name, role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (username) DO UPDATE SET full_name = EXCLUDED.full_name
		RETURNING id
	`, "guru1", "guru1@school.com", hash, "Bapak Guru Budi", "guru", true).Scan(&id)
	
	if err != nil {
		// If exists and we can't return ID directly because of no update, let's select it
		db.QueryRow("SELECT id FROM users WHERE username = 'guru1'").Scan(&id)
	}
	return id
}

func createClasses(db *sql.DB, teacherID string) []string {
	var classIDs []string
	academicYear := "2025/2026"
	
	for i := 1; i <= 6; i++ {
		className := fmt.Sprintf("Kelas %d", i)
		var id string
		
		err := db.QueryRow(`
			INSERT INTO classes (class_name, teacher_id, academic_year, capacity)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (class_name, academic_year) DO UPDATE SET teacher_id = EXCLUDED.teacher_id
			RETURNING id
		`, className, teacherID, academicYear, 30).Scan(&id)

		if err != nil {
			db.QueryRow("SELECT id FROM classes WHERE class_name = $1 AND academic_year = $2", className, academicYear).Scan(&id)
		}
		classIDs = append(classIDs, id)
	}
	return classIDs
}

type studentData struct {
	ID      string
	ClassID string
}

func createStudents(db *sql.DB, classID string, classNum int, count int) []studentData {
	var students []studentData
	for i := 1; i <= count; i++ {
		nisn := fmt.Sprintf("100%d%02d%02d", classNum, i, rand.Intn(99))
		fullName := fmt.Sprintf("Siswa %d Kelas %d", i, classNum)
		qrData := fmt.Sprintf("%s|%s|%s", nisn, fullName, classID)
		gender := "laki-laki"
		if i%2 == 0 {
			gender = "perempuan"
		}

		var id string
		// Adding UUID directly if needed or let DB generate
		idStr := uuid.New().String()
		err := db.QueryRow(`
			INSERT INTO students (id, nisn, full_name, class_id, gender, qr_code_data, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (nisn, class_id) DO NOTHING
			RETURNING id
		`, idStr, nisn, fullName, classID, gender, qrData, true).Scan(&id)

		if err != nil { // Could be no rows returned because of conflict
			db.QueryRow("SELECT id FROM students WHERE nisn = $1 AND class_id = $2", nisn, classID).Scan(&id)
		}
		
		students = append(students, studentData{ID: id, ClassID: classID})
	}
	return students
}

func createAttendances(db *sql.DB, students []studentData, teacherID string) {
	// Seed 1 month of data: May 2026
	startDate := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	
	// Create a transaction for faster bulk insert
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	
	stmt, err := tx.Prepare(`
		INSERT INTO attendances (student_id, class_id, attendance_date, status, recorded_by, is_manual)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (student_id, class_id, attendance_date) DO NOTHING
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for day := 0; day < 31; day++ {
		currentDate := startDate.AddDate(0, 0, day)
		// Skip weekends (optional, but realistic)
		if currentDate.Weekday() == time.Saturday || currentDate.Weekday() == time.Sunday {
			continue
		}

		for _, student := range students {
			status := getRandomStatus()
			_, err = stmt.Exec(student.ID, student.ClassID, currentDate, status, teacherID, true)
			if err != nil {
				log.Printf("Error inserting attendance: %v", err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Error committing transaction: %v", err)
	}
}

func getRandomStatus() string {
	r := rand.Intn(100)
	if r < 85 {
		return "hadir"
	} else if r < 90 {
		return "izin"
	} else if r < 95 {
		return "sakit"
	}
	return "tanpa_keterangan"
}
