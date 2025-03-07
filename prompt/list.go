package prompt

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ListPrompt represents a list selection prompt
type ListPrompt struct {
	*BasePrompt
	Options     []string
	Selected    int
	PageSize    int
	ShowNumbers bool
}

// NewList creates a new list selection prompt
func NewList(question string, options []string) *ListPrompt {
	return &ListPrompt{
		BasePrompt:  NewBasePrompt(question),
		Options:     options,
		Selected:    0,
		PageSize:    7,
		ShowNumbers: true,
	}
}

type listModel struct {
	prompt *ListPrompt
	cursor int
	offset int
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.offset {
					m.offset = m.cursor
				}
			}
		case "down", "j":
			if m.cursor < len(m.prompt.Options)-1 {
				m.cursor++
				if m.cursor >= m.offset+m.prompt.PageSize {
					m.offset = m.cursor - m.prompt.PageSize + 1
				}
			}
		case "enter":
			return m, tea.Quit
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m listModel) View() string {
	var b strings.Builder

	// Render question
	question := m.prompt.Style.QuestionStyle.Render(m.prompt.Question)
	b.WriteString(fmt.Sprintf("%s\n", question))

	// Calculate visible options
	endIdx := m.offset + m.prompt.PageSize
	if endIdx > len(m.prompt.Options) {
		endIdx = len(m.prompt.Options)
	}
	visibleOptions := m.prompt.Options[m.offset:endIdx]

	// Render options
	for i, option := range visibleOptions {
		cursor := " "
		if m.offset+i == m.cursor {
			cursor = ">"
		}

		prefix := ""
		if m.prompt.ShowNumbers {
			prefix = fmt.Sprintf("%d) ", m.offset+i+1)
		}

		optionStyle := m.prompt.Style.UnfocusedStyle
		if m.offset+i == m.cursor {
			optionStyle = m.prompt.Style.FocusedStyle
		}

		b.WriteString(fmt.Sprintf("%s %s%s\n", cursor, prefix, optionStyle.Render(option)))
	}

	return b.String()
}

// Run starts the list selection prompt and returns the selected option
func (p *ListPrompt) Run() (interface{}, error) {
	if len(p.Options) == 0 {
		return "", fmt.Errorf("no options provided")
	}

	model := listModel{
		prompt: p,
		cursor: p.Selected,
	}

	p.Program = tea.NewProgram(model)
	m, err := p.Program.Run()
	if err != nil {
		return "", err
	}

	finalModel := m.(listModel)
	return p.Options[finalModel.cursor], nil
}

// GetQuestion returns the prompt question
func (p *ListPrompt) GetQuestion() string {
	return p.Question
}

// Validate is not used for list prompts
func (p *ListPrompt) Validate(input string) error {
	return nil
}
