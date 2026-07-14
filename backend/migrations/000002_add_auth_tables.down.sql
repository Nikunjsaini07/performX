DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS user_otps;
ALTER TABLE users DROP COLUMN IF EXISTS password_hash;
ALTER TABLE users DROP COLUMN IF EXISTS email_verified;
