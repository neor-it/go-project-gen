// internal/cli/wizard.go - Interactive CLI wizard for project configuration
package cli

import (
	"fmt"
	//"path/filepath"
	//"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/username/goprojectgen/internal/config"
	"github.com/username/goprojectgen/internal/logger"
)

// Wizard represents the interactive CLI wizard
type Wizard struct {
	log logger.Logger
}

// NewWizard creates a new wizard
func NewWizard(log logger.Logger) *Wizard {
	return &Wizard{
		log: log,
	}
}

// Run runs the wizard and returns the project configuration
func (w *Wizard) Run() (config.ProjectConfig, error) {
	w.log.Info("Starting interactive project configuration wizard")

	// Define project configuration
	var projectCfg config.ProjectConfig

	// Ask for username
	username := ""
	prompt := &survey.Input{
		Message: "GitHub username or organization:",
		Help:    "This will be used to create the module path (e.g., github.com/username/project-name)",
	}
	if err := survey.AskOne(prompt, &username, survey.WithValidator(survey.Required)); err != nil {
		return projectCfg, err
	}
	projectCfg.Username = username

	// Ask for project name
	projectName := ""
	prompt = &survey.Input{
		Message: "Project name:",
		Help:    "This will be used as the directory name and in the module path",
	}
	if err := survey.AskOne(prompt, &projectName, survey.WithValidator(survey.Required)); err != nil {
		return projectCfg, err
	}
	projectCfg.ProjectName = projectName

	// Create module name
	projectCfg.ModuleName = fmt.Sprintf("github.com/%s/%s", username, projectName)

	// Ask for components
	components := []string{}
	componentsPrompt := &survey.MultiSelect{
		Message: "Select components to include:",
		Options: []string{
			"HTTP (Gin)",
			"PostgreSQL",
			"Docker",
			"Kubernetes",
			"CI/CD",
		},
		Default: []string{"HTTP (Gin)"},
	}
	if err := survey.AskOne(componentsPrompt, &components); err != nil {
		return projectCfg, err
	}

	// Set components
	projectCfg.Components = config.Components{
		HTTP:       contains(components, "HTTP (Gin)"),
		Postgres:   contains(components, "PostgreSQL"),
		Docker:     contains(components, "Docker"),
		Kubernetes: contains(components, "Kubernetes"),
		CICD:       contains(components, "CI/CD"),
	}

	// Print configuration
	w.log.Info("Project configuration",
		"username", projectCfg.Username,
		"projectName", projectCfg.ProjectName,
		"moduleName", projectCfg.ModuleName,
		"http", projectCfg.Components.HTTP,
		"postgres", projectCfg.Components.Postgres,
		"docker", projectCfg.Components.Docker,
		"kubernetes", projectCfg.Components.Kubernetes,
		"cicd", projectCfg.Components.CICD,
	)

	// Ask for confirmation
	confirmed := false
	confirmPrompt := &survey.Confirm{
		Message: "Confirm project configuration?",
		Default: true,
	}
	if err := survey.AskOne(confirmPrompt, &confirmed); err != nil {
		return projectCfg, err
	}

	if !confirmed {
		return w.Run()
	}

	return projectCfg, nil
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
