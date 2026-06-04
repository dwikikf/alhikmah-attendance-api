## Relevant Files

### New Files
- `migrations/002_schema_v2.up.sql` - Migration baru: split `classes.class_name` → `room_name + grade + section`, tabel `student_enrollments` (dengan `end_reason`), alter `attendances` (kolom `subject`, partial unique index), alter `class_teachers` (kolom `subject` dan `role`).
- `migrations/002_schema_v2.down.sql` - Rollback migration.
- `core/domain/enrollment.go` - Struct `StudentEnrollment` dan interface `EnrollmentRepository`, `EnrollmentService`.
- `core/domain/class_teacher.go` - Struct `ClassTeacher` dan interface `ClassTeacherRepository`, `ClassTeacherService`.
- `core/repository/enrollment_postgres.go` - Implementasi repository `student_enrollments`.
- `core/repository/class_teacher_postgres.go` - Implementasi repository `class_teachers`.
- `core/service/enrollment_service.go` - Business logic enrollment, promote, dan transfer siswa.
- `core/service/class_teacher_service.go` - Business logic assign/unassign guru muatan lokal.
- `core/handler/enrollment_handler.go` - HTTP handler untuk enrollment, promote, dan transfer.
- `core/handler/class_teacher_handler.go` - HTTP handler untuk manage guru muatan lokal.

### Modified Files
- `migrations/001_init_schema.up.sql` - Tidak diubah (sudah live), semua perubahan via migration baru.
- `core/domain/class.go` - Ganti field `ClassName string` → `RoomName string`, `Grade int`, `Section *int`. Tambah helper method `DisplayName() string`.
- `core/domain/attendance.go` - Tambah field `Subject *string` di struct `Attendance`.
- `core/domain/student.go` - Field `ClassID` tetap ada untuk backward compat (akan deprecated setelah enrollment aktif).
- `core/repository/class_postgres.go` - Update semua query: ganti `class_name` → `room_name, grade, section`. Update `IsTeacherResponsibleForStudent` agar cek `class_teachers` juga. Update `Create` dan `Update`.
- `core/repository/attendance_postgres.go` - Update `MarkAttendance`, `GetByClassAndDate` untuk kolom `subject`.
- `core/service/class_service.go` - Update validasi Create/Update untuk field `room_name`, `grade`, `section`.
- `core/service/attendance_service.go` - Update `ProcessQRScan` dan `ProcessManualAttendance` terima param `subject`, update cache key.
- `core/service/report_service.go` - Update query laporan filter by `subject` dan `academic_year`, tambah validasi akses historis.
- `core/handler/attendance_handler.go` - Tambah param `subject` di endpoint scan dan manual.
- `core/handler/report_handler.go` - Tambah query param `subject` dan `academic_year`.
- `cmd/api/main.go` - Register route baru untuk enrollment, transfer, dan class_teachers.
- `cmd/seed/main.go` - Update seed: pakai `room_name/grade/section`, tambah seed `student_enrollments`, contoh `class_teachers` muatan lokal.

### Test Files
- `core/service/enrollment_service_test.go` - Unit test enrollment, promote, transfer.
- `core/service/class_teacher_service_test.go` - Unit test class teacher service.
- `core/service/attendance_service_test.go` - Update test untuk subject-aware attendance.
- `core/service/class_service_test.go` - Update test untuk schema baru.

### Notes
- **URUTAN PENTING:** Task 1 (migration) → Task 2 (update domain Class) → Task 3 (enrollment) → Task 4 (class teacher) → Task 5 (attendance) → Task 6 (reports). Jangan lewati urutan ini karena ada ketergantungan.
- **Display name kelas** digenerate di Go: `"Kelas {Grade} {RoomName}"` atau `"Kelas {Grade} {RoomName} {Section}"` jika `Section != nil`. Jangan simpan display name di DB.
- **`section` bersifat opsional (NULL):** jika hanya ada 1 kelas di grade+room tersebut, section = NULL. Jika ada Madinah 1 & Madinah 2, baru section = 1 dan 2.
- **Promote ≠ bulk flat:** Endpoint promote menerima array mapping per-siswa `[{student_id, target_class_id}]`. Sistem default-suggest `target_class_id` = kelas dengan `room_name` sama, `grade+1`, `section` sama. Admin bisa ubah individual.
- **Transfer (pindah seksi):** Endpoint terpisah, bisa terjadi kapan saja dalam TA yang sama maupun beda TA. `end_reason = 'transferred'`.
- **`end_reason` values:** `'promoted'` (naik kelas), `'transferred'` (pindah seksi/kelas), `'graduated'` (lulus), `'dropped'` (keluar sekolah).
- Cache key QR scan: `"qr_scan_" + nisn + "_" + subject` (jika subject kosong pakai `"reguler"`).
- Endpoint promote dan class_teachers: admin only.
- Partial unique index: `uq_attendance_regular` untuk `subject IS NULL`, `uq_attendance_subject` untuk `subject IS NOT NULL`.
- Gunakan `docker compose up -d --build` untuk rebuild lokal.

