package cmd

import (
	"fmt"
	"strings"

	"github.com/foxycorps/stack/pkg/git"
	"github.com/foxycorps/stack/pkg/stack"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:     "add [name]",
	Short:   "Add a new layer to the stack",
	GroupID: "STACK",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := git.NewRepository(".")
		if err != nil {
			return err
		}

		currentBranch, err := repo.CurrentBranch()
		if err != nil {
			return err
		}
		stackName := strings.Split(currentBranch, "_")[0]
		branches, err := stack.ReadStackFile(repo.Path)
		if err != nil {
			return err
		}
		order := len(branches) - 1
		newBranchName := fmt.Sprintf("%s_%s_%02d", stackName, args[0], order)

		// Check if the current branch is at the top of the stack
		if currentBranch != branches[len(branches)-1] {
			// Find the position to insert the new branch
			insertPosition := 0
			for i, branch := range branches {
				if branch == currentBranch {
					insertPosition = i + 1
					break
				}
			}

			// Insert the new branch at the appropriate position
			err = stack.InsertBranch(repo.Path, newBranchName, insertPosition)
			if err != nil {
				return err
			}

			// Rename subsequent branches to reflect the updated number suffixes
			for i := insertPosition + 1; i < len(branches); i++ {
				oldBranchName := branches[i]
				newSuffix := fmt.Sprintf("%02d", i)
				parts := strings.Split(oldBranchName, "_")
				parts[len(parts)-1] = newSuffix
				newBranchName := strings.Join(parts, "_")
				err = repo.RenameBranch(oldBranchName, newBranchName)
				if err != nil {
					return err
				}
			}
		} else {
			// Current branch is at the top, proceed with the usual branch creation
			err = repo.CreateBranch(newBranchName)
			if err != nil {
				return err
			}

			err = stack.UpdateStackFile(repo.Path, newBranchName)
			if err != nil {
				return err
			}
		}

		err = repo.SwitchBranch(newBranchName)
		if err != nil {
			return err
		}

		fmt.Printf("New branch '%s' created successfully\n", newBranchName)

		return nil
	},
}
