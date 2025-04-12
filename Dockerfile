# Build stage
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
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/goprojectgen

# Final stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/bin/goprojectgen .

# Set environment variables
ENV TZ=UTC

# Expose port
EXPOSE 8080

# Run application
ENTRYPOINT ["/app/goprojectgen"]