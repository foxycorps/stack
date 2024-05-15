package cmd

import (
	"fmt"

	"github.com/foxycorps/stack/pkg/git"
	"github.com/foxycorps/stack/pkg/stack"
	"github.com/foxycorps/stack/pkg/utils"
	"github.com/spf13/cobra"
)

/**
* Up and down are technically inverse because of how we store the stack file.
**/

var DownCmd = &cobra.Command{
	Use:     "down",
	Short:   "Move down the stack",
	GroupID: "NAVIGATION",
	RunE: func(cmd *cobra.Command, args []string) error {

		repo, err := git.NewRepository(".")
		if err != nil {
			return err
		}

		currentBranch, err := repo.CurrentBranch()
		if err != nil {
			return err
		}

		// stackName := strings.SplitN(currentBranch, "_", 1)[0]
		// branches, err := repo.ListBranches(stackName)
		branches, err := stack.ReadStackFile(repo.Path)
		if err != nil {
			return err
		}

		currentIndex := utils.CurrentIndex(branches, currentBranch)
		if currentIndex == -1 {
			return fmt.Errorf("Current branch is not in the stack")
		}

		nextBranch := branches[currentIndex+1]
		err = repo.SwitchBranch(nextBranch)
		if err != nil {
			return err
		}
		fmt.Printf("Moved to '%s'\n", nextBranch)

		return nil
	},
}
