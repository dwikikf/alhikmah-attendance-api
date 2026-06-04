## Relevant Files

### Backend (core/*)

**Domain Layer**
- `core/domain/class.go` - Entity Class + Repository Interface + Service Interface (split ke file terpisah)
- `core/domain/student.go` - Entity Student + Repository Interface + Service Interface
- `core/domain/user.go` - Entity User + Repository Interface + Service Interface
- `core/domain/attendance.go` - Entity Attendance + Repository Interface + Service Interface
- `core/domain/enrollment.go` - Entity Enrollment + Repository Interface + Service Interface
- `core/domain/class_teacher.go` - Entity ClassTeacher + Repository Interface + Service Interface
- `core/domain/report.go` - Entity Report + Repository Interface + Service Interface

**DTO Layer (Input/Output per use-case)**
- `core/dto/class_dto.go` - CreateClassRequest, UpdateClassRequest, ClassResponse
- `core/dto/student_dto.go` - CreateStudentRequest, UpdateStudentRequest, StudentResponse
- `core/dto/user_dto.go` - CreateUserRequest, UpdateUserRequest, UserResponse
- `core/dto/attendance_dto.go` - ScanQRRequest, ManualAttendanceRequest, AttendanceResponse
- `core/dto/enrollment_dto.go` - EnrollRequest, PromoteRequest, TransferRequest, EnrollmentResponse
- `core/dto/class_teacher_dto.go` - AssignTeacherRequest, ClassTeacherResponse
- `core/dto/report_dto.go` - ReportQueryRequest, DailyReportResponse, MonthlyReportResponse, SemesterReportResponse

**Repository Layer**
- `core/repository/class_postgres.go` - Implementasi ClassRepository
- `core/repository/student_postgres.go` - Implementasi StudentRepository
- `core/repository/user_postgres.go` - Implementasi UserRepository
- `core/repository/attendance_postgres.go` - Implementasi AttendanceRepository
- `core/repository/enrollment_postgres.go` - Implementasi EnrollmentRepository
- `core/repository/class_teacher_postgres.go` - Implementasi ClassTeacherRepository
- `core/repository/report_postgres.go` - Implementasi ReportRepository

**Service / UseCase Layer**
- `core/service/class_service.go` - Implementasi ClassService
- `core/service/student_service.go` - Implementasi StudentService
- `core/service/user_service.go` - Implementasi UserService
- `core/service/attendance_service.go` - Implementasi AttendanceService
- `core/service/enrollment_service.go` - Implementasi EnrollmentService
- `core/service/class_teacher_service.go` - Implementasi ClassTeacherService
- `core/service/report_service.go` - Implementasi ReportService

**Handler / Delivery Layer**
- `core/handler/class_handler.go` - Handler HTTP untuk Class (tidak boleh import domain)
- `core/handler/student_handler.go` - Handler HTTP untuk Student (tidak boleh import domain)
- `core/handler/user_handler.go` - Handler HTTP untuk User (tidak boleh import domain)
- `core/handler/attendance_handler.go` - Handler HTTP untuk Attendance (tidak boleh import domain)
- `core/handler/enrollment_handler.go` - Handler HTTP untuk Enrollment (tidak boleh import domain)
- `core/handler/class_teacher_handler.go` - Handler HTTP untuk ClassTeacher (tidak boleh import domain)
- `core/handler/report_handler.go` - Handler HTTP untuk Report (tidak boleh import domain)
- `core/handler/auth.go` - Handler HTTP untuk Auth

**Test Files**
- `core/service/*_service_test.go` - Unit test service (masing-masing entity)
- `core/handler/*_handler_test.go` - Unit test handler (masing-masing entity)

---

### Notes

- Struktur target adalah **Clean Architecture** (mendekati), bukan Hexagonal murni.
- Aturan utama: **Dependency hanya boleh ke dalam** (Handler → Service → Repository → Domain)
- `domain` hanya berisi **Entity struct** dan **Interface** (Repository + Service). Tidak ada logic bisnis.
- `dto` berisi **Request struct** dan **Response struct** untuk setiap endpoint.
- **Handler tidak boleh mengimport `domain`** secara langsung untuk membuat entity (`domain.Class{}`). Semua mapping dilakukan di Service.
- **Service tidak boleh mengimport `handler` atau `dto`** dari handler. Service menerima parameter primitif atau domain entity.
- Repository hanya berinteraksi dengan database, tidak ada logic bisnis.
- Gunakan `go test ./...` untuk menjalankan semua test.

---

## Instructions for Completing Tasks

**IMPORTANT:** Setiap task yang sudah selesai, ubah `- [ ]` menjadi `- [x]`. Lakukan setelah setiap sub-task selesai, bukan hanya setelah parent task selesai.

---

## Tasks

- [x] 0.0 Buat Feature Branch
  - [x] 0.1 Buat dan checkout branch baru: `git checkout -b refactor/clean-architecture`

- [x] 1.0 Refactor Domain Layer — Pisahkan Interface dari Entity
  - [x] 1.1 Pelajari isi setiap file di `core/domain/` untuk memahami apa yang perlu dipisah
  - [x] 1.2 Pastikan setiap domain file **hanya berisi**: (a) Entity struct, (b) Repository Interface, (c) Service/UseCase Interface
  - [x] 1.3 Hapus semua logic bisnis yang ada di domain (jika ada) — pindahkan ke service
  - [x] 1.4 Verifikasi tidak ada import ke package lain (selain standard library) di dalam `core/domain/`
  - [x] 1.5 Jalankan `go build ./...` untuk memastikan tidak ada compile error

