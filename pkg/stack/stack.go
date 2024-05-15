package stack

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const stackFile = "/.stack"

func stackFilePath(root string) string {
	path, _ := filepath.Abs(filepath.Join(root, stackFile))
	return path
}

func InitializeStackFile(root string, baseBranch string) error {
	path := stackFilePath(root)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, baseBranch+"\n")
	return err
}

func UpdateStackFile(root string, branchName string) error {
	path := stackFilePath(root)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, branchName+"\n")
	return err
}

func InsertBranch(root string, branchName string, position int) error {
	path := stackFilePath(root)
	branches, err := ReadStackFile(root)
	if err != nil {
		return err
	}

	if position < 0 || position > len(branches) {
		return fmt.Errorf("invalid position provided")
	}

	newBranchList := branches[:position]
	newBranchList = append(newBranchList, branchName)
	newBranchList = append(newBranchList, branches[position:]...)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, strings.Join(newBranchList, "\n"))
	return err
}

func ReadStackFile(root string) ([]string, error) {
	path := stackFilePath(root)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := make([]string, 0)
	for _, line := range strings.Split(string(data), "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, nil
}

func GetParents(root string) (map[string]string, error) {
	branches, err := ReadStackFile(root)
	if err != nil {
		return nil, err
	}

	parents := make(map[string]string)
	for i := 1; i < len(branches); i++ {
		parents[branches[i]] = branches[i-1]
	}
	return parents, nil
}
