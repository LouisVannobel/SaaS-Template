package routes

import (
	"net/http"

	"github.com/LouisVannobel/SaaS-Template/backend/api/handlers"
	"github.com/LouisVannobel/SaaS-Template/backend/api/middleware"
	"github.com/LouisVannobel/SaaS-Template/backend/db"
	"github.com/gorilla/mux"
)

// SetupRouter configures and returns the API router
func SetupRouter(database *db.DB) *mux.Router {
	// Create router
	router := mux.NewRouter()

	// Create handlers
	userHandler := handlers.NewUserHandler(database)
	taskHandler := handlers.NewTaskHandler(database)

	// Auth routes (public)
	router.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", userHandler.Login).Methods("POST")

	// API routes (protected)
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.AuthMiddleware)

	// User routes
	apiRouter.HandleFunc("/users/profile", userHandler.GetProfile).Methods("GET")
	apiRouter.HandleFunc("/users/profile", userHandler.UpdateProfile).Methods("PUT")

	// Task routes
	apiRouter.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	apiRouter.HandleFunc("/tasks", taskHandler.GetAllTasks).Methods("GET")
	apiRouter.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.GetTask).Methods("GET")
	apiRouter.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.UpdateTask).Methods("PUT")
	apiRouter.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.DeleteTask).Methods("DELETE")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`)) 
	}).Methods("GET")

	return router
}
