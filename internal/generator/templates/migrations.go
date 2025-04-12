// internal/generator/templates/migrations.go - Templates for migration files
package templates

// MigrationsScriptTemplate returns the content of the migrations.sh script
func MigrationsScriptTemplate() string {
	return `#!/bin/sh
# scripts/migrate.sh - Database migrations runner

# Change to project root directory
cd "$(dirname "$0")/.." || exit 1

# Parse arguments
COMMAND="up"
STEPS=0
ENV_FILE=".env"

print_usage() {
  echo "Usage: $0 [options]"
  echo "Options:"
  echo "  -c, --command=COMMAND  Migration command (up, down, version) [default: up]"
  echo "  -s, --steps=STEPS      Number of migrations to apply (0 means all) [default: 0]"
  echo "  -e, --env=ENV_FILE     Path to .env file [default: .env]"
  echo "  -h, --help             Show this help message"
}

while [ $# -gt 0 ]; do
  case "$1" in
    -c=*|--command=*)
      COMMAND="${1#*=}"
      shift
      ;;
    -s=*|--steps=*)
      STEPS="${1#*=}"
      shift
      ;;
    -e=*|--env=*)
      ENV_FILE="${1#*=}"
      shift
      ;;
    -h|--help)
      print_usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      print_usage
      exit 1
      ;;
  esac
done

# Run migrations tool
go run ./scripts/migtool/migrations.go -command="$COMMAND" -steps="$STEPS" -env="$ENV_FILE"
`
}

// MigrationToolTemplate returns the content of the migrations tool
func MigrationToolTemplate() string {
	return `// scripts/migtool/migrations.go - Database migrations tool
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/joho/godotenv"

	"{{ .ModuleName }}/internal/migrations"
)

func main() {
	// Define flags
	var (
		command = flag.String("command", "up", "Migration command (up, down, version)")
		steps   = flag.Int("steps", 0, "Number of migrations to apply (0 means all)")
		env     = flag.String("env", ".env", "Path to .env file")
	)

	flag.Parse()

	// Load environment variables from .env file
	if err := godotenv.Load(*env); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Get database connection string from environment
	connString := os.Getenv("DB_CONNECTION_STRING")
	if connString == "" {
		fmt.Println("Error: DB_CONNECTION_STRING environment variable is not set")
		os.Exit(1)
	}

	// Get migrations directory from environment or use default
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		// Use embedded migrations
		if err := runEmbeddedMigrations(connString, *command, *steps); err != nil {
			if err != migrate.ErrNoChange {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}
	} else {
		// Use file-based migrations
		sourceURL := fmt.Sprintf("file://%s", filepath.Clean(migrationsDir))
		
		// Create migrate instance
		m, err := migrate.New(sourceURL, connString)
		if err != nil {
			fmt.Printf("Error: Failed to create migrate instance: %v\n", err)
			os.Exit(1)
		}
		
		// Set logger
		m.Log = &migrationLogger{}
		
		// Execute migration command
		if err := executeMigrationCommand(m, *command, *steps); err != nil {
			if err != migrate.ErrNoChange {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}
	}
}

// runEmbeddedMigrations runs migrations from embedded filesystem
func runEmbeddedMigrations(connString, command string, steps int) error {
	// Create migrations source
	migrations, err := migrations.GetFS()
	if err != nil {
		return fmt.Errorf("failed to access embedded migrations: %w", err)
	}

	// Create source instance
	d, err := iofs.New(migrations, ".")
	if err != nil {
		return fmt.Errorf("failed to create migrations source: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithSourceInstance("iofs", d, connString)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Set logger
	m.Log = &migrationLogger{}

	// Execute migration command
	return executeMigrationCommand(m, command, steps)
}

// executeMigrationCommand executes the migration command
func executeMigrationCommand(m *migrate.Migrate, command string, steps int) error {
	switch strings.ToLower(command) {
	case "up":
		if steps > 0 {
			err := m.Steps(steps)
			if err != nil && err != migrate.ErrNoChange {
				return fmt.Errorf("failed to apply migrations: %w", err)
			}
			fmt.Printf("Successfully applied %d migrations\n", steps)
		} else {
			err := m.Up()
			if err != nil && err != migrate.ErrNoChange {
				return fmt.Errorf("failed to apply all migrations: %w", err)
			}
			fmt.Println("Successfully applied all migrations")
		}

	case "down":
		if steps > 0 {
			err := m.Steps(-steps)
			if err != nil && err != migrate.ErrNoChange {
				return fmt.Errorf("failed to rollback migrations: %w", err)
			}
			fmt.Printf("Successfully rolled back %d migrations\n", steps)
		} else {
			err := m.Down()
			if err != nil && err != migrate.ErrNoChange {
				return fmt.Errorf("failed to rollback all migrations: %w", err)
			}
			fmt.Println("Successfully rolled back all migrations")
		}

	case "version":
		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			return fmt.Errorf("failed to get migration version: %w", err)
		}
		fmt.Printf("Current migration version: %d (dirty: %v)\n", version, dirty)

	default:
		return fmt.Errorf("unknown command: %s", command)
	}

	return nil
}

// Custom logger for migrations
type migrationLogger struct{}

func (l *migrationLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *migrationLogger) Verbose() bool {
	return true
}
`
}

