package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LouisVannobel/SaaS-Template/backend/api/handlers"
	"github.com/LouisVannobel/SaaS-Template/backend/api/middleware"
	"github.com/LouisVannobel/SaaS-Template/backend/db"
	"github.com/LouisVannobel/SaaS-Template/backend/utils/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger" // Fiber logger middleware
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Initialize logger
	logger.InitLogger()

	// Get port from environment variable or use default
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to database
	database, err := db.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default 500 statuscode
			code := fiber.StatusInternalServerError

			// Check if it's a Fiber error
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Return JSON error response
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173", // Frontend URL
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))
	app.Use(fiberlogger.New())

	// Setup routes
	setupRoutes(app, database)

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

// setupRoutes configures all the routes for our application
func setupRoutes(app *fiber.App, database *db.DB) {
	// Create handlers
	userHandler := handlers.NewFiberUserHandler(database)
	taskHandler := handlers.NewFiberTaskHandler(database)

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// API routes
	api := app.Group("/api")

	// Public routes (no auth required)
	api.Post("/register", userHandler.Register)
	api.Post("/login", userHandler.Login)

	// Protected routes (auth required)
	// Create a protected group
	protected := api.Group("/")
	protected.Use(middleware.JWTProtected())

	// User routes
	protected.Get("/users/profile", userHandler.GetProfile)
	protected.Put("/users/profile", userHandler.UpdateProfile)

	// Task routes
	protected.Post("/tasks", taskHandler.CreateTask)
	protected.Get("/tasks", taskHandler.GetAllTasks)
	protected.Get("/tasks/:id", taskHandler.GetTask)
	protected.Put("/tasks/:id", taskHandler.UpdateTask)
	protected.Delete("/tasks/:id", taskHandler.DeleteTask)
}
