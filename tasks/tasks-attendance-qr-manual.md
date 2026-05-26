## Relevant Files

- `internal/domain/attendance.go` - Domain models for attendance.
- `internal/dto/attendance.go` - Data Transfer Objects for QR and Manual requests.
- `internal/repository/attendance_repository.go` - Database operations for inserting attendance.
- `internal/repository/class_repository.go` - Database operations for validating teacher class assignments.
- `internal/repository/student_repository.go` - Database operations for validating student NISN.
- `internal/service/attendance_service.go` - Business logic including cache validation for QR and single-entry manual attendance.
- `internal/handler/attendance_handler.go` - API route handlers for QR and manual endpoints.
- `pkg/cache/cache.go` - Setup for temporary in-memory caching for QR code scans.

### Notes

- Unit tests should typically be placed alongside the code files they are testing.
- Use `go test ./...` to run tests since this is a Go project.

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [ ]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

Example:
- `- [ ] 1.1 Read file` → `- [x] 1.1 Read file` (after completing)

## Tasks

- [x] 0.0 Create feature branch
  - [x] 0.1 Create and checkout a new branch (e.g., `git checkout -b feature/attendance-qr-manual`)
- [x] 1.0 Setup In-Memory Cache Mechanism
  - [x] 1.1 Create `pkg/cache/cache.go` to implement an in-memory cache (e.g., using a map with mutex or a library like `go-cache`) to store scanned QR codes temporarily.
  - [x] 1.2 Implement `Set(key string, value interface{}, expiration time.Duration)` method.
  - [x] 1.3 Implement `Get(key string) (interface{}, bool)` method.
- [x] 2.0 Update DTOs for Attendance Requests
  - [x] 2.1 Add `QRScanRequest` DTO in `internal/dto/attendance.go` containing `NISN` and `TeacherID`.
  - [x] 2.2 Add `ManualAttendanceRequest` DTO containing `StudentID`, `Status` (Present, Absent, etc.), and `TeacherID`.
  - [x] 2.3 Add validation tags to DTOs.
- [x] 3.0 Implement Repository Methods for Teacher Validation and Student Lookup
  - [x] 3.1 In `internal/repository/class_repository.go`, add `IsTeacherResponsibleForStudent` to check if the student belongs to a class assigned to the teacher.
  - [x] 3.2 In `internal/repository/student_repository.go`, add a method to find student details by NISN (`FindByNISN`).
  - [x] 3.3 Ensure `attendance_repository.go` has a method to insert a single attendance record.
- [x] 4.0 Implement Service Logic for QR Code Attendance (with caching)
  - [x] 4.1 In `internal/service/attendance_service.go`, inject the Cache and required repositories.
  - [x] 4.2 Create `ProcessQRScan` method.
  - [x] 4.3 Inside `ProcessQRScan`, check if the NISN exists in the cache to prevent duplicate processing within a short time. Return a "success" or "ignored" response immediately if cached.
  - [x] 4.4 If not in cache, validate via `IsTeacherResponsibleForStudent`. Return error if invalid.
  - [x] 4.5 Look up the student ID via NISN.
  - [x] 4.6 Insert the attendance record into the database via the repository.
  - [x] 4.7 Save the NISN to the cache with a short expiration (e.g., 10 seconds).
- [x] 5.0 Implement Service Logic for Manual Attendance (one-by-one)
  - [x] 5.1 In `internal/service/attendance_service.go`, create `ProcessManualAttendance` method.
  - [x] 5.2 Inside `ProcessManualAttendance`, look up student by ID to get NISN (or just use ID if teacher validation supports it).
  - [x] 5.3 Validate via `IsTeacherResponsibleForStudent`. Return error if invalid.
  - [x] 5.4 Insert the attendance record into the database with the given status.
- [x] 6.0 Create Handlers and Register Routes
  - [x] 6.1 In `internal/handler/attendance_handler.go`, create `HandleQRScan` endpoint to parse request and call `ProcessQRScan`.
  - [x] 6.2 Create `HandleManualAttendance` endpoint to parse request and call `ProcessManualAttendance`.
  - [x] 6.3 Register new endpoints in the router (e.g., `POST /api/attendance/qr` and `POST /api/attendance/manual`).
