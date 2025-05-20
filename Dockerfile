# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build tools if necessary (e.g., for CGO_ENABLED=0 builds or specific dependencies)
# RUN apk add --no-cache gcc musl-dev

# Copy go.mod and go.sum files to download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the application
# Disabling CGO for a smaller, static binary if not needed, adjust as necessary
# Ensure your main package is in cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/billing-mcp-server ./cmd/main.go

# Stage 2: Create the final lightweight image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/billing-mcp-server /app/billing-mcp-server

# Copy migration files
# Ensure this path matches your project structure
COPY ./database/migrations ./database/migrations

# .config.yaml will be mounted via docker-compose, so no need to copy it here
# If you prefer to bake it in, uncomment the next line:
# COPY .config.yaml .config.yaml

# Expose the port the application listens on (should match config.yaml and docker-compose)
EXPOSE 8080

# Command to run the application
# The application will read CONFIG_PATH from environment variable set in docker-compose
ENTRYPOINT ["/app/billing-mcp-server"]
