package cmd

import (
	"fmt"

	"github.com/foxycorps/stack/pkg/git"
	"github.com/foxycorps/stack/pkg/stack"
	"github.com/spf13/cobra"
)

var ParentsCmd = &cobra.Command{
	Use:   "parents",
	Short: "List the whole stack as a parent map",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := git.NewRepository(".")
		if err != nil {
			return err
		}

		branchParentMap, err := stack.GetParents(repo.Path)
		if err != nil {
			return err
		}

		fmt.Println("Branch Parent Map:", branchParentMap)

		return nil
	},
}
