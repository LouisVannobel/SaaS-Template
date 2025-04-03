package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/LouisVannobel/SaaS-Template/backend/db"
	"github.com/LouisVannobel/SaaS-Template/backend/models"
	"github.com/LouisVannobel/SaaS-Template/backend/utils/logger"
	"github.com/gorilla/mux"
)

// TaskHandler handles HTTP requests related to tasks
type TaskHandler struct {
	taskRepo *models.TaskRepository
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(database *db.DB) *TaskHandler {
	return &TaskHandler{
		taskRepo: models.NewTaskRepository(database),
	}
}

// CreateTask handles task creation
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userIDValue := r.Context().Value("userID")
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDValue.(int)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	// Parse request body
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if task.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Set user ID from context
	task.UserID = userID

	// Set default status if not provided
	if task.Status == "" {
		task.Status = "pending"
	}

	// Save task to database
	if err := h.taskRepo.Create(&task); err != nil {
		logger.Error("Failed to create task: %v", err)
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	// Return created task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// GetAllTasks retrieves all tasks for the authenticated user
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userIDValue := r.Context().Value("userID")
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDValue.(int)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	// Get tasks from database
	tasks, err := h.taskRepo.GetAllByUserID(userID)
	if err != nil {
		logger.Error("Failed to get tasks: %v", err)
		http.Error(w, "Failed to get tasks", http.StatusInternalServerError)
		return
	}

	// Return tasks
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// GetTask retrieves a specific task by ID
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userIDValue := r.Context().Value("userID")
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDValue.(int)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	// Get task ID from URL parameters
	vars := mux.Vars(r)
	taskIDStr := vars["id"]
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Get task from database
	task, err := h.taskRepo.GetByID(taskID, userID)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Return task
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// UpdateTask updates a specific task
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userIDValue := r.Context().Value("userID")
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDValue.(int)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	// Get task ID from URL parameters
	vars := mux.Vars(r)
	taskIDStr := vars["id"]
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set task ID and user ID
	task.ID = taskID
	task.UserID = userID

	// Update task in database
	if err := h.taskRepo.Update(&task); err != nil {
		logger.Error("Failed to update task: %v", err)
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	// Return updated task
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// DeleteTask deletes a specific task
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userIDValue := r.Context().Value("userID")
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDValue.(int)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	// Get task ID from URL parameters
	vars := mux.Vars(r)
	taskIDStr := vars["id"]
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Delete task from database
	if err := h.taskRepo.Delete(taskID, userID); err != nil {
		logger.Error("Failed to delete task: %v", err)
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusNoContent)
}
