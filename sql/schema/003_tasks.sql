-- +goose Up
CREATE TABLE
    tasks (
        id UUID PRIMARY KEY,
        title VARCHAR(200) NOT NULL UNIQUE,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP NULL,
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE SET NULL
    );

-- +goose Down
DROP TABLE tasks;