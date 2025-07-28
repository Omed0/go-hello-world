package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/omed0/go-hello-world/internal/database"
)

// Task represents a task in the system
type Task struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsCompleted bool       `json:"is_completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	UserID      uuid.UUID  `json:"user_id"`
}

// DatabaseTaskToTask converts a database task to a task model
func DatabaseTaskToTask(dbTask database.Task) Task {
	task := Task{
		ID:          dbTask.ID,
		Title:       dbTask.Title,
		Description: dbTask.Description,
		IsCompleted: dbTask.IsCompleted,
		CreatedAt:   dbTask.CreatedAt,
		UpdatedAt:   dbTask.UpdatedAt,
		UserID:      dbTask.UserID,
	}

	// Handle nullable deleted_at field
	if dbTask.DeletedAt.Valid {
		task.DeletedAt = &dbTask.DeletedAt.Time
	}

	return task
}

// DatabaseTasksToTasks converts a slice of database tasks to task models
func DatabaseTasksToTasks(dbTasks []database.Task) []Task {
	tasks := make([]Task, len(dbTasks))
	for i, dbTask := range dbTasks {
		tasks[i] = DatabaseTaskToTask(dbTask)
	}
	return tasks
}
