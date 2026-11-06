package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables
type Config struct {
	Port            string
	GitProvider     string
	GitHubToken     string
	GitHubOwner     string
	GitHubRepo      string
	GitLabURL       string
	GitLabToken     string
	GitLabProjectID string
	ConfigDirPath   string
	SchemaDirPath   string
}

// Load loads configuration from environment variables and optional .env file
func Load() (*Config, error) {
	// Load .env file if present (ignored in production if file doesn't exist)
	_ = godotenv.Load()

	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		GitProvider:     getEnv("GIT_PROVIDER", "github"),
		GitHubToken:     os.Getenv("GITHUB_TOKEN"),
		GitHubOwner:     os.Getenv("GITHUB_OWNER"),
		GitHubRepo:      os.Getenv("GITHUB_REPO"),
		GitLabURL:       getEnv("GITLAB_URL", "https://gitlab.com"),
		GitLabToken:     os.Getenv("GITLAB_TOKEN"),
		GitLabProjectID: os.Getenv("GITLAB_PROJECT_ID"),
		ConfigDirPath:   getEnv("CONFIG_DIR_PATH", "configs/"),
		SchemaDirPath:   getEnv("SCHEMA_DIR_PATH", "schemas/"),
	}

	// Basic validation based on selected provider
	if cfg.GitProvider == "github" && (cfg.GitHubOwner == "" || cfg.GitHubRepo == "") {
		return nil, fmt.Errorf("github provider selected, but GITHUB_OWNER or GITHUB_REPO is missing")
	}

	if cfg.GitProvider == "gitlab" && cfg.GitLabProjectID == "" {
		return nil, fmt.Errorf("gitlab provider selected, but GITLAB_PROJECT_ID is missing")
	}

	return cfg, nil
}

// getEnv helper function to retrieve environment variables with a default fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
