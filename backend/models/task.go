package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/LouisVannobel/SaaS-Template/backend/db"
)

// Task represents a task in the system
type Task struct {
	ID          int          `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      string       `json:"status"`
	DueDate     *time.Time   `json:"due_date,omitempty"`
	UserID      int          `json:"user_id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// TaskRepository handles database operations for tasks
type TaskRepository struct {
	DB *db.DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(database *db.DB) *TaskRepository {
	return &TaskRepository{DB: database}
}

// Create adds a new task to the database
func (r *TaskRepository) Create(task *Task) error {
	query := `
		INSERT INTO tasks (title, description, status, due_date, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.DB.QueryRow(
		query,
		task.Title,
		task.Description,
		task.Status,
		task.DueDate,
		task.UserID,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)

	return err
}

// GetByID retrieves a task by ID and user ID
func (r *TaskRepository) GetByID(id, userID int) (*Task, error) {
	task := &Task{}

	query := `
		SELECT id, title, description, status, due_date, user_id, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND user_id = $2
	`

	err := r.DB.QueryRow(query, id, userID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.DueDate,
		&task.UserID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	return task, nil
}

// GetAllByUserID retrieves all tasks for a specific user
func (r *TaskRepository) GetAllByUserID(userID int) ([]*Task, error) {
	query := `
		SELECT id, title, description, status, due_date, user_id, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*Task{}

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.DueDate,
			&task.UserID,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Update updates a task
func (r *TaskRepository) Update(task *Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, due_date = $4, updated_at = NOW()
		WHERE id = $5 AND user_id = $6
		RETURNING updated_at
	`

	result, err := r.DB.Exec(
		query,
		task.Title,
		task.Description,
		task.Status,
		task.DueDate,
		task.ID,
		task.UserID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("task not found or not owned by user")
	}

	return nil
}

// Delete removes a task
func (r *TaskRepository) Delete(id, userID int) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`

	result, err := r.DB.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("task not found or not owned by user")
	}

	return nil
}
