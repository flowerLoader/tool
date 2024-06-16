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
	zone "github.com/lrstanley/bubblezone"
)

type welcomeComponent struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	// Components
	filterInput FormField

	// Optimization for re-use
	sortedGameNames []string

	// State
	cursorPos int
	cursorMax int
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

	return tea.Batch(
		textinput.Blink,
		c.filterInput.Cursor.SetMode(cursor.CursorBlink),
	)
}

func (c *welcomeComponent) Update(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 1)
	c.filterInput, cmds[0] = c.filterInput.Update(msg)

	prevPos := c.cursorPos
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			c.setCursorPos(c.cursorPos - 1)

		case "down":
			c.setCursorPos(c.cursorPos + 1)

		case "enter":
			cmds = append(cmds, c.handleSubmit())
		}

	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			c.setCursorPos(c.cursorPos - 1)
		case tea.MouseButtonWheelDown:
			c.setCursorPos(c.cursorPos + 1)
		}

		for i := 0; i <= c.cursorMax; i++ {
			if zone.Get(fmt.Sprintf("cursor%d", i)).InBounds(msg) {
				c.cursorPos = i
				if msg.Action == tea.MouseActionRelease || msg.Button == tea.MouseButtonLeft {
					cmds = append(cmds, c.handleSubmit())
				}
				break
			}
		}
	}

	if prevPos != c.cursorPos {
		if c.cursorPos == 0 {
			c.filterInput.TextStyle = theme.Gloss(PrimaryStyle)
			cmds = append(cmds, c.filterInput.Focus())
		} else {
			c.filterInput.TextStyle = theme.Gloss(DefaultStyle)
			c.filterInput.Blur()
		}
	}

	return tea.Batch(cmds...)
}

func (c *welcomeComponent) handleSubmit() tea.Cmd {
	switch c.cursorPos {
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

func (c *welcomeComponent) renderCursor(pos int, after string) string {
	if c.cursorMax < pos {
		c.cursorMax = pos
	}

	elements := make([]string, 1)
	if c.cursorPos == pos {
		elements[0] = theme.Gloss(PrimaryStyle).Render(" → ")
		elements = append(elements, theme.Gloss(PrimaryStyle).Render(" "))
		elements = append(elements, theme.Gloss(PrimaryStyle).Render(after))
	} else {
		elements[0] = theme.Gloss(DefaultStyle).Render("   ")
		elements = append(elements, theme.Gloss(DefaultStyle).Render(" "))
		elements = append(elements, theme.Gloss(DefaultStyle).Render(after))
	}

	return zone.Mark(
		fmt.Sprintf("cursor%d", pos),
		lipgloss.JoinHorizontal(lipgloss.Left, elements...),
	)
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
		Background(lipgloss.Color(theme.Default.Background)).
		Render(lipgloss.JoinVertical(lipgloss.Left, elements...))
}

func (c *welcomeComponent) setCursorPos(pos int) {
	if pos < 0 {
		pos = 0
	}

	if pos > c.cursorMax {
		pos = c.cursorMax
	}

	c.cursorPos = pos
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
		innerBoxStyle.Render(c.renderCursor(0, c.renderInput())),

		// Help Text
		innerBoxStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				c.renderCursor(1, "Add Unsupported Game (Advanced)"),
				c.renderCursor(2, "Manage Environments            "),
				c.renderCursor(3, "TUI Settings                   "),
				c.renderCursor(4, "Quit                           "),
			),
		),
	)
}
