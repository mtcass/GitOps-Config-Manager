package main

import (
	"log"

	"github.com/mtcass/gitops-config-manager/internal/config"
	"github.com/mtcass/gitops-config-manager/internal/gitprovider"
	"github.com/mtcass/gitops-config-manager/internal/server"
)

func main() {
	// 1. Load application configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Initialize the GitHub client provider
	gitClient := gitprovider.NewGitHubClient(cfg.GitHubToken, cfg.GitHubOwner, cfg.GitHubRepo)

	// 3. Initialize and start the HTTP server
	srv := server.NewServer(cfg, gitClient)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
