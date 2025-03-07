package prompt

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ConfirmPrompt represents a yes/no confirmation prompt
type ConfirmPrompt struct {
	*BasePrompt
	DefaultValue bool
}

// NewConfirm creates a new confirmation prompt
func NewConfirm(question string) *ConfirmPrompt {
	return &ConfirmPrompt{
		BasePrompt: NewBasePrompt(question),
	}
}

// Model represents the state of the confirmation prompt
type confirmModel struct {
	prompt *ConfirmPrompt
	value  bool
}

func (m confirmModel) Init() tea.Cmd {
	return nil
}

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			m.value = true
			return m, tea.Quit
		case "n", "N":
			m.value = false
			return m, tea.Quit
		case "enter":
			m.value = m.prompt.DefaultValue
			return m, tea.Quit
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	var b strings.Builder

	question := m.prompt.Style.QuestionStyle.Render(m.prompt.Question)
	defaultValue := "(y/N)"
	if m.prompt.DefaultValue {
		defaultValue = "(Y/n)"
	}

	b.WriteString(fmt.Sprintf("%s %s: ", question, defaultValue))
	return b.String()
}

// Run starts the confirmation prompt and returns the user's choice
func (p *ConfirmPrompt) Run() (interface{}, error) {
	model := confirmModel{
		prompt: p,
		value:  p.DefaultValue,
	}

	p.Program = tea.NewProgram(model)
	m, err := p.Program.Run()
	if err != nil {
		return false, err
	}

	finalModel := m.(confirmModel)
	return finalModel.value, nil
}

// GetQuestion returns the prompt question
func (p *ConfirmPrompt) GetQuestion() string {
	return p.Question
}

// Validate is not used for confirmation prompts
func (p *ConfirmPrompt) Validate(input string) error {
	return nil
}
