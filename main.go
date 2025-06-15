package main

import (
	"embed"
	"fmt"
	"git.homelab.lan/denga/go-real-world-example/api"
	"git.homelab.lan/denga/go-real-world-example/internal/auth"
	"git.homelab.lan/denga/go-real-world-example/internal/db"
	"git.homelab.lan/denga/go-real-world-example/internal/handlers"
	"git.homelab.lan/denga/go-real-world-example/internal/middleware"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

//go:embed openapi.yml
var openAPISpec embed.FS

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
	// Create a new router
	r := chi.NewRouter()

	// Add middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Create auth config
	authConfig := auth.DefaultConfig()

	// Add auth middleware
	r.Use(middleware.Auth(authConfig))

	// Serve the OpenAPI spec
	r.Get("/openapi.yml", func(w http.ResponseWriter, r *http.Request) {
		specBytes, err := openAPISpec.ReadFile("openapi.yml")
		if err != nil {
			http.Error(w, "Could not read OpenAPI spec", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/yaml")
		w.Write(specBytes)
	})

	// Initialize in-memory database
	db := db.NewInMemoryDB()

	// Create API handlers
	handler := handlers.NewHandler(db, authConfig)

	// Register API handlers
	apiHandler := api.HandlerFromMux(handler, r)

	// Mount the API handler
	r.Mount("/", apiHandler)

	// Determine port for HTTP service
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
