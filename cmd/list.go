package cmd

import (
	"strings"

	"github.com/foxycorps/stack/pkg/git"
	"github.com/foxycorps/stack/pkg/github"
	"github.com/foxycorps/stack/pkg/stack"
	"github.com/foxycorps/stack/pkg/ui"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all branches in stack",
	GroupID: "INFORMATION",
	RunE: func(cmd *cobra.Command, args []string) error {

		repo, err := git.NewRepository(".")
		if err != nil {
			return err
		}

		// We are only going to list the branches we actually know of.
		// So we will not do any fetches or anything.
		branches, err := stack.ReadStackFile(repo.Path)
		if err != nil {
			return err
		}

		currentBranch, err := repo.CurrentBranch()
		if err != nil {
			return err
		}

		client := github.NewClient(repo)

		var BranchInfo []ui.BranchInformation
		var stackOwner string
		for _, branch := range branches[1:] {
			// We will try and get a PR for this branch.
			pr, err := client.FindPrForBranch(branch)
			if err != nil {
				return err
			}

			owner, err := repo.GetOwner(branch)
			if err != nil {
				return err
			}

			if strings.HasSuffix(branch, ".base") {
				stackOwner = owner
			}

			NewBranchInfo := ui.BranchInformation{
				Name:        branch,
				HasPR:       pr != nil,
				NeedsReview: false, // TODO: revert this and fix the `IsPRAwaitingReview` function.
				Owner:       owner,
			}

			if pr != nil {
				NewBranchInfo.Number = pr.GetNumber()
				approved, approver, err := client.IsPRApproved(pr)
				if err != nil {
					return err
				}
				if approved {
					NewBranchInfo.IsApproved = true
					NewBranchInfo.Approver = approver
				}
			}

			BranchInfo = append(BranchInfo, NewBranchInfo)
		}

		err = ui.ListStack(currentBranch, BranchInfo, stackOwner)
		if err != nil {
			return err
		}

		repo.SwitchBranch(currentBranch) // Switch back to the branch we were on.

		return nil
	},
}
