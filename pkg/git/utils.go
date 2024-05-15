package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func (r *Repository) Execute(args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", strings.TrimSuffix(r.Path, "/.git")}, args...)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute git command '%s': %s", strings.Join(args, " "), string(output))
	}
	return strings.TrimSpace(string(output)), nil
}

func (r *Repository) ExecuteStrings(args ...string) ([]string, error) {
	cmd := exec.Command("git", append([]string{"-C", strings.TrimSuffix(r.Path, "/.git")}, args...)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute git command '%s': %s", strings.Join(args, " "), string(output))
	}

	// Split the output into lines and trim whitespace for each line.
	lines := strings.Split(string(output), "\n")
	trimmedLines := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			trimmedLines = append(trimmedLines, trimmedLine)
		}
	}
	return trimmedLines, nil
}
