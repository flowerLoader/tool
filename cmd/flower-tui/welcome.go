package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
)

type welcomeComponent struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	// Components
	nav         *navComponent
	filterInput FormField

	// Optimization for re-use
	sortedGameNames []string
}

const minHeight = 18   // # of terminal lines reserved for header and footer
const spacingRatio = 4 // # of terminal lines per 1 spacing line

func (c *welcomeComponent) Init() tea.Cmd {
	c.filterInput = NewFormField("", "? ", "Type to filter games...", "")
	c.filterInput.ShowSuggestions = true
	c.filterInput.Focus()

	c.sortedGameNames = make([]string, 0)
	for _, game := range config.Games {
		for _, names := range game.Meta.Name {
			c.sortedGameNames = append(c.sortedGameNames, names)
			break
		}
	}
	slices.Sort(c.sortedGameNames)
	c.filterInput.SetSuggestions(c.sortedGameNames)

	c.nav = NewNavComponent(theme.Gloss(PrimaryStyle), []Item{
		{Func: c.renderInput},
		{Name: "Add Unsupported Game (Advanced)"},
		{Name: "Manage Environments"},
		{Name: "TUI Settings"},
		{Name: "Quit"},
	}, func(position int) tea.Cmd {
		switch position {
		case 0:
			reactea.SetCurrentRoute(fmt.Sprintf("game/%s", c.inputAutocomplete(c.filterInput.Value())))
		case 1:
			reactea.SetCurrentRoute("game/unsupported")
		case 2:
			reactea.SetCurrentRoute("manage-environments")
		case 3:
			reactea.SetCurrentRoute("settings")
		case 4:
			return reactea.Destroy
		}

		return nil
	})

	return tea.Batch(
		textinput.Blink,
		c.filterInput.Cursor.SetMode(cursor.CursorBlink),
	)
}

func (c *welcomeComponent) Update(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 1)
	c.filterInput, cmds[0] = c.filterInput.Update(msg)

	prevPos := c.nav.position
	if cmd := c.nav.Update(msg); cmd != nil {
		return cmd // return early if nav component handled the message
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			cmds = append(cmds, c.nav.onClick(c.nav.position))
		}
	}

	if prevPos != c.nav.position {
		if c.nav.position == 0 {
			c.filterInput.TextStyle = theme.Gloss(PrimaryStyle)
			cmds = append(cmds, c.filterInput.Focus())
		} else {
			c.filterInput.TextStyle = theme.Gloss(DefaultStyle)
			c.filterInput.Blur()
		}
	}

	return tea.Batch(cmds...)
}

func (c *welcomeComponent) inputAutocomplete(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))

	for _, name := range c.sortedGameNames {
		if strings.Contains(strings.ToLower(name), text) {
			return name
		}
	}

	return ""
}

func (c *welcomeComponent) renderInput() string {
	elements := []string{
		c.filterInput.View(),
	}

	if guess := c.inputAutocomplete(c.filterInput.Value()); guess == "" {
		elements = append(elements, theme.Gloss(ErrorStyle).Render(
			"Game not found. Please check your spelling and try again.",
		))
	} else {
		elements = append(elements, theme.Gloss(SecondaryStyle).Render(
			fmt.Sprintf("Press Enter to select: %s", guess),
		))
	}

	return lipgloss.NewStyle().
		Width(c.filterInput.Width).
		Background(lipgloss.Color(theme.Styles[DefaultStyle].Background)).
		Render(lipgloss.JoinVertical(lipgloss.Left, elements...))
}

func (c *welcomeComponent) Render(width, height int) string {
	spacing := 0
	if height > minHeight {
		spacing = (height - minHeight) / spacingRatio
	}

	var innerBoxStyle = theme.Gloss(BorderStyle).
		Width(width-10).
		Margin(spacing, 4, 0).
		Padding(spacing/2, 0).
		Border(lipgloss.RoundedBorder(), true)

	c.filterInput.Width = width - 36

	return lipgloss.JoinVertical(
		lipgloss.Top,

		// Filter Input Entry
		innerBoxStyle.Render(c.nav.Render(0, width-10, 1)),

		// Help Text
		innerBoxStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				c.nav.Render(1, width-10, 1),
				c.nav.Render(2, width-10, 1),
				c.nav.Render(3, width-10, 1),
				c.nav.Render(4, width-10, 1),
			),
		),
	)
}
