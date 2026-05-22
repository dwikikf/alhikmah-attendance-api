ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL;
ALTER TABLE students ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL;
ALTER TABLE classes ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL;

CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_students_deleted_at ON students(deleted_at);
CREATE INDEX idx_classes_deleted_at ON classes(deleted_at);
