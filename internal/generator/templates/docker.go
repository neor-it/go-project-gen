// internal/generator/templates/docker.go - Templates for Docker files
package templates

import "github.com/username/goprojectgen/internal/config"

// DockerfileTemplate returns the content of the Dockerfile
func DockerfileTemplate(cfg config.ProjectConfig) string {
	return `# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/` + cfg.ProjectName + ` main.go

# Final stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/bin/` + cfg.ProjectName + ` .

# Copy configuration
COPY --from=builder /app/config/config.yaml /app/config/

# Set environment variables
ENV TZ=UTC

# Expose port
EXPOSE 8080

# Run application
CMD ["./` + cfg.ProjectName + `"]
`
}

// DockerComposeTemplate returns the content of the docker-compose.yml file
func DockerComposeTemplate(cfg config.ProjectConfig) string {
	// Base docker-compose.yml
	compose := `version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: ` + cfg.ProjectName + `
    restart: unless-stopped
    environment:
      - TZ=UTC
    ports:
      - "8080:8080"
`

	// Add Postgres service if needed
	if cfg.Components.Postgres {
		compose += `
  postgres:
    image: postgres:16-alpine
    container_name: ` + cfg.ProjectName + `-postgres
    restart: unless-stopped
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=` + cfg.ProjectName + `
      - TZ=UTC
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
`
	}

	return compose
}

// DockerignoreTemplate returns the content of the .dockerignore file
func DockerignoreTemplate() string {
	return `# Git
.git
.gitignore

# Docker
Dockerfile
docker-compose.yml
.dockerignore

# IDE
.idea
.vscode

# Binaries
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Tests
*_test.go
*.test

# Build
.build/

# Misc
*.md
LICENSE
README.md

# Environment variables
.env
.env.local

# Temporary files
tmp/
temp/

# Log files
*.log

# OS specific files
.DS_Store
`
}
