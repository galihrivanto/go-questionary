package prompt

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// TextPrompt represents a text input prompt
type TextPrompt struct {
	*BasePrompt
	textInput textinput.Model
}

// NewText creates a new text input prompt
func NewText(question string) *TextPrompt {
	ti := textinput.New()
	ti.Focus()

	return &TextPrompt{
		BasePrompt: NewBasePrompt(question),
		textInput:  ti,
	}
}

// Model represents the state of the text prompt
type textModel struct {
	prompt    *TextPrompt
	textInput textinput.Model
	err       error
}

func (m textModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m textModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			value := m.textInput.Value()
			if m.prompt.Validator != nil {
				if err := m.prompt.Validator(value); err != nil {
					m.err = err
					return m, nil
				}
			}
			return m, tea.Quit
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m textModel) View() string {
	var b strings.Builder

	question := m.prompt.Style.QuestionStyle.Render(m.prompt.Question)
	b.WriteString(fmt.Sprintf("%s\n", question))

	if m.err != nil {
		b.WriteString(m.prompt.Style.ErrorStyle.Render(m.err.Error()) + "\n")
	}

	b.WriteString(m.textInput.View())
	return b.String()
}

// Run starts the text prompt and returns the user's input
func (p *TextPrompt) Run() (interface{}, error) {
	if p.Default != "" {
		p.textInput.SetValue(p.Default)
	}

	model := textModel{
		prompt:    p,
		textInput: p.textInput,
	}

	p.Program = tea.NewProgram(model)
	m, err := p.Program.Run()
	if err != nil {
		return "", err
	}

	finalModel := m.(textModel)
	return finalModel.textInput.Value(), nil
}

// GetQuestion returns the prompt question
func (p *TextPrompt) GetQuestion() string {
	return p.Question
}

// Validate checks if the input is valid
func (p *TextPrompt) Validate(input string) error {
	if p.Validator != nil {
		return p.Validator(input)
	}
	return nil
}
