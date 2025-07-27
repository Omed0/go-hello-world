package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/omed0/go-hello-world/handlers"

	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: failed to load .env file: %v", err)
	}

	// Get PORT from environment
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is required")
	}

	// Initialize database connection
	apiCfg, err := handlers.NewApiConfig()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer handlers.CloseDB()

	// Create router
	router := chi.NewRouter()

	// CORS configuration
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// API v1 routes
	v1Router := chi.NewRouter()

	// Health and error endpoints
	v1Router.Get("/healthz", handlers.HandlerReadiness)
	v1Router.Get("/err", handlers.HandlerErr)

	// User endpoints
	v1Router.Post("/user", apiCfg.HandlerCreateUser)
	v1Router.Get("/user", apiCfg.HandlerGetUser)

	// Task endpoints
	v1Router.Post("/tasks", apiCfg.HandlerCreateTask)
	v1Router.Get("/tasks", apiCfg.HandlerGetTasks)
	v1Router.Get("/tasks/search", apiCfg.HandlerSearchTasks)
	v1Router.Get("/tasks/{taskId}", apiCfg.HandlerGetTask)
	v1Router.Put("/tasks/{taskId}", apiCfg.HandlerUpdateTask)
	v1Router.Delete("/tasks/{taskId}", apiCfg.HandlerDeleteTask)

	router.Mount("/v1", v1Router)

	// Create server
	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
