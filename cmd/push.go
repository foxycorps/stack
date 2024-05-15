package cmd

import (
	"fmt"

	"github.com/foxycorps/stack/pkg/git"
	"github.com/foxycorps/stack/pkg/stack"
	"github.com/spf13/cobra"
)

var PushCmd = &cobra.Command{
	Use:     "push",
	Short:   "Pushes the local database to the remote database",
	GroupID: "STACK",
	RunE: func(cmd *cobra.Command, args []string) error {

		repo, err := git.NewRepository(".")
		if err != nil {
			return err
		}

		// We will go through the known branches and push them to the remote.
		branches, err := stack.ReadStackFile(repo.Path)
		if err != nil {
			return err
		}

		for _, branch := range branches[1:] { // We will skip trunk branch
			fmt.Printf("Pushing branch %s\n", branch)
			err = repo.Push(branch)
			if err != nil {
				return err
			}
		}

		return nil
	},
}