---

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [ ]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

Update the file after completing each sub-task, not just after completing an entire parent task.

---

## Tasks

- [x] 0.0 Create feature branch
  - [x] 0.1 Buat dan checkout branch baru: `git checkout -b feature/schema-v2-enrollment`

- [x] 1.0 Database Migration — Schema Baru
  - [x] 1.1 Buat file `migrations/002_schema_v2.up.sql`:
    ```sql
    -- 1. Split class_name → room_name + grade + section di tabel classes
    ALTER TABLE classes ADD COLUMN room_name VARCHAR(50);
    ALTER TABLE classes ADD COLUMN grade SMALLINT;
    ALTER TABLE classes ADD COLUMN section SMALLINT DEFAULT NULL;

    -- Isi room_name dari class_name yang ada (data migration — sesuaikan manual jika perlu)
    -- Contoh sederhana: salin class_name ke room_name dulu, admin akan edit via UI
    UPDATE classes SET room_name = class_name, grade = 1 WHERE room_name IS NULL;

    -- Terapkan NOT NULL setelah data diisi
    ALTER TABLE classes ALTER COLUMN room_name SET NOT NULL;
    ALTER TABLE classes ALTER COLUMN grade SET NOT NULL;

    -- Hapus constraint lama, ganti dengan constraint baru
    ALTER TABLE classes DROP CONSTRAINT IF EXISTS classes_class_name_academic_year_key;
    ALTER TABLE classes ADD CONSTRAINT classes_room_grade_section_year_key
      UNIQUE(room_name, grade, section, academic_year);

    -- Opsional: hapus class_name setelah data migration selesai (bisa di migration berikutnya)
    -- ALTER TABLE classes DROP COLUMN class_name;

    -- 2. Tabel student_enrollments
    CREATE TABLE student_enrollments (
      id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      student_id   UUID NOT NULL REFERENCES students(id),
      class_id     UUID NOT NULL REFERENCES classes(id),
      academic_year VARCHAR(9) NOT NULL,
      enrolled_at  TIMESTAMP DEFAULT NOW(),
      ended_at     TIMESTAMP NULL,
      end_reason   VARCHAR(20) NULL,
      -- end_reason: 'promoted' | 'transferred' | 'graduated' | 'dropped'
      UNIQUE(student_id, class_id, academic_year)
    );
    CREATE INDEX idx_enrollments_student_id ON student_enrollments(student_id);
    CREATE INDEX idx_enrollments_class_id ON student_enrollments(class_id);
    CREATE INDEX idx_enrollments_academic_year ON student_enrollments(academic_year);
    CREATE INDEX idx_enrollments_active ON student_enrollments(student_id) WHERE ended_at IS NULL;

    -- Seed enrollment aktif dari data students.class_id yang sudah ada
    INSERT INTO student_enrollments (student_id, class_id, academic_year)
    SELECT s.id, s.class_id, c.academic_year
    FROM students s
    JOIN classes c ON s.class_id = c.id
    WHERE s.deleted_at IS NULL
    ON CONFLICT DO NOTHING;

    -- 3. Tambah kolom subject ke attendances
    ALTER TABLE attendances ADD COLUMN subject VARCHAR(100) DEFAULT NULL;

    -- 4. Hapus unique constraint lama attendances, ganti dengan partial index
    ALTER TABLE attendances
      DROP CONSTRAINT attendances_student_id_class_id_attendance_date_key;
    CREATE UNIQUE INDEX uq_attendance_regular
      ON attendances(student_id, class_id, attendance_date)
      WHERE subject IS NULL;
    CREATE UNIQUE INDEX uq_attendance_subject
      ON attendances(student_id, class_id, attendance_date, subject)
      WHERE subject IS NOT NULL;

    -- 5. Tambah kolom ke class_teachers
    ALTER TABLE class_teachers ADD COLUMN subject VARCHAR(100) DEFAULT NULL;
    ALTER TABLE class_teachers ADD COLUMN role VARCHAR(20) DEFAULT 'subject_teacher' NOT NULL;
    -- role: 'homeroom' untuk wali kelas, 'subject_teacher' untuk muatan lokal
    ```
  - [x] 1.2 Buat file `migrations/002_schema_v2.down.sql`:
    ```sql
    -- Rollback class_teachers
    ALTER TABLE class_teachers DROP COLUMN IF EXISTS subject;
    ALTER TABLE class_teachers DROP COLUMN IF EXISTS role;

    -- Rollback attendances
    DROP INDEX IF EXISTS uq_attendance_subject;
    DROP INDEX IF EXISTS uq_attendance_regular;
    ALTER TABLE attendances ADD CONSTRAINT attendances_student_id_class_id_attendance_date_key
      UNIQUE(student_id, class_id, attendance_date);
    ALTER TABLE attendances DROP COLUMN IF EXISTS subject;

    -- Rollback student_enrollments
    DROP TABLE IF EXISTS student_enrollments CASCADE;

    -- Rollback classes
    ALTER TABLE classes DROP CONSTRAINT IF EXISTS classes_room_grade_section_year_key;
    ALTER TABLE classes ADD CONSTRAINT classes_class_name_academic_year_key
      UNIQUE(class_name, academic_year);
    ALTER TABLE classes DROP COLUMN IF EXISTS section;
    ALTER TABLE classes DROP COLUMN IF EXISTS grade;
    ALTER TABLE classes DROP COLUMN IF EXISTS room_name;
    ```
  - [x] 1.3 Jalankan migration di lokal: `migrate -path migrations -database "$DATABASE_URL" up`
  - [x] 1.4 Verifikasi di psql:
    - `\d classes` → pastikan ada kolom `room_name`, `grade`, `section`
    - `\d student_enrollments` → pastikan tabel ada dengan kolom `ended_at` dan `end_reason`
    - `\d attendances` → pastikan ada kolom `subject`
    - `\d class_teachers` → pastikan ada kolom `subject` dan `role`

