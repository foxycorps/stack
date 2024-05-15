package cmd

import (
	"slices"
	"strings"

	"github.com/foxycorps/stack/pkg/git"
	"github.com/foxycorps/stack/pkg/stack"
	"github.com/spf13/cobra"
)

var HomeCmd = &cobra.Command{
	Use:   "home",
	Short: "Navigates to the root of the stack",
	RunE: func(cmd *cobra.Command, args []string) error {

		repo, err := git.NewRepository(".")
		if err != nil {
			return err
		}

		err = repo.FetchAll()
		if err != nil {
			return err
		}

		branchName, err := repo.CurrentBranch()
		if err != nil {
			return err
		}

		stackName := strings.Split(branchName, ".")[0]

		branches, err := stack.ReadStackFile(repo.Path)
		if err != nil {
			return err
		}

		if slices.Contains(branches, stackName+".base") {
			repo.SwitchBranch(stackName + ".base")
		} else {
			repo.SwitchBranch("main")
		}

		return nil
	},
}
