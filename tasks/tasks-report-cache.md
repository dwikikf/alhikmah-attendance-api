## Relevant Files

- `internal/domain/report.go` - Report domain. Tambahkan interface `ReportCacheRepository`.
- `internal/repository/report_cache_postgres.go` - **[NEW]** Repository untuk menyimpan/mengambil cache report dari tabel `reports`.
- `internal/service/report_service.go` - Report service. Tambahkan logika cache: cek cache → jika ada return cache → jika tidak ada generate lalu simpan ke cache.
- `cmd/api/main.go` - Inject `ReportCacheRepository` ke `ReportService`.

### Notes

- Tabel `reports` sudah ada di database (lihat `migrations/001_init_schema.up.sql`, baris 103-117). Tabel ini menyimpan:
  - `report_type` (ENUM: 'harian', 'mingguan', 'bulanan', 'semesteran')
  - `class_id` (UUID)
  - `period_start` dan `period_end` (DATE)
  - `generated_by` (UUID) - user yang generate (diambil dari JWT claims)
  - `report_data` (JSONB) - data report yang di-cache
  - Constraint `UNIQUE(report_type, class_id, period_start, period_end)`
- **Strategi Cache:**
  1. Saat request report masuk, cek apakah ada record di tabel `reports` dengan key yang cocok (`report_type`, `class_id`, `period_start`, `period_end`).
  2. Jika ada DAN `generated_at` masih dalam batas waktu (misalnya < 1 jam), return data dari cache (`report_data` JSONB).
  3. Jika tidak ada atau sudah expired, generate report seperti biasa, lalu simpan hasilnya ke tabel `reports` (INSERT ... ON CONFLICT UPDATE).
- **Cache TTL:** 1 jam (3600 detik). Ini bisa dikonfigurasi lewat environment variable `REPORT_CACHE_TTL_SECONDS` (optional).
- Untuk rebuild dan test: `docker compose up -d --build`
- **PENTING:** `generated_by` memerlukan `userID` dari context JWT. Handler perlu mengirim `userID` ke service layer.

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [ ]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

Example:
- `- [ ] 1.1 Read file` → `- [x] 1.1 Read file` (after completing)

Update the file after completing each sub-task, not just after completing an entire parent task.

## Tasks

- [x] 1.0 Definisikan Cache Interface di Domain
  - [x] 1.1 Buka `internal/domain/report.go`. Tambahkan interface baru `ReportCacheRepository`:
    ```go
    type ReportCacheRepository interface {
        Get(reportType, classID, periodStart, periodEnd string) (json.RawMessage, error)
        Set(reportType, classID, periodStart, periodEnd, generatedBy string, data json.RawMessage) error
    }
    ```
  - [x] 1.2 Tambahkan import `encoding/json` di file tersebut.

- [x] 2.0 Implementasi Cache Repository
  - [x] 2.1 Buat file baru `internal/repository/report_cache_postgres.go`.
  - [x] 2.2 Implementasikan struct `reportCachePostgres` yang menerima `*sql.DB`.
  - [x] 2.3 Implementasikan method `Get(reportType, classID, periodStart, periodEnd string) (json.RawMessage, error)`:
    - Query: `SELECT report_data FROM reports WHERE report_type = $1 AND class_id = $2 AND period_start = $3 AND period_end = $4 AND generated_at > NOW() - INTERVAL '1 hour'`
    - Jika `sql.ErrNoRows`, return `nil, nil` (cache miss, bukan error).
    - Jika ada data, return `report_data`.
  - [x] 2.4 Implementasikan method `Set(reportType, classID, periodStart, periodEnd, generatedBy string, data json.RawMessage) error`:
    - Query: `INSERT INTO reports (report_type, class_id, period_start, period_end, generated_by, report_data) VALUES ($1::report_type_enum, $2, $3, $4, $5, $6) ON CONFLICT (report_type, class_id, period_start, period_end) DO UPDATE SET report_data = EXCLUDED.report_data, generated_by = EXCLUDED.generated_by, generated_at = NOW()`
  - [x] 2.5 Buat fungsi constructor `NewReportCacheRepository(db *sql.DB) domain.ReportCacheRepository`.

- [x] 3.0 Update Report Service untuk Menggunakan Cache
  - [x] 3.1 Buka `internal/service/report_service.go`. Update struct `reportService` untuk menyertakan field `cache domain.ReportCacheRepository`.
  - [x] 3.2 Update constructor `NewReportService` untuk menerima parameter `cacheRepo domain.ReportCacheRepository` tambahan.
  - [x] 3.3 Update method `GetDailyReport(classID, dateStr string) (*domain.DailyReport, error)`:
    - Sebelum generate, cek cache: `cached, err := s.cache.Get("harian", classID, dateStr, dateStr)`. Jika `cached != nil`, unmarshal ke `*domain.DailyReport` dan return.
    - Setelah generate report, marshal hasilnya ke JSON dan simpan: `data, _ := json.Marshal(report); s.cache.Set("harian", classID, dateStr, dateStr, "", data)`.
    - **Catatan:** Untuk saat ini, `generatedBy` bisa dikosongkan ("") atau di-hardcode. Idealnya, service menerima `userID` dari handler.
  - [x] 3.4 Update method `GetMonthlyReport(classID, monthStr string) (*domain.MonthlyReport, error)`:
    - Hitung `periodStart` (hari pertama bulan) dan `periodEnd` (hari terakhir bulan) dari `monthStr`.
    - Cek cache: `cached, err := s.cache.Get("bulanan", classID, periodStart, periodEnd)`.
    - Jika cache miss, generate lalu simpan.
  - [x] 3.5 Update method `GetSemesterReport(classID, academicYear string, semester int) (*domain.SemesterReport, error)`:
    - Hitung `periodStart` dan `periodEnd` berdasarkan semester (semester 1: Juli-Desember, semester 2: Januari-Juni).
    - Cek cache: `cached, err := s.cache.Get("semesteran", classID, periodStart, periodEnd)`.
    - Jika cache miss, generate lalu simpan.

- [x] 4.0 Update Dependency Injection di `main.go`
  - [x] 4.1 Buka `cmd/api/main.go`. Setelah `reportRepo := repository.NewReportRepository(db)`, tambahkan:
    ```go
    reportCacheRepo := repository.NewReportCacheRepository(db)
    ```
  - [x] 4.2 Update pemanggilan `NewReportService` untuk menyertakan `reportCacheRepo`:
    ```go
    reportService := service.NewReportService(reportRepo, studentRepo, reportCacheRepo)
    ```

- [x] 5.0 Build dan Verifikasi
  - [x] 5.1 Jalankan `docker compose up -d --build` untuk memastikan kode berhasil di-compile.
  - [x] 5.2 Test `GET /api/v1/reports/daily?class_id=...&date=2026-05-21` dua kali berturut-turut. Request kedua seharusnya lebih cepat karena menggunakan cache.
  - [x] 5.3 Verifikasi di database bahwa record baru muncul di tabel `reports` setelah request pertama:
    ```sql
    SELECT id, report_type, class_id, period_start, period_end, generated_at FROM reports;
    ```
