# Dockerfile - Updated for better volume handling
# Use latest Go version for build environment
FROM golang:1.23-alpine

# Install necessary build tools
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o goprojectgen .

# Create and set permissions for /output directory
RUN mkdir -p /output && chmod 777 /output

# Set the entrypoint script
ENTRYPOINT ["/app/goprojectgen"]

# Default command if none is provided
CMD []