package prompt

import (
	"chronicl/internal/config"
	"errors"
	"fmt"
	"os"

	"github.com/cqroot/prompt"
	// "github.com/cqroot/prompt/choose"
	// "github.com/cqroot/prompt/input"
)

func Check(err error) {
	if err != nil {
		if errors.Is(err, prompt.ErrUserQuit) {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		} else {
			panic(err)
		}
	}
}

// GetUserInput prompts for commit type, scope, and message
func GetUserInput(commitTypes []config.CommitType, scopes []string) (string, string, string) {
	display := make([]string, 0)
	dispToKeyMap := make(map[string]string)
	for _, obj := range commitTypes {
		str := obj.Key + " | " + obj.Description
		display = append(display, str)
		dispToKeyMap[str] = obj.Key
	}

	commitType, err := prompt.New().Ask("Select commit type:").Choose(display)
	Check(err)
	commitType = dispToKeyMap[commitType]

	// commitType := goprompter.Choose("Select commit type:", commitTypes)
	if commitType == "" {
		fmt.Println("No commit type selected. Aborting.")
		return "", "", ""
	}

	var scope string
	if len(scopes) > 0 {
		scopes = append(scopes, "(custom)", "(none)")
		scope, err = prompt.New().Ask("Select scope (optional):").Choose(scopes)
		Check(err)
		if scope == "(custom)" {
			scope, err = prompt.New().Ask("Enter custom scope:").Input("blah blah")
			Check(err)
		} else if scope == "(none)" {
			scope = ""
		}
	} else {
		scope, err = prompt.New().Ask("Enter scope (optional):").Input("blah blah")
	}

	message, err := prompt.New().Ask("Enter commit message").Input("bleh bleh bleh")
	Check(err)
	if message == "" {
		fmt.Println("Empty commit message. Aborting.")
		return "", "", ""
	}

	return commitType, scope, message
}

func Confirm() bool {
	chose, err := prompt.New().Ask("Message ok?").Choose([]string{"yes", "no"})
	Check(err)
	return chose == "yes"
}
