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
	dot := ColorDisabled.Render(" • ")

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Left, " ",
			ColorPrimary.Render("Ctrl+C"), " ", ColorSecondary.Render("Quit"),
			dot,
			ColorPrimary.Render("↑↓"), " ", ColorSecondary.Render("Navigate"),
			dot,
			ColorPrimary.Render("Enter"), " ", ColorSecondary.Render("Select Highlighted Option"),
		),
	)
}
