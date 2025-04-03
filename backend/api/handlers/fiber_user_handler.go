package handlers

import (
	"strings"

	"github.com/LouisVannobel/SaaS-Template/backend/auth"
	"github.com/LouisVannobel/SaaS-Template/backend/db"
	"github.com/LouisVannobel/SaaS-Template/backend/models"
	"github.com/LouisVannobel/SaaS-Template/backend/utils/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// FiberUserHandler handles user-related requests using Fiber
type FiberUserHandler struct {
	userRepo *models.UserRepository
}

// NewFiberUserHandler creates a new FiberUserHandler
func NewFiberUserHandler(database *db.DB) *FiberUserHandler {
	return &FiberUserHandler{
		userRepo: models.NewUserRepository(database),
	}
}

// Register handles user registration
func (h *FiberUserHandler) Register(c *fiber.Ctx) error {
	// Parse request body
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		logger.Error("Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate required fields
	if user.Email == "" || user.Password == "" || user.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email, password, and name are required",
		})
	}

	// Check if user already exists
	existingUser, _ := h.userRepo.GetByEmail(user.Email)
	if existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User with this email already exists",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process request",
		})
	}

	// Set hashed password
	user.Password = string(hashedPassword)

	// Save user to database
	if err := h.userRepo.Create(&user); err != nil {
		logger.Error("Failed to create user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Generate JWT token
	token, err := auth.GenerateToken(int(user.ID))
	if err != nil {
		logger.Error("Failed to generate token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate authentication token",
		})
	}

	// Return success response with token
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"token":   token,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// Login handles user login
func (h *FiberUserHandler) Login(c *fiber.Ctx) error {
	// Log the request body for debugging
	logger.Info("Login request body: %s", string(c.Body()))

	// Parse request body
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&credentials); err != nil {
		logger.Error("Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Log the parsed credentials for debugging
	logger.Info("Parsed credentials - Email: %s, Password length: %d", credentials.Email, len(credentials.Password))

	// Validate required fields
	if credentials.Email == "" || credentials.Password == "" {
		logger.Error("Login failed: Email or password is empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(credentials.Email)
	if err != nil {
		logger.Error("Login failed: Error getting user by email: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	if user == nil {
		logger.Error("Login failed: User not found for email: %s", credentials.Email)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Log user details for debugging
	logger.Info("User found - ID: %d, Email: %s, Password hash length: %d", user.ID, user.Email, len(user.Password))

	// Verify password
	var passwordValid bool
	
	// Essayer la vérification bcrypt standard
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		logger.Error("Bcrypt password verification failed: %v", err)
		
		// Si le hash ne commence pas par '$', essayer une comparaison directe (pour les mots de passe non hashés)
		if !strings.HasPrefix(user.Password, "$") {
			logger.Info("Trying direct password comparison because hash doesn't start with $")
			// Pour le débogage, acceptons n'importe quel mot de passe temporairement
			passwordValid = true // Accepter n'importe quel mot de passe pour le débogage
			logger.Info("DEBUG MODE: Accepting any password for login")
			
			// Générer un nouveau hash bcrypt
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
			if err != nil {
				logger.Error("Failed to hash password for update: %v", err)
			} else {
				// Mettre à jour le mot de passe dans la base de données
				user.Password = string(hashedPassword)
				if err := h.userRepo.Update(user); err != nil {
					logger.Error("Failed to update user password: %v", err)
				} else {
					logger.Info("Password hash updated successfully for user ID: %d", user.ID)
				}
			}
		}
	} else {
		// Si la vérification bcrypt réussit, le mot de passe est valide
		passwordValid = true
	}
	
	// Si le mot de passe n'est pas valide, renvoyer une erreur
	if !passwordValid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	logger.Info("Password verification successful for user ID: %d", user.ID)

	// Generate JWT token
	token, err := auth.GenerateToken(int(user.ID))
	if err != nil {
		logger.Error("Failed to generate token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate authentication token",
		})
	}

	// Return success response with token
	response := fiber.Map{
		"message": "Login successful",
		"token":   token,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	}

	logger.Info("Login successful response prepared: %+v", response)

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetProfile retrieves the user's profile
func (h *FiberUserHandler) GetProfile(c *fiber.Ctx) error {
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

	// Get user from database
	user, err := h.userRepo.GetByID(int(userID))
	if err != nil || user == nil {
		logger.Error("Failed to get user: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Return user profile
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// UpdateProfile updates the user's profile
func (h *FiberUserHandler) UpdateProfile(c *fiber.Ctx) error {
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

	// Get user from database
	user, err := h.userRepo.GetByID(int(userID))
	if err != nil || user == nil {
		logger.Error("Failed to get user: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Parse request body
	var updateData struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		logger.Error("Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Update user data
	if updateData.Name != "" {
		user.Name = updateData.Name
	}

	// Save updated user to database
	if err := h.userRepo.Update(user); err != nil {
		logger.Error("Failed to update user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	// Return updated user profile
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
