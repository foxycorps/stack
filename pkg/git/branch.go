package git

import (
	"fmt"
	"strings"

	"github.com/foxycorps/stack/pkg/utils"
)

func (r *Repository) CreateBranch(branchName string) error {
	_, err := r.Execute("switch", "-c", branchName)
	return err
}

func (r *Repository) SwitchBranch(branchName string) error {
	_, err := r.Execute("switch", branchName)
	return err
}

func (r *Repository) CurrentBranch() (string, error) {
	output, err := r.Execute("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %v", err)
	}

	return output, nil
}

func (r *Repository) ListBranches(stackName string) ([]string, error) {
	output, err := r.ExecuteStrings("for-each-ref", "--sort=creatordate", "--format=%(refname:short)", "refs/heads/")
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %v", err)
	}

	return utils.FilterArrayByValue(output, stackName+"_"), nil
}

func (r *Repository) ListRemoteBranches(stackName string) ([]string, error) {
	output, err := r.ExecuteStrings("for-each-ref", "--sort=creatordate", "--format=%(refname:short)", "refs/remotes/origin/")
	if err != nil {
		return nil, fmt.Errorf("failed to list remote branches: %v", err)
	}

	fmt.Println("output", output)

	return utils.FilterArrayByValue(output, stackName+"_"), nil
}

func (r *Repository) TrackRemoteBranch(branchName string) error {
	_, err := r.Execute("branch", "--track", branchName, fmt.Sprintf("origin/%s", branchName))
	return err
}

func (r *Repository) RebaseBranch(branchName, parentName string) error {
	// Switch to the branch we want to rebase.
	err := r.SwitchBranch(branchName)
	if err != nil {
		return err
	}

	_, err = r.Execute("rebase", parentName, branchName)
	if err != nil {
		// We will try and cancel the rebase.
		_, err = r.Execute("rebase", "--abort")
		if err != nil {
			return fmt.Errorf("failed to rebase branch: %v", err)
		}

		return fmt.Errorf("failed to rebase branch: %v", err)
	}
	return err
}

func (r *Repository) Push(branchName string) error {
	_, err := r.Execute("push", "--set-upstream", "origin", "--force-with-lease", branchName)
	if err != nil {
		// If the error has todo with the permissions, we will say that.
		if strings.Contains(err.Error(), "403") {
			return fmt.Errorf("it appears you do not have write access to this repo: %v", err)
		}
		return fmt.Errorf("failed to push branch to remote: %v", err)
	}

	return nil
}

func (r *Repository) Pull(branchName string) error {
	_, err := r.Execute("pull", "origin", branchName, "--rebase")
	if err != nil {
		return fmt.Errorf("failed to pull branch from remote: %v", err)
	}

	return nil
}

func (r *Repository) BranchExists(branchName string) (bool, error) {
	_, err := r.Execute("show-ref", "--verify", "--quiet", "refs/heads/"+branchName)
	if err != nil {
		// If the command fails, it means the branch doesn't exist locally.
		return false, nil
	}

	// If the command succeeds, it means the branch exists locally.
	return true, nil
}

func (r *Repository) GetOwner(branchName string) (string, error) {
	output, err := r.Execute("log", "--pretty=format:%an", branchName)
	if err != nil {
		return "", fmt.Errorf("failed to get owner of branch: %v", err)
	}

	// Split the output by newlines
	lines := strings.Split(output, "\n")

	// Return the first line (name of the person who committed the first commit)
	return lines[0], nil
}

func (r *Repository) FindRemoteStacks() ([]string, error) {
	// List all the `*.base` branches from remote.
	output, err := r.ExecuteStrings("ls-remote", "--heads", "origin")
	if err != nil {
		return nil, fmt.Errorf("failed to list remote branches: %v", err)
	}

	// Filter the branches by the stack name.
	stacks := []string{}
	for _, branch := range output {
		if strings.Contains(branch, ".base") {
			stacks = append(stacks, strings.TrimSuffix(strings.Split(branch, "/")[2], ".base"))
		}
	}
	return stacks, nil
}

func (r *Repository) RenameBranch(oldBranchName, newBranchName string) error {
	err := r.SwitchBranch(oldBranchName)
	if err != nil {
		return fmt.Errorf("failed to switch to branch %s: %v", oldBranchName, err)
	}
	_, err = r.Execute("branch", "-m", newBranchName)
	if err != nil {
		return fmt.Errorf("failed to rename branch %s to %s: %v", oldBranchName, newBranchName, err)
	}
	return nil
}
