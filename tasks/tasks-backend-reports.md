# Task Implementasi: Attendance Reports & Specific Date Attendance

Dokumen ini berisi panduan teknis mendetail untuk mengimplementasikan fitur laporan absensi (Daily, Monthly, Semester) dan pengambilan data absensi untuk tanggal spesifik. Panduan ini dirancang agar dapat diimplementasikan secara langsung oleh model AI/developer.

## 1. Persiapan Struktur Domain & Interface
**Target File:** `internal/domain/report.go` (Buat baru jika belum ada) & `internal/domain/attendance.go`

*   **Action:** Definisikan struktur data untuk balasan (response) laporan dan tambahkan metode di interface repository/service.
*   **Struktur Data:**
    *   `DailyReportResponse`: Menyimpan daftar siswa beserta status absensinya (Hadir, Izin, Sakit, Alpa) pada satu hari tertentu.
    *   `MonthlyReportResponse`: Menyimpan ringkasan total (Hadir, Izin, Sakit, Alpa) per siswa selama satu bulan penuh.
    *   `SemesterReportResponse`: Menyimpan ringkasan total per siswa berdasarkan rentang tanggal semester (Ganjil/Genap) pada tahun ajaran tertentu.
*   **Update Interface `AttendanceRepository` & `AttendanceService` (di `internal/domain/attendance.go`):**
    *   Tambahkan fungsi: `GetByClassAndDate(classID string, date time.Time) ([]*Attendance, error)`

## 2. Implementasi Endpoint: Attendance by Date
**Target:** `internal/repository/attendance_postgres.go`, `internal/service/attendance_service.go`, `internal/handler/attendance_handler.go`, `cmd/api/main.go`

*   **Endpoint:** `GET /api/v1/attendances/:class_id/:date`
*   **Repository:** Implementasikan `GetByClassAndDate`. Lakukan `SELECT` ke tabel `attendances` dengan filter `class_id = $1 AND attendance_date = $2`.
*   **Service:** Panggil fungsi repository tersebut.
*   **Handler:** Buat metode `GetByClassAndDate` di `AttendanceHandler`. Ambil parameter URL `:class_id` dan `:date` (konversi string `YYYY-MM-DD` ke `time.Time`). Panggil service dan return JSON.
*   **Routing (main.go):** Tambahkan `protected.GET("/attendances/:class_id/:date", attendanceHandler.GetByClassAndDate)`

## 3. Implementasi Endpoint: Daily Report
**Target:** `internal/handler/report_handler.go` (Buat baru), `internal/service/report_service.go`, `internal/repository/report_postgres.go`

*   **Endpoint:** `GET /api/v1/reports/daily?class_id=...&date=YYYY-MM-DD`
*   **Query SQL Logika (Repository):** Lakukan `LEFT JOIN` antara tabel `students` (difilter berdasarkan `class_id`) dan tabel `attendances` (difilter berdasarkan `attendance_date`). Ini memastikan semua siswa di kelas tersebut muncul di laporan, meskipun mereka belum diabsen (status bisa direturn null atau "belum_absen").
*   **Handler:** Ambil query parameter `class_id` dan `date`. Validasi format `date`.

## 4. Implementasi Endpoint: Monthly Report
**Target:** `internal/handler/report_handler.go`, `internal/service/report_service.go`, `internal/repository/report_postgres.go`

*   **Endpoint:** `GET /api/v1/reports/monthly?class_id=...&month=YYYY-MM`
*   **Query SQL Logika (Repository):** Lakukan agregasi `COUNT` berdasarkan `status`. Filter rentang waktu dari hari pertama di bulan tersebut hingga hari terakhir.
    *   *Contoh filter:* `attendance_date >= '2026-05-01' AND attendance_date < '2026-06-01'`
    *   *Group by:* `student_id`. Return daftar siswa beserta total masing-masing status.

## 5. Implementasi Endpoint: Semester Report
**Target:** `internal/handler/report_handler.go`, `internal/service/report_service.go`, `internal/repository/report_postgres.go`

*   **Endpoint:** `GET /api/v1/reports/semester?class_id=...&semester=1&academic_year=2024/2025`
*   **Logika Penentuan Tanggal (Service):**
    *   Jika `semester=1` (Ganjil) dan `academic_year=2024/2025`: Rentang waktu biasanya Juli 2024 hingga Desember 2024.
    *   Jika `semester=2` (Genap) dan `academic_year=2024/2025`: Rentang waktu biasanya Januari 2025 hingga Juni 2025.
    *   Tugas *Service* adalah menerjemahkan `semester` dan `academic_year` menjadi `startDate` dan `endDate`.
*   **Query SQL Logika (Repository):** Mirip dengan Monthly Report, gunakan agregasi `COUNT` dengan filter rentang `startDate` dan `endDate` yang diberikan oleh Service.

## 6. Update Routing Utama
**Target:** `cmd/api/main.go`

Daflarkan service dan handler baru untuk laporan:
```go
reportRepo := repository.NewReportRepository(db)
reportService := service.NewReportService(reportRepo, studentRepo)
reportHandler := handler.NewReportHandler(reportService)

// Di dalam blok protected routes:
protected.GET("/reports/daily", reportHandler.GetDailyReport)
protected.GET("/reports/monthly", reportHandler.GetMonthlyReport)
protected.GET("/reports/semester", reportHandler.GetSemesterReport)
```
