package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
)

type controlsFooterProps struct {
}

type controlsFooterComponent struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[controlsFooterProps]
}

func (c *controlsFooterComponent) Init(props *controlsFooterProps) tea.Cmd {
	return nil
}

func (c *controlsFooterComponent) Render(width, height int) string {
	dot := TextDisabled.Render(" • ")

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Left, " ",
			TextMain.Render("Ctrl+C"), " ", TextDark.Render("Quit"),
			dot,
			TextMain.Render("↑↓"), " ", TextDark.Render("Navigate"),
			dot,
			TextMain.Render("Enter"), " ", TextDark.Render("Select Highlighted Option"),
		),
	)
}
