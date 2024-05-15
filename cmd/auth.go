package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/foxycorps/stack/pkg/git"
	"github.com/foxycorps/stack/pkg/github"
	"github.com/foxycorps/stack/pkg/keyring"
	"github.com/spf13/cobra"
)

var AuthCmd = &cobra.Command{
	Use:   "auth [token]",
	Short: "Authenticate with github",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var token string

		if len(args) > 0 {
			token = args[0]
		} else {
			// Read token from stdin
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				bytes, _ := io.ReadAll(os.Stdin)
				token = strings.TrimSpace(string(bytes))
			} else {
				return fmt.Errorf("no token provided")
			}
		}

		err := keyring.Set(token)
		if err != nil {
			return fmt.Errorf("failed to set github token: %v", err)
		}

		repo, err := git.NewRepository(".")
		if err != nil {
			return err
		}

		client := github.NewClient(repo)
		if !client.ValidateToken() {
			// We are going to remove the token from the keyring if it is invalid.
			keyring.Delete()
			return fmt.Errorf("invalid github token")
		}

		fmt.Println("Authenticated with github.")

		return nil
	},
}
