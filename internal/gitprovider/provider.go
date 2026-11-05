package gitprovider

// ConfigFile represents a configuration file fetched from the repository
type ConfigFile struct {
	Path    string
	Content string
	SHA     string
}

// GitProvider defines the standard interface for Git platform operations
type GitProvider interface {
	// ListConfigFiles lists configuration files inside a specific directory path
	ListConfigFiles(dirPath string) ([]ConfigFile, error)

	// GetConfigFile retrieves the content and SHA of a specific configuration file
	GetConfigFile(filePath string) (*ConfigFile, error)

	// CreateBranchAndPR creates a new work branch, commits changes to a file, and opens a Pull Request
	CreateBranchAndPR(branchName, commitMessage, prTitle, prBody, filePath string, content []byte) error
}
