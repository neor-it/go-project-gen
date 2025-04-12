// internal/config/config.go - Configuration structures for the project generator
package config

// Config represents the main configuration for the generator
type Config struct {
	// Is the generator running in interactive mode
	IsInteractive bool
	// Output directory for the generated project
	OutputDir string
	// Configuration for the project to be generated
	ProjectConfig ProjectConfig
}

// ProjectConfig represents the configuration for the project to be generated
type ProjectConfig struct {
	// Username (for module name, e.g., github.com/username/projectname)
	Username string
	// Project name
	ProjectName string
	// Module name (e.g., github.com/username/projectname)
	ModuleName string
	// Components to include in the project
	Components Components
}

// Components represents the components to include in the project
type Components struct {
	// Include HTTP server with Gin
	HTTP bool
	// Include PostgreSQL database
	Postgres bool
	// Include Docker support
	Docker bool
	// Include Kubernetes manifests
	Kubernetes bool
	// Include CI/CD configuration
	CICD bool
}

// ParseArgs parses command line arguments
func ParseArgs(args []string) (*Config, error) {
	// Default configuration with interactive mode
	cfg := &Config{
		IsInteractive: true,
		OutputDir:     ".",
	}

	// TODO: Add argument parsing logic if needed
	// For now, just return the default configuration

	return cfg, nil
}
