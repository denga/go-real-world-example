# Go Real World Example

This project is a real-world example of a RESTful API implementation in Go, following the [RealWorld](https://github.com/gothinkster/realworld) API specification. It demonstrates how to build a backend application using modern Go practices and libraries, along with a modern frontend using Next.js and chadcn UI.

## Features

### Backend
- RESTful API implementation based on the RealWorld API specification
- Built with the [Chi](https://github.com/go-chi/chi) router for HTTP routing
- API code generation using [oapi-codegen](https://github.com/deepmap/oapi-codegen)
- JWT-based authentication system
- In-memory database for data storage
- OpenAPI specification embedded in the binary
- Password hashing with bcrypt
- Middleware for request authentication

### Frontend
- Built with [Next.js](https://nextjs.org/) and [chadcn UI](https://ui.shadcn.com/)
- Dark mode support
- Mobile responsive design
- Integration with backend API
- Static export that can be embedded in the backend binary

## Project Structure

```
.
├── api/                  # API-related code
│   ├── api.gen.go        # Generated API code
│   ├── config.yaml       # Configuration for oapi-codegen
│   └── gen.go            # Go generate directive
├── frontend/             # Frontend application
│   ├── app/              # Next.js app directory
│   ├── components/       # React components
│   │   ├── theme-provider.tsx  # Dark mode provider
│   │   ├── theme-toggle.tsx    # Dark mode toggle
│   │   └── ui/           # UI components from shadcn/ui
│   ├── lib/              # Utility functions and API client
│   ├── public/           # Static assets
│   ├── next.config.ts    # Next.js configuration
│   └── package.json      # Frontend dependencies
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

### Backend
- Go 1.24 or later

### Frontend
- Node.js 18.0.0 or later
- npm 9.0.0 or later

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

### Continuous Integration

This project uses Gitea Actions for continuous integration. The CI workflow automatically runs tests on push to the main branch and on pull requests.

You can see the CI workflow configuration in the `.gitea/workflows/go-test.yml` file.

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

### Frontend Development

The frontend is built with Next.js and chadcn UI. To work on the frontend:

1. Navigate to the frontend directory:
   ```
   cd frontend
   ```

2. Install dependencies:
   ```
   npm install
   ```

3. Start the development server:
   ```
   npm run dev
   ```

4. Open [http://localhost:3000](http://localhost:3000) in your browser to see the frontend.

### Building and Embedding the Frontend

To build the frontend and embed it in the backend binary:

1. Build the frontend:
   ```
   cd frontend
   npm run build
   ```
   This will create a static export in the `frontend/dist` directory.

2. Update the main.go file to embed and serve the frontend (see the frontend/README.md for detailed instructions).

3. Build the backend:
   ```
   go build
   ```

### Regenerating API Code

If you make changes to the OpenAPI specification, you can regenerate the API code using:
```
go generate ./...
```

## License

This project is open source and available under the [MIT License](LICENSE).
