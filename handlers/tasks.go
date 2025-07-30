package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/omed0/go-hello-world/internal/auth"
	"github.com/omed0/go-hello-world/internal/database"
	"github.com/omed0/go-hello-world/models"
)

// Constants for better maintainability
const (
	maxTitleLength = 255
	maxDescLength  = 2255
	defaultLimit   = 10
	maxLimit       = 100
)

// Error messages as constants
const (
	errRequiredFields    = "Task title and description are required"
	errTitleTooLong      = "Task title must be less than 255 characters"
	errDescTooLong       = "Task description must be less than 2255 characters"
	errInvalidTitle      = "Task title can only contain alphanumeric characters and spaces"
	errInvalidDesc       = "Task description can only contain alphanumeric characters, spaces, and punctuation (.,!?-)"
	errInvalidJSON       = "Invalid JSON format"
	errInvalidTaskID     = "Invalid task ID"
	errTaskNotFound      = "Task not found"
	errAccessDenied      = "Access denied"
	errCreateTaskFailed  = "Failed to create task"
	errUpdateTaskFailed  = "Failed to update task"
	errDeleteTaskFailed  = "Failed to delete task"
	errSearchTasksFailed = "Failed to search tasks"
	errUserNotFound      = "User not found in context"
)

// Compile regex patterns once at package level
var (
	regexTitlePattern = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
	regexDescPattern  = regexp.MustCompile(`^[a-zA-Z0-9\s\.,!?-]+$`)
)

// CreateTaskRequest represents the request body for creating a task
type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required,max=255"`
	Description string `json:"description,omitempty" validate:"max=2255"`
}

// UpdateTaskRequest represents the request body for updating a task
type UpdateTaskRequest struct {
	Title       string `json:"title" validate:"required,max=255"`
	Description string `json:"description,omitempty" validate:"max=2255"`
	IsCompleted *bool  `json:"is_completed,omitempty"`
}

