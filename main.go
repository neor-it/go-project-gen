// main.go - Main entry point for the Go Project Generator
package main

import (
	"fmt"
	"os"

	"github.com/neor-it/go-project-gen/internal/cli"
	"github.com/neor-it/go-project-gen/internal/config"
	"github.com/neor-it/go-project-gen/internal/generator"
	"github.com/neor-it/go-project-gen/internal/logger"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()
	log.Info("Starting Go Project Generator")

	// Check if /output directory exists and is writable when running in Docker
	outputDir := "."
	if _, err := os.Stat("/output"); err == nil {
		// We're inside Docker with mounted volume
		outputDir = "/output"
		log.Info("Using Docker volume output directory", "path", outputDir)
	}

	// Parse command line arguments
	cfg, err := config.ParseArgs(os.Args[1:])
	if err != nil {
		log.Fatal("Failed to parse arguments", "error", err)
	}

	// Set the output directory
	cfg.OutputDir = outputDir

	// Run CLI wizard if no configuration file provided
	if cfg.IsInteractive {
		wizard := cli.NewWizard(log)
		projectCfg, err := wizard.Run()
		if err != nil {
			log.Fatal("Failed to run wizard", "error", err)
		}
		cfg.ProjectConfig = projectCfg
	}

	// Generate project
	gen := generator.NewGenerator(log, cfg)
	if err := gen.Generate(); err != nil {
		log.Fatal("Failed to generate project", "error", err)
	}

	// Show success message with correct path information
	projectPath := fmt.Sprintf("%s/%s", outputDir, cfg.ProjectConfig.ProjectName)
	if outputDir == "/output" {
		// When running in Docker, show the path relative to the user's current directory
		projectPath = cfg.ProjectConfig.ProjectName
	}

	fmt.Println("✅ Project successfully generated!")
	fmt.Printf("📂 Location: %s\n", projectPath)
}
