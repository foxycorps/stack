package cmd

import (
	"github.com/foxycorps/stack/pkg/git"
	"github.com/foxycorps/stack/pkg/stack"
	"github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
	Use:     "start [name]",
	Short:   "Start a new stack",
	GroupID: "STACK",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		stackName := args[0]

		repo, err := git.NewRepository(".")
		if err != nil {
			return err
		}

		// First we will checkout the main branch
		err = repo.FetchAll()
		if err != nil {
			return err
		}
		repo.SwitchBranch("main")

		branch_name := stackName + "_base"
		err = repo.CreateBranch(branch_name)
		if err != nil {
			return err
		}

		stack.InitializeStackFile(repo.Path, "main")
		stack.UpdateStackFile(repo.Path, branch_name)

		return nil
	},
}
