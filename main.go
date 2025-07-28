package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/omed0/go-hello-world/handlers"
	"github.com/omed0/go-hello-world/internal/config"
	"github.com/omed0/go-hello-world/internal/middleware"

	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: failed to load .env file: %v", err)
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Validate required configuration
	if cfg.DatabaseURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}

	// Initialize database connection
	apiCfg, err := handlers.NewApiConfig()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer handlers.CloseDB()

	// Create router
	router := chi.NewRouter()

	// Global middleware
	router.Use(middleware.Recovery)      // Panic recovery
	router.Use(middleware.RequestLogger) // Request logging

	// CORS configuration
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// API v1 routes
	v1Router := chi.NewRouter()

	// Public endpoints (no authentication required)
	v1Router.Get("/healthz", handlers.HandlerReadiness)
	v1Router.Get("/err", handlers.HandlerErr)
	v1Router.Post("/user", apiCfg.HandlerCreateUser)
	v1Router.Post("/login", apiCfg.HandlerLogin)

	// Protected endpoints (authentication required)
	v1Router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(apiCfg))

		// User endpoints
		r.Get("/user", apiCfg.HandlerGetUser)
		r.Put("/user", apiCfg.HandlerUpdateUser)

		// Organization endpoints
		r.Post("/organizations", apiCfg.HandlerCreateOrganization)
		r.Get("/organizations/{orgId}", apiCfg.HandlerGetOrganization)
		r.With(middleware.RequireRole(apiCfg, "admin", "owner")).Put("/organizations/{orgId}", apiCfg.HandlerUpdateOrganization)
		r.With(middleware.RequireRole(apiCfg, "owner")).Delete("/organizations/{orgId}", apiCfg.HandlerDeleteOrganization)
		r.Get("/organizations/{orgId}/users", apiCfg.HandlerGetOrganizationUsers)

		// Task endpoints
		r.Post("/tasks", apiCfg.HandlerCreateTask)
		r.Get("/tasks", apiCfg.HandlerGetTasks)
		r.Get("/tasks/search", apiCfg.HandlerSearchTasks)
		r.Get("/tasks/{taskId}", apiCfg.HandlerGetTask)
		r.Put("/tasks/{taskId}", apiCfg.HandlerUpdateTask)
		r.Delete("/tasks/{taskId}", apiCfg.HandlerDeleteTask)
		r.Patch("/tasks/{taskId}/complete", apiCfg.HandlerToggleTaskCompletion)
	})

	router.Mount("/v1", v1Router)

	// Create server with configuration-based timeouts
	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + cfg.Port,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block until a signal is received
	log.Println("Shutting down server...")

	// Graceful shutdown with configurable timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ServerTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
