-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email,hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: CreateChirp :one
INSERT INTO chirps(id,created_at,updated_at,body,user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    now(),
    $1,
    $2
)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps ORDER BY created_at ASC;

-- name: GetChirp :one
SELECT * FROM chirps WHERE id = $1;

-- name: GetChirpsByUserId :many
SELECT * FROM chirps WHERE user_id = $1 ORDER BY created_at ASC;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token,user_id,expires_at,revoked_at)
VALUES(
    $1,
    $2,
    $3,
    NULL
)
RETURNING *;


-- name: GetUserFromRefreshToken :one
SELECT user_id from refresh_tokens where token = $1 AND revoked_at IS NULL;


-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = NOW(),
updated_at = NOW()
WHERE token = $1;


-- name: UpdateUserDetails :one
UPDATE users SET email = $2, hashed_password = $3
WHERE id = $1
RETURNING *;

-- name: DeleteChirp :exec
DELETE FROM chirps WHERE id = $1;

-- name: UpgradeUser :exec
UPDATE users SET is_chirpy_red = true
WHERE id = $1;