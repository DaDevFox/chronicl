package prompt

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/cqroot/prompt/choose"
	"github.com/cqroot/prompt/input"

	"chronicl/internal/config"
	"chronicl/internal/message"

	"github.com/DaDevFox/hof"
)

type stage int

const (
	stageCommitType stage = iota
	stageScope
	stageCustomScope
	stageCommitAnnotations
	stageScopeAnnotations
	stageMessage
	stageDone
)

type commitForm struct {
	stage  stage
	config *config.Config

	commitType  *choose.Model
	promptScope bool
	scope       *choose.Model
	customScope *input.Model

	promptCommitAnnotations bool
	commitAnnotations       *choose.Model
	promptScopeAnotations   bool
	scopeAnnotations        *choose.Model

	message *input.Model

	// The actual values accumulated so far.
	commitTypeValue string
	scopeValue      string

	commitAnnotationsValue string
	scopeAnnotationsValue  string

	messageValue string

	result commitResult
	err    error

	width int
}

type commitResult struct {
	commitType           string
	scope                string
	scopeAnnotation      string
	commitTypeAnnotation string
	message              string
}

var (
	previewStyle = lipgloss.NewStyle().
			Padding(1, 1)

	activeLabelStyle = lipgloss.NewStyle().
				Bold(true)

	dimStyle = lipgloss.NewStyle().
			Faint(true)
)

func newCommitForm(cfg *config.Config) *commitForm {
	// Commit types.
	commitChoices := hof.MapToArray(cfg.CommitTypes, func(c config.MessageFormatterConfig) choose.Choice {
		return choose.Choice{
			Text: c.Key,
			Note: c.Description,
		}
	})

	// Scope choices.
	scopeChoices := hof.MapToArray(cfg.Scopes, func(c config.MessageFormatterConfig) choose.Choice {
		return choose.Choice{
			Text: c.Key,
			Note: c.Description,
		}
	})

	if len(scopeChoices) > 0 {
		scopeChoices = append(scopeChoices, choose.Choice{Text: "(custom)", Note: "write free text scope not in the list"},
			choose.Choice{Text: "(none)", Note: "skip scope"})
	}

	// These are initially empty. We replace the relevant models when the
	// corresponding stage becomes active.
	m := &commitForm{
		stage:  stageCommitType,
		config: cfg,

		commitType: choose.New(commitChoices),

		scope: choose.New(scopeChoices),

		customScope: input.New(""),

		message: input.New(""),

		// Placeholder models; replaced on stage transitions.
		commitAnnotations: choose.NewWithStrings([]string{}),
		scopeAnnotations:  choose.NewWithStrings([]string{}),

		result: commitResult{},

		width: 0,
	}

	// Determine whether scope stage should be shown at all.
	m.promptScope = len(scopeChoices) > 0

	return m
}

// buildCommitAnnotations (re)builds the commit-annotation model based on the
// currently selected commit type. Call this only on stage transitions, never
// inside Update().
func (m *commitForm) buildCommitAnnotations() {
	annotations, prompt := m.config.CommitTypeToAnnotations[m.commitTypeValue]
	m.promptCommitAnnotations = prompt
	m.commitAnnotations = choose.New(hof.MapToArray(
		annotations,
		func(c config.MessageFormatterConfig) choose.Choice {
			return choose.Choice{
				Text: c.Key,
				Note: c.Description,
			}
		},
	))
}

// buildScopeAnnotations (re)builds the scope-annotation model based on the
// currently selected commit type (scope annotations are keyed by commit type
// in the config).
func (m *commitForm) buildScopeAnnotations() {
	annotations, prompt := m.config.ScopeToAnnotations[m.commitTypeValue]
	m.promptScopeAnotations = prompt
	m.scopeAnnotations = choose.New(hof.MapToArray(
		annotations,
		func(c config.MessageFormatterConfig) choose.Choice {
			return choose.Choice{
				Text: c.Key,
				Note: c.Description,
			}
		},
	))
}

func (m commitForm) Init() tea.Cmd {
	return m.commitType.Init()
}

func (m commitForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = sizeMsg.Width
	}

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "ctrl+c", "esc", "q":
			m.err = tea.ErrProgramKilled
			return m, tea.Quit

		case "enter":
			return m.confirm()
		}
	}

	var cmd tea.Cmd
	switch m.stage {
	case stageCommitType:
		var model tea.Model
		model, cmd = m.commitType.Update(msg)
		*m.commitType = model.(choose.Model)

		m.commitTypeValue = m.commitType.Data()

	case stageScope:
		var model tea.Model
		model, cmd = m.scope.Update(msg)
		*m.scope = model.(choose.Model)

		m.scopeValue = m.scope.Data()

	case stageCustomScope:
		var model tea.Model
		model, cmd = m.customScope.Update(msg)
		*m.customScope = model.(input.Model)

		m.scopeValue = m.customScope.Data()

	case stageCommitAnnotations:
		var model tea.Model
		model, cmd = m.commitAnnotations.Update(msg)
		*m.commitAnnotations = model.(choose.Model)

		m.commitAnnotationsValue = m.commitAnnotations.Data()

	case stageScopeAnnotations:
		var model tea.Model
		model, cmd = m.scopeAnnotations.Update(msg)
		*m.scopeAnnotations = model.(choose.Model)

		m.scopeAnnotationsValue = m.scopeAnnotations.Data()

	case stageMessage:
		var model tea.Model
		model, cmd = m.message.Update(msg)
		*m.message = model.(input.Model)

		m.messageValue = m.message.Data()
	}

	return m, cmd
}

func (m commitForm) confirm() (tea.Model, tea.Cmd) {
	switch m.stage {
	case stageCommitType:
		m.commitTypeValue = m.commitType.Data()

		if m.commitTypeValue == "" {
			return m, nil
		}

		// Build annotation models now that we know the commit type.
		m.buildCommitAnnotations()
		m.buildScopeAnnotations()

		if m.promptScope {
			m.stage = stageScope
		} else {
			m.stage = m.nextStageAfterScope()
		}

	case stageScope:
		m.scopeValue = m.scope.Data()

		switch m.scopeValue {
		case "(custom)":
			m.scopeValue = ""
			m.stage = stageCustomScope

		case "(none)":
			m.scopeValue = ""
			m.stage = m.nextStageAfterScope()

		default:
			m.stage = m.nextStageAfterScope()
		}

	case stageCustomScope:
		m.scopeValue = m.customScope.Data()
		m.stage = m.nextStageAfterScope()

	case stageCommitAnnotations:
		m.commitAnnotationsValue = m.commitAnnotations.Data()
		if m.promptScopeAnotations {
			m.stage = stageScopeAnnotations
		} else {
			m.stage = stageMessage
		}

	case stageScopeAnnotations:
		m.scopeAnnotationsValue = m.scopeAnnotations.Data()
		m.stage = stageMessage

	case stageMessage:
		m.messageValue = m.message.Data()

		m.result = commitResult{
			commitType:           m.commitTypeValue,
			scope:                m.scopeValue,
			commitTypeAnnotation: m.commitAnnotationsValue,
			scopeAnnotation:      m.scopeAnnotationsValue,
			message:              m.messageValue,
		}

		m.stage = stageDone
		return m, tea.Quit
	}

	return m, m.enterCurrentStage()
}

func (m *commitForm) nextStageAfterScope() stage {
	if m.promptCommitAnnotations {
		return stageCommitAnnotations
	}
	if m.promptScopeAnotations {
		return stageScopeAnnotations
	}

	return stageMessage
}

func (m *commitForm) enterCurrentStage() tea.Cmd {
	switch m.stage {
	case stageCommitType:
		return m.commitType.Init()

	case stageScope:
		return m.scope.Init()

	case stageCustomScope:
		return m.customScope.Init()

	case stageCommitAnnotations:
		return m.commitAnnotations.Init()

	case stageScopeAnnotations:
		return m.scopeAnnotations.Init()

	case stageMessage:
		return m.message.Init()
	}

	return nil
}

func (m commitForm) View() string {
	var b strings.Builder

	// Wrap the preview at the terminal width (minus padding/border).
	previewWidth := m.width - 6 // border (2) + padding (4)
	if previewWidth < 15 {
		previewWidth = 15 // sane minimum
	}

	preview := previewStyle.Width(previewWidth).Render(m.preview())

	b.WriteString(preview)
	b.WriteString("\n\n")

	switch m.stage {
	case stageCommitType:
		b.WriteString(activeLabelStyle.Render("Commit type"))
		b.WriteString("\n")
		b.WriteString(m.commitType.View())

	case stageScope:
		b.WriteString(activeLabelStyle.Render("Scope"))
		b.WriteString("\n")
		b.WriteString(m.scope.View())

	case stageCustomScope:
		b.WriteString(activeLabelStyle.Render("Custom scope"))
		b.WriteString("\n")
		b.WriteString(m.customScope.View())

	case stageCommitAnnotations:
		b.WriteString(activeLabelStyle.Render("Commit annotations"))
		b.WriteString("\n")
		b.WriteString(m.commitAnnotations.View())

	case stageScopeAnnotations:
		b.WriteString(activeLabelStyle.Render("Scope annotations"))
		b.WriteString("\n")
		b.WriteString(m.scopeAnnotations.View())

	case stageMessage:
		b.WriteString(activeLabelStyle.Render("Commit message"))
		b.WriteString("\n")
		b.WriteString(m.message.View())
	}

	b.WriteString("\n\n")
	b.WriteString(dimStyle.Render("↑/↓ move • space select • enter confirm • esc quit"))

	return b.String()
}

// TODO: add {author} substitution
func (m commitForm) subHelper(format string) string {
	ret := format
	if strings.Trim(m.scopeValue, " ") != "" {
		ret = strings.ReplaceAll(ret, "{scope}", m.scopeValue)
	} else {
		// TODO: find containing [] surrounding {scope} and delete it
	}
	if strings.Trim(m.messageValue, " ") != "" {
		ret = strings.ReplaceAll(ret, "{message}", m.messageValue)
	} else {
		// TODO: ditto
	}

	return ret
}

func (m commitForm) subWithConfigFromList(data []config.MessageFormatterConfig, matching string) (string, error) {
	datum, found := hof.Find(data, func(c config.MessageFormatterConfig) bool { return c.Key == matching })
	if !found && strings.Trim(matching, " ") != "" {
		return "", fmt.Errorf("couldn't find datum config for datum selected among config data it was supposedly selected from; datum: %s, data: %s", matching, data)
	}
	return m.subHelper(datum.HelperFormatting), nil
}

func (m commitForm) preview() string {
	ret, err := message.Serialize(
		m.commitTypeValue,
		m.scopeValue,
		m.scopeAnnotationsValue,
		m.commitAnnotationsValue,
		m.messageValue,
		m.config,
	)

	// TODO: pass err up as validation hint, block commit until validated correctly
	if err != nil {
		return fmt.Sprintf("error: %s", err)
	}

	retHelpers := []string{}
	toAppend, err := m.subWithConfigFromList(m.config.CommitTypes, m.commitTypeValue)
	if err == nil {

		retHelpers = append(retHelpers, toAppend)
	}
	toAppend, err = m.subWithConfigFromList(m.config.Scopes, m.scopeValue)
	if err == nil {
		retHelpers = append(retHelpers, toAppend)
	}

	commitTypeAnotations, ok := m.config.CommitTypeToAnnotations[m.commitTypeValue]
	if ok {
		toAppend, err = m.subWithConfigFromList(commitTypeAnotations, m.commitAnnotationsValue)
		if err == nil {
			retHelpers = append(retHelpers, toAppend)
		}
	}

	scopeAnnotationsValue, ok := m.config.ScopeToAnnotations[m.scopeValue]
	if ok {
		toAppend, err = m.subWithConfigFromList(scopeAnnotationsValue, m.scopeAnnotationsValue)
		if err == nil {
			retHelpers = append(retHelpers, toAppend)
		}
	}

	retHelpers = hof.FilterToArray(retHelpers, func(s string) bool {
		return strings.Trim(s, " ") != ""
	})
	return fmt.Sprintf("%s\n\n%s", ret, strings.Join(retHelpers, "\n"))
}

func GetUserInputV2(promptConfig *config.Config) (
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
