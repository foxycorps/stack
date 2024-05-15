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

var UpCmd = &cobra.Command{
	Use:     "up",
	Short:   "Move up the stack",
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

		branches, err := stack.ReadStackFile(repo.Path)
		if err != nil {
			return err
		}

		fmt.Printf("Branches: %v\n", branches)
		currentIndex := utils.CurrentIndex(branches, currentBranch)

		if currentIndex == -1 {
			return fmt.Errorf("current branch is not in the stack")
		}

		nextBranch := branches[currentIndex-1]
		err = repo.SwitchBranch(nextBranch)
		if err != nil {
			return err
		}
		fmt.Printf("Moved to '%s'\n", nextBranch)

		return nil
	},
}