- [x] 2.0 Update Domain & Repository Class — Schema Baru
  - [x] 2.1 Update `core/domain/class.go`:
    - Ganti field `ClassName string` dengan `RoomName string`, `Grade int`, `Section *int`
    - Tambah helper method:
      ```go
      func (c *Class) DisplayName() string {
          if c.Section != nil {
              return fmt.Sprintf("Kelas %d %s %d", c.Grade, c.RoomName, *c.Section)
          }
          return fmt.Sprintf("Kelas %d %s", c.Grade, c.RoomName)
      }
      ```
    - Tambah field `DisplayName string` (computed, tidak disimpan di DB) untuk JSON response
  - [x] 2.2 Update `core/repository/class_postgres.go`:
    - `GetAll`: ganti `c.class_name` → `c.room_name, c.grade, c.section` di SELECT dan WHERE
    - `GetByID`: idem
    - `Create`: ganti `class_name` → `room_name, grade, section` di INSERT
    - `Update`: ganti `class_name` → `room_name, grade, section` di UPDATE SET
    - Tambah ORDER BY default: `ORDER BY c.grade, c.room_name, c.section`
    - Update `IsTeacherResponsibleForStudent` — tambah cek via `class_teachers`:
      ```sql
      SELECT EXISTS (
        SELECT 1 FROM student_enrollments se
        JOIN classes c ON se.class_id = c.id
        WHERE se.student_id = $1 AND c.teacher_id = $2
          AND se.ended_at IS NULL AND c.deleted_at IS NULL
      ) OR EXISTS (
        SELECT 1 FROM student_enrollments se
        JOIN class_teachers ct ON se.class_id = ct.class_id
        WHERE se.student_id = $1 AND ct.teacher_id = $2
          AND se.ended_at IS NULL
      )
      ```
      *Catatan: setelah `student_enrollments` aktif, gunakan join via enrollment, bukan `students.class_id`*
  - [x] 2.3 Update `core/service/class_service.go`:
    - Validasi `Create`: `RoomName` tidak boleh kosong, `Grade` harus 1–12, `Section` opsional tapi jika ada harus > 0
    - Validasi `Update`: idem
  - [x] 2.4 Update `core/handler/class_handler.go`:
    - Request body Create/Update: ganti `class_name` → `room_name` (string), `grade` (int), `section` (int, opsional)
    - Response: sertakan `display_name` (panggil `class.DisplayName()`) di setiap response kelas
  - [x] 2.5 Update `core/service/class_service_test.go` — sesuaikan semua test dengan field baru

