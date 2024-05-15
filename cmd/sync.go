package cmd

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/foxycorps/stack/pkg/git"
	"github.com/foxycorps/stack/pkg/stack"
	"github.com/spf13/cobra"
)

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncs the local database with the remote database",
	RunE: func(cmd *cobra.Command, args []string) error {

		repo, err := git.NewRepository(".")
		if err != nil {
			return err
		}

		err = repo.FetchAll()
		if err != nil {
			return err
		}

		currentBranch, err := repo.CurrentBranch()
		if err != nil {
			return err
		}

		stackName := strings.Split(currentBranch, ".")[0]

		// Getting the branches we know.
		branches, err := stack.ReadStackFile(repo.Path)
		if err != nil {
			return err
		}

		// Getting the branches on remote (incase they added to remote)
		remoteBranches, err := repo.ListRemoteBranches(stackName)
		if err != nil {
			return err
		}

		for _, remoteBranch := range remoteBranches {
			cleanedName := strings.TrimPrefix(remoteBranch, "origin/")
			if slices.Contains(branches, cleanedName) {
				continue
			}
			err = repo.TrackRemoteBranch(cleanedName)
			if err != nil {
				return err
			}
			err = stack.UpdateStackFile(repo.Path, cleanedName)
			if err != nil {
				return err
			}
		}

		sort.Slice(branches, func(i, j int) bool {
			orderI := getOrderFromBranchName(branches[i])
			orderJ := getOrderFromBranchName(branches[j])
			return orderI < orderJ
		})

		var prev string = branches[0] // Previous branch
		for _, branch := range branches[1:] {
			if err := repo.RebaseBranch(branch, prev); err != nil {
				return fmt.Errorf("failed to rebase branch '%s' onto '%s': %v\n", branch, prev, err)
			}
			prev = branch
		}

		// Getting the branches that are not on remote anymore
		// Finally switching back to the original branch.
		err = repo.SwitchBranch(currentBranch)

		return err
	},
}
