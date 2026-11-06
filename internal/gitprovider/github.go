package gitprovider

import (
	"context"
	"fmt"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

// GitHubClient implements GitProvider interface using GitHub REST API
type GitHubClient struct {
	client *github.Client
	owner  string
	repo   string
	ctx    context.Context
}

// NewGitHubClient creates a new authenticated GitHub client instance
func NewGitHubClient(token, owner, repo string) *GitHubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &GitHubClient{
		client: client,
		owner:  owner,
		repo:   repo,
		ctx:    ctx,
	}
}

// ListConfigFiles lists all JSON configuration files in the specified directory
func (g *GitHubClient) ListConfigFiles(dirPath string) ([]ConfigFile, error) {
	_, directoryContent, _, err := g.client.Repositories.GetContents(g.ctx, g.owner, g.repo, dirPath, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get directory content from github: %w", err)
	}

	var files []ConfigFile
	for _, item := range directoryContent {
		if item.GetType() == "file" {
			files = append(files, ConfigFile{
				Path: item.GetPath(),
				SHA:  item.GetSHA(),
			})
		}
	}

	return files, nil
}

// GetConfigFile fetches the content of a single configuration file
func (g *GitHubClient) GetConfigFile(filePath string) (*ConfigFile, error) {
	fileContent, _, _, err := g.client.Repositories.GetContents(g.ctx, g.owner, g.repo, filePath, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file content: %w", err)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return nil, fmt.Errorf("failed to decode file content: %w", err)
	}

	return &ConfigFile{
		Path:    fileContent.GetPath(),
		Content: content,
		SHA:     fileContent.GetSHA(),
	}, nil
}

// CreateBranchAndPR handles the full GitOps write flow: branch creation, file commit, and PR opening
func (g *GitHubClient) CreateBranchAndPR(branchName, commitMessage, prTitle, prBody, filePath string, content []byte) error {
	// 1. Get the default branch reference (e.g., 'main' or 'dev') to base the new branch on
	repoInfo, _, err := g.client.Repositories.Get(g.ctx, g.owner, g.repo)
	if err != nil {
		return fmt.Errorf("failed to fetch repository info: %w", err)
	}
	defaultBranch := repoInfo.GetDefaultBranch()

	baseRef, _, err := g.client.Git.GetRef(g.ctx, g.owner, g.repo, "refs/heads/"+defaultBranch)
	if err != nil {
		return fmt.Errorf("failed to get base branch reference: %w", err)
	}

	// 2. Create the new work branch
	newRefName := "refs/heads/" + branchName
	newRef := &github.Reference{
		Ref: github.String(newRefName),
		Object: &github.GitObject{
			SHA: baseRef.Object.SHA,
		},
	}
	_, _, err = g.client.Git.CreateRef(g.ctx, g.owner, g.repo, newRef)
	if err != nil {
		return fmt.Errorf("failed to create branch %s: %w", branchName, err)
	}

	// 3. Get existing file SHA if it already exists (needed for updating files via API)
	existingFile, _ := g.GetConfigFile(filePath)
	var fileSHA *string
	if existingFile != nil {
		fileSHA = &existingFile.SHA
	}

	// 4. Create or update file commit on the new branch
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(commitMessage),
		Content: content,
		Branch:  github.String(branchName),
		SHA:     fileSHA,
	}

	_, _, err = g.client.Repositories.CreateFile(g.ctx, g.owner, g.repo, filePath, opts)
	if err != nil {
		return fmt.Errorf("failed to commit file changes: %w", err)
	}

	// 5. Open the Pull Request
	newPR := &github.NewPullRequest{
		Title:               github.String(prTitle),
		Head:                github.String(branchName),
		Base:                github.String(defaultBranch),
		Body:                github.String(prBody),
		MaintainerCanModify: github.Bool(true),
	}

	_, _, err = g.client.PullRequests.Create(g.ctx, g.owner, g.repo, newPR)
	if err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}

	return nil
}
