package prompt

import (
	"chronicl/internal/config"
	"errors"
	"fmt"
	"os"

	"github.com/DaDevFox/hof"
	"github.com/cqroot/prompt"
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
func GetUserInput(promptConfig *config.Config) (string, string, []string, string) {
	display := make([]string, 0)
	dispToKeyMap := make(map[string]string)
	for _, obj := range promptConfig.CommitTypes {
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
		return "", "", nil, ""
	}

	var scope string
	scopeKeys := hof.MapToArray(promptConfig.Scopes, func(s config.MessageFormatterConfig) string {
		return s.Key
	})

	if len(scopeKeys) > 0 {
		scopeKeys = append(scopeKeys, "(custom)", "(none)")
		scope, err = prompt.New().Ask("Select scope (optional):").Choose(scopeKeys)
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

	annotationsForCurrCommit := []string{}
	annotationConfigsForCurrCommit, present := promptConfig.CommitTypeToAnnotations[commitType]
	if present {
		annotationsForCurrCommit, err = prompt.New().Ask("Annotate (optional):").MultiChoose(hof.MapToArray(annotationConfigsForCurrCommit, func(c config.MessageFormatterConfig) string {
			return c.Key
		}))
	}

	annotationsForCurrScope := []string{}
	annotationConfigsForCurrScope, present := promptConfig.ScopeToAnnotations[commitType]
	if present {
		annotationsForCurrScope, err = prompt.New().Ask("Annotate (optional):").MultiChoose(hof.MapToArray(annotationConfigsForCurrScope, func(c config.MessageFormatterConfig) string {
			return c.Key
		}))
	}

	message, err := prompt.New().Ask("Enter commit message").Input("bleh bleh bleh")
	Check(err)
	if message == "" && promptConfig.AbortOnEmptyCommit {
		fmt.Println("Empty commit message. Aborting.")
		return "", "", nil, ""
	}

	return commitType, scope, append(annotationsForCurrCommit, annotationsForCurrScope...), message
}

func Confirm() bool {
	chose, err := prompt.New().Ask("Message ok?").Choose([]string{"yes", "no"})
	Check(err)
	return chose == "yes"
}
