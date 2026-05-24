CREATE TYPE user_role AS ENUM ('admin', 'teacher');
CREATE TYPE gender_type AS ENUM ('laki-laki', 'perempuan');
CREATE TYPE attendance_status AS ENUM ('hadir', 'izin', 'sakit', 'tanpa_keterangan');
CREATE TYPE report_type_enum AS ENUM ('harian', 'mingguan', 'bulanan', 'semesteran');

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username VARCHAR(50) UNIQUE NOT NULL,
  email VARCHAR(100) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  full_name VARCHAR(150) NOT NULL,
  role user_role DEFAULT 'teacher',
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL DEFAULT NULL,
  last_login TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

CREATE TABLE classes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  class_name VARCHAR(50) NOT NULL,
  teacher_id UUID NOT NULL REFERENCES users(id),
  academic_year VARCHAR(9) NOT NULL,
  capacity INT DEFAULT 30,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL DEFAULT NULL,
  UNIQUE(class_name, academic_year)
);

CREATE INDEX idx_classes_teacher_id ON classes(teacher_id);
CREATE INDEX idx_classes_academic_year ON classes(academic_year);
CREATE INDEX idx_classes_deleted_at ON classes(deleted_at);

CREATE TABLE students (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  nisn VARCHAR(10) UNIQUE NOT NULL,
  full_name VARCHAR(150) NOT NULL,
  class_id UUID NOT NULL REFERENCES classes(id),
  date_of_birth DATE,
  gender gender_type,
  photo_url VARCHAR(500),
  qr_code_data TEXT NOT NULL,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL DEFAULT NULL,
  UNIQUE(nisn, class_id)
);

CREATE INDEX idx_students_nisn ON students(nisn);
CREATE INDEX idx_students_class_id ON students(class_id);
CREATE INDEX idx_students_is_active ON students(is_active);
CREATE INDEX idx_students_deleted_at ON students(deleted_at);

CREATE TABLE attendances (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  student_id UUID NOT NULL REFERENCES students(id),
  class_id UUID NOT NULL REFERENCES classes(id),
  attendance_date DATE NOT NULL,
  status attendance_status DEFAULT 'hadir',
  recorded_by UUID NOT NULL REFERENCES users(id),
  recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  scanned_at TIMESTAMP,
  notes TEXT,
  is_manual BOOLEAN DEFAULT false,
  UNIQUE(student_id, class_id, attendance_date)
);

CREATE INDEX idx_attendances_student_id ON attendances(student_id);
CREATE INDEX idx_attendances_class_id ON attendances(class_id);
CREATE INDEX idx_attendances_attendance_date ON attendances(attendance_date);
CREATE INDEX idx_attendances_status ON attendances(status);
CREATE INDEX idx_attendances_recorded_by ON attendances(recorded_by);
CREATE INDEX idx_attendances_composite ON attendances(class_id, attendance_date);

CREATE TABLE attendance_audits (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  attendance_id UUID NOT NULL REFERENCES attendances(id),
  old_status VARCHAR(50),
  new_status VARCHAR(50) NOT NULL,
  changed_by UUID NOT NULL REFERENCES users(id),
  changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  reason TEXT
);

CREATE INDEX idx_attendance_audits_attendance_id ON attendance_audits(attendance_id);
CREATE INDEX idx_attendance_audits_changed_by ON attendance_audits(changed_by);
CREATE INDEX idx_attendance_audits_changed_at ON attendance_audits(changed_at);

CREATE TABLE class_teachers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  teacher_id UUID NOT NULL REFERENCES users(id),
  class_id UUID NOT NULL REFERENCES classes(id),
  academic_year VARCHAR(9) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(teacher_id, class_id, academic_year)
);

CREATE INDEX idx_class_teachers_teacher_id ON class_teachers(teacher_id);
CREATE INDEX idx_class_teachers_class_id ON class_teachers(class_id);

CREATE TABLE reports (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  report_type report_type_enum NOT NULL,
  class_id UUID NOT NULL REFERENCES classes(id),
  period_start DATE NOT NULL,
  period_end DATE NOT NULL,
  generated_by UUID NOT NULL REFERENCES users(id),
  generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  report_data JSONB,
  UNIQUE(report_type, class_id, period_start, period_end)
);

CREATE INDEX idx_reports_class_id ON reports(class_id);
CREATE INDEX idx_reports_report_type ON reports(report_type);
CREATE INDEX idx_reports_period ON reports(period_start, period_end);
