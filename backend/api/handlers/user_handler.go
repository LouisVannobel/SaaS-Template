package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/LouisVannobel/SaaS-Template/backend/auth"
	"github.com/LouisVannobel/SaaS-Template/backend/db"
	"github.com/LouisVannobel/SaaS-Template/backend/models"
	"github.com/LouisVannobel/SaaS-Template/backend/utils/logger"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userRepo *models.UserRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler(database *db.DB) *UserHandler {
	return &UserHandler{
		userRepo: models.NewUserRepository(database),
	}
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" || req.Name == "" {
		http.Error(w, "Email, password, and name are required", http.StatusBadRequest)
		return
	}

	// Create user object
	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	// Save user to database
	if err := h.userRepo.Create(user); err != nil {
		logger.Error("Failed to create user: %v", err)
		// Log l'erreur complète pour le débogage
		log.Printf("Detailed error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		logger.Error("Failed to generate token: %v", err)
		http.Error(w, "Failed to generate authentication token", http.StatusInternalServerError)
		return
	}

	// Return token and user info
	response := AuthResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Login handles user login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Check password
	hashedPassword := models.HashPassword(req.Password)
	if user.Password != hashedPassword {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		logger.Error("Failed to generate token: %v", err)
		http.Error(w, "Failed to generate authentication token", http.StatusInternalServerError)
		return
	}

	// Clear password before sending response
	user.Password = ""

	// Return token and user info
	response := AuthResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetProfile retrieves the user's profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
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

	// Get user from database
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		logger.Error("Failed to get user profile: %v", err)
		http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
		return
	}

	// Return user info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateProfile updates the user's profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
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
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set user ID from context
	user.ID = userID

	// Update user in database
	if err := h.userRepo.Update(&user); err != nil {
		logger.Error("Failed to update user profile: %v", err)
		http.Error(w, "Failed to update user profile", http.StatusInternalServerError)
		return
	}

	// Get updated user from database
	updatedUser, err := h.userRepo.GetByID(userID)
	if err != nil {
		logger.Error("Failed to get updated user profile: %v", err)
		http.Error(w, "Failed to get updated user profile", http.StatusInternalServerError)
		return
	}

	// Return updated user info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}
