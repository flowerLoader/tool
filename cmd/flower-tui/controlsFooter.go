package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
)

type controlsFooterComponent struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]
}

func (c *controlsFooterComponent) Init() tea.Cmd {
	return nil
}

func (c *controlsFooterComponent) Render(width, height int) string {
	dot := theme.Gloss(DisabledStyle).Render(" • ")
	space := theme.Gloss(DisabledStyle).Render(" ")

	return theme.Gloss(DefaultStyle).Width(width).Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left, space,
			theme.Gloss(PrimaryStyle).Render("Ctrl+C"), space,
			theme.Gloss(SecondaryStyle).Render("Quit"),
			dot,
			theme.Gloss(PrimaryStyle).Render("↑↓"), space,
			theme.Gloss(SecondaryStyle).Render("Navigate"),
			dot,
			theme.Gloss(PrimaryStyle).Render("Enter"), space,
			theme.Gloss(SecondaryStyle).Render("Select Highlighted Option"),
		),
	)
}