// ToggleCompletionRequest represents the request body for toggling task completion
type ToggleCompletionRequest struct {
	IsCompleted bool `json:"is_completed"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

// validateTaskInput validates and sanitizes task input
func validateTaskInput(title, description string) error {
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)

	if title == "" {
		return &ValidationError{Message: errRequiredFields}
	}

	if len(title) > maxTitleLength {
		return &ValidationError{Message: errTitleTooLong}
	}

	if len(description) > maxDescLength {
		return &ValidationError{Message: errDescTooLong}
	}

	if !regexTitlePattern.MatchString(title) {
		return &ValidationError{Message: errInvalidTitle}
	}

	if description != "" && !regexDescPattern.MatchString(description) {
		return &ValidationError{Message: errInvalidDesc}
	}

	return nil
}

// parseTaskID parses and validates task ID from string
func parseTaskID(taskIDStr string) (uuid.UUID, error) {
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return uuid.Nil, &ValidationError{Message: errInvalidTaskID}
	}
	return taskID, nil
}

// checkTaskOwnership verifies task ownership without additional DB query
func checkTaskOwnership(task database.Task, userID uuid.UUID) error {
	if task.UserID != userID {
		return &ValidationError{Message: errAccessDenied}
	}
	return nil
}

// filterTasksByQuery filters tasks by search query
func filterTasksByQuery(tasks []database.Task, query string) []database.Task {
	if query == "" {
		return tasks
	}

	filtered := make([]database.Task, 0, len(tasks))
	queryLower := strings.ToLower(query)

	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Title), queryLower) ||
			strings.Contains(strings.ToLower(task.Description), queryLower) {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

// decodeAndValidateTaskRequest decodes and validates task request
func decodeAndValidateTaskRequest(r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return &ValidationError{Message: errInvalidJSON}
	}

	switch req := dst.(type) {
	case *CreateTaskRequest:
		req.Title = strings.TrimSpace(req.Title)
		req.Description = strings.TrimSpace(req.Description)
		return validateTaskInput(req.Title, req.Description)
	case *UpdateTaskRequest:
		req.Title = strings.TrimSpace(req.Title)
		req.Description = strings.TrimSpace(req.Description)
		return validateTaskInput(req.Title, req.Description)
	case *ToggleCompletionRequest:
		// No validation needed for boolean
		return nil
	}

	return nil
}

// HandlerCreateTask creates a new task
func (api *ApiConfig) HandlerCreateTask(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, errUserNotFound)
		return
	}

	var params CreateTaskRequest
	if err := decodeAndValidateTaskRequest(r, &params); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := api.Queries.CreateTask(r.Context(), database.CreateTaskParams{
		ID:          uuid.New(),
		Title:       params.Title,
		Description: params.Description,
		UserID:      userID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, errCreateTaskFailed)
		return
	}

	RespondWithJSON(w, http.StatusCreated, models.DatabaseTaskToTask(task))
}

// HandlerGetTasks gets all tasks for the authenticated user
func (api *ApiConfig) HandlerGetTasks(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, errUserNotFound)
		return
	}

	tasks, err := api.Queries.GetTasksByUserId(r.Context(), userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, errSearchTasksFailed)
		return
	}

	RespondWithJSON(w, http.StatusOK, models.DatabaseTasksToTasks(tasks))
}

// HandlerGetTask gets a specific task by ID
func (api *ApiConfig) HandlerGetTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := parseTaskID(chi.URLParam(r, "taskId"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, errUserNotFound)
		return
	}

	task, err := api.Queries.GetTaskById(r.Context(), taskID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, errTaskNotFound)
		return
	}

	// Verify task ownership
	if err := checkTaskOwnership(task, userID); err != nil {
		RespondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, models.DatabaseTaskToTask(task))
}

// HandlerUpdateTask updates a task
func (api *ApiConfig) HandlerUpdateTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := parseTaskID(chi.URLParam(r, "taskId"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, errUserNotFound)
		return
	}

	// Get task first to check ownership
	task, err := api.Queries.GetTaskById(r.Context(), taskID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, errTaskNotFound)
		return
	}

	// Verify task ownership
	if err := checkTaskOwnership(task, userID); err != nil {
		RespondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	var params UpdateTaskRequest
	if err := decodeAndValidateTaskRequest(r, &params); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update task title and description
	updatedTask, err := api.Queries.UpdateTaskPartial(r.Context(), database.UpdateTaskPartialParams{
		Column1: params.Title,
		Column2: params.Description,
		ID:      taskID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, errUpdateTaskFailed)
		return
	}

	// If completion status is being updated, handle it separately
	if params.IsCompleted != nil {
		if *params.IsCompleted && !task.IsCompleted {
			// Mark as completed
			updatedTask, err = api.Queries.CompleteTask(r.Context(), taskID)
		} else if !*params.IsCompleted && task.IsCompleted {
			// Mark as incomplete
			updatedTask, err = api.Queries.UndoCompleteTask(r.Context(), taskID)
		}
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, errUpdateTaskFailed)
			return
		}
	}

	RespondWithJSON(w, http.StatusOK, models.DatabaseTaskToTask(updatedTask))
}

// HandlerDeleteTask soft deletes a task
func (api *ApiConfig) HandlerDeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := parseTaskID(chi.URLParam(r, "taskId"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, errUserNotFound)
		return
	}

	// Get task first to check ownership
	task, err := api.Queries.GetTaskById(r.Context(), taskID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, errTaskNotFound)
		return
	}

	// Verify task ownership
	if err := checkTaskOwnership(task, userID); err != nil {
		RespondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	// Soft delete the task
	if _, err := api.Queries.SoftDeleteTask(r.Context(), taskID); err != nil {
		RespondWithError(w, http.StatusInternalServerError, errDeleteTaskFailed)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandlerSearchTasks searches tasks with filters
func (api *ApiConfig) HandlerSearchTasks(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, errUserNotFound)
		return
	}

	// Parse query parameters
	query := strings.TrimSpace(r.URL.Query().Get("query"))
	limit := parseLimit(r.URL.Query().Get("limit"))

	tasks, err := api.Queries.GetTasksByUserId(r.Context(), userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, errSearchTasksFailed)
		return
	}

	// Filter by query if provided
	filteredTasks := filterTasksByQuery(tasks, query)

	// Apply limit
	if limit > 0 && len(filteredTasks) > limit {
		filteredTasks = filteredTasks[:limit]
	}

	RespondWithJSON(w, http.StatusOK, models.DatabaseTasksToTasks(filteredTasks))
}

// HandlerToggleTaskCompletion toggles the completion status of a task
func (api *ApiConfig) HandlerToggleTaskCompletion(w http.ResponseWriter, r *http.Request) {
	taskID, err := parseTaskID(chi.URLParam(r, "taskId"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, errUserNotFound)
		return
	}

	// Get task first to check ownership
	task, err := api.Queries.GetTaskById(r.Context(), taskID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, errTaskNotFound)
		return
	}

	// Verify task ownership
	if err := checkTaskOwnership(task, userID); err != nil {
		RespondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	var params ToggleCompletionRequest
	if err := decodeAndValidateTaskRequest(r, &params); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var updatedTask database.Task

	// Toggle completion based on request
	if params.IsCompleted && !task.IsCompleted {
		// Mark as completed
		updatedTask, err = api.Queries.CompleteTask(r.Context(), taskID)
	} else if !params.IsCompleted && task.IsCompleted {
		// Mark as incomplete
		updatedTask, err = api.Queries.UndoCompleteTask(r.Context(), taskID)
	} else {
		// No change needed, return current state
		RespondWithJSON(w, http.StatusOK, models.DatabaseTaskToTask(task))
		return
	}

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, errUpdateTaskFailed)
		return
	}

	RespondWithJSON(w, http.StatusOK, models.DatabaseTaskToTask(updatedTask))
}
