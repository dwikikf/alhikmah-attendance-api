ALTER TABLE users DROP COLUMN IF EXISTS password_reset_token;
ALTER TABLE users DROP COLUMN IF EXISTS reset_token_expires;

ALTER TYPE user_role RENAME VALUE 'teacher' TO 'guru';
-- Note: 'principal' value cannot be dropped from the enum natively without recreating the type.