- [x] 3.0 Student Enrollment — Domain, Repository, Service, Handler
  - [x] 3.1 Buat `core/domain/enrollment.go`:
    ```go
    type StudentEnrollment struct {
        ID           string
        StudentID    string
        StudentName  string  // join dari students
        ClassID      string
        ClassDisplay string  // display name kelas
        AcademicYear string
        EnrolledAt   time.Time
        EndedAt      *time.Time
        EndReason    *string  // 'promoted' | 'transferred' | 'graduated' | 'dropped'
    }

    type PromoteItem struct {
        StudentID     string  // siswa yang akan dipromote
        TargetClassID string  // kelas tujuan (bisa beda section)
    }

    type EnrollmentRepository interface {
        Enroll(e *StudentEnrollment) error
        GetActiveByStudentID(studentID string) (*StudentEnrollment, error)
        GetActiveByClassID(classID string) ([]*StudentEnrollment, error)
        GetHistoryByStudentID(studentID string) ([]*StudentEnrollment, error)
        EndEnrollment(studentID, classID string, reason string) error
        BulkEnroll(items []PromoteItem, academicYear string) (int, error)
    }

    type EnrollmentService interface {
        Enroll(studentID, classID, academicYear string) error
        PromoteClass(items []PromoteItem, academicYear string) (int, error)
        TransferStudent(studentID, fromClassID, toClassID, academicYear string) error
        GetActiveByStudentID(studentID string) (*StudentEnrollment, error)
        GetActiveByClassID(classID string) ([]*StudentEnrollment, error)
        GetHistoryByStudentID(studentID string) ([]*StudentEnrollment, error)
    }
    ```
  - [x] 3.2 Buat `core/repository/enrollment_postgres.go`:
    - `Enroll`: INSERT ke `student_enrollments`, return error jika sudah ada enrollment aktif di TA yang sama
    - `GetActiveByStudentID`: SELECT WHERE `ended_at IS NULL` JOIN ke `classes` untuk display name
    - `GetActiveByClassID`: SELECT semua siswa aktif di kelas (JOIN ke `students`)
    - `GetHistoryByStudentID`: SELECT semua baris, ORDER BY `enrolled_at DESC`
    - `EndEnrollment(studentID, classID, reason string)`:
      ```sql
      UPDATE student_enrollments
      SET ended_at = NOW(), end_reason = $3
      WHERE student_id = $1 AND class_id = $2 AND ended_at IS NULL
      ```
    - `BulkEnroll(items []PromoteItem, academicYear string) (int, error)`:
      — gunakan `sql.Tx`, loop INSERT per item, return jumlah berhasil
  - [x] 3.3 Buat `core/service/enrollment_service.go`:
    - `Enroll`: cek siswa dan kelas ada, cek tidak double aktif di TA yang sama, panggil `repo.Enroll`
    - `PromoteClass(items []PromoteItem, academicYear string)`:
      1. Validasi `academicYear` format `YYYY/YYYY`
      2. Validasi semua `TargetClassID` ada di DB
      3. Jalankan dalam 1 transaksi:
         - Loop setiap item: panggil `EndEnrollment(studentID, currentClassID, 'promoted')`
         - Panggil `BulkEnroll(items, academicYear)`
      4. Return jumlah siswa yang berhasil dipromote
    - `TransferStudent(studentID, fromClassID, toClassID, academicYear string)`:
      1. Cek enrollment aktif siswa di `fromClassID`
      2. `EndEnrollment(studentID, fromClassID, 'transferred')`
      3. `Enroll(studentID, toClassID, academicYear)`
  - [x] 3.4 Buat `core/service/enrollment_service_test.go`:
    - Test `Enroll` berhasil
    - Test `Enroll` gagal jika sudah ada enrollment aktif di TA yang sama
    - Test `PromoteClass` dengan 3 siswa, 2 ke kelas sama, 1 ke kelas berbeda (beda section)
    - Test `TransferStudent` berhasil — enrollment lama tutup, enrollment baru buka
  - [x] 3.5 Buat `core/handler/enrollment_handler.go`:
    - `POST /classes/:class_id/enroll` — enroll 1 siswa:
      ```json
      { "student_id": "UUID", "academic_year": "2026/2027" }
      ```
    - `POST /classes/:class_id/promote` — promote per-siswa (bukan bulk flat):
      ```json
      {
        "academic_year": "2027/2028",
        "items": [
          { "student_id": "UUID-A", "target_class_id": "UUID-KELAS-2-MADINAH-1" },
          { "student_id": "UUID-B", "target_class_id": "UUID-KELAS-2-MADINAH-2" }
        ]
      }
      ```
      Return: `{ "success": true, "data": { "promoted_count": 30 } }`
    - `POST /students/:student_id/transfer` — pindah seksi/kelas:
      ```json
      { "from_class_id": "UUID", "to_class_id": "UUID", "academic_year": "2026/2027" }
      ```
    - `GET /students/:student_id/enrollments` — history enrollment siswa
    - `GET /classes/:class_id/enrollments` — list siswa aktif di kelas (via enrollment)
  - [x] 3.6 Register route di `cmd/api/main.go`:
    - `POST /classes/:class_id/enroll` — admin only
    - `POST /classes/:class_id/promote` — admin only
    - `POST /students/:student_id/transfer` — admin only
    - `GET /students/:student_id/enrollments` — admin + teacher (kelas sendiri)
    - `GET /classes/:class_id/enrollments` — admin + teacher (kelas sendiri)

