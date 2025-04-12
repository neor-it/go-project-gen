// internal/generator/templates/main.go - Templates for main files
package templates

import (
	"github.com/username/goprojectgen/internal/config"
)

// MainTemplate returns the content of the main.go file
func MainTemplate(cfg config.ProjectConfig) string {
	imports := `
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"` + cfg.ModuleName + `/internal/app"
	"` + cfg.ModuleName + `/internal/config"
	"` + cfg.ModuleName + `/internal/logger"
`

	return `// main.go - Main entry point for the ` + cfg.ProjectName + ` service
package main

import (` + imports + `)

func main() {
	// Create context that listens for termination signals
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Initialize logger
	log := logger.NewLogger()
	log.Info("Starting ` + cfg.ProjectName + ` service")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", "error", err)
	}

	// Create and start application
	application, err := app.NewApp(log, cfg)
	if err != nil {
		log.Fatal("Failed to create application", "error", err)
	}

	// Start the application
	if err := application.Start(ctx); err != nil {
		log.Fatal("Failed to start application", "error", err)
	}

	// Wait for termination signal
	<-ctx.Done()
	log.Info("Shutting down...")

	// Create a new context for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	// Stop the application
	if err := application.Stop(shutdownCtx); err != nil {
		log.Error("Error during shutdown", "error", err)
	}

	log.Info("Service stopped")
}
`
}

// GoModTemplate returns the content of the go.mod file
func GoModTemplate(moduleName string) string {
	return `module ` + moduleName + `

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	go.uber.org/zap v1.26.0
)

require (
	github.com/bytedance/sonic v1.10.2 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.17.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.6 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.1.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.7.0 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
google.golang.org/protobuf/proto v1.32.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
`
}

// GitignoreTemplate returns the content of the .gitignore file
func GitignoreTemplate() string {
	return `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
*.db

# Test binary, built with "go test -c"
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
vendor/

# Go workspace file
go.work

# IDE specific files
.idea/
.vscode/
*.swp
*.swo

# Environment variables
.env
.env.local

# Build directory
/build/
/bin/

# Log files
*.log

# OS specific files
.DS_Store

# Temporary files
tmp/
temp/
`
}

// ReadmeTemplate returns the content of the README.md file
func ReadmeTemplate(cfg config.ProjectConfig) string {
	components := ""

	if cfg.Components.HTTP {
		components += "- HTTP API (Gin)\n"
	}
	if cfg.Components.Postgres {
		components += "- PostgreSQL database\n"
	}
	if cfg.Components.Docker {
		components += "- Docker support\n"
	}
	if cfg.Components.Kubernetes {
		components += "- Kubernetes deployment\n"
	}
	if cfg.Components.CICD {
		components += "- CI/CD pipeline\n"
	}

	return `# ` + cfg.ProjectName + `

## Overview

This is a Go service generated with Go Project Generator.

## Components

` + components + `

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

1. Clone the repository:

   ` + "```bash" + `
   git clone https://github.com/` + cfg.Username + `/` + cfg.ProjectName + `.git
   cd ` + cfg.ProjectName + `
   ` + "```" + `

2. Install dependencies:

   ` + "```bash" + `
   go mod download
   ` + "```" + `

3. Build the application:

   ` + "```bash" + `
   go build -o bin/` + cfg.ProjectName + ` main.go
   ` + "```" + `

4. Run the application:

   ` + "```bash" + `
   ./bin/` + cfg.ProjectName + `
   ` + "```" + `

## Configuration

The application is configured using environment variables or a configuration file.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
`
}

// ConfigTemplate returns the content of the config.go file
func ConfigTemplate() string {
	return `// internal/config/config.go - Configuration loading and parsing
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	// Server configuration
	Server struct {
		Port         int           ` + "`mapstructure:\"port\" yaml:\"port\"`" + `
		ReadTimeout  time.Duration ` + "`mapstructure:\"read_timeout\" yaml:\"read_timeout\"`" + `
		WriteTimeout time.Duration ` + "`mapstructure:\"write_timeout\" yaml:\"write_timeout\"`" + `
	} ` + "`mapstructure:\"server\" yaml:\"server\"`" + `

	// Database configuration
	Database struct {
		Host     string ` + "`mapstructure:\"host\" yaml:\"host\"`" + `
		Port     int    ` + "`mapstructure:\"port\" yaml:\"port\"`" + `
		User     string ` + "`mapstructure:\"user\" yaml:\"user\"`" + `
		Password string ` + "`mapstructure:\"password\" yaml:\"password\"`" + `
		Name     string ` + "`mapstructure:\"name\" yaml:\"name\"`" + `
		SSLMode  string ` + "`mapstructure:\"ssl_mode\" yaml:\"ssl_mode\"`" + `
	} ` + "`mapstructure:\"database\" yaml:\"database\"`" + `

	// Logging configuration
	Logging struct {
		Level  string ` + "`mapstructure:\"level\" yaml:\"level\"`" + `
		Format string ` + "`mapstructure:\"format\" yaml:\"format\"`" + `
	} ` + "`mapstructure:\"logging\" yaml:\"logging\"`" + `

	// Shutdown timeout
	ShutdownTimeout time.Duration ` + "`mapstructure:\"shutdown_timeout\" yaml:\"shutdown_timeout\"`" + `
}

// LoadConfig loads the configuration from environment variables or a file
func LoadConfig() (*Config, error) {
	var config Config

	// Set default values
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "10s")
	viper.SetDefault("server.write_timeout", "10s")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.name", "postgres")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("shutdown_timeout", "5s")

	// Read from environment variables
	viper.AutomaticEnv()

	// Look for config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/app")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// ConnectionString returns the database connection string
func (c *Config) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}
`
}

