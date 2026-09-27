package cmd

import (
	"fmt"
	"os/exec"

	"chronicl/internal/config"
	"chronicl/internal/git"
	"chronicl/internal/message"
	"chronicl/internal/prompt"

	"github.com/spf13/cobra"
)

// RootCmd is the main CLI command
var RootCmd = &cobra.Command{
	Use:   "chronicl",
	Short: "chronicl commit/vcs analysis tool",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		commitType, scope, commitAnnotations, scopeAnotations, messageText := prompt.GetUserInputV2(cfg)
		if commitType == "" || messageText == "" {
			fmt.Println("Commit aborted.")
			return
		}

		commitMsg, err := message.Serialize(commitType, scope, commitAnnotations, scopeAnotations, messageText, cfg)
		if err != nil {
			fmt.Errorf("invalid message\n")
		}
		fmt.Println("\nGenerated commit message:", commitMsg)

		if !cfg.AutoCommit {
			confirm := prompt.Confirm()
			if !confirm {
				fmt.Println("Commit aborted.")
				return
			}
		} else {
			fmt.Println("Autocomitting... (turn this off in your chronicl config if undesired)")
		}

		if cfg.CommitCommandFormat == "" {
			if err := git.Commit(commitMsg); err != nil {
				fmt.Printf("Git commit failed:%s\n", err)
			}
			// TODO: security! sanitize for multiple %s or other specifiers
		} else {
			out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf(cfg.CommitCommandFormat, messageText)).Output()
			if err != nil {
				fmt.Printf("Commit failed: %s\n", err)
			} else {
				fmt.Printf("success\n%s", out)
			}
		}
	},
}
