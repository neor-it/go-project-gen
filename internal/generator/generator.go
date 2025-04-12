// generator.go - Обновляем метод createStandardStructure
// internal/generator/generator.go - Updated with conditional scripts directory creation
package generator

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/neor-it/go-project-gen/internal/config"
	"github.com/neor-it/go-project-gen/internal/generator/templates"
	"github.com/neor-it/go-project-gen/internal/logger"
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
		"outputDir", g.config.OutputDir,
	)

	// Check if output directory is writable
	testFile := filepath.Join(g.config.OutputDir, ".test-write-permission")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("output directory %s is not writable: %w", g.config.OutputDir, err)
	}

	// Clean up test file
	os.Remove(testFile)

	// Create project directory
	projectDir := filepath.Join(g.config.OutputDir, g.config.ProjectConfig.ProjectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	g.log.Info("Project directory created", "path", projectDir)

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

	// Run go mod tidy to update dependencies
	if err := g.runGoModTidy(projectDir); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	return nil
}

// runGoModTidy runs go mod tidy in the project directory
func (g *Generator) runGoModTidy(projectDir string) error {
	g.log.Info("Running go mod tidy in the project directory")

	// Create command to run go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	g.log.Info("Successfully ran go mod tidy")
	return nil
}

// createStandardStructure creates the standard Go project structure
func (g *Generator) createStandardStructure(projectDir string) error {
	// Create base directories
	dirs := []string{
		"internal",
		"internal/app",
		"internal/config",
		"internal/logger",
		"pkg",
	}

	// Add scripts directories only if Postgres is selected
	if g.config.ProjectConfig.Components.Postgres {
		dirs = append(dirs,
			"scripts",
			"scripts/migtool",
			"scripts/modelgen",
		)
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

	// Create config files - use dynamic template generation
	configContent := templates.ConfigTemplate(g.config.ProjectConfig)
	if err := os.WriteFile(filepath.Join(projectDir, "internal/config/config.go"), []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to create config.go file: %w", err)
	}

	// Create .env and .env.example files
	envContent := g.generateEnvFile()
	if err := os.WriteFile(filepath.Join(projectDir, ".env.example"), []byte(envContent), 0644); err != nil {
		return fmt.Errorf("failed to create .env.example file: %w", err)
	}

	if err := os.WriteFile(filepath.Join(projectDir, ".env"), []byte(envContent), 0644); err != nil {
		return fmt.Errorf("failed to create .env file: %w", err)
	}

	// Create logger files
	loggerContent := templates.LoggerTemplate()
	if err := g.writeTemplateFile(filepath.Join(projectDir, "internal/logger/logger.go"), loggerContent); err != nil {
		return fmt.Errorf("failed to create logger.go file: %w", err)
	}

	// Create app files
	appContent := templates.AppTemplate(g.config.ProjectConfig)
	if err := g.writeTemplateFile(filepath.Join(projectDir, "internal/app/app.go"), appContent); err != nil {
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

	// Generate migrations files
	if g.config.ProjectConfig.Components.Postgres {
		if err := g.generateMigrationsFiles(projectDir); err != nil {
			return fmt.Errorf("failed to generate migrations files: %w", err)
		}
	}

	// Generate Docker files
	if g.config.ProjectConfig.Components.Docker {
		if err := g.generateDockerFiles(projectDir); err != nil {
			return fmt.Errorf("failed to generate Docker files: %w", err)
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

// writeFile writes raw content to a file without template processing
func (g *Generator) writeFile(path, content string) error {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// writeTemplateFile writes a template file with the given content
func (g *Generator) writeTemplateFile(path, content string) error {
	tmpl, err := template.New(filepath.Base(path)).Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	data := map[string]interface{}{
		"ModuleName":  g.config.ProjectConfig.ModuleName,
		"ProjectName": g.config.ProjectConfig.ProjectName,
		"Username":    g.config.ProjectConfig.Username,
		"Components":  g.config.ProjectConfig.Components,
		"Timestamp":   time.Now().Format(time.RFC3339),
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
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
	if err := g.writeTemplateFile(filepath.Join(projectDir, "internal/api/server.go"), serverContent); err != nil {
		return fmt.Errorf("failed to create server.go file: %w", err)
	}

	handlersContent := templates.APIHandlersTemplate()
	if err := g.writeTemplateFile(filepath.Join(projectDir, "internal/api/handlers/handlers.go"), handlersContent); err != nil {
		return fmt.Errorf("failed to create handlers.go file: %w", err)
	}

	middlewareContent := templates.APIMiddlewareTemplate()
	if err := g.writeTemplateFile(filepath.Join(projectDir, "internal/api/middleware/middleware.go"), middlewareContent); err != nil {
		return fmt.Errorf("failed to create middleware.go file: %w", err)
	}

	routesContent := templates.APIRoutesTemplate()
	if err := g.writeTemplateFile(filepath.Join(projectDir, "internal/api/routes/routes.go"), routesContent); err != nil {
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
	if err := g.writeTemplateFile(filepath.Join(projectDir, "internal/db/db.go"), dbContent); err != nil {
		return fmt.Errorf("failed to create db.go file: %w", err)
	}

	modelsContent := templates.UserModelTemplate()
	if err := g.writeTemplateFile(filepath.Join(projectDir, "internal/db/models/users.go"), modelsContent); err != nil {
		return fmt.Errorf("failed to create models.go file: %w", err)
	}

	reposContent := templates.DBRepositoriesTemplate()
	if err := g.writeTemplateFile(filepath.Join(projectDir, "internal/db/repositories/repositories.go"), reposContent); err != nil {
		return fmt.Errorf("failed to create repositories.go file: %w", err)
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

func (g *Generator) generateEnvFile() string {
	env := `# Server Configuration
SERVER_PORT=8080
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s

# Logging Configuration
LOGGING_LEVEL=info

# Application Configuration
SHUTDOWN_TIMEOUT=5s
`

	// Add database configuration if PostgreSQL is selected
	if g.config.ProjectConfig.Components.Postgres {
		// Base connection string uses localhost for direct development
		env += `
# Database Configuration for local development
# DB_CONNECTION_STRING=postgres://postgres:postgres@localhost:5432/` + g.config.ProjectConfig.ProjectName + `?sslmode=disable
`

		// If Docker is also selected, add a commented Docker-specific connection string as reference
		if g.config.ProjectConfig.Components.Docker {
			env += `
# Database Configuration for Docker environment:
DB_CONNECTION_STRING=postgres://postgres:postgres@postgres:5432/` + g.config.ProjectConfig.ProjectName + `?sslmode=disable
`
		}
	}

	// Add Docker configuration if Docker is selected
	if g.config.ProjectConfig.Components.Docker {
		env += `
# Docker Configuration
DOCKER_REGISTRY=` + g.config.ProjectConfig.Username + `
`
	}

	// Add CI/CD configuration if CI/CD is selected
	if g.config.ProjectConfig.Components.CICD {
		env += `
# CI/CD Configuration
CI_ENABLE_TESTS=true
CI_ENABLE_LINTING=true
`
	}

	return env
}

// generateMigrationsFiles generates the migration-specific files
func (g *Generator) generateMigrationsFiles(projectDir string) error {
	g.log.Info("Generating migrations files")

	// Create directories
	dirs := []string{
		"scripts/migtool",
		"scripts/modelgen",
		"internal/migrations",
		"internal/migrations/sql",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create migration tool files
	migrationToolContent := templates.MigrationToolTemplate()
	if err := g.writeTemplateFile(filepath.Join(projectDir, "scripts/migtool/migrations.go"), migrationToolContent); err != nil {
		return fmt.Errorf("failed to create migrations tool file: %w", err)
	}

	// Create model generator tool - Using our new comprehensive template
	// Use writeFile directly as modelgen.go content should not be templated here.
	modelGenContent := templates.ModelGeneratorFullTemplate()
	if err := g.writeFile(filepath.Join(projectDir, "scripts/modelgen/modelgen.go"), modelGenContent); err != nil {
		return fmt.Errorf("failed to create model generator file: %w", err)
	}

	// Create migration package file
	migrationPackageContent := templates.MigrationsPackageTemplate()
	if err := g.writeTemplateFile(filepath.Join(projectDir, "internal/migrations/migrations.go"), migrationPackageContent); err != nil {
		return fmt.Errorf("failed to create migrations package file: %w", err)
	}

	// Create initial migration files
	migrationUpContent := templates.MigrationFileTemplate()
	if err := os.WriteFile(filepath.Join(projectDir, "internal/migrations/sql", "001_init.up.sql"), []byte(migrationUpContent), 0644); err != nil {
		return fmt.Errorf("failed to create migration up file: %w", err)
	}

	migrationDownContent := templates.MigrationDownFileTemplate()
	if err := os.WriteFile(filepath.Join(projectDir, "internal/migrations/sql", "001_init.down.sql"), []byte(migrationDownContent), 0644); err != nil {
		return fmt.Errorf("failed to create migration down file: %w", err)
	}

	// Create migration script file
	scriptContent := templates.MigrationsScriptTemplate()
	scriptFile := filepath.Join(projectDir, "scripts/migrate.sh")
	if err := os.WriteFile(scriptFile, []byte(scriptContent), 0755); err != nil {
		return fmt.Errorf("failed to create migration script file: %w", err)
	}

	// Create model generator script file
	modelGenScriptContent := templates.ModelGeneratorScriptTemplate()
	modelGenScriptFile := filepath.Join(projectDir, "scripts/generate_models.sh")
	if err := os.WriteFile(modelGenScriptFile, []byte(modelGenScriptContent), 0755); err != nil {
		return fmt.Errorf("failed to create model generator script file: %w", err)
	}

	// Make scripts executable
	if err := os.Chmod(scriptFile, 0755); err != nil {
		return fmt.Errorf("failed to make migration script executable: %w", err)
	}

	if err := os.Chmod(modelGenScriptFile, 0755); err != nil {
		return fmt.Errorf("failed to make model generator script executable: %w", err)
	}

	return nil
}
