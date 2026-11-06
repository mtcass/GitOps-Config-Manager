package main

import (
	"fmt"
	"log"

	"github.com/mtcass/gitops-config-manager/internal/config"
	"github.com/mtcass/gitops-config-manager/internal/gitprovider"
)

func main() {
	// 1. Load configuration from environment / .env
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("Initializing GitHub client for owner: %s, repo: %s...\n", cfg.GitHubOwner, cfg.GitHubRepo)

	// 2. Initialize authenticated GitHub client
	client := gitprovider.NewGitHubClient(cfg.GitHubToken, cfg.GitHubOwner, cfg.GitHubRepo)

	// 3. Try to list files in the configured config directory
	fmt.Printf("Listing files in directory: %s\n", cfg.ConfigDirPath)
	files, err := client.ListConfigFiles(cfg.ConfigDirPath)
	if err != nil {
		log.Fatalf("Failed to list config files: %v", err)
	}

	fmt.Printf("Success! Found %d configuration file(s):\n", len(files))
	for _, file := range files {
		fmt.Printf(" - Path: %s (SHA: %s)\n", file.Path, file.SHA)
	}
}
