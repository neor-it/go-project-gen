// internal/generator/templates/templates.go - Centralized template interface
package templates

import (
	"github.com/neor-it/go-project-gen/internal/config"
)

// AllTemplates implements the interface for all available templates
type AllTemplates struct {
	Config    ConfigTemplates
	API       APITemplates
	DB        DBTemplates
	Migration MigrationTemplates
	Docker    DockerTemplates
	Main      MainTemplates
	Logger    LoggerTemplates
	CICD      CICDTemplates
}

// ConfigTemplates interface represents templates for configuration
type ConfigTemplates interface {
	ConfigTemplate(config.ProjectConfig) string
}

// APITemplates interface contains methods for generating API templates
type APITemplates interface {
	APIServerTemplate() string
	APIHandlersTemplate() string
	APIMiddlewareTemplate() string
	APIRoutesTemplate() string
}

// DBTemplates interface contains methods for generating database templates
type DBTemplates interface {
	DBTemplate() string
	DBModelsTemplate() string
	DBRepositoriesTemplate() string
}

// MigrationTemplates interface represents templates for migrations
type MigrationTemplates interface {
	MigrationsScriptTemplate() string
	MigrationToolTemplate() string
	MigrationsPackageTemplate() string
	ModelGeneratorScriptTemplate() string
	ModelGeneratorFullTemplate() string
	MigrationFileTemplate() string
	MigrationDownFileTemplate() string
}

// DockerTemplates represents templates for Docker
type DockerTemplates interface {
	DockerfileTemplate(config.ProjectConfig) string
	DockerComposeTemplate(config.ProjectConfig) string
	DockerignoreTemplate() string
}

// MainTemplates represents templates for main application files
type MainTemplates interface {
	MainTemplate(config.ProjectConfig) string
	GoModTemplate(string) string
	GitignoreTemplate() string
	ReadmeTemplate(config.ProjectConfig) string
	AppTemplate(config.ProjectConfig) string
}

// LoggerTemplates represents templates for logging
type LoggerTemplates interface {
	LoggerTemplate() string
}

// CICDTemplates represents templates for CI/CD
type CICDTemplates interface {
	GitHubWorkflowTemplate(config.ProjectConfig) string
}
