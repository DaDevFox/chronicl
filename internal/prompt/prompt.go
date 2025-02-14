package prompt

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// GetUserInput prompts for commit type, scope, and message
func GetUserInput(commitTypes, scopes []string) (string, string, string) {
	var commitType, scope, message string

	survey.AskOne(&survey.Select{
		Message: "Select commit type:",
		Options: commitTypes,
	}, &commitType)

	if len(scopes) > 0 {
		survey.AskOne(&survey.Select{
			Message: "Select scope (optional):",
			Options: append(scopes, "(custom)"),
		}, &scope)

		if scope == "(custom)" {
			survey.AskOne(&survey.Input{Message: "Enter custom scope:"}, &scope)
		}
	} else {
		survey.AskOne(&survey.Input{Message: "Enter scope (optional):"}, &scope)
	}

	survey.AskOne(&survey.Input{Message: "Enter commit message:"}, &message)

	return commitType, scope, message
}

