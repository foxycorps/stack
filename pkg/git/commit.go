package git

import (
	"fmt"
	"strings"
)

func (r *Repository) FindLastCommit() (string, error) {
	output, err := r.Execute("rev-parse", "HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to find last commit: %v", err)
	}

	return output, nil
}

func (r *Repository) CreateCommit(message string) error {
	_, err := r.Execute("commit", "-m", message)
	if err != nil {
		return fmt.Errorf("failed to create commit: %v", err)
	}
	return nil
}

func (r *Repository) GetCommitsOnBranch(branch string) ([]string, error) {
	output, err := r.Execute("log", "--format=%H", branch)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits on branch: %v", err)
	}
	commits := strings.Split(output, "\n")
	return commits, nil
}

func (r *Repository) GetFirstCommitNotOnParent(branch, parentBranch string) (string, error) {
	// Get the commits on the branch
	commits, err := r.GetCommitsOnBranch(branch)
	if err != nil {
		return "", fmt.Errorf("failed to get commits on branch: %v", err)
	}

	// Get the commits on the parent branch
	parentCommits, err := r.GetCommitsOnBranch(parentBranch)
	if err != nil {
		return "", fmt.Errorf("failed to get commits on parent branch: %v", err)
	}

	// Find the first commit on the branch that is not on the parent branch
	for _, commit := range commits {
		found := false
		for _, parentCommit := range parentCommits {
			if commit == parentCommit {
				found = true
				break
			}
		}
		if !found {
			return commit, nil
		}
	}

	return "", fmt.Errorf("no commit found on branch that is not on parent branch")
}

func (r *Repository) GetCommit(sha string) (string, error) {
	output, err := r.Execute("show", "--format=%s%n%b", sha)
	if err != nil {
		return "", fmt.Errorf("failed to get commit: %v", err)
	}
	lines := strings.SplitN(output, "\n", 2)
	return strings.TrimSpace(lines[0]), nil
}

func (r *Repository) AddAll() error {
	_, err := r.Execute("add", ".")
	if err != nil {
		return fmt.Errorf("failed to add all files: %v", err)
	}
	return nil
}
