package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pratham15541/go-crud/internal/config"
	"github.com/pratham15541/go-crud/internal/database"
	"github.com/pratham15541/go-crud/internal/handlers"
	"github.com/pratham15541/go-crud/internal/middleware"
	"github.com/pratham15541/go-crud/internal/repository"
	"github.com/pratham15541/go-crud/internal/services"
)

// @title Go CRUD API
// @version 1.0
// @description A production-ready CRUD API built with Go
// @contact.name Pratham Parikh
// @contact.email prathamparikh94@gmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	healthHandler := handlers.NewHealthHandler(db)

	// Setup router
	router := mux.NewRouter()

	// Add middleware
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CORSMiddleware)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Health check
	api.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")

	// User routes
	userRoutes := api.PathPrefix("/users").Subrouter()
	userRoutes.HandleFunc("", userHandler.GetUsers).Methods("GET")
	userRoutes.HandleFunc("", userHandler.CreateUser).Methods("POST")
	userRoutes.HandleFunc("/{id:[0-9]+}", userHandler.GetUser).Methods("GET")
	userRoutes.HandleFunc("/{id:[0-9]+}", userHandler.UpdateUser).Methods("PUT")
	userRoutes.HandleFunc("/{id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")

	// Create server
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: router,
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}