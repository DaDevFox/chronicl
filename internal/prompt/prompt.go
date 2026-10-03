package prompt

import (
	"fmt"

	"chronicl/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

func GetUserInput(promptConfig *config.Config) (
	string,
	string,
	string,
	string,
	string,
) {
	form := newCommitForm(promptConfig)

	finalModel, err := tea.NewProgram(form).Run()
	Check(err)

	result := finalModel.(commitForm).result

	if result.commitType == "" {
		fmt.Println("No commit type selected. Aborting.")
		return "", "", "", "", ""
	}

	if result.message == "" && promptConfig.AbortOnEmptyCommit {
		fmt.Println("Empty commit message. Aborting.")
		return "", "", "", "", ""
	}

	return result.commitType,
		result.scope,
		result.commitTypeAnnotation,
		result.scopeAnnotation,
		result.message
}
