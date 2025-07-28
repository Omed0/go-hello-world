-- name: CreateUser :one
INSERT INTO users (id, username, password_hash, age, gender, role, organization_id, api_key) 
VALUES ($1, $2, $3, $4, $5, COALESCE($6, 'user'), COALESCE($7, '00000000-0000-0000-0000-000000000000'),
    encode(sha256(random()::text::bytea), 'hex')
)
RETURNING *;

-- name: CreateUserWithPassword :one
INSERT INTO users (id, username, password_hash, age, gender, role, organization_id, api_key) 
VALUES ($1, $2, $3, $4, $5, COALESCE($6, 'user'), COALESCE($7, '00000000-0000-0000-0000-000000000000'),
    encode(sha256(random()::text::bytea), 'hex')
)
RETURNING *;

-- name: GetAllUsers :many
SELECT u.*, o.name as organization_name 
FROM users u 
LEFT JOIN organizations o ON u.organization_id = o.id;

-- name: GetUserByUsername :one
SELECT u.*, o.name as organization_name 
FROM users u 
LEFT JOIN organizations o ON u.organization_id = o.id 
WHERE u.username = $1;

-- name: GetUserByUsernameAndPassword :one
SELECT u.*, o.name as organization_name 
FROM users u 
LEFT JOIN organizations o ON u.organization_id = o.id 
WHERE u.username = $1 AND u.password_hash = $2;

-- name: GetUserByID :one
SELECT u.*, o.name as organization_name 
FROM users u 
LEFT JOIN organizations o ON u.organization_id = o.id 
WHERE u.id = $1;

-- name: GetUserByAPIKey :one
SELECT u.*, o.name as organization_name 
FROM users u 
LEFT JOIN organizations o ON u.organization_id = o.id 
WHERE u.api_key = $1;

-- name: GetUsersByOrganization :many
SELECT u.*, o.name as organization_name 
FROM users u 
LEFT JOIN organizations o ON u.organization_id = o.id 
WHERE u.organization_id = $1;

-- name: UpdateUser :one
UPDATE users
SET username = COALESCE($2, username), 
    age = COALESCE($3, age),
    gender = COALESCE($4, gender),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET password_hash = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserRole :one
UPDATE users
SET role = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserOrganization :one
UPDATE users
SET organization_id = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users WHERE id = $1
RETURNING *;    

