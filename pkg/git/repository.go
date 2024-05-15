package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type Repository struct {
	Path string
}

func NewRepository(path string) (*Repository, error) {

	// Check if the provided path is a valid Git repository
	location, err := exec.Command("git", "-C", path, "rev-parse", "--git-dir").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("invalid Git repository: %s", path)
	}

	abs, err := filepath.Abs(strings.TrimSpace(string(location)))
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %v", err)
	}

	return &Repository{
		Path: abs,
	}, nil
}

func (r *Repository) Owner() (string, error) {
	output, err := r.Execute("config", "--get", "remote.origin.url")
	if err != nil {
		return "", fmt.Errorf("failed to get repository owner: %v", err)
	}

	remoteURL := strings.TrimSpace(output)
	parts := strings.Split(remoteURL, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid remote URL format")
	}

	owner := parts[len(parts)-2]
	return owner, nil
}

func (r *Repository) Name() (string, error) {
	output, err := r.Execute("config", "--get", "remote.origin.url")
	if err != nil {
		return "", fmt.Errorf("failed to get repository name: %v", err)
	}

	remoteURL := strings.TrimSpace(output)
	parts := strings.Split(remoteURL, "/")
	if len(parts) < 1 {
		return "", fmt.Errorf("invalid remote URL format")
	}

	name := strings.TrimSuffix(parts[len(parts)-1], ".git")
	return name, nil
}

func (r *Repository) FetchAll() error {
	_, err := r.Execute("fetch", "--all")
	if err != nil {
		return fmt.Errorf("failed to fetch all: %v", err)
	}
	return nil
}
