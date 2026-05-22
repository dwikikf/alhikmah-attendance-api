ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'principal';
ALTER TYPE user_role RENAME VALUE 'guru' TO 'teacher';

ALTER TABLE users ADD COLUMN IF NOT EXISTS password_reset_token VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS reset_token_expires TIMESTAMP;
