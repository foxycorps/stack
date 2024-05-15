package cmd

import (
	"fmt"
	"slices"

	"github.com/foxycorps/stack/pkg/git"
	"github.com/foxycorps/stack/pkg/stack"
	"github.com/foxycorps/stack/pkg/ui"
	"github.com/spf13/cobra"
)

var CheckoutCmd = &cobra.Command{
	Use:     "checkout",
	Aliases: []string{"co"},
	Short:   "Checkout a branch",
	GroupID: "NAVIGATION",
	RunE: func(cmd *cobra.Command, args []string) error {

		repo, err := git.NewRepository(".")
		if err != nil {
			return err
		}

		knownBranches, err := stack.ReadStackFile(repo.Path)
		if err != nil {
			return err
		}

		var selectedBranch string

		if len(args) == 0 {
			// There was no arg provided... so we will ask them to select a branch
			selectedBranch, err = ui.SelectBranch(knownBranches)
			if err != nil {
				return err
			}
		} else {
			selectedBranch = args[0]
			if !slices.Contains(knownBranches, selectedBranch) {
				return fmt.Errorf("branch %s is not known", selectedBranch)
			}
		}

		err = repo.SwitchBranch(selectedBranch)

		return err
	},
}
