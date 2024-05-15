package git

import (
	"os"
	"path/filepath"
)

func (r *Repository) HasPendingChanges() (bool, error) {
	output, err := r.Execute("status", "--porcelain")
	if err != nil {
		return false, err
	}

	return len(output) > 0, nil
}

func (r *Repository) FindPRTemplate() (string, error) {
	templatePaths := []string{
		"PULL_REQUEST_TEMPLATE.md",
		".github/PULL_REQUEST_TEMPLATE.md",
		"docs/PULL_REQUEST_TEMPLATE.md",
	}

	for _, path := range templatePaths {
		templatePath := filepath.Join(r.Path, path)
		if _, err := os.Stat(templatePath); err == nil {
			content, err := os.ReadFile(templatePath)
			if err != nil {
				return "", err
			}
			return string(content), nil
		}
	}

	return "", nil
}
