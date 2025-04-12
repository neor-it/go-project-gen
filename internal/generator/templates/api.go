// internal/generator/templates/api.go - Templates for API files
package templates

// APIServerTemplate returns the content of the server.go file
func APIServerTemplate() string {
	return `// internal/api/server.go - HTTP server implementation
package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"

	"{{ .ModuleName }}/internal/api/middleware"
	"{{ .ModuleName }}/internal/api/routes"
	"{{ .ModuleName }}/internal/config"
	"{{ .ModuleName }}/internal/logger"
)

// Server represents the HTTP server
type Server struct {
	log    logger.Logger
	cfg    *config.Config
	router *gin.Engine
	server *http.Server
}

// NewServer creates a new HTTP server
func NewServer(log logger.Logger, cfg *config.Config, dependencies ...interface{}) (*Server, error) {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger(log))
	router.Use(middleware.Recovery(log))
	router.Use(cors.Default())

	// Add pprof endpoints in debug mode
	pprof.Register(router)

	// Register routes
	routes.RegisterRoutes(router, log, dependencies...)

	// Create server
	server := &Server{
		log:    log,
		cfg:    cfg,
		router: router,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
			Handler:      router,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
		},
	}

	return server, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.log.Info("Starting HTTP server", "port", s.cfg.Server.Port)

	// Start server in a goroutine
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("Failed to start HTTP server", "error", err)
		}
	}()

	return nil
}

// Stop stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("Stopping HTTP server")

	// Shutdown server
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	return nil
}
`
}

// APIHandlersTemplate returns the content of the handlers.go file
func APIHandlersTemplate() string {
	return `// internal/api/handlers/handlers.go - HTTP request handlers
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"{{ .ModuleName }}/internal/logger"
)

// Handler represents a HTTP handler
type Handler struct {
	log logger.Logger
}

// NewHandler creates a new handler
func NewHandler(log logger.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

// HealthCheck handles the health check endpoint
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// Status handles the status endpoint
func (h *Handler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"version": "1.0.0",
	})
}
`
}

// APIMiddlewareTemplate returns the content of the middleware.go file
func APIMiddlewareTemplate() string {
	return `// internal/api/middleware/middleware.go - HTTP middleware
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"{{ .ModuleName }}/internal/logger"
)

// Logger returns a middleware that logs HTTP requests
func Logger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Info("HTTP request",
			"status", statusCode,
			"method", method,
			"path", path,
			"ip", clientIP,
			"latency", latency,
			"user_agent", c.Request.UserAgent(),
		)
	}
}

// Recovery returns a middleware that recovers from panics
func Recovery(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log error
				log.Error("Panic recovered", "error", err)

				// Return error response
				c.AbortWithStatus(500)
			}
		}()

		c.Next()
	}
}

// RequestID returns a middleware that adds a request ID to the context
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add request ID to context
		c.Next()
	}
}
`
}

// APIRoutesTemplate returns the content of the routes.go file
func APIRoutesTemplate() string {
	return `// internal/api/routes/routes.go - HTTP routes
package routes

import (
	"github.com/gin-gonic/gin"

	"{{ .ModuleName }}/internal/api/handlers"
	"{{ .ModuleName }}/internal/logger"
)

// RegisterRoutes registers the HTTP routes
func RegisterRoutes(router *gin.Engine, log logger.Logger, dependencies ...interface{}) {
	// Create handlers
	handler := handlers.NewHandler(log)

	// Register top-level routes
	router.GET("/health", handler.HealthCheck)
	router.GET("/status", handler.Status)

	// Register API v1 routes with TODO
	v1 := router.Group("/api/v1")
	{
		// TODO: Add API v1 routes here
	}
}
`
}
