package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/LouisVannobel/SaaS-Template/backend/db"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository handles database operations for users
type UserRepository struct {
	DB *db.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(database *db.DB) *UserRepository {
	return &UserRepository{DB: database}
}

// HashPassword creates a SHA256 hash of the password
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// Create adds a new user to the database
func (r *UserRepository) Create(user *User) error {
	// Hash the password before storing
	hashedPassword := HashPassword(user.Password)

	// SQL query to insert a new user
	query := `
		INSERT INTO users (email, password, name, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	// Execute the query
	err := r.DB.QueryRow(query, user.Email, hashedPassword, user.Name).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// Clear password from the returned user object for security
	user.Password = ""
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int) (*User, error) {
	user := &User{}

	query := `SELECT id, email, name, created_at, updated_at FROM users WHERE id = $1`
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*User, error) {
	user := &User{}

	query := `SELECT id, email, password, name, created_at, updated_at FROM users WHERE email = $1`
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// Update updates a user's information
func (r *UserRepository) Update(user *User) error {
	query := `
		UPDATE users
		SET name = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING updated_at
	`

	err := r.DB.QueryRow(query, user.Name, user.ID).Scan(&user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(userID int, password string) error {
	hashedPassword := HashPassword(password)

	query := `UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.DB.Exec(query, hashedPassword, userID)
	return err
}
