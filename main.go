package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
)

// CommitType represents a conventional commit type
var CommitTypes = []string{
	"feat", "fix", "docs", "style", "refactor", "perf", "test", "chore", "ci",
}

func main() {
	var commitType, scope, description, confirm string

	// Select commit type
	promptType := &survey.Select{
		Message: "Select commit type:",
		Options: CommitTypes,
	}
	survey.AskOne(promptType, &commitType)

	// Enter scope (optional)
	promptScope := &survey.Input{
		Message: "Enter scope (optional):",
	}
	survey.AskOne(promptScope, &scope)

	// Enter commit message
	promptMessage := &survey.Input{
		Message: "Enter commit message:",
	}
	survey.AskOne(promptMessage, &description)

	// Construct commit message
	var commitMsg string
	if scope != "" {
		commitMsg = fmt.Sprintf("%s(%s): %s", commitType, scope, description)
	} else {
		commitMsg = fmt.Sprintf("%s: %s", commitType, description)
	}

	fmt.Printf("\nGenerated commit message:\n%s\n", commitMsg)

	// Confirm commit
	promptConfirm := &survey.Select{
		Message: "Commit message looks good?",
		Options: []string{"Yes", "No"},
	}
	survey.AskOne(promptConfirm, &confirm)

	if confirm == "No" {
		fmt.Println("Commit aborted.")
		os.Exit(0)
	}

	// Run git commit
	cmd := exec.Command("git", "commit", "-m", commitMsg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Commit successful!")
}


