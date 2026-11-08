package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mtcass/gitops-config-manager/internal/config"
	"github.com/mtcass/gitops-config-manager/internal/gitprovider"
)

func main() {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("Initializing GitHub client for owner: %s, repo: %s...\n", cfg.GitHubOwner, cfg.GitHubRepo)
	client := gitprovider.NewGitHubClient(cfg.GitHubToken, cfg.GitHubOwner, cfg.GitHubRepo)

	// 2. Validation 1: List configuration files
	files := runListValidation(client, cfg.ConfigDirPath)
	if len(files) == 0 {
		log.Fatalf("Aborting: No configuration files found.")
	}

	// 3. Validation 2: Read specific file content
	targetPath := files[0].Path
	fileContent := runGetFileValidation(client, targetPath)

	// 4. Validation 3: GitOps write flow (Controlled by environment variable to avoid unwanted PRs)
	if os.Getenv("RUN_WRITE_TEST") == "true" {
		runGitOpsWriteValidation(client, targetPath, fileContent)
	} else {
		fmt.Println("\n[3/3] Write validation skipped. Set RUN_WRITE_TEST=true to test branch creation and PR opening.")
	}
}

// runListValidation validates directory listing
func runListValidation(client *gitprovider.GitHubClient, dirPath string) []gitprovider.ConfigFile {
	fmt.Printf("\n[1/3] Validating ListConfigFiles in '%s'...\n", dirPath)
	files, err := client.ListConfigFiles(dirPath)
	if err != nil {
		log.Fatalf("Failed to list config files: %v", err)
	}

	fmt.Printf("✔ Success! Found %d file(s):\n", len(files))
	for _, f := range files {
		fmt.Printf("  - %s (SHA: %s)\n", f.Path, f.SHA)
	}
	return files
}

// runGetFileValidation validates single file content reading
func runGetFileValidation(client *gitprovider.GitHubClient, filePath string) string {
	fmt.Printf("\n[2/3] Validating GetConfigFile for '%s'...\n", filePath)
	file, err := client.GetConfigFile(filePath)
	if err != nil {
		log.Fatalf("Failed to get file content: %v", err)
	}

	fmt.Printf("✔ Success! Content length: %d bytes\n", len(file.Content))
	return file.Content
}

// runGitOpsWriteValidation validates branch creation, commit, and PR opening
func runGitOpsWriteValidation(client *gitprovider.GitHubClient, filePath, currentContent string) {
	timestamp := time.Now().Format("20060102-150405")
	branchName := fmt.Sprintf("feat/smoke-test-%s", timestamp)
	commitMessage := "chore: automated smoke test validation"
	prTitle := "Automated Smoke Test PR"
	prBody := "Validating GitOps create branch and PR pipeline."
	updatedContent := []byte(currentContent + "\n# Smoke test timestamp: " + timestamp)

	fmt.Printf("\n[3/3] Validating GitOps flow (Branch: %s)...\n", branchName)
	err := client.CreateBranchAndPR(branchName, commitMessage, prTitle, prBody, filePath, updatedContent)
	if err != nil {
		log.Fatalf("Failed to execute GitOps flow: %v", err)
	}

	fmt.Println("✔ Success! Branch created, file committed, and Pull Request opened.")
}
