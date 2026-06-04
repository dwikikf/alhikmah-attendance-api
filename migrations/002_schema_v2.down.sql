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
