package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	zone "github.com/lrstanley/bubblezone"
)

type settingsComponent struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	// Components
	borderColor    FormField
	primaryColor   FormField
	secondaryColor FormField
	disabledColor  FormField
	errorColor     FormField

	// State
	cursorPos int
	cursorMax int
}

func (c *settingsComponent) Init() tea.Cmd {
	fields := []struct {
		field *FormField
		label string
		def   string
	}{
		{&c.borderColor, "Border", theme.Styles[BorderStyle].Foreground},
		{&c.primaryColor, "Primary", theme.Styles[PrimaryStyle].Foreground},
		{&c.secondaryColor, "Secondary", theme.Styles[SecondaryStyle].Foreground},
		{&c.disabledColor, "Disabled", theme.Styles[DisabledStyle].Foreground},
		{&c.errorColor, "Error", theme.Styles[ErrorStyle].Foreground},
	}

	cmds := []tea.Cmd{}
	for _, f := range fields {
		*f.field = NewFormField(f.label, "", fmt.Sprintf("(default: %s)", f.def), f.def)
		cmds = append(cmds, (*f.field).Focus(), (*f.field).Cursor.SetMode(cursor.CursorBlink))
	}

	c.setCursorPos(0)

	return tea.Batch(cmds...)
}

func (c *settingsComponent) Update(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	var cmd tea.Cmd

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

	fields := []*FormField{
		&c.borderColor,
		&c.primaryColor,
		&c.secondaryColor,
		&c.disabledColor,
		&c.errorColor,
	}

	for i, field := range fields {
		if c.cursorPos == i {
			*field, cmd = field.Update(msg)
			styleType := StyleType(strings.ToLower(field.Label))
			style := theme.Styles[styleType]
			style.Foreground = field.Value()
			theme.Styles[styleType] = style
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

func (c *settingsComponent) handleSubmit() tea.Cmd {
	if c.cursorPos == 5 { // Back to Main Menu
		reactea.SetCurrentRoute("default")
	}

	return nil
}

func (c *settingsComponent) renderCursor(pos int, after string) string {
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

func (c *settingsComponent) setCursorPos(pos int) {
	if pos < 0 {
		pos = 0
	}

	if pos > c.cursorMax {
		pos = c.cursorMax
	}

	c.cursorPos = pos

	fields := []*FormField{
		&c.borderColor,
		&c.primaryColor,
		&c.secondaryColor,
		&c.disabledColor,
		&c.errorColor,
	}

	for i, field := range fields {
		if c.cursorPos == i {
			field.TextStyle = theme.Gloss(PrimaryStyle)
			field.Focus()
		} else {
			field.TextStyle = theme.Gloss(DefaultStyle)
			field.Blur()
		}
	}
}

func (c *settingsComponent) Render(width, height int) string {
	spacing := 0
	if height > minHeight {
		spacing = (height - minHeight) / spacingRatio
	}

	var innerBoxStyle = theme.Gloss(BorderStyle).
		Width(width-10).
		Margin(spacing, 4, 0).
		Padding(spacing/2, 0).
		Border(lipgloss.RoundedBorder(), true)

	return lipgloss.JoinVertical(
		lipgloss.Top,

		// Input Entries
		innerBoxStyle.Render(lipgloss.JoinVertical(
			lipgloss.Left,
			c.renderCursor(0, c.borderColor.View()),
			c.renderCursor(1, c.primaryColor.View()),
			c.renderCursor(2, c.secondaryColor.View()),
			c.renderCursor(3, c.disabledColor.View()),
			c.renderCursor(4, c.errorColor.View()),
		)),

		// Help Text
		innerBoxStyle.Render(c.renderCursor(5, "Back to Main Menu")),
	)
}
