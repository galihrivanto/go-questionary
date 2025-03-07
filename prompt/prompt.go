package prompt

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Style defines the visual appearance of prompts
type Style struct {
	QuestionStyle  lipgloss.Style
	AnswerStyle    lipgloss.Style
	SelectionStyle lipgloss.Style
	ErrorStyle     lipgloss.Style
	FocusedStyle   lipgloss.Style
	UnfocusedStyle lipgloss.Style
}

// DefaultStyle returns the default styling for prompts
func DefaultStyle() Style {
	return Style{
		QuestionStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("69")),
		AnswerStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("39")),
		SelectionStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("212")),
		ErrorStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
		FocusedStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("212")),
		UnfocusedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
	}
}

// Prompt represents a generic command line prompt
type Prompt interface {
	// Run starts the prompt and returns the result
	Run() (interface{}, error)

	// GetQuestion returns the prompt question
	GetQuestion() string

	// Validate checks if the input is valid
	Validate(input string) error
}

// BasePrompt contains common functionality for all prompts
type BasePrompt struct {
	Question  string
	Default   string
	Validator func(string) error
	Style     Style
	Program   *tea.Program
}

// NewBasePrompt creates a new base prompt with default styling
func NewBasePrompt(question string) *BasePrompt {
	return &BasePrompt{
		Question: question,
		Style:    DefaultStyle(),
	}
}
