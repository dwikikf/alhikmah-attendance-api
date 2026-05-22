## Relevant Files

- `migrations/002_add_soft_delete.up.sql` - Migration file to add `deleted_at` column to users, students, classes tables.
- `migrations/002_add_soft_delete.down.sql` - Rollback migration to remove `deleted_at` column.
- `internal/domain/user.go` - User struct and interfaces. Add `DeletedAt` field, add `SoftDelete` method to repository/service.
- `internal/domain/student.go` - Student struct and interfaces. Add `DeletedAt` field, add `SoftDelete` method to repository/service.
- `internal/domain/class.go` - Class struct and interfaces. Add `DeletedAt` field, add `SoftDelete` method to repository/service.
- `internal/repository/user_postgres.go` - User repository. Add soft delete query, update all SELECT queries to exclude deleted records.
- `internal/repository/student_postgres.go` - Student repository. Add soft delete query, update all SELECT queries to exclude deleted records.
- `internal/repository/class_postgres.go` - Class repository. Add soft delete query, update all SELECT queries to exclude deleted records.
- `internal/service/user_service.go` - User service. Add `SoftDelete` method.
- `internal/service/student_service.go` - Student service. Add `SoftDelete` method.
- `internal/service/class_service.go` - Class service. Add `SoftDelete` method.
- `internal/handler/user_handler.go` - User handler. Add `DELETE /users/:user_id` endpoint.
- `internal/handler/student_handler.go` - Student handler. Add `DELETE /students/:student_id` endpoint.
- `internal/handler/class_handler.go` - Class handler. Add `DELETE /classes/:class_id` endpoint.
- `cmd/api/main.go` - Register the new DELETE routes.

### Notes

- Soft delete menggunakan pola `deleted_at TIMESTAMP NULL`. Jika `deleted_at IS NOT NULL`, record dianggap "deleted".
- JANGAN gunakan hard delete (`DELETE FROM ...`). Seluruh penghapusan harus menggunakan `UPDATE ... SET deleted_at = NOW()`.
- Semua query SELECT yang sudah ada harus ditambahkan filter `AND deleted_at IS NULL` agar record yang sudah dihapus tidak muncul.
- Kolom `is_active` yang sudah ada tetap digunakan untuk fitur non-aktifkan (bukan hapus). Soft delete adalah layer terpisah.
- Untuk rebuild dan test: `docker compose up -d --build`

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [ ]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

Example:
- `- [ ] 1.1 Read file` → `- [x] 1.1 Read file` (after completing)

Update the file after completing each sub-task, not just after completing an entire parent task.

## Tasks

- [x] 1.0 Buat Migration untuk Soft Delete
  - [x] 1.1 Buat file `migrations/002_add_soft_delete.up.sql` yang menambahkan kolom `deleted_at TIMESTAMP NULL DEFAULT NULL` ke tabel `users`, `students`, dan `classes`.
  - [x] 1.2 Buat file `migrations/002_add_soft_delete.down.sql` yang menghapus kolom `deleted_at` dari ketiga tabel tersebut.
  - [x] 1.3 Tambahkan index pada kolom `deleted_at` untuk ketiga tabel: `CREATE INDEX idx_users_deleted_at ON users(deleted_at);` (dan seterusnya).

- [x] 2.0 Update Domain Structs dan Interfaces
  - [x] 2.1 Buka `internal/domain/user.go`. Tambahkan field `DeletedAt *time.Time \`json:"deleted_at,omitempty"\`` pada struct `User`. Tambahkan method `SoftDelete(id string) error` pada interface `UserRepository` dan `UserService`.
  - [x] 2.2 Buka `internal/domain/student.go`. Tambahkan field `DeletedAt *time.Time \`json:"deleted_at,omitempty"\`` pada struct `Student`. Tambahkan method `SoftDelete(id string) error` pada interface `StudentRepository` dan `StudentService`.
  - [x] 2.3 Buka `internal/domain/class.go`. Tambahkan field `DeletedAt *time.Time \`json:"deleted_at,omitempty"\`` pada struct `Class`. Tambahkan method `SoftDelete(id string) error` pada interface `ClassRepository` dan `ClassService`.

