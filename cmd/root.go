package cmd

import (
	"fmt"

	"chronicl/internal/config"
	"chronicl/internal/git"
	"chronicl/internal/prompt"
	"github.com/spf13/cobra"
)

// RootCmd is the main CLI command
var RootCmd = &cobra.Command{
	Use:   "rgrc",
	Short: "rgrc - Fast and secure commit tool",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		commitType, scope, message := prompt.GetUserInput(cfg.CommitTypes, cfg.Scopes)
		if commitType == "" || message == "" {
			fmt.Println("Commit aborted.")
			return
		}

		commitMsg := fmt.Sprintf("%s(%s): %s", commitType, scope, message)
		fmt.Println("\nGenerated commit message:", commitMsg)

		if !cfg.AutoCommit {
			confirm := prompt.Confirm()
			if !confirm {
				fmt.Println("Commit aborted.")
				return
			}
		} else {
			fmt.Println("Autocomitting... (turn this off in your chronicl config)")
		}

		if err := git.Commit(commitMsg); err != nil {
			fmt.Println("Git commit failed:", err)
		} else {
			fmt.Println("Commit successful!")
		}
	},
}
