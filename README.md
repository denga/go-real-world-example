# Go Real World Example

This project is a real-world example of a RESTful API implementation in Go, following the [RealWorld](https://github.com/gothinkster/realworld) API specification. It demonstrates how to build a backend application using modern Go practices and libraries.

## Features

- RESTful API implementation based on the RealWorld API specification
- Built with the [Chi](https://github.com/go-chi/chi) router for HTTP routing
- API code generation using [oapi-codegen](https://github.com/deepmap/oapi-codegen)
- JWT-based authentication system
- In-memory database for data storage
- OpenAPI specification embedded in the binary
- Password hashing with bcrypt
- Middleware for request authentication

## Project Structure

```
.
├── api/                  # API-related code
│   ├── api.gen.go        # Generated API code
│   ├── config.yaml       # Configuration for oapi-codegen
│   └── gen.go            # Go generate directive
├── internal/             # Internal application code
│   ├── auth/             # Authentication functionality
│   │   └── auth.go       # JWT token generation and validation
│   ├── db/               # Database implementation
│   │   └── db.go         # In-memory database
│   ├── handlers/         # API handlers
│   │   └── handlers.go   # Implementation of API endpoints
│   ├── middleware/       # HTTP middleware
│   │   └── auth.go       # Authentication middleware
│   └── util/             # Utility functions
│       └── slug.go       # Slug generation for articles
├── go.mod                # Go module file
├── go.sum                # Go module checksum
├── main.go               # Application entry point
├── openapi.yml           # OpenAPI specification
└── README.md             # This file
```

## Requirements

- Go 1.24 or later

## Getting Started

### Installation

1. Clone the repository:
   ```
   git clone git.homelab.lan/denga/go-real-world-example.git
   cd go-real-world-example
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

### Running the Application

Start the server:
```
go run main.go
```

The server will start on port 8080 by default. You can change the port by setting the `PORT` environment variable.

### Running Tests

To run all tests:
```
go test git.homelab.lan/denga/go-real-world-example/internal/...
```

To run tests for a specific package:
```
go test git.homelab.lan/denga/go-real-world-example/internal/util
go test git.homelab.lan/denga/go-real-world-example/internal/auth
go test git.homelab.lan/denga/go-real-world-example/internal/db
go test git.homelab.lan/denga/go-real-world-example/internal/middleware
go test git.homelab.lan/denga/go-real-world-example/internal/handlers
```

To run tests with verbose output:
```
go test -v git.homelab.lan/denga/go-real-world-example/internal/...
```

### API Documentation

The OpenAPI specification is available at `/openapi.yml` when the server is running. You can use tools like [Swagger UI](https://swagger.io/tools/swagger-ui/) to explore the API.

#### API Endpoints

The API implements the [RealWorld API specification](https://github.com/gothinkster/realworld/tree/main/api) with the following main endpoints:

- **Authentication**:
  - `POST /api/users/login` - Login for existing user
  - `POST /api/users` - Register a new user

- **User**:
  - `GET /api/user` - Get current user
  - `PUT /api/user` - Update user

- **Profiles**:
  - `GET /api/profiles/:username` - Get a profile
  - `POST /api/profiles/:username/follow` - Follow a user
  - `DELETE /api/profiles/:username/follow` - Unfollow a user

- **Articles**:
  - `GET /api/articles` - List articles
  - `GET /api/articles/feed` - Feed articles
  - `GET /api/articles/:slug` - Get an article
  - `POST /api/articles` - Create an article
  - `PUT /api/articles/:slug` - Update an article
  - `DELETE /api/articles/:slug` - Delete an article

- **Comments**:
  - `GET /api/articles/:slug/comments` - Get comments for an article
  - `POST /api/articles/:slug/comments` - Add a comment to an article
  - `DELETE /api/articles/:slug/comments/:id` - Delete a comment

- **Favorites**:
  - `POST /api/articles/:slug/favorite` - Favorite an article
  - `DELETE /api/articles/:slug/favorite` - Unfavorite an article

- **Tags**:
  - `GET /api/tags` - Get tags

#### Authentication

The API uses JWT tokens for authentication. To authenticate, you need to:

1. Register a user or login with existing credentials
2. Include the received token in the `Authorization` header of your requests:
   ```
   Authorization: Token <your-token>
   ```

## Development

### Regenerating API Code

If you make changes to the OpenAPI specification, you can regenerate the API code using:
```
go generate ./...
```

## License

This project is open source and available under the [MIT License](LICENSE).