- [x] 3.0 Update Repository Layer - Tambah Filter `deleted_at IS NULL`
  - [x] 3.1 Buka `internal/repository/user_postgres.go`. Update SEMUA query `SELECT` yang ada agar menambahkan kondisi `AND u.deleted_at IS NULL` (atau `AND deleted_at IS NULL` jika tanpa alias). Ini mencakup `GetByID`, `GetAll`, dan query login di `internal/handler/auth.go` (jika ada query langsung di sana).
  - [x] 3.2 Tambahkan method `SoftDelete(id string) error` pada `userPostgres` yang menjalankan: `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`.
  - [x] 3.3 Buka `internal/repository/student_postgres.go`. Update SEMUA query `SELECT` yang ada agar menambahkan kondisi `AND s.deleted_at IS NULL`. Ini mencakup `Create`, `GetByID`, `GetByClassID`, `GetByNISN`, dan `GetAll`.
  - [x] 3.4 Tambahkan method `SoftDelete(id string) error` pada `studentPostgres`.
  - [x] 3.5 Buka `internal/repository/class_postgres.go`. Update SEMUA query `SELECT` yang ada agar menambahkan kondisi `AND c.deleted_at IS NULL`. Ini mencakup `GetAll` dan `GetByID`.
  - [x] 3.6 Tambahkan method `SoftDelete(id string) error` pada `classPostgres`.

- [x] 4.0 Update Service Layer
  - [x] 4.1 Buka `internal/service/user_service.go`. Tambahkan method `SoftDelete(id string) error` yang memanggil `repo.SoftDelete(id)`.
  - [x] 4.2 Buka `internal/service/student_service.go`. Tambahkan method `SoftDelete(id string) error` yang memanggil `repo.SoftDelete(id)`.
  - [x] 4.3 Buka `internal/service/class_service.go`. Tambahkan method `SoftDelete(id string) error` yang memanggil `repo.SoftDelete(id)`. **Penting:** Sebelum menghapus kelas, periksa apakah masih ada siswa aktif di kelas tersebut. Jika ada, kembalikan error "cannot delete class with active students".

- [x] 5.0 Update Handler Layer - Tambah Endpoint DELETE
  - [x] 5.1 Buka `internal/handler/user_handler.go`. Tambahkan method `Delete(c *gin.Context)` yang mengambil `user_id` dari URL param, lalu memanggil `service.SoftDelete(userID)`. Return `{ "success": true, "message": "User deleted successfully" }`.
  - [x] 5.2 Buka `internal/handler/student_handler.go`. Tambahkan method `Delete(c *gin.Context)` yang mengambil `student_id` dari URL param, lalu memanggil `service.SoftDelete(studentID)`. Return `{ "success": true, "message": "Student deleted successfully" }`.
  - [x] 5.3 Buka `internal/handler/class_handler.go`. Tambahkan method `Delete(c *gin.Context)` yang mengambil `class_id` dari URL param, lalu memanggil `service.SoftDelete(classID)`. Return `{ "success": true, "message": "Class deleted successfully" }`.

- [x] 6.0 Register Routes di `main.go`
  - [x] 6.1 Buka `cmd/api/main.go`. Tambahkan route berikut di dalam blok `protected`:
    - `protected.DELETE("/users/:user_id", middleware.RoleMiddleware("admin"), userHandler.Delete)`
    - `protected.DELETE("/students/:student_id", middleware.RoleMiddleware("admin"), studentHandler.Delete)`
    - `protected.DELETE("/classes/:class_id", middleware.RoleMiddleware("admin"), classHandler.Delete)`

- [x] 7.0 Build dan Verifikasi
  - [x] 7.1 Jalankan `docker compose up -d --build` untuk memastikan kode berhasil di-compile dan migrasi `002` berjalan otomatis.
  - [x] 7.2 Verifikasi bahwa endpoint `GET /api/v1/students`, `GET /api/v1/users`, `GET /api/v1/classes` tetap berfungsi normal (record yang dihapus tidak muncul).
