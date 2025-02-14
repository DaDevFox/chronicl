package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"chronicl/internal/config"
	"chronicl/internal/git"
	"chronicl/internal/prompt"
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
		commitMsg := fmt.Sprintf("%s(%s): %s", commitType, scope, message)

		fmt.Println("\nGenerated commit message:", commitMsg)

		if !cfg.AutoCommit {
			var confirm string
			survey.AskOne(&survey.Select{
				Message: "Commit message looks good?",
				Options: []string{"Yes", "No"},
			}, &confirm)

			if confirm == "No" {
				fmt.Println("Commit aborted.")
				return
			}
		}

		if err := git.Commit(commitMsg); err != nil {
			fmt.Println("Git commit failed:", err)
		} else {
			fmt.Println("Commit successful!")
		}
	},
}

