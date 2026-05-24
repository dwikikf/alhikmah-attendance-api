BEGIN;

UPDATE users SET role = 'teacher' WHERE role::text = 'principal';

ALTER TYPE user_role RENAME TO user_role_old;
CREATE TYPE user_role AS ENUM ('admin', 'teacher');
ALTER TABLE users ALTER COLUMN role DROP DEFAULT;
ALTER TABLE users ALTER COLUMN role TYPE user_role USING role::text::user_role;
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'teacher';
DROP TYPE user_role_old;

COMMIT;
