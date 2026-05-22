## Relevant Files

- `internal/domain/class.go` - Class struct dan interfaces. Tambahkan method `Create` dan `Update` ke `ClassRepository` dan `ClassService`.
- `internal/repository/class_postgres.go` - Class repository. Implementasi `Create` dan `Update` dengan query SQL.
- `internal/service/class_service.go` - Class service. Implementasi `Create` dan `Update` logic.
- `internal/handler/class_handler.go` - Class handler. Tambahkan handler `Create` dan `Update`.
- `cmd/api/main.go` - Register route `POST /classes` dan `PUT /classes/:class_id`.

### Notes

- Berdasarkan PRD (prd.md), `POST /classes` hanya bisa diakses oleh **admin**.
- Berdasarkan PRD, `PUT /classes/:class_id` bisa diakses oleh **admin** (update info kelas).
- Constraint `UNIQUE(class_name, academic_year)` sudah ada di database. Jika ada duplikat, return error yang jelas.
- Untuk rebuild dan test: `docker compose up -d --build`
- Field yang diterima saat Create class (berdasarkan PRD):
  ```json
  {
    "class_name": "2B",
    "teacher_id": "UUID",
    "academic_year": "2024/2025",
    "capacity": 32,
    "description": "Kelas 2 Angkatan B"
  }
  ```
- Field yang diterima saat Update class (berdasarkan PRD):
  ```json
  {
    "class_name": "2B Updated",
    "teacher_id": "UUID",
    "capacity": 35,
    "description": "Deskripsi baru"
  }
  ```

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [x]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

Example:
- `- [x] 1.1 Read file` → `- [x] 1.1 Read file` (after completing)

Update the file after completing each sub-task, not just after completing an entire parent task.

## Tasks

- [x] 1.0 Update Domain Interface untuk Class
  - [x] 1.1 Buka `internal/domain/class.go`. Tambahkan method berikut ke interface `ClassRepository`:
    - `Create(class *Class) error`
    - `Update(class *Class) error`
  - [x] 1.2 Tambahkan method berikut ke interface `ClassService`:
    - `Create(class *Class) error`
    - `Update(class *Class) error`

- [x] 2.0 Implementasi Repository Layer
  - [x] 2.1 Buka `internal/repository/class_postgres.go`. Tambahkan method `Create(class *Class) error` yang menjalankan query:
    ```sql
    INSERT INTO classes (class_name, teacher_id, academic_year, capacity, description)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, created_at, updated_at
    ```
    Scan hasil `RETURNING` ke `class.ID`, `class.CreatedAt`, `class.UpdatedAt`.
  - [x] 2.2 Tambahkan method `Update(class *Class) error` yang menjalankan query:
    ```sql
    UPDATE classes SET class_name = $1, teacher_id = $2, capacity = $3, description = $4, updated_at = NOW()
    WHERE id = $5
    RETURNING updated_at
    ```
    Scan `updated_at` ke `class.UpdatedAt`. Jika `RowsAffected == 0`, return error "class not found".

- [x] 3.0 Implementasi Service Layer
  - [x] 3.1 Buka `internal/service/class_service.go`. Tambahkan method `Create(class *Class) error`:
    - Validasi bahwa `class.ClassName` tidak kosong.
    - Validasi bahwa `class.TeacherID` tidak kosong.
    - Validasi bahwa `class.AcademicYear` tidak kosong dan format-nya valid (misal `2024/2025`).
    - Jika `class.Capacity` kosong/0, set default ke `30`.
    - Panggil `repo.Create(class)`.
  - [x] 3.2 Tambahkan method `Update(class *Class) error`:
    - Validasi bahwa `class.ID` tidak kosong.
    - Panggil `repo.Update(class)`.

- [x] 4.0 Implementasi Handler Layer
  - [x] 4.1 Buka `internal/handler/class_handler.go`. Tambahkan method `Create(c *gin.Context)`:
    - Bind JSON request body ke struct: `class_name` (required), `teacher_id` (required), `academic_year` (required), `capacity` (optional), `description` (optional).
    - Panggil `service.Create(class)`.
    - Return `201 Created` dengan `{ "success": true, "data": class }`.
  - [x] 4.2 Tambahkan method `Update(c *gin.Context)`:
    - Ambil `class_id` dari URL param `c.Param("class_id")`.
    - Bind JSON request body.
    - Set `class.ID = classID`.
    - Panggil `service.Update(class)`.
    - Return `200 OK` dengan `{ "success": true, "data": class, "message": "Class updated successfully" }`.

- [x] 5.0 Register Routes
  - [x] 5.1 Buka `cmd/api/main.go`. Tambahkan route berikut di dalam blok `protected`:
    - `protected.POST("/classes", middleware.RoleMiddleware("admin"), classHandler.Create)`
    - `protected.PUT("/classes/:class_id", middleware.RoleMiddleware("admin"), classHandler.Update)`

- [x] 6.0 Build dan Verifikasi
  - [x] 6.1 Jalankan `docker compose up -d --build` untuk memastikan kode berhasil di-compile.
  - [x] 6.2 Test `POST /api/v1/classes` dengan payload JSON yang valid dan pastikan mendapat response `201`.
  - [x] 6.3 Test `PUT /api/v1/classes/:class_id` dengan payload update dan pastikan data berubah.
