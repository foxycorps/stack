package cmd

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/foxycorps/stack/pkg/git"
	"github.com/foxycorps/stack/pkg/stack"
	"github.com/foxycorps/stack/pkg/ui"
	"github.com/spf13/cobra"
)

var UseCmd = &cobra.Command{
	Use:     "use [name]",
	Short:   "Fetch a remote stack",
	GroupID: "COLLABORATION",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := git.NewRepository(".")
		if err != nil {
			return err
		}

		// Ensuring we are up to date with the remote.
		err = repo.FetchAll()
		if err != nil {
			return err
		}

		// Getting the stack name to use.
		var stackName string
		if len(args) == 1 {
			stackName = args[0]
		}

		// If the stackName is not provided, we will make a call and get all the `*_base` branches.
		// We will then prompt the user to select the stack they want to use.
		if stackName == "" {
			// We will list all the remote stacks
			remoteStacks, err := repo.FindRemoteStacks()
			if err != nil {
				return err
			}
			stackName, err = ui.SelectStack(remoteStacks)
			if err != nil {
				return err
			}
		}
		fmt.Println("Stack name:", stackName)
		// Getting all remote branches for the stack.
		branchNames, err := repo.ListRemoteBranches(stackName)
		if err != nil {
			return err
		}
		fmt.Println("Branch names:", branchNames)

		// Starting the stack file.
		err = stack.InitializeStackFile(repo.Path, "main")
		if err != nil {
			return err
		}

		// Sorting the branches into the correct order (flipping them)
		sort.Slice(branchNames, func(i, j int) bool {
			orderI := getOrderFromBranchName(branchNames[i])
			orderJ := getOrderFromBranchName(branchNames[j])
			return orderI < orderJ
		})

		// Tracking each remote branch locally.
		for _, branchName := range branchNames {
			cleanedName := strings.TrimPrefix(branchName, "origin/")
			exists, err := repo.BranchExists(cleanedName)
			if err != nil {
				return err
			}
			// Branch doesn't exist locally, we will track it.
			if !exists {
				err = repo.TrackRemoteBranch(cleanedName)
				if err != nil {
					return err
				}
			} else {
				// Pulling the latest code.
				err = repo.Pull(cleanedName)
				if err != nil {
					return err
				}
			}

			// Updating the stack file with the tracked branch.
			err = stack.UpdateStackFile(repo.Path, cleanedName)
			if err != nil {
				return err
			}
		}

		// Switching to the first branch in the stack.
		repo.SwitchBranch(fmt.Sprintf("%s.base", stackName))

		fmt.Printf("Now using '%s%s%s' stack\n", ui.PastelBlue, stackName, ui.Reset)

		fmt.Println("Base branch:", fmt.Sprintf("%s%s%s", ui.PastelBlue, stackName+".base", ui.Reset))
		fmt.Println("Found branches:", branchNames)
		return nil
	},
}

func getOrderFromBranchName(branchName string) int {
	parts := strings.Split(branchName, "_")
	if len(parts) < 2 {
		return 0
	}
	order, _ := strconv.Atoi(parts[len(parts)-1])
	return order
}