- [x] 4.0 Class Teachers (Muatan Lokal) — Domain, Repository, Service, Handler
  - [x] 4.1 Buat `core/domain/class_teacher.go`:
    ```go
    type ClassTeacher struct {
        ID           string
        TeacherID    string
        TeacherName  string
        ClassID      string
        ClassDisplay string  // display name kelas
        AcademicYear string
        Subject      string  // nama mapel, e.g. "PJOK", "Bahasa Arab"
        Role         string  // 'homeroom' | 'subject_teacher'
        CreatedAt    time.Time
    }

    type ClassTeacherRepository interface {
        Assign(ct *ClassTeacher) error
        Unassign(teacherID, classID, subject string) error
        GetByClassID(classID string) ([]*ClassTeacher, error)
        GetByTeacherID(teacherID string) ([]*ClassTeacher, error)
        GetSubjectAssignments(teacherID string) ([]*ClassTeacher, error)
    }

    type ClassTeacherService interface {
        Assign(ct *ClassTeacher) error
        Unassign(teacherID, classID, subject string) error
        GetByClassID(classID string) ([]*ClassTeacher, error)
        GetSubjectAssignments(teacherID string) ([]*ClassTeacher, error)
    }
    ```
  - [x] 4.2 Buat `core/repository/class_teacher_postgres.go`:
    - `Assign`: INSERT ke `class_teachers` dengan `role='subject_teacher'`
    - `Unassign`: DELETE WHERE `teacher_id=$1 AND class_id=$2 AND subject=$3`
    - `GetByClassID`: SELECT JOIN `users` (untuk `teacher_name`) dan `classes` — include homeroom dari `classes.teacher_id` dengan `role='homeroom'` dan semua baris di `class_teachers`
    - `GetByTeacherID`: SELECT semua assignment guru ini
    - `GetSubjectAssignments`: SELECT untuk mobile — hanya yang `role='subject_teacher'`, include display name kelas
  - [x] 4.3 Buat `core/service/class_teacher_service.go`:
    - `Assign`: validasi user ada dan `role='teacher'`, validasi kelas ada, `subject` tidak boleh kosong, panggil `repo.Assign`
    - `Unassign`: validasi assignment ada, panggil `repo.Unassign`
    - `GetSubjectAssignments`: untuk mobile app saat load ScannerScreen
  - [x] 4.4 Buat `core/service/class_teacher_service_test.go`:
    - Test `Assign` berhasil
    - Test `Assign` gagal jika user bukan role teacher
    - Test `Assign` gagal jika subject kosong
    - Test `Unassign` berhasil
  - [x] 4.5 Buat `core/handler/class_teacher_handler.go`:
    - `POST /classes/:class_id/teachers` — assign guru muatan lokal (admin only):
      ```json
      { "teacher_id": "UUID", "subject": "PJOK", "academic_year": "2026/2027" }
      ```
    - `DELETE /classes/:class_id/teachers/:teacher_id?subject=PJOK` — unassign (admin only)
    - `GET /classes/:class_id/teachers` — list semua guru di kelas (admin + teacher)
    - `GET /teachers/:teacher_id/subjects` — list assignment mapel guru (untuk mobile, admin + teacher sendiri)
  - [x] 4.6 Register semua route di `cmd/api/main.go`