// LoggerTemplate returns the content of the logger.go file
func LoggerTemplate() string {
	return `// internal/logger/logger.go - Logger implementation
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger interface defines the methods that the logger should implement
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	With(keysAndValues ...interface{}) Logger
}

// ZapLogger implements the Logger interface using Zap
type ZapLogger struct {
	logger *zap.SugaredLogger
}

// NewLogger creates a new logger
func NewLogger() Logger {
	// Create encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create JSON encoder
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// Create core
	core := zapcore.NewCore(
		jsonEncoder,
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zapcore.InfoLevel),
	)

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	defer logger.Sync()

	// Return sugared logger
	return &ZapLogger{
		logger: logger.Sugar(),
	}
}

// Debug logs a debug message
func (l *ZapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

// Info logs an info message
func (l *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

// Warn logs a warning message
func (l *ZapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

// Error logs an error message
func (l *ZapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

// Fatal logs a fatal message and exits
func (l *ZapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}

// With returns a logger with the given key-value pairs
func (l *ZapLogger) With(keysAndValues ...interface{}) Logger {
	return &ZapLogger{
		logger: l.logger.With(keysAndValues...),
	}
}
`
}

// AppTemplate returns the content of the app.go file
func AppTemplate(cfg config.ProjectConfig) string {
	imports := `
	"context"

	"` + cfg.ModuleName + `/internal/config"
	"` + cfg.ModuleName + `/internal/logger"
`

	// Add HTTP import
	if cfg.Components.HTTP {
		imports += `	"` + cfg.ModuleName + `/internal/api"
`
	}

	// Add DB import
	if cfg.Components.Postgres {
		imports += `	"` + cfg.ModuleName + `/internal/db"
`
	}

	// App struct
	appStruct := `
// App represents the application
type App struct {
	log logger.Logger
	cfg *config.Config
`

	// Add HTTP field
	if cfg.Components.HTTP {
		appStruct += `	server *api.Server
`
	}

	// Add DB field
	if cfg.Components.Postgres {
		appStruct += `	db *db.Database
`
	}

	appStruct += `}
`

	// NewApp function
	newApp := `
// NewApp creates a new application
func NewApp(log logger.Logger, cfg *config.Config) (*App, error) {
	app := &App{
		log: log,
		cfg: cfg,
	}

`

	// Add DB initialization
	if cfg.Components.Postgres {
		newApp += `	// Initialize database
	db, err := db.NewDatabase(log, cfg.ConnectionString())
	if err != nil {
		return nil, err
	}
	app.db = db

`
	}

	// Add HTTP initialization
	if cfg.Components.HTTP {
		newApp += `	// Initialize HTTP server
	server, err := api.NewServer(log, cfg`

		if cfg.Components.Postgres {
			newApp += `, db`
		}

		newApp += `)
	if err != nil {
		return nil, err
	}
	app.server = server

`
	}

	newApp += `	return app, nil
}
`

	// Start function
	start := `
// Start starts the application
func (a *App) Start(ctx context.Context) error {
	a.log.Info("Starting application")

`

	// Add DB start
	if cfg.Components.Postgres {
		start += `	// Start database
	if err := a.db.Connect(); err != nil {
		return err
	}

`
	}

	// Add HTTP start
	if cfg.Components.HTTP {
		start += `	// Start HTTP server
	if err := a.server.Start(); err != nil {
		return err
	}

`
	}

	start += `	return nil
}
`

	// Stop function
	stop := `
// Stop stops the application
func (a *App) Stop(ctx context.Context) error {
	a.log.Info("Stopping application")

`

	// Add HTTP stop
	if cfg.Components.HTTP {
		stop += `	// Stop HTTP server
	if err := a.server.Stop(ctx); err != nil {
		return err
	}

`
	}

	// Add DB stop
	if cfg.Components.Postgres {
		stop += `	// Close database connection
	if err := a.db.Close(); err != nil {
		return err
	}

`
	}

	stop += `	return nil
}
`

	return `// internal/app/app.go - Application initialization and lifecycle management
package app

import (` + imports + `)
` + appStruct + newApp + start + stop
}
