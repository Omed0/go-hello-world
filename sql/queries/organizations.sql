-- name: CreateOrganization :one
INSERT INTO organizations (id, name, description, settings) 
VALUES ($1, $2, $3, COALESCE($4, '{}'))
RETURNING *;

-- name: GetAllOrganizations :many
SELECT * FROM organizations WHERE deleted_at IS NULL;

-- name: GetOrganizationByID :one
SELECT * FROM organizations WHERE id = $1 AND deleted_at IS NULL;

-- name: GetOrganizationByName :one
SELECT * FROM organizations WHERE name = $1 AND deleted_at IS NULL;

-- name: UpdateOrganization :one
UPDATE organizations
SET name = CASE 
    WHEN $2 != '' THEN $2 
    ELSE name 
END,
description = COALESCE($3, description),
settings = COALESCE($4, settings),
updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteOrganization :one
UPDATE organizations
SET deleted_at = NOW()
WHERE id = $1
RETURNING *;

-- name: HardDeleteOrganization :one
DELETE FROM organizations WHERE id = $1
RETURNING *;
