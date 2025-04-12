// main.go - Main entry point for the Go Project Generator
package main

import (
	"fmt"
	"os"

	"github.com/username/goprojectgen/internal/cli"
	"github.com/username/goprojectgen/internal/config"
	"github.com/username/goprojectgen/internal/generator"
	"github.com/username/goprojectgen/internal/logger"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()
	log.Info("Starting Go Project Generator")

	// Parse command line arguments
	cfg, err := config.ParseArgs(os.Args[1:])
	if err != nil {
		log.Fatal("Failed to parse arguments", "error", err)
	}

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

	fmt.Println("✅ Project successfully generated!")
	fmt.Printf("📂 Location: %s\n", cfg.OutputDir)
}
