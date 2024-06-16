package main

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FormField struct {
	textinput.Model

	Label string
}

func NewFormField(
	label string,
	prompt string,
	placeholder string,
	value string,
) FormField {
	m := textinput.New()
	m.CharLimit = 250
	m.Cursor.BlinkSpeed = time.Second / 3
	m.Placeholder = placeholder
	m.Prompt = prompt
	m.Width = 30
	m.SetValue(value)

	return FormField{
		Model: m,

		Label: label,
	}
}

func (f FormField) Update(msg tea.Msg) (model FormField, cmd tea.Cmd) {
	f.Model, cmd = f.Model.Update(msg)
	return f, cmd
}

func (f FormField) View() string {
	elements := make([]string, 0)

	if f.Label != "" {
		elements = append(elements, lipgloss.NewStyle().Width(20).Render(f.Label))
	}

	elements = append(elements, f.Model.View())

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		elements...,
	)
}
