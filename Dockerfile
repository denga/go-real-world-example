FROM node:18-alpine AS frontend-builder
WORKDIR /app

# Copy frontend files
COPY frontend/ ./frontend/

# Install frontend dependencies and build
WORKDIR /app/frontend
RUN npm ci
RUN npm run build

# Go build stage
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app

# Install required packages
RUN apk add --no-cache git

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Copy the built frontend from the frontend-builder stage
COPY --from=frontend-builder /app/frontend/dist /app/frontend/dist

# Install codegen
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Build the Go application
RUN go generate ./api
RUN CGO_ENABLED=0 GOOS=linux go build -o realworld-app

# Final stage
FROM alpine:latest
WORKDIR /app

# Add ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from the build stage
COPY --from=backend-builder /app/realworld-app .

# Expose the application port
EXPOSE 8080

# Set environment variables
ENV PORT=8080

# Run the application
CMD ["./realworld-app"]