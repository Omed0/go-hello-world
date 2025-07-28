-- name: CreateTask :one
INSERT INTO tasks (id, title, description, user_id) 
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAllTasks :many
SELECT * FROM tasks WHERE deleted_at IS NULL ORDER BY created_at DESC; 

-- name: GetTaskById :one
SELECT * FROM tasks WHERE id = $1 AND deleted_at IS NULL;

-- name: GetTasksByUserId :many
SELECT * FROM tasks WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC;

-- name: GetDeletedTasksByUserId :many
SELECT * FROM tasks WHERE user_id = $1 AND deleted_at IS NOT NULL ORDER BY updated_at DESC;

-- name: CompleteTask :one
UPDATE tasks
SET is_completed = TRUE, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UndoCompleteTask :one
UPDATE tasks
SET is_completed = FALSE, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateTaskPartial :one
UPDATE tasks
SET
  title = COALESCE(NULLIF($1, ''), title),
  description = COALESCE(NULLIF($2, ''), description),
  updated_at = NOW()
WHERE id = $3 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteTask :one
UPDATE tasks
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: RestoreTask :one
UPDATE tasks
SET deleted_at = NULL, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NOT NULL
RETURNING *;

-- name: HardDeleteTask :one
DELETE FROM tasks
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: HandleSearchTasks :many
SELECT * FROM tasks
WHERE (title ILIKE '%' || $1 || '%')
AND ($2::int IS NULL OR user_id = $2)
AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT COALESCE(NULLIF($3, 0), 10) 
OFFSET COALESCE(NULLIF(($4 - 1) * COALESCE(NULLIF($3, 0), 10), -10), 0); 

