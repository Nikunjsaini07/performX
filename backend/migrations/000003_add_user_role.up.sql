-- Add role column to users table using pre-existing user_role ENUM
ALTER TABLE users ADD COLUMN IF NOT EXISTS role user_role NOT NULL DEFAULT 'USER';

CREATE INDEX IF NOT EXISTS users_role_idx ON users(role);
