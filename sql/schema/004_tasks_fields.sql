-- +goose Up
ALTER TABLE tasks
ADD COLUMN description TEXT DEFAULT '',
ADD COLUMN is_completed BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE tasks SET description = '' WHERE description IS NULL;

ALTER TABLE tasks ALTER COLUMN description SET NOT NULL;


-- +goose Down
ALTER TABLE tasks
DROP COLUMN description,
DROP COLUMN is_completed;