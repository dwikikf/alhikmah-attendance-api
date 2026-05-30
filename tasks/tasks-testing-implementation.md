## Relevant Files

### API (Go)
- `internal/domain/attendance.go` - Interface `AttendanceRepository` dan `AttendanceService` yang akan di-mock.
- `internal/domain/user.go` - Interface `UserRepository` dan `UserService`.
- `internal/domain/student.go` - Interface `StudentRepository`.
- `internal/domain/class.go` - Interface `ClassRepository`.
- `internal/service/attendance_service.go` - Business logic absensi utama (QR Scan, Manual).
- `internal/service/attendance_service_test.go` - **[BUAT BARU]** Unit tests untuk `attendanceService`.
- `internal/service/user_service.go` - Business logic user (create, update, validasi role).
- `internal/service/user_service_test.go` - **[BUAT BARU]** Unit tests untuk `userService`.
- `internal/service/student_service.go` - Business logic siswa.
- `internal/service/student_service_test.go` - **[BUAT BARU]** Unit tests untuk `studentService`.
- `internal/handler/auth.go` - Handler login, refresh, logout.
- `internal/handler/auth_test.go` - **[BUAT BARU]** Handler tests dengan `httptest`.
- `internal/handler/attendance_handler.go` - Handler attendance.
- `internal/handler/attendance_handler_test.go` - **[BUAT BARU]** Handler tests dengan `httptest`.
- `internal/mocks/` - **[BUAT BARU]** Direktori berisi mock yang di-generate oleh `mockery`.
- `pkg/utils/` - Helper functions (hash, format).

### Notes

- Jalankan semua tests: `go test ./... -v`
- Jalankan test satu package: `go test ./internal/service/... -v`
- Generate mocks: `go run github.com/vektra/mockery/v2@latest --all --dir internal/domain --output internal/mocks --outpkg mocks`
- Lihat coverage: `go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out`

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [ ]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

Update the file after completing each sub-task, not just after completing an entire parent task.

## Tasks

- [x] 0.0 Create feature branch
  - [x] 0.1 Buat dan checkout branch baru: `git checkout -b test/unit-tests`

- [x] 1.0 Setup testing infrastructure & dependencies
  - [x] 1.1 Tambahkan `github.com/stretchr/testify` ke `go.mod` dengan perintah: `go get github.com/stretchr/testify`
  - [x] 1.2 Pastikan `mockery` tersedia: `go run github.com/vektra/mockery/v2@latest --version`
  - [x] 1.3 Generate semua mock dari interfaces di `internal/domain/`: `go run github.com/vektra/mockery/v2@latest --all --dir internal/domain --output internal/mocks --outpkg mocks`
  - [x] 1.4 Verifikasi folder `internal/mocks/` terbuat dan berisi file mock (misal: `MockAttendanceRepository.go`, `MockUserRepository.go`, dll.)

- [x] 2.0 Implement Service Layer Unit Tests (`internal/service`)
  - [x] 2.1 Buat file `internal/service/attendance_service_test.go`
  - [x] 2.2 Tulis test `TestProcessQRScan_DuplicateCachePrevention`: pastikan error dikembalikan jika NISN sudah ada di cache
  - [x] 2.3 Tulis test `TestProcessQRScan_StudentNotFound`: mock `GetByNISN` mengembalikan error, pastikan service mengembalikan error "Siswa tidak ditemukan"
  - [x] 2.4 Tulis test `TestProcessQRScan_UnauthorizedTeacher`: role bukan admin dan `IsTeacherResponsibleForStudent` return false, pastikan error "Guru tidak memiliki akses"
  - [x] 2.5 Tulis test `TestProcessQRScan_AlreadyScannedToday`: mock `GetByClassAndDate` return attendance yang sudah ada untuk studentID tersebut
  - [x] 2.6 Tulis test `TestProcessQRScan_Success`: semua mock return data valid, pastikan `MarkAttendance` terpanggil sekali dan tidak ada error
  - [x] 2.7 Tulis test `TestProcessManualAttendance_InvalidStatus`: kirim status selain `hadir/izin/sakit/tanpa_keterangan`, pastikan error
  - [x] 2.8 Tulis test `TestProcessManualAttendance_Success`: semua mock valid, pastikan attendance berhasil tersimpan
  - [x] 2.9 Buat file `internal/service/user_service_test.go`
  - [x] 2.10 Tulis test `TestUserCreate_MissingRequiredFields`: buat user tanpa username/email/nama, pastikan error "missing required fields"
  - [x] 2.11 Tulis test `TestUserCreate_InvalidRole`: buat user dengan role selain `admin`/`teacher`, pastikan error "invalid role"
  - [x] 2.12 Tulis test `TestUserCreate_Success`: semua field valid, mock `repo.Create` return nil, pastikan tidak ada error
  - [x] 2.13 Tulis test `TestUserUpdate_MissingID`: update tanpa ID, pastikan error "missing user ID"
  - [x] 2.14 Tulis test `TestGetAll_PaginationDefaults`: kirim `page=0` dan `limit=200`, pastikan fallback ke page=1 dan limit=10

- [x] 3.0 Implement Handler Layer Tests (`internal/handler`)
  - [x] 3.1 Buat file `internal/handler/auth_test.go`
  - [x] 3.2 Setup helper `newTestRouter()` yang membuat Gin router dengan mode `test` menggunakan `gin.SetMode(gin.TestMode)`
  - [x] 3.3 Tulis test `TestLogin_BadRequest`: kirim body JSON yang tidak valid (misal: kosong), ekspektasi response `400 Bad Request`
  - [x] 3.4 Tulis test `TestLogout_Success`: panggil `POST /logout`, ekspektasi response `200 OK` dan cookie `refresh_token` ter-clear (MaxAge = -1)
  - [x] 3.5 Tulis test `TestRefresh_NoCookie`: panggil `POST /refresh` tanpa cookie, ekspektasi response `401 Unauthorized`
  - [x] 3.6 Tulis test `TestParseTokenDuration_ValidInput`: unit test untuk helper `parseTokenDuration` dengan input `"1h"`, `"24h"`, `""`, dan nilai tidak valid
  - [x] 3.7 Buat file `internal/handler/attendance_handler_test.go`
  - [x] 3.8 Tulis test `TestScanQR_MissingNISN`: kirim request tanpa body NISN, ekspektasi `400`
  - [x] 3.9 Tulis test `TestManualAttendance_InvalidStatus`: kirim status yang salah, ekspektasi `400` dengan pesan error yang sesuai

- [x] 4.0 Verify test coverage & run all tests
  - [x] 4.1 Jalankan `go test ./... -v` dan pastikan semua test lulus (PASS)
  - [x] 4.2 Jalankan `go test ./... -coverprofile=coverage.out` untuk generate coverage report
  - [x] 4.3 Jalankan `go tool cover -func=coverage.out` untuk melihat persentase coverage per fungsi
  - [x] 4.4 Pastikan coverage pada package `internal/service` minimal **70%**
  - [x] 4.5 Fix test yang gagal jika ada
