-- name: CreateUser :one
INSERT INTO users (id, username, created_at, updated_at, api_key) 
VALUES ($1, $2, $3, $4,
    encode(sha256(random()::text::bytea), 'hex')
)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE api_key = $1;