// MigrationsPackageTemplate returns the content of the migrations package
func MigrationsPackageTemplate() string {
	return `// internal/migrations/migrations.go - Embedded SQL migrations
package migrations

import (
	"embed"
	"io/fs"
)

//go:embed sql/*.sql
var Migrations embed.FS

// GetFS returns a filesystem with SQL migrations
func GetFS() (fs.FS, error) {
	return fs.Sub(Migrations, "sql")
}
`
}

// ModelGeneratorScriptTemplate returns the content of the model generator script
func ModelGeneratorScriptTemplate() string {
	return `#!/bin/sh
# scripts/generate_models.sh - Database models generator from SQL migrations

# Change to project root directory
cd "$(dirname "$0")/.." || exit 1

# Parse arguments
ENV_FILE=".env"
OUTPUT_DIR="internal/db/models"
MIGRATIONS_DIR="internal/migrations/sql"
FROM_MIGRATIONS=true

print_usage() {
  echo "Usage: $0 [options]"
  echo "Options:"
  echo "  -e, --env=ENV_FILE     Path to .env file [default: .env]"
  echo "  -o, --output=DIR       Output directory for models [default: internal/db/models]"
  echo "  -m, --migrations=DIR   Directory with migration files [default: internal/migrations/sql]"
  echo "  -d, --from-db          Generate models from database instead of migrations"
  echo "  -h, --help             Show this help message"
}

while [ $# -gt 0 ]; do
  case "$1" in
    -e=*|--env=*)
      ENV_FILE="${1#*=}"
      shift
      ;;
    -o=*|--output=*)
      OUTPUT_DIR="${1#*=}"
      shift
      ;;
    -m=*|--migrations=*)
      MIGRATIONS_DIR="${1#*=}"
      shift
      ;;
    -d|--from-db)
      FROM_MIGRATIONS=false
      shift
      ;;
    -h|--help)
      print_usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      print_usage
      exit 1
      ;;
  esac
done

# Check if .env file exists
if [ ! -f "$ENV_FILE" ]; then
  echo "Warning: $ENV_FILE file not found"
fi

# Run model generator tool
echo "Generating models from migrations..."
if [ "$FROM_MIGRATIONS" = true ]; then
  go run ./scripts/modelgen/modelgen.go -env="$ENV_FILE" -output="$OUTPUT_DIR" -migrations="$MIGRATIONS_DIR" -from-migrations=true
else
  go run ./scripts/modelgen/modelgen.go -env="$ENV_FILE" -output="$OUTPUT_DIR" -from-migrations=false
fi
`
}

// MigrationFileTemplate returns the content of the initial migration file
func MigrationFileTemplate() string {
	return `-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
`
}

// MigrationDownFileTemplate returns the content of the initial down migration file
func MigrationDownFileTemplate() string {
	return `-- Drop indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;

-- Drop tables
DROP TABLE IF EXISTS users;
`
}