- [x] 2.0 Refactor DTO Layer — Pisahkan Request dan Response per Entity
  - [x] 2.1 Buat/update `core/dto/class_dto.go`: tambah `ClassResponse` struct yang merepresentasikan output API (field yang dikembalikan ke FE)
  - [x] 2.2 Buat/update `core/dto/student_dto.go`: tambah `StudentResponse` struct
  - [x] 2.3 Buat/update `core/dto/user_dto.go`: tambah `UserResponse` struct
  - [x] 2.4 Buat/update `core/dto/attendance_dto.go`: tambah `AttendanceResponse` struct
  - [x] 2.5 Buat/update `core/dto/enrollment_dto.go`: tambah `EnrollmentResponse` dan `EnrollmentHistoryResponse` struct
  - [x] 2.6 Buat `core/dto/class_teacher_dto.go`: `AssignTeacherRequest`, `ClassTeacherResponse`
  - [x] 2.7 Buat/update `core/dto/report_dto.go`: `DailyReportResponse`, `MonthlyReportResponse`, `SemesterReportResponse`
  - [x] 2.8 Pastikan semua DTO **tidak mengimport domain** — hanya tipe primitif (string, int, bool, time.Time, dll)
  - [x] 2.9 Jalankan `go build ./...`

- [x] 3.0 Refactor Service Layer — Service Menerima DTO, Mengembalikan DTO
  - [x] 3.1 Update signature interface `ClassService` di `core/domain/class.go` agar metode Create dan Update menerima DTO bukan `*domain.Class`
  - [x] 3.2 Update implementasi `classService` di `core/service/class_service.go`: mapping dari DTO ke entity di dalam service, panggil repo
  - [x] 3.3 Update signature interface dan implementasi `StudentService` (sama seperti class)
  - [x] 3.4 Update signature interface dan implementasi `UserService`
  - [x] 3.5 Update signature interface dan implementasi `AttendanceService`
  - [x] 3.6 Update signature interface dan implementasi `EnrollmentService`
  - [x] 3.7 Update signature interface dan implementasi `ClassTeacherService`
  - [x] 3.8 Update signature interface dan implementasi `ReportService`
  - [x] 3.9 Pastikan setiap service method yang mengembalikan data, mengembalikan **DTO Response** bukan `*domain.Entity`
  - [x] 3.10 Jalankan `go test ./core/service/...` dan pastikan semua test lulus

- [x] 4.0 Refactor Handler Layer — Handler Hanya Tahu DTO, Tidak Import Domain
  - [x] 4.1 Update `core/handler/class_handler.go`: hapus semua `domain.Class{}` struct literal — cukup bind ke DTO request, serahkan ke service
  - [x] 4.2 Update `core/handler/student_handler.go`: hapus semua `domain.Student{}` struct literal
  - [x] 4.3 Update `core/handler/user_handler.go`: hapus semua `domain.User{}` struct literal
  - [x] 4.4 Update `core/handler/enrollment_handler.go`: hapus semua `domain.PromoteItem{}` dan entity domain lainnya
  - [x] 4.5 Update `core/handler/class_teacher_handler.go`: hapus semua `domain.ClassTeacher{}` struct literal
  - [x] 4.6 Update `core/handler/report_handler.go` dan `report_exporter.go`: pastikan hanya pakai DTO response
  - [x] 4.7 Update `core/handler/attendance_handler.go`: pastikan tidak mengimport domain entity secara langsung
  - [x] 4.8 Hapus import `"alhikmah-attendance-api/core/domain"` dari semua handler file
  - [x] 4.9 Jalankan `go build ./...` — tidak boleh ada compile error
  - [x] 4.10 Jalankan `go test ./core/handler/...` dan pastikan semua test lulus

- [x] 5.0 Update Dependency Injection di `cmd/api/main.go`
  - [x] 5.1 Pastikan semua constructor (`NewXxxHandler`, `NewXxxService`, `NewXxxRepository`) masih compatible dengan signature yang baru
  - [x] 5.2 Sesuaikan pemanggilan di `main.go` jika ada perubahan constructor parameter
  - [x] 5.3 Jalankan `go build -o api-server cmd/api/main.go` dan pastikan build sukses

- [x] 6.0 Update dan Tambah Unit Tests
  - [x] 6.1 Update `core/service/*_service_test.go` sesuai signature service yang baru (mock, input DTO, output DTO)
  - [x] 6.2 Update `core/handler/*_handler_test.go` sesuai perubahan handler (tidak pakai domain)
  - [x] 6.3 Tambah test case untuk error path (validasi DTO, not found, conflict, dll)
  - [x] 6.4 Jalankan `go test ./...` — semua test harus lulus (target: 0 FAIL)

- [x] 7.0 Final Verification & Cleanup
  - [x] 7.1 Jalankan `grep -rn '"alhikmah-attendance-api/core/domain"' core/handler/` — harus kosong (0 hasil)
  - [x] 7.2 Jalankan `go vet ./...` — tidak ada warning/error
  - [x] 7.3 Jalankan `go test ./...` — semua PASS
  - [x] 7.4 Test manual: Start server (`go run cmd/api/main.go`) dan verifikasi endpoint CRUD Class, Student, User, Attendance tidak error
  - [x] 7.5 Commit perubahan ke branch `refactor/clean-architecture`
