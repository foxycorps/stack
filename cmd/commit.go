package cmd

import (
	"github.com/foxycorps/stack/pkg/git"
	"github.com/spf13/cobra"
)

var CommitCmd = &cobra.Command{
	Use:     "commit [message]",
	Short:   "Commit the changes to the current branch",
	GroupID: "STACK",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		repo, err := git.NewRepository(".")
		if err != nil {
			return err
		}

		commitMessage := args[0]

		err = repo.AddAll()
		if err != nil {
			return err
		}

		err = repo.CreateCommit(commitMessage)
		if err != nil {
			return err
		}

		return nil
	},
}
