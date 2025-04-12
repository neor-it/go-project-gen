// internal/generator/generator.go - Project structure generator
package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/username/goprojectgen/internal/config"
	"github.com/username/goprojectgen/internal/generator/templates"
	"github.com/username/goprojectgen/internal/logger"
)

// Generator represents the project generator
type Generator struct {
	log    logger.Logger
	config *config.Config
}

// NewGenerator creates a new generator
func NewGenerator(log logger.Logger, cfg *config.Config) *Generator {
	return &Generator{
		log:    log,
		config: cfg,
	}
}

// Generate generates the project structure
func (g *Generator) Generate() error {
	g.log.Info("Generating project structure",
		"projectName", g.config.ProjectConfig.ProjectName,
		"moduleName", g.config.ProjectConfig.ModuleName,
	)

	// Create project directory
	projectDir := filepath.Join(g.config.OutputDir, g.config.ProjectConfig.ProjectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create standard Go project structure
	if err := g.createStandardStructure(projectDir); err != nil {
		return fmt.Errorf("failed to create standard structure: %w", err)
	}

	// Generate project-specific files
	if err := g.generateProjectFiles(projectDir); err != nil {
		return fmt.Errorf("failed to generate project files: %w", err)
	}

	// Generate component-specific files
	if err := g.generateComponentFiles(projectDir); err != nil {
		return fmt.Errorf("failed to generate component files: %w", err)
	}

	return nil
}

// createStandardStructure creates the standard Go project structure
func (g *Generator) createStandardStructure(projectDir string) error {
	// Create directories
	dirs := []string{
		"cmd",
		"internal",
		"internal/app",
		"internal/config",
		"internal/logger",
		"pkg",
		"scripts",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateProjectFiles generates the project-specific files
func (g *Generator) generateProjectFiles(projectDir string) error {
	g.log.Info("Generating project files")

	// Create go.mod file
	goModContent := templates.GoModTemplate(g.config.ProjectConfig.ModuleName)
	if err := os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("failed to create go.mod file: %w", err)
	}

	// Create main.go file
	mainContent := templates.MainTemplate(g.config.ProjectConfig)
	if err := os.WriteFile(filepath.Join(projectDir, "main.go"), []byte(mainContent), 0644); err != nil {
		return fmt.Errorf("failed to create main.go file: %w", err)
	}

	// Create .gitignore file
	gitignoreContent := templates.GitignoreTemplate()
	if err := os.WriteFile(filepath.Join(projectDir, ".gitignore"), []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore file: %w", err)
	}

	// Create README.md file
	readmeContent := templates.ReadmeTemplate(g.config.ProjectConfig)
	if err := os.WriteFile(filepath.Join(projectDir, "README.md"), []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to create README.md file: %w", err)
	}

	// Create config files
	configContent := templates.ConfigTemplate()
	if err := os.WriteFile(filepath.Join(projectDir, "internal/config/config.go"), []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to create config.go file: %w", err)
	}

	// Create logger files
	loggerContent := templates.LoggerTemplate()
	if err := os.WriteFile(filepath.Join(projectDir, "internal/logger/logger.go"), []byte(loggerContent), 0644); err != nil {
		return fmt.Errorf("failed to create logger.go file: %w", err)
	}

	// Create app files
	appContent := templates.AppTemplate(g.config.ProjectConfig)
	if err := os.WriteFile(filepath.Join(projectDir, "internal/app/app.go"), []byte(appContent), 0644); err != nil {
		return fmt.Errorf("failed to create app.go file: %w", err)
	}

	return nil
}

// generateComponentFiles generates the component-specific files
func (g *Generator) generateComponentFiles(projectDir string) error {
	g.log.Info("Generating component files",
		"http", g.config.ProjectConfig.Components.HTTP,
		"postgres", g.config.ProjectConfig.Components.Postgres,
		"docker", g.config.ProjectConfig.Components.Docker,
	)

	// Generate HTTP files
	if g.config.ProjectConfig.Components.HTTP {
		if err := g.generateHTTPFiles(projectDir); err != nil {
			return fmt.Errorf("failed to generate HTTP files: %w", err)
		}
	}

	// Generate PostgreSQL files
	if g.config.ProjectConfig.Components.Postgres {
		if err := g.generatePostgresFiles(projectDir); err != nil {
			return fmt.Errorf("failed to generate PostgreSQL files: %w", err)
		}
	}

	// Generate Docker files
	if g.config.ProjectConfig.Components.Docker {
		if err := g.generateDockerFiles(projectDir); err != nil {
			return fmt.Errorf("failed to generate Docker files: %w", err)
		}
	}

	// Generate Kubernetes files
	if g.config.ProjectConfig.Components.Kubernetes {
		if err := g.generateKubernetesFiles(projectDir); err != nil {
			return fmt.Errorf("failed to generate Kubernetes files: %w", err)
		}
	}

	// Generate CI/CD files
	if g.config.ProjectConfig.Components.CICD {
		if err := g.generateCICDFiles(projectDir); err != nil {
			return fmt.Errorf("failed to generate CI/CD files: %w", err)
		}
	}

	return nil
}

// generateHTTPFiles generates the HTTP-specific files
func (g *Generator) generateHTTPFiles(projectDir string) error {
	g.log.Info("Generating HTTP files")

	// Create directories
	dirs := []string{
		"internal/api",
		"internal/api/handlers",
		"internal/api/middleware",
		"internal/api/routes",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create API files
	serverContent := templates.APIServerTemplate()
	if err := os.WriteFile(filepath.Join(projectDir, "internal/api/server.go"), []byte(serverContent), 0644); err != nil {
		return fmt.Errorf("failed to create server.go file: %w", err)
	}

	handlersContent := templates.APIHandlersTemplate()
	if err := os.WriteFile(filepath.Join(projectDir, "internal/api/handlers/handlers.go"), []byte(handlersContent), 0644); err != nil {
		return fmt.Errorf("failed to create handlers.go file: %w", err)
	}

	middlewareContent := templates.APIMiddlewareTemplate()
	if err := os.WriteFile(filepath.Join(projectDir, "internal/api/middleware/middleware.go"), []byte(middlewareContent), 0644); err != nil {
		return fmt.Errorf("failed to create middleware.go file: %w", err)
	}

	routesContent := templates.APIRoutesTemplate()
	if err := os.WriteFile(filepath.Join(projectDir, "internal/api/routes/routes.go"), []byte(routesContent), 0644); err != nil {
		return fmt.Errorf("failed to create routes.go file: %w", err)
	}

	return nil
}

// generatePostgresFiles generates the PostgreSQL-specific files
func (g *Generator) generatePostgresFiles(projectDir string) error {
	g.log.Info("Generating PostgreSQL files")

	// Create directories
	dirs := []string{
		"internal/db",
		"internal/db/migrations",
		"internal/db/models",
		"internal/db/repositories",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create DB files
	dbContent := templates.DBTemplate()
	if err := os.WriteFile(filepath.Join(projectDir, "internal/db/db.go"), []byte(dbContent), 0644); err != nil {
		return fmt.Errorf("failed to create db.go file: %w", err)
	}

	modelsContent := templates.DBModelsTemplate()
	if err := os.WriteFile(filepath.Join(projectDir, "internal/db/models/models.go"), []byte(modelsContent), 0644); err != nil {
		return fmt.Errorf("failed to create models.go file: %w", err)
	}

	reposContent := templates.DBRepositoriesTemplate()
	if err := os.WriteFile(filepath.Join(projectDir, "internal/db/repositories/repositories.go"), []byte(reposContent), 0644); err != nil {
		return fmt.Errorf("failed to create repositories.go file: %w", err)
	}

	// Create migration files
	migrationContent := templates.DBMigrationTemplate()
	if err := os.WriteFile(filepath.Join(projectDir, "internal/db/migrations/000001_init.up.sql"), []byte(migrationContent), 0644); err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "internal/db/migrations/000001_init.down.sql"), []byte("-- Revert initial migration"), 0644); err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}

	return nil
}

// generateDockerFiles generates the Docker-specific files
func (g *Generator) generateDockerFiles(projectDir string) error {
	g.log.Info("Generating Docker files")

	// Create Dockerfile
	dockerfileContent := templates.DockerfileTemplate(g.config.ProjectConfig)
	if err := os.WriteFile(filepath.Join(projectDir, "Dockerfile"), []byte(dockerfileContent), 0644); err != nil {
		return fmt.Errorf("failed to create Dockerfile: %w", err)
	}

	// Create docker-compose.yml
	composeContent := templates.DockerComposeTemplate(g.config.ProjectConfig)
	if err := os.WriteFile(filepath.Join(projectDir, "docker-compose.yml"), []byte(composeContent), 0644); err != nil {
		return fmt.Errorf("failed to create docker-compose.yml: %w", err)
	}

	// Create .dockerignore
	dockerignoreContent := templates.DockerignoreTemplate()
	if err := os.WriteFile(filepath.Join(projectDir, ".dockerignore"), []byte(dockerignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .dockerignore: %w", err)
	}

	return nil
}

// generateKubernetesFiles generates the Kubernetes-specific files
func (g *Generator) generateKubernetesFiles(projectDir string) error {
	g.log.Info("Generating Kubernetes files")

	// Create directory
	if err := os.MkdirAll(filepath.Join(projectDir, "deployments/kubernetes"), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create Kubernetes manifests
	deploymentContent := templates.KubernetesDeploymentTemplate(g.config.ProjectConfig)
	if err := os.WriteFile(filepath.Join(projectDir, "deployments/kubernetes/deployment.yaml"), []byte(deploymentContent), 0644); err != nil {
		return fmt.Errorf("failed to create deployment.yaml: %w", err)
	}

	serviceContent := templates.KubernetesServiceTemplate(g.config.ProjectConfig)
	if err := os.WriteFile(filepath.Join(projectDir, "deployments/kubernetes/service.yaml"), []byte(serviceContent), 0644); err != nil {
		return fmt.Errorf("failed to create service.yaml: %w", err)
	}

	configMapContent := templates.KubernetesConfigMapTemplate(g.config.ProjectConfig)
	if err := os.WriteFile(filepath.Join(projectDir, "deployments/kubernetes/configmap.yaml"), []byte(configMapContent), 0644); err != nil {
		return fmt.Errorf("failed to create configmap.yaml: %w", err)
	}

	return nil
}

// generateCICDFiles generates the CI/CD-specific files
func (g *Generator) generateCICDFiles(projectDir string) error {
	g.log.Info("Generating CI/CD files")

	// Create directory
	if err := os.MkdirAll(filepath.Join(projectDir, ".github/workflows"), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create GitHub Actions workflow
	workflowContent := templates.GitHubWorkflowTemplate(g.config.ProjectConfig)
	if err := os.WriteFile(filepath.Join(projectDir, ".github/workflows/main.yml"), []byte(workflowContent), 0644); err != nil {
		return fmt.Errorf("failed to create main.yml: %w", err)
	}

	return nil
}
