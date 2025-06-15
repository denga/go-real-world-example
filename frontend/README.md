# Frontend for Go Real World Example

This is the frontend for the Go Real World Example project. It's built with Next.js and chadcn UI.

## Features

- Dark mode support
- Mobile responsive design
- Integration with backend API
- Built with chadcn UI components

## Getting Started

First, run the development server:

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## Building for Production

To build the frontend for production:

```bash
npm run build
```

This will create a static export in the `dist` directory as configured in `next.config.ts`.

## Embedding in the Backend Binary

To embed the frontend in the backend binary, follow these steps:

1. Build the frontend:

```bash
cd frontend
npm run build
```

2. Update the main.go file to embed and serve the frontend:

```go
package main

import (
    "embed"
    "fmt"
    "io/fs"
    "log"
    "net/http"
    "os"
    "strings"

    // ... other imports
)

//go:embed openapi.yml
var openAPISpec embed.FS

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
    // Create a new router
    r := chi.NewRouter()

    // ... other middleware and setup

    // Serve the frontend
    frontendRoot, err := fs.Sub(frontendFS, "frontend/dist")
    if err != nil {
        log.Fatal(err)
    }

    // Create a file server handler for the frontend
    fileServer := http.FileServer(http.FS(frontendRoot))

    // Serve the frontend at the root path
    r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
        // If the path is for the API, don't serve the frontend
        if strings.HasPrefix(r.URL.Path, "/api") {
            http.NotFound(w, r)
            return
        }

        // Check if the file exists
        _, err := fs.Stat(frontendRoot, strings.TrimPrefix(r.URL.Path, "/"))
        if err != nil {
            // If the file doesn't exist, serve the index.html file
            r.URL.Path = "/"
        }

        fileServer.ServeHTTP(w, r)
    })

    // ... rest of the code
}
```

3. Build the backend:

```bash
go build
```

This will create a binary that includes the embedded frontend.

## Configuration

The frontend is configured to connect to the backend API at `http://localhost:8080/api` by default. You can change this by setting the `NEXT_PUBLIC_API_URL` environment variable.
