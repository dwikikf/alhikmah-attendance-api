package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"alhikmah-attendance-api/config"
	"alhikmah-attendance-api/pkg/database"

	"github.com/brianvoe/gofakeit/v6"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Inisialisasi GoFakeIt v6
	gofakeit.Seed(0)

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

	fmt.Println("Seeding database menggunakan brianvoe/gofakeit/v6...")

	seedDataWithGofakeit(db)

	fmt.Println("Seeding selesai dengan sukses!")
}

func seedDataWithGofakeit(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Could not begin transaction: %v", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
			log.Fatalf("Transaction failed: %v", err)
		} else {
			err = tx.Commit()
			if err != nil {
				log.Fatalf("Could not commit transaction: %v", err)
			}
		}
	}()

	// Hapus semua data yang ada sebelumnya agar bersih dan tidak conflict
	_, err = tx.Exec(`TRUNCATE TABLE reports, attendance_audits, attendances, students, class_teachers, classes, users RESTART IDENTITY CASCADE;`)
	if err != nil {
		log.Printf("Warning: Could not truncate tables: %v", err)
	}

	// 1. Buat 1 Superadmin
	createAdmin(tx)

	// 2. Buat 10 Guru menggunakan database internal GoFakeIt
	var teachers []string
	for i := 0; i < 10; i++ {
		fn := gofakeit.FirstName()
		ln := gofakeit.LastName()
		fullName := fn + " " + ln
		
		// Username tanpa spasi, tanpa tanda petik, dan huruf kecil
		username := strings.ToLower(fn + ln)
		username = strings.ReplaceAll(username, " ", "")
		username = strings.ReplaceAll(username, "'", "")

		id := createTeacher(tx, username, fullName)
		teachers = append(teachers, id)
	}

	// 3. Buat 6 Kelas
	classNames := []string{
		"Kelas 1 ( Al-Khawarizmi )",
		"Kelas 2 ( Ibnu Sina )",
		"Kelas 3 ( Al-Kindi )",
		"Kelas 4 ( Al-Farabi )",
		"Kelas 5 ( Ibnu Khaldun )",
		"Kelas 6 ( Jabir bin Hayyan )",
	}

	var classes []string
	academicYear := "2025/2026"

	for i, className := range classNames {
		teacherID := teachers[i] // 6 guru pertama jadi wali kelas

		classID := gofakeit.UUID()
		_, err = tx.Exec(`
			INSERT INTO classes (id, class_name, teacher_id, academic_year, capacity, description)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, classID, className, teacherID, academicYear, 40, gofakeit.Sentence(5))
		if err != nil {
			panic(err)
		}
		classes = append(classes, classID)

		// Hubungkan ke class_teachers
		_, err = tx.Exec(`
			INSERT INTO class_teachers (teacher_id, class_id, academic_year)
			VALUES ($1, $2, $3)
		`, teacherID, classID, academicYear)
		if err != nil {
			panic(err)
		}
	}

	// 4. Buat 30-40 Siswa per Kelas dan Data Kehadiran Juli 2025 sampai Juni 2026
	startDate := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)
	var dates []time.Time

	for d := startDate; d.Before(endDate) || d.Equal(endDate); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			dates = append(dates, d)
		}
	}

	statuses := []string{"hadir", "hadir", "hadir", "hadir", "hadir", "hadir", "hadir", "hadir", "izin", "sakit", "tanpa_keterangan"}
	studentCounter := 1

	for classIdx, classID := range classes {
		numStudents := gofakeit.Number(30, 40)
		fmt.Printf("Seeding class %d with %d students using GoFakeIt v6...\n", classIdx+1, numStudents)
		
		for i := 0; i < numStudents; i++ {
			studentID := gofakeit.UUID()
			nisn := fmt.Sprintf("100%07d", studentCounter)
			
			// Otomatis generate nama lengkap dari DB internal GoFakeIt
			fullName := gofakeit.Name()
			
			gender := "laki-laki"
			if gofakeit.Bool() {
				gender = "perempuan"
			}
			
			dob := gofakeit.DateRange(time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2015, 12, 31, 0, 0, 0, 0, time.UTC))

			_, err = tx.Exec(`
				INSERT INTO students (id, nisn, full_name, class_id, date_of_birth, gender, photo_url, qr_code_data, is_active)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`, studentID, nisn, fullName, classID, dob, gender, "", nisn, true)
			if err != nil {
				panic(err)
			}
			studentCounter++

			// Kehadiran per siswa untuk tiap hari kerja
			for _, date := range dates {
				status := statuses[gofakeit.Number(0, len(statuses)-1)]
				attendanceID := gofakeit.UUID()
				teacherID := teachers[classIdx]

				_, err = tx.Exec(`
					INSERT INTO attendances (id, student_id, class_id, attendance_date, status, recorded_by, recorded_at, is_manual)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				`, attendanceID, studentID, classID, date, status, teacherID, date.Add(time.Hour*8), gofakeit.Bool())
				if err != nil {
					panic(err)
				}

				// Fake Audit secara random (2% chance)
				if gofakeit.Number(1, 100) <= 2 {
					oldStatus := "tanpa_keterangan"
					reason := gofakeit.Sentence(4) 
					_, err = tx.Exec(`
						INSERT INTO attendance_audits (attendance_id, old_status, new_status, changed_by, changed_at, reason)
						VALUES ($1, $2, $3, $4, $5, $6)
					`, attendanceID, oldStatus, status, teacherID, date.Add(time.Hour*10), reason)
					if err != nil {
						panic(err)
					}
				}
			}
		}

		// Fake Report untuk bulan Juli 2025 sampai Juni 2026
		reportDates := []time.Time{
			time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		}

		for _, startPeriod := range reportDates {
			endPeriod := startPeriod.AddDate(0, 1, -1)
			
			teacherID := teachers[classIdx]

			reportData := map[string]interface{}{
				"total_students": numStudents,
				"total_days":     20,
				"summary":        "Report generated by GoFakeIt - " + gofakeit.Sentence(6),
			}
			reportJSON, _ := json.Marshal(reportData)

			_, err = tx.Exec(`
				INSERT INTO reports (report_type, class_id, period_start, period_end, generated_by, generated_at, report_data)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`, "bulanan", classID, startPeriod, endPeriod, teacherID, endPeriod.Add(time.Hour*24), reportJSON)
			if err != nil {
				panic(err)
			}
		}
	}
}

func createAdmin(tx *sql.Tx) {
	password := "superadmin"
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	hash := string(bytes)
	id := gofakeit.UUID()

	_, err := tx.Exec(`
		INSERT INTO users (id, username, email, password_hash, full_name, role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (username) DO UPDATE SET 
			full_name = EXCLUDED.full_name,
			password_hash = EXCLUDED.password_hash,
			role = EXCLUDED.role,
			email = EXCLUDED.email
	`, id, "superadmin", "superadmin@alhikmah.sch.id", hash, "Super Admin", "admin", true)

	if err != nil {
		panic(fmt.Errorf("error creating admin: %v", err))
	}
}

func createTeacher(tx *sql.Tx, username, fullName string) string {
	password := "password123"
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	hash := string(bytes)
	
	id := gofakeit.UUID()
	email := username + "@alhikmah.sch.id"

	_, err := tx.Exec(`
		INSERT INTO users (id, username, email, password_hash, full_name, role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (username) DO UPDATE SET 
			full_name = EXCLUDED.full_name,
			password_hash = EXCLUDED.password_hash,
			role = EXCLUDED.role,
			email = EXCLUDED.email
		RETURNING id
	`, id, username, email, hash, fullName, "teacher", true)

	if err != nil {
		errQuery := tx.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&id)
		if errQuery != nil {
			panic(fmt.Errorf("error creating teacher %s: %v", username, err))
		}
	}
	
	return id
}