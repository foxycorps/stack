package ui

import (
	"fmt"
	"strings"
)

type BranchInformation struct {
	Name        string
	Number      int
	HasPR       bool
	NeedsReview bool
	Owner       string
	IsApproved  bool
	Approver    string
}

func ListStack(currentBranch string, branches []BranchInformation, stackOwner string) error {
	ownerLink := CreateHyperlink(fmt.Sprintf("https://github.com/%s", stackOwner), stackOwner)
	stackName := strings.Split(currentBranch, "_")[0]
	fmt.Printf("Current stack: %s%s%s%s\t[Owner: %s%s%s]\n", White, Bold, stackName, Reset, PastelPurple, ownerLink, Reset)
	for _, branchInfo := range branches {
		var prefix string = "    "
		var Color string = Dim
		if branchInfo.Name == currentBranch {
			prefix = "  * "
			Color = White
		}
		var prStatus = ""
		if branchInfo.HasPR {
			if branchInfo.NeedsReview {
				prStatus = fmt.Sprintf("[%s%s%s%s]", PastelRed, Bold, "Needs review", Reset)
			} else if branchInfo.IsApproved {
				prStatus = fmt.Sprintf("[%s%s%s%s by %s%s%s]", PastelGreen, Bold, "Approved", Reset, PastelPurple, CreateHyperlink(fmt.Sprintf("https://github.com/%s", branchInfo.Approver), branchInfo.Approver), Reset)
			} else {
				prStatus = fmt.Sprintf("[PR #%s%s%d%s]", PastelGreen, Bold, branchInfo.Number, Reset)
			}
		}
		cleanedName := strings.Replace(branchInfo.Name, ".", "/", 1)
		fmt.Printf("%s%s%s%s\t%s\n", Color, prefix, cleanedName, Reset, prStatus)
	}

	return nil
}
