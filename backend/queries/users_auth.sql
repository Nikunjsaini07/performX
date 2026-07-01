
-- name: VerifyUserEmail :exec
UPDATE users
SET email_verified = TRUE, updated_at = now()
WHERE email = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2, updated_at = now()
WHERE email = $1;

-- name: UpdateUserProfile :one
UPDATE users
SET display_name = COALESCE(sqlc.narg('display_name'), display_name),
    bio = COALESCE(sqlc.narg('bio'), bio),
    avatar_url = COALESCE(sqlc.narg('avatar_url'), avatar_url),
    updated_at = now()
WHERE id = $1
RETURNING id, username, display_name, email, bio, avatar_url, email_verified, created_at, updated_at;

-- name: CreateOTP :exec
INSERT INTO user_otps (email, otp_code, purpose, expires_at)
VALUES ($1, $2, $3, $4);

-- name: GetValidOTP :one
SELECT id, email, otp_code, purpose, expires_at, created_at
FROM user_otps
WHERE email = $1 AND otp_code = $2 AND purpose = $3 AND expires_at > now()
LIMIT 1;

-- name: DeleteOTP :exec
DELETE FROM user_otps
WHERE id = $1;

-- name: CleanExpiredOTPs :exec
DELETE FROM user_otps
WHERE expires_at < now();

-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (user_id, token, expires_at)
VALUES ($1, $2, $3);

-- name: GetRefreshToken :one
SELECT id, user_id, token, expires_at, created_at
FROM refresh_tokens
WHERE token = $1 AND expires_at > now()
LIMIT 1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token = $1;

-- name: DeleteUserRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;

-- name: CleanExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < now();
