package handlers

import (
	"strconv"

	"github.com/LouisVannobel/SaaS-Template/backend/db"
	"github.com/LouisVannobel/SaaS-Template/backend/models"
	"github.com/LouisVannobel/SaaS-Template/backend/utils/logger"
	"github.com/gofiber/fiber/v2"
)

// FiberTaskHandler handles task-related requests using Fiber
type FiberTaskHandler struct {
	taskRepo *models.TaskRepository
}

// NewFiberTaskHandler creates a new FiberTaskHandler
func NewFiberTaskHandler(database *db.DB) *FiberTaskHandler {
	return &FiberTaskHandler{
		taskRepo: models.NewTaskRepository(database),
	}
}

// CreateTask handles task creation
func (h *FiberTaskHandler) CreateTask(c *fiber.Ctx) error {
	// Get user ID from context (set by JWTProtected middleware)
	userIDValue := c.Locals("userID")
	// Log the type and value for debugging
	logger.Info("UserID type: %T, value: %v", userIDValue, userIDValue)
	
	var userID int
	switch v := userIDValue.(type) {
	case int:
		userID = v
	case float64:
		userID = int(v)
	case uint:
		userID = int(v)
	default:
		logger.Error("Invalid userID type: %T", userIDValue)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Invalid user ID",
		})
	}

	// Parse request body
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		logger.Error("Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate required fields
	if task.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	// Set user ID
	task.UserID = userID

	// Save task to database
	if err := h.taskRepo.Create(&task); err != nil {
		logger.Error("Failed to create task: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create task",
		})
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Task created successfully",
		"task":    task,
	})
}

// GetAllTasks retrieves all tasks for the authenticated user
func (h *FiberTaskHandler) GetAllTasks(c *fiber.Ctx) error {
	// Get user ID from context (set by JWTProtected middleware)
	userIDValue := c.Locals("userID")
	// Log the type and value for debugging
	logger.Info("UserID type: %T, value: %v", userIDValue, userIDValue)
	
	var userID int
	switch v := userIDValue.(type) {
	case int:
		userID = v
	case float64:
		userID = int(v)
	case uint:
		userID = int(v)
	default:
		logger.Error("Invalid userID type: %T", userIDValue)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Invalid user ID",
		})
	}

	// Get tasks from database
	tasks, err := h.taskRepo.GetAllByUserID(userID)
	if err != nil {
		logger.Error("Failed to get tasks: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve tasks",
		})
	}

	// Return tasks
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tasks": tasks,
	})
}

// GetTask retrieves a specific task by ID
func (h *FiberTaskHandler) GetTask(c *fiber.Ctx) error {
	// Get user ID from context (set by JWTProtected middleware)
	userIDValue := c.Locals("userID")
	// Log the type and value for debugging
	logger.Info("UserID type: %T, value: %v", userIDValue, userIDValue)
	
	var userID int
	switch v := userIDValue.(type) {
	case int:
		userID = v
	case float64:
		userID = int(v)
	case uint:
		userID = int(v)
	default:
		logger.Error("Invalid userID type: %T", userIDValue)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Invalid user ID",
		})
	}

	// Get task ID from URL parameter
	taskIDStr := c.Params("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
		})
	}

	// Get task from database
	task, err := h.taskRepo.GetByID(int(taskID), userID)
	if err != nil || task == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	// The task belongs to the authenticated user (already checked in GetByID)

	// Return task
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"task": task,
	})
}

// UpdateTask updates a specific task by ID
func (h *FiberTaskHandler) UpdateTask(c *fiber.Ctx) error {
	// Get user ID from context (set by JWTProtected middleware)
	userIDValue := c.Locals("userID")
	// Log the type and value for debugging
	logger.Info("UserID type: %T, value: %v", userIDValue, userIDValue)
	
	var userID int
	switch v := userIDValue.(type) {
	case int:
		userID = v
	case float64:
		userID = int(v)
	case uint:
		userID = int(v)
	default:
		logger.Error("Invalid userID type: %T", userIDValue)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Invalid user ID",
		})
	}

	// Get task ID from URL parameter
	taskIDStr := c.Params("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
		})
	}

	// Get existing task from database
	existingTask, err := h.taskRepo.GetByID(int(taskID), userID)
	if err != nil || existingTask == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	// The task belongs to the authenticated user (already checked in GetByID)

	// Parse request body
	var updatedTask models.Task
	if err := c.BodyParser(&updatedTask); err != nil {
		logger.Error("Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate required fields
	if updatedTask.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	// Update task fields
	existingTask.Title = updatedTask.Title
	existingTask.Description = updatedTask.Description
	existingTask.Status = updatedTask.Status

	// Save updated task to database
	if err := h.taskRepo.Update(existingTask); err != nil {
		logger.Error("Failed to update task: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update task",
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task updated successfully",
		"task":    existingTask,
	})
}

// DeleteTask deletes a specific task by ID
func (h *FiberTaskHandler) DeleteTask(c *fiber.Ctx) error {
	// Get user ID from context (set by JWTProtected middleware)
	userIDValue := c.Locals("userID")
	// Log the type and value for debugging
	logger.Info("UserID type: %T, value: %v", userIDValue, userIDValue)
	
	var userID int
	switch v := userIDValue.(type) {
	case int:
		userID = v
	case float64:
		userID = int(v)
	case uint:
		userID = int(v)
	default:
		logger.Error("Invalid userID type: %T", userIDValue)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Invalid user ID",
		})
	}

	// Get task ID from URL parameter
	taskIDStr := c.Params("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
		})
	}

	// Get existing task from database
	existingTask, err := h.taskRepo.GetByID(int(taskID), userID)
	if err != nil || existingTask == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	// The task belongs to the authenticated user (already checked in GetByID)

	// Delete task from database
	if err := h.taskRepo.Delete(int(taskID), userID); err != nil {
		logger.Error("Failed to delete task: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete task",
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task deleted successfully",
	})
}
