package main

import (
	//i have a package name user

	"log"
	"net/http"
	"os"

	"github.com/omed0/go-hello-world/internal/api"
	"github.com/omed0/go-hello-world/user"
	"github.com/omed0/go-hello-world/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("You must setup PORT Env on .env file")
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	db, er := api.GetDB()
	if er != nil {
		log.Fatalf("Failed to connect to the database: %v", er)
	}
	defer db.Close()

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", utils.HandlerReadiness)
	v1Router.Get("/err", utils.HandlerErr)
	v1Router.Post("/user", user.HandlerCreateUser)
	v1Router.Get("/user", user.HandlerGetUser)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + PORT,
	}

	// Print a greeting message
	log.Printf("Server Startup with PORT: %s", PORT)

	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
