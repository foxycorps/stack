package cmd

import (
	"fmt"

	"github.com/foxycorps/stack/pkg/git"
	"github.com/foxycorps/stack/pkg/github"
	"github.com/spf13/cobra"
)

var reject, approve, comment bool

var ReviewCmd = &cobra.Command{
	Use:     "review [comment]",
	Short:   "Review a pull request",
	Args:    cobra.ExactArgs(1),
	GroupID: "COLLABORATION",
	RunE: func(cmd *cobra.Command, args []string) error {

		if !reject && !approve && !comment {
			return fmt.Errorf("must specify --reject, --approve, or --comment")
		}

		repo, err := git.NewRepository(".")
		if err != nil {
			return err
		}

		currentBranch, err := repo.CurrentBranch()
		if err != nil {
			return err
		}

		client := github.NewClient(repo)
		pr, err := client.FindPrForBranch(currentBranch)
		if pr == nil {
			return fmt.Errorf("no pull request found for branch %s", currentBranch)
		}
		if err != nil {
			return err
		}

		sha, _ := repo.FindLastCommit()

		var event string = "COMMENT"

		switch {
		case reject:
			event = "REQUEST_CHANGES"
		case approve:
			event = "APPROVE"
		case comment:
			event = "COMMENT"
		}

		err = client.WriteReviewComment(pr, args[0], event, sha)

		return err
	},
}

func init() {
	ReviewCmd.Flags().BoolVarP(&reject, "reject", "r", false, "Reject the pull request")
	ReviewCmd.Flags().BoolVarP(&approve, "approve", "a", false, "Approve the pull request")
	ReviewCmd.Flags().BoolVarP(&comment, "comment", "c", false, "Comment on the pull request")
}