- [ ] 5.0 Update Attendance — Subject-Aware
  - [ ] 5.1 Update `core/domain/attendance.go` — tambah field `Subject *string`
  - [ ] 5.2 Update `core/repository/attendance_postgres.go`:
    - `MarkAttendance`: tambah kolom `subject` di INSERT
    - `GetByClassAndDate`: tambah kolom `subject` di SELECT + Scan, tambah opsional filter `WHERE subject IS NULL` atau `WHERE subject = $n`
  - [ ] 5.3 Update `core/service/attendance_service.go`:
    - `ProcessQRScan(nisn, teacherID, role, subject string) error`:
      - Cache key: `"qr_scan_" + nisn + "_" + coalesce(subject, "reguler")`
      - Cek duplikat harus include filter `subject`
      - Set `attendance.Subject` sebelum `MarkAttendance`
    - `ProcessManualAttendance`: tambah param `subject string`
  - [ ] 5.4 Update `core/handler/attendance_handler.go`:
    - `POST /attendances/scan`: tambah field `subject` (opsional):
      ```json
      { "nisn": "1234567890", "subject": "PJOK" }
      ```
    - `POST /attendances/manual`: tambah field `subject`
  - [ ] 5.5 Update `core/service/attendance_service_test.go` — update semua test yang terpengaruh

- [ ] 6.0 Update Reports — Filter Subject & Akses Historis
  - [ ] 6.1 Update `core/domain/report.go` — tambah `Subject *string` di struct parameter/filter
  - [ ] 6.2 Update `core/service/report_service.go`:
    - Semua query `attendances` tambah kondisi filter `subject`:
      - Jika `subject` kosong → `WHERE subject IS NULL` (default: reguler)
      - Jika `subject` diisi → `WHERE subject = $n`
    - Validasi akses historis: guru (muatan lokal) boleh akses laporan jika ada di `class_teachers` untuk `class_id` + `academic_year` tersebut
  - [ ] 6.3 Update `core/handler/report_handler.go`:
    - Tambah query param `subject` dan `academic_year`:
      ```
      GET /reports?class_id=xxx&academic_year=2026/2027&subject=PJOK
      ```
  - [ ] 6.4 Update `core/service/report_service_test.go`

- [x] 7.0 Update Seed Data
  - [x] 7.1 Update `cmd/seed/main.go`:
    - Ganti `class_name` → insert dengan `room_name`, `grade`, `section` yang realistis (contoh: `room_name='Aqoba', grade=1`, `room_name='Madinah', grade=1, section=1`, dst.)
    - Tambah seed `student_enrollments` untuk semua siswa seed
    - Tambah contoh `class_teachers` dengan `subject='PJOK'` dan `role='subject_teacher'`

- [x] 8.0 Build dan Verifikasi
  - [x] 8.1 Jalankan `docker compose up -d --build` — pastikan tidak ada error compile
  - [x] 8.2 Jalankan semua unit test: `go test ./...`
  - [x] 8.3 Test `GET /classes` — pastikan response sertakan `display_name` seperti `"Kelas 1 Aqoba"`
  - [x] 8.4 Test `POST /classes` — buat kelas dengan `room_name="Madinah", grade=1, section=1`
  - [x] 8.5 Test `POST /classes/:id/promote` — kirim array mapping per-siswa, pastikan enrollment lama tutup (`ended_at` terisi, `end_reason='promoted'`), enrollment baru terbuka
  - [x] 8.6 Test `POST /students/:id/transfer` — pindah siswa dari Madinah 1 ke Madinah 2, pastikan history enrollment benar
  - [x] 8.7 Test QR scan reguler (subject kosong) + QR scan PJOK di hari yang sama → keduanya berhasil
  - [x] 8.8 Test scan PJOK kedua kali di hari yang sama → error "sudah absen"
  - [x] 8.9 Test `GET /reports?subject=PJOK` — hanya tampil data absensi PJOK
