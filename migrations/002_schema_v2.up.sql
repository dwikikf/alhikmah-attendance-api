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
