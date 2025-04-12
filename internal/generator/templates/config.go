// internal/generator/templates/config.go - Templates for configuration files
package templates

import (
	"github.com/neor-it/go-project-gen/internal/config"
)

// ConfigTemplate returns the content of the config.go file
func ConfigTemplate(projectCfg config.ProjectConfig) string {
	baseConfig := `// internal/config/config.go - Configuration loading and parsing
package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	// Server configuration
	Server struct {
		Port         int           ` + "`mapstructure:\"port\"`" + `
		ReadTimeout  time.Duration ` + "`mapstructure:\"read_timeout\"`" + `
		WriteTimeout time.Duration ` + "`mapstructure:\"write_timeout\"`" + `
	} ` + "`mapstructure:\"server\"`" + `

`

	// Add Database configuration if Postgres is enabled
	if projectCfg.Components.Postgres {
		baseConfig += `	// Database configuration
	Database struct {
		ConnectionString string ` + "`mapstructure:\"connection_string\"`" + `
	} ` + "`mapstructure:\"database\"`" + `

`
	}

	// Add Logging configuration
	baseConfig += `	// Logging configuration
	Logging struct {
		Level string ` + "`mapstructure:\"level\"`" + `
	} ` + "`mapstructure:\"logging\"`" + `

	// Shutdown timeout
	ShutdownTimeout time.Duration ` + "`mapstructure:\"shutdown_timeout\"`" + `
}

// LoadConfig loads the configuration from environment variables or .env file
func LoadConfig() (*Config, error) {
	var config Config

	// Load .env file if it exists
	_ = godotenv.Load()

	// Set default values and override with environment variables
	
	// Server configuration
	config.Server.Port = getEnvInt("SERVER_PORT", 8080)
	config.Server.ReadTimeout = getEnvDuration("SERVER_READ_TIMEOUT", 10*time.Second)
	config.Server.WriteTimeout = getEnvDuration("SERVER_WRITE_TIMEOUT", 10*time.Second)
	
`

	// Add Database configuration loading if Postgres is enabled
	if projectCfg.Components.Postgres {
		baseConfig += `	// Database configuration
	config.Database.ConnectionString = getEnvString("DB_CONNECTION_STRING", "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable")
	
`
	}

	baseConfig += `	// Logging configuration
	config.Logging.Level = getEnvString("LOGGING_LEVEL", "info")
	
	// Shutdown timeout
	config.ShutdownTimeout = getEnvDuration("SHUTDOWN_TIMEOUT", 5*time.Second)

	return &config, nil
}
`

	// Add ConnectionString method if Postgres is enabled
	if projectCfg.Components.Postgres {
		baseConfig += `
// ConnectionString returns the database connection string
func (c *Config) ConnectionString() string {
	return c.Database.ConnectionString
}
`
	}

	baseConfig += `
// GetLogLevel returns the configured log level
func (c *Config) GetLogLevel() string {
	return c.Logging.Level
}

// getEnvString gets a string value from environment variable or returns the default
func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt gets an integer value from environment variable or returns the default
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvDuration gets a duration value from environment variable or returns the default
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
`

	return baseConfig
}
