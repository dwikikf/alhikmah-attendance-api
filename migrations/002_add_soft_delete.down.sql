DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_students_deleted_at;
DROP INDEX IF EXISTS idx_classes_deleted_at;

ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE students DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE classes DROP COLUMN IF EXISTS deleted_at;
