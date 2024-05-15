package cmd

import (
	"fmt"

	"github.com/foxycorps/stack/pkg/git"
	"github.com/foxycorps/stack/pkg/github"
	"github.com/foxycorps/stack/pkg/stack"
	"github.com/foxycorps/stack/pkg/ui"
	"github.com/spf13/cobra"
)

var SubmitCmd = &cobra.Command{
	Use:     "submit",
	Short:   "Submit stack as pull requests",
	GroupID: "COLLABORATION",
	RunE: func(cmd *cobra.Command, args []string) error {

		repo, err := git.NewRepository(".")
		if err != nil {
			return err
		}
		client := github.NewClient(repo)
		branches, err := stack.ReadStackFile(repo.Path)
		if err != nil {
			return err
		}

		parentMap, err := stack.GetParents(repo.Path)
		if err != nil {
			return err
		}

		for _, branch := range branches[1:] {
			pr, err := client.FindPrForBranch(branch)
			if err != nil {
				return err
			}
			if pr != nil {
				continue
			}
			markCommit, err := repo.GetFirstCommitNotOnParent(branch, parentMap[branch])
			if err != nil {
				return err
			}
			firstCommit, err := repo.GetCommit(markCommit)
			if err != nil {
				return err
			}

			prDetails := *ui.CreatePR(branch, firstCommit)
			if !prDetails.Create {
				continue
			}

			pr, err = client.CreatePullRequest(branch, parentMap[branch], prDetails.Title, prDetails.Body)
			if err != nil {
				return err
			}

			fmt.Printf("Submitting PR for '%s%s%s'... Success! View it here: %s\n", ui.White, branch, ui.Reset, pr.GetHTMLURL())
		}

		return nil
	},
}
