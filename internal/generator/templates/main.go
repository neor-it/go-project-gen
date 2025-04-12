// internal/generator/templates/main.go - Templates for main files
package templates

import (
	"github.com/neor-it/go-project-gen/internal/config"
)

// MainTemplate returns the content of the main.go file
func MainTemplate(cfg config.ProjectConfig) string {
	imports := `
	"context"
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
	
	// Set log level from configuration
	log.SetLevel(cfg.GetLogLevel())

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

go 1.23

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/gin-contrib/cors v1.5.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/golang-migrate/migrate/v4 v4.17.0
	github.com/spf13/viper v1.18.2
	go.uber.org/zap v1.26.0
	github.com/gertd/go-pluralize v0.2.1
	github.com/iancoleman/strcase v0.3.0
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
	if cfg.Components.CICD {
		components += "- CI/CD pipeline\n"
	}

	migrationsSection := ""
	modelsSection := ""

	if cfg.Components.Postgres {
		migrationsSection = `## Database Migrations

This project uses Go-based migrations with [golang-migrate](https://github.com/golang-migrate/migrate). Migration files are stored in the 'internal/migrations/sql' directory using the format 'NNN_description.(up|down).sql'.

### Running Migrations

To apply migrations:

` + "```bash" + `
# Apply all pending migrations
./scripts/migrate.sh

# Apply specific number of migrations
./scripts/migrate.sh --steps=1

# Rollback migrations
./scripts/migrate.sh --command=down

# Check current migration version
./scripts/migrate.sh --command=version
` + "```" + `

### Creating New Migrations

To create a new migration:

1. Create two new files in 'internal/migrations/sql/' using sequential numbering:
   - 'NNN_description.up.sql' - Forward migration
   - 'NNN_description.down.sql' - Rollback migration

Example:

` + "```sql" + `
-- 002_add_posts_table.up.sql
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
` + "```" + `

` + "```sql" + `
-- 002_add_posts_table.down.sql
DROP TABLE IF EXISTS posts;
` + "```" + `

`

		modelsSection = `## Database Models

This project can automatically generate Go struct models from your database schema.

### Generating Models

After applying migrations, you can generate models with:

` + "```bash" + `
# Generate all models
./scripts/generate_models.sh

# Specify output directory
./scripts/generate_models.sh --output=internal/custom/models
` + "```" + `

The generator creates type-safe Go structs with appropriate field types and struct tags for database models.

Models will be placed in 'internal/db/models/' by default.

`
	}

	postgresPrereq := ""
	if cfg.Components.Postgres {
		postgresPrereq = "- PostgreSQL"
	}

	postgresSetup := ""
	if cfg.Components.Postgres {
		postgresSetup = `
4. Set up database:

   ` + "```bash" + `
   # Run PostgreSQL (if using Docker)
   docker-compose up -d postgres

   # Apply database migrations
   ./scripts/migrate.sh
   ` + "```" + `
`
	}

	apiSection := ""
	if cfg.Components.HTTP {
		apiSection = `│   ├── api/             # HTTP API implementation
│   │   ├── handlers/    # HTTP request handlers
│   │   ├── middleware/  # HTTP middleware
│   │   └── routes/      # HTTP route definitions`
	}

	dbSection := ""
	if cfg.Components.Postgres {
		dbSection = `│   ├── db/              # Database code
│   │   ├── models/      # Database models
│   │   └── repositories/ # Data access layer
│   ├── migrations/      # Database migrations
│   │   └── sql/         # SQL migration files`
	}

	scriptsSection := ""
	if cfg.Components.Postgres {
		scriptsSection = `│   ├── migrate.sh       # Database migration script
│   ├── generate_models.sh # Model generation script
│   ├── migtool/         # Migration tool implementation
│   └── modelgen/        # Model generator implementation`
	}

	dockerSection := ""
	if cfg.Components.Docker {
		dockerSection = `├── Dockerfile           # Docker build file
├── docker-compose.yml   # Docker Compose file`
	}

	// Add Docker Compose section for running app with Docker
	dockerComposeSection := ""
	if cfg.Components.Docker {
		dockerComposeSection = `## Running with Docker Compose

This project includes Docker support for easy deployment and development.

### Prerequisites

- Docker
- Docker Compose

### Setup

1. Make sure Docker is installed and running on your system.

2. Update your .env file for Docker environment:

   ` + "```bash" + `
   # Create a copy of the example file if you haven't done so
   cp .env.example .env
   ` + "```" + `

3. Important settings for Docker environment in .env:
`
		// Add database specific settings if Postgres is included
		if cfg.Components.Postgres {
			dockerComposeSection += `
   ` + "```bash" + `
   # Use the service name as the hostname (not localhost):
   DB_CONNECTION_STRING=postgres://postgres:postgres@postgres:5432/` + cfg.ProjectName + `?sslmode=disable
   ` + "```" + `
`
		}

		dockerComposeSection += `
4. Build and start the containers:

   ` + "```bash" + `
   # Start all services in detached mode
   docker-compose up -d
   
   # View logs from all containers
   docker-compose logs -f
   
   # Stop all services
   docker-compose down
   ` + "```" + `
`
		// Add database migration info if Postgres is included
		if cfg.Components.Postgres {
			dockerComposeSection += `
5. Run database migrations from within the container:

   ` + "```bash" + `
   # Connect to the application container
   docker-compose exec app sh
   
   # Inside the container, run migrations
   ./scripts/migrate.sh
   ` + "```" + `
`
		}

		dockerComposeSection += `
### Rebuilding and Updating

When you make changes to the application:

` + "```bash" + `
# Rebuild the application container
docker-compose build app

# Restart with the updated image
docker-compose up -d
` + "```" + `

### Accessing the Application

Once the containers are running:

- The HTTP API will be available at: http://localhost:8080
`
		if cfg.Components.Postgres {
			dockerComposeSection += `- PostgreSQL will be available at: localhost:5432
`
		}
	}

	return `# ` + cfg.ProjectName + `

## Overview

This is a Go service generated with Go Project Generator.

## Components

` + components + `

## Getting Started

### Prerequisites

- Go 1.23 or higher
- Git
` + postgresPrereq + `

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

3. Set up environment variables:

   ` + "```bash" + `
   cp .env.example .env
   # Edit .env file with your configuration
   ` + "```" + `
` + postgresSetup + `
5. Build the application:

   ` + "```bash" + `
   go build -o bin/` + cfg.ProjectName + ` main.go
   ` + "```" + `

6. Run the application:

   ` + "```bash" + `
   ./bin/` + cfg.ProjectName + `
   ` + "```" + `

` + dockerComposeSection + `
## Project Structure

` + "```" + `
├── internal/            # Private application code
│   ├── app/             # Application initialization
│   ├── config/          # Configuration handling
│   ├── logger/          # Logging implementation
` + apiSection + `
` + dbSection + `
├── pkg/                 # Public libraries
├── scripts/             # Utility scripts
` + scriptsSection + `
├── main.go              # Application entry point
├── go.mod               # Go module file
├── go.sum               # Go module checksums
` + dockerSection + `
├── .env.example         # Example environment file
├── .env                 # Environment file (git-ignored)
└── README.md            # This file
` + "```" + `

## Configuration

The application is configured using environment variables in the .env file.

` + migrationsSection + modelsSection + `
## License

This project is licensed under the MIT License - see the LICENSE file for details.
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
