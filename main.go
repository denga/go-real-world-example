package main

import (
	"embed"
	"fmt"
	"git.homelab.lan/denga/go-real-world-example/api"
	"git.homelab.lan/denga/go-real-world-example/internal/auth"
	"git.homelab.lan/denga/go-real-world-example/internal/db"
	"git.homelab.lan/denga/go-real-world-example/internal/handlers"
	"git.homelab.lan/denga/go-real-world-example/internal/middleware"
	iofs "io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

//go:embed openapi.yml
var openAPISpec embed.FS

//go:embed frontend/dist/*
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

	// Create a separate router for API routes
	apiRouter := chi.NewRouter()

	// Add auth middleware only to API routes
	apiRouter.Use(middleware.Auth(authConfig))

	// Serve the OpenAPI spec
	apiRouter.Get("/openapi.yml", func(w http.ResponseWriter, r *http.Request) {
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
	apiHandler := api.HandlerFromMux(handler, apiRouter)

	// Mount the API handler to the main router without additional wrapping
	r.Mount("/api", apiHandler)

	// Serve the embedded frontend files
	frontendRoot, err := iofs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}

	// Create a custom file server with correct MIME types
	fileServer := http.FileServer(http.FS(frontendRoot))

	// Create a handler that sets the correct MIME types
	fileServerWithMIME := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set correct MIME types based on file extension
		if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		} else if strings.HasSuffix(r.URL.Path, ".json") {
			w.Header().Set("Content-Type", "application/json")
		} else if strings.HasSuffix(r.URL.Path, ".svg") {
			w.Header().Set("Content-Type", "image/svg+xml")
		} else if strings.HasSuffix(r.URL.Path, ".png") {
			w.Header().Set("Content-Type", "image/png")
		} else if strings.HasSuffix(r.URL.Path, ".jpg") || strings.HasSuffix(r.URL.Path, ".jpeg") {
			w.Header().Set("Content-Type", "image/jpeg")
		} else if strings.HasSuffix(r.URL.Path, ".woff2") {
			w.Header().Set("Content-Type", "font/woff2")
		} else if strings.HasSuffix(r.URL.Path, ".woff") {
			w.Header().Set("Content-Type", "font/woff")
		}
		fileServer.ServeHTTP(w, r)
	})

	// Serve all files with proper MIME types
	r.Handle("/*", fileServerWithMIME)

	// Determine port for HTTP service
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
