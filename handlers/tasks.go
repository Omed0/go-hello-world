package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/omed0/go-hello-world/internal/auth"
	"github.com/omed0/go-hello-world/internal/database"
	"github.com/omed0/go-hello-world/models"
)

// CreateTaskRequest represents the request body for creating a task
type CreateTaskRequest struct {
	Title string `json:"title"`
}

// UpdateTaskRequest represents the request body for updating a task
type UpdateTaskRequest struct {
	Title string `json:"title"`
}

// SearchTasksRequest represents query parameters for searching tasks
type SearchTasksRequest struct {
	Query  string `json:"query,omitempty"`
	UserID string `json:"user_id,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Page   int    `json:"page,omitempty"`
}

// HandlerCreateTask creates a new task
func (api *ApiConfig) HandlerCreateTask(w http.ResponseWriter, r *http.Request) {
	// Get user from API key
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := api.Queries.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	var params CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Validate and sanitize input
	params.Title = strings.TrimSpace(params.Title)
	if params.Title == "" {
		RespondWithError(w, http.StatusBadRequest, "Task title is required")
		return
	}
	if len(params.Title) > 255 {
		RespondWithError(w, http.StatusBadRequest, "Task title must be less than 255 characters")
		return
	}

	task, err := api.Queries.CreateTask(r.Context(), database.CreateTaskParams{
		ID:     uuid.New(),
		Title:  params.Title,
		UserID: user.ID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create task: "+err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, models.DatabaseTaskToTask(task))
}

// HandlerGetTasks gets all tasks for the authenticated user
func (api *ApiConfig) HandlerGetTasks(w http.ResponseWriter, r *http.Request) {
	// Get user from API key
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := api.Queries.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	tasks, err := api.Queries.GetTasksByUserId(r.Context(), user.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to get tasks: "+err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, models.DatabaseTasksToTasks(tasks))
}

// HandlerGetTask gets a specific task by ID
func (api *ApiConfig) HandlerGetTask(w http.ResponseWriter, r *http.Request) {
	taskIDStr := chi.URLParam(r, "taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Get user from API key
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := api.Queries.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	task, err := api.Queries.GetTaskById(r.Context(), taskID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Task not found")
		return
	}

	// Check if the task belongs to the user
	if task.UserID != user.ID {
		RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	RespondWithJSON(w, http.StatusOK, models.DatabaseTaskToTask(task))
}

// HandlerUpdateTask updates a task
func (api *ApiConfig) HandlerUpdateTask(w http.ResponseWriter, r *http.Request) {
	taskIDStr := chi.URLParam(r, "taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Get user from API key
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := api.Queries.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	// Check if task exists and belongs to user
	existingTask, err := api.Queries.GetTaskById(r.Context(), taskID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Task not found")
		return
	}

	if existingTask.UserID != user.ID {
		RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	var params UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Validate and sanitize input
	params.Title = strings.TrimSpace(params.Title)
	if params.Title == "" {
		RespondWithError(w, http.StatusBadRequest, "Task title is required")
		return
	}
	if len(params.Title) > 255 {
		RespondWithError(w, http.StatusBadRequest, "Task title must be less than 255 characters")
		return
	}

	task, err := api.Queries.UpdateTask(r.Context(), database.UpdateTaskParams{
		ID:     taskID,
		Title:  params.Title,
		UserID: user.ID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to update task: "+err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, models.DatabaseTaskToTask(task))
}

// HandlerDeleteTask soft deletes a task
func (api *ApiConfig) HandlerDeleteTask(w http.ResponseWriter, r *http.Request) {
	taskIDStr := chi.URLParam(r, "taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Get user from API key
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := api.Queries.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	// Check if task exists and belongs to user
	existingTask, err := api.Queries.GetTaskById(r.Context(), taskID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Task not found")
		return
	}

	if existingTask.UserID != user.ID {
		RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	_, err = api.Queries.SoftDeleteTask(r.Context(), taskID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to delete task: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandlerSearchTasks searches tasks with filters
func (api *ApiConfig) HandlerSearchTasks(w http.ResponseWriter, r *http.Request) {
	// Get user from API key
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := api.Queries.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	// For now, just return user's tasks since the search query has parameter type issues
	// TODO: Fix the SQL query to properly handle UUID parameters
	tasks, err := api.Queries.GetTasksByUserId(r.Context(), user.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to search tasks: "+err.Error())
		return
	}

	// Filter by query if provided
	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("query")))
	if query != "" {
		filteredTasks := make([]database.Task, 0)
		for _, task := range tasks {
			if strings.Contains(strings.ToLower(task.Title), query) {
				filteredTasks = append(filteredTasks, task)
			}
		}
		tasks = filteredTasks
	}

	RespondWithJSON(w, http.StatusOK, models.DatabaseTasksToTasks(tasks))
}
