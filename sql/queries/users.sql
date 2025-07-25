-- name: CreateUser :one
INSERT INTO users (id, username, api_key) 
VALUES ($1, $2, 
    encode(sha256(random()::text::bytea), 'hex')
)
RETURNING *;


-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE api_key = $1;

-- name: UpdateUser :one
UPDATE users
SET username = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;    

-- name: DeleteUser :one
DELETE FROM users WHERE id = $1
RETURNING *;    

