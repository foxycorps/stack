package cmd

import (
	"fmt"
	"os"

	"github.com/foxycorps/stack/pkg/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stack",
	Short: "Stack is a tool for managing git workflow",
	Long:  `Stack is a tool for managing git workflow`,
}

func init() {

	rootCmd.AddGroup(
		&cobra.Group{
			ID:    "STACK",
			Title: fmt.Sprintf("%s%sStack Commands:%s", ui.White, ui.Bold, ui.Reset),
		},
	)

	rootCmd.AddGroup(
		&cobra.Group{
			ID:    "INFORMATION",
			Title: fmt.Sprintf("%s%sInformation Commands:%s", ui.White, ui.Bold, ui.Reset),
		},
	)

	rootCmd.AddGroup(
		&cobra.Group{
			ID:    "NAVIGATION",
			Title: fmt.Sprintf("%s%sNavigation Commands:%s", ui.White, ui.Bold, ui.Reset),
		},
	)

	rootCmd.AddGroup(
		&cobra.Group{
			ID:    "COLLABORATION",
			Title: fmt.Sprintf("%s%sCollaboration Commands:%s", ui.White, ui.Bold, ui.Reset),
		},
	)

	rootCmd.AddCommand(StartCmd)
	rootCmd.AddCommand(AddCmd)
	rootCmd.AddCommand(UpCmd)
	rootCmd.AddCommand(DownCmd)
	rootCmd.AddCommand(AuthCmd)
	rootCmd.AddCommand(PushCmd)
	rootCmd.AddCommand(ListCmd)
	rootCmd.AddCommand(CommitCmd)
	rootCmd.AddCommand(SyncCmd)
	rootCmd.AddCommand(UseCmd)
	rootCmd.AddCommand(CheckoutCmd)
	rootCmd.AddCommand(SubmitCmd)
	rootCmd.AddCommand(ReviewCmd)
	rootCmd.AddCommand(HomeCmd)

	rootCmd.AddCommand(ParentsCmd)

	rootCmd.SetUsageTemplate(ui.ColorHeadings(rootCmd.UsageTemplate()))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
