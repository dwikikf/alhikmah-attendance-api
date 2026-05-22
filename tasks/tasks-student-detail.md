## Relevant Files

- `internal/domain/student.go` - Student struct dan interfaces. Method `GetByID` sudah ada di interface. Perlu dicek implementasinya.
- `internal/repository/student_postgres.go` - Student repository. Pastikan `GetByID` mengembalikan semua field termasuk `class_name`.
- `internal/service/student_service.go` - Student service. Method `GetByID` sudah ada.
- `internal/handler/student_handler.go` - Student handler. Tambahkan method `GetByID` handler.
- `cmd/api/main.go` - Register route `GET /students/:student_id`.

### Notes

- Berdasarkan PRD (prd.md bagian `GET /students/{student_id}`), response harus menyertakan `class_name` (diambil via JOIN dengan tabel `classes`).
- Response yang diharapkan oleh PRD:
  ```json
  {
    "success": true,
    "data": {
      "id": "UUID",
      "nisn": "1234567890",
      "full_name": "Ahmad Rizki Pratama",
      "class_id": "UUID",
      "class_name": "1A",
      "date_of_birth": "2018-03-15",
      "gender": "laki-laki",
      "photo_url": "https://...",
      "qr_code_data": "1234567890|Ahmad Rizki Pratama|1A",
      "is_active": true,
      "created_at": "2024-01-20T10:00:00Z"
    }
  }
  ```
- Perhatikan bahwa struct `Student` saat ini di `domain/student.go` **belum memiliki field `ClassName`**. Field ini perlu ditambahkan sebagai derived field (didapatkan dari JOIN, bukan disimpan di tabel students).
- Untuk rebuild dan test: `docker compose up -d --build`

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [ ]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

Example:
- `- [ ] 1.1 Read file` → `- [x] 1.1 Read file` (after completing)

Update the file after completing each sub-task, not just after completing an entire parent task.

## Tasks

- [x] 1.0 Update Domain Struct
  - [x] 1.1 Buka `internal/domain/student.go`. Tambahkan field `ClassName string \`json:"class_name,omitempty"\`` pada struct `Student` (sebagai derived field, taruh setelah `ClassID` atau di bagian "Derived fields").

- [x] 2.0 Update Repository - Perbaiki `GetByID` Query
  - [x] 2.1 Buka `internal/repository/student_postgres.go`. Cari method `GetByID`. Update query-nya agar melakukan JOIN dengan tabel `classes` dan mengambil `class_name`:
    ```sql
    SELECT s.id, s.nisn, s.full_name, s.class_id, c.class_name,
           s.date_of_birth, s.gender, s.photo_url, s.qr_code_data,
           s.is_active, s.created_at, s.updated_at
    FROM students s
    JOIN classes c ON s.class_id = c.id
    WHERE s.id = $1
    ```
  - [x] 2.2 Update `rows.Scan()` agar meng-scan `c.class_name` ke `student.ClassName`.
  - [x] 2.3 Pastikan handler nullable fields (`date_of_birth`, `gender`, `photo_url`) menggunakan `sql.NullString` agar tidak crash jika bernilai NULL.

- [x] 3.0 Tambahkan Handler untuk `GET /students/:student_id`
  - [x] 3.1 Buka `internal/handler/student_handler.go`. Tambahkan method `GetByID(c *gin.Context)`:
    - Ambil `student_id` dari URL param: `c.Param("student_id")`.
    - Panggil `service.GetByID(studentID)`.
    - Jika error (misalnya student not found), return `404 Not Found` dengan `{ "error": "Student not found" }`.
    - Jika sukses, return `200 OK` dengan `{ "success": true, "data": student }`.

- [x] 4.0 Register Route
  - [x] 4.1 Buka `cmd/api/main.go`. Tambahkan route di dalam blok `protected`:
    - `protected.GET("/students/:student_id", studentHandler.GetByID)`
    - **PENTING:** Pastikan route ini ditambahkan **setelah** route `protected.GET("/students", ...)` agar tidak terjadi konflik routing di Gin.

- [x] 5.0 Build dan Verifikasi
  - [x] 5.1 Jalankan `docker compose up -d --build` untuk memastikan kode berhasil di-compile.
  - [x] 5.2 Test `GET /api/v1/students/:student_id` dengan salah satu student ID dari database (bisa didapat dari `GET /api/v1/students?page=1&limit=1`) dan pastikan response menyertakan `class_name`.
