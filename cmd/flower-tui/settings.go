package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
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
	c.borderColor = NewFormField("Border", "#", "(default: 103)", theme.Border.Foreground)
	c.primaryColor = NewFormField("Primary", "#", "(default: 176)", theme.Primary.Foreground)
	c.secondaryColor = NewFormField("Secondary", "#", "(default: 96)", theme.Secondary.Foreground)
	c.disabledColor = NewFormField("Disabled", "#", "(default: 243)", theme.Disabled.Foreground)
	c.errorColor = NewFormField("Error", "#", "(default: 9)", theme.Error.Foreground)

	c.setCursorPos(0)

	return tea.Batch(
		textinput.Blink,
		c.borderColor.Focus(),
		c.borderColor.Cursor.SetMode(cursor.CursorBlink),
		c.primaryColor.Cursor.SetMode(cursor.CursorBlink),
		c.secondaryColor.Cursor.SetMode(cursor.CursorBlink),
		c.disabledColor.Cursor.SetMode(cursor.CursorBlink),
		c.errorColor.Cursor.SetMode(cursor.CursorBlink),
	)
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

	switch c.cursorPos {
	case 0:
		c.borderColor, cmd = c.borderColor.Update(msg)
		theme.Border.Foreground = c.borderColor.Value()
	case 1:
		c.primaryColor, cmd = c.primaryColor.Update(msg)
		theme.Primary.Foreground = c.primaryColor.Value()
	case 2:
		c.secondaryColor, cmd = c.secondaryColor.Update(msg)
		theme.Secondary.Foreground = c.secondaryColor.Value()
	case 3:
		c.disabledColor, cmd = c.disabledColor.Update(msg)
		theme.Disabled.Foreground = c.disabledColor.Value()
	case 4:
		c.errorColor, cmd = c.errorColor.Update(msg)
		theme.Error.Foreground = c.errorColor.Value()
	}
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (c *settingsComponent) handleSubmit() tea.Cmd {
	switch c.cursorPos {
	case 0: // Border Color
	case 1: // Primary Color
	case 2: // Secondary Color
	case 3: // Disabled Color
	case 4: // Error Color
	case 5: // Back to Main Menu
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

	//
	// Update styles according to cursor position
	//

	if c.cursorPos == c.cursorMax {
		c.borderColor.TextStyle = theme.Gloss(DefaultStyle)
		c.primaryColor.TextStyle = theme.Gloss(DefaultStyle)
		c.secondaryColor.TextStyle = theme.Gloss(DefaultStyle)
		c.disabledColor.TextStyle = theme.Gloss(DefaultStyle)
		c.errorColor.TextStyle = theme.Gloss(DefaultStyle)
	}

	if c.cursorPos == 0 {
		c.borderColor.TextStyle = theme.Gloss(PrimaryStyle)
		c.borderColor.Focus()
	} else {
		c.borderColor.TextStyle = theme.Gloss(DefaultStyle)
		c.borderColor.Blur()
	}

	if c.cursorPos == 1 {
		c.primaryColor.TextStyle = theme.Gloss(PrimaryStyle)
		c.primaryColor.Focus()
	} else {
		c.primaryColor.TextStyle = theme.Gloss(DefaultStyle)
		c.primaryColor.Blur()
	}

	if c.cursorPos == 2 {
		c.secondaryColor.TextStyle = theme.Gloss(PrimaryStyle)
		c.secondaryColor.Focus()
	} else {
		c.secondaryColor.TextStyle = theme.Gloss(DefaultStyle)
		c.secondaryColor.Blur()
	}

	if c.cursorPos == 3 {
		c.disabledColor.TextStyle = theme.Gloss(PrimaryStyle)
		c.disabledColor.Focus()
	} else {
		c.disabledColor.TextStyle = theme.Gloss(DefaultStyle)
		c.disabledColor.Blur()
	}

	if c.cursorPos == 4 {
		c.errorColor.TextStyle = theme.Gloss(PrimaryStyle)
		c.errorColor.Focus()
	} else {
		c.errorColor.TextStyle = theme.Gloss(DefaultStyle)
		c.errorColor.Blur()
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
		innerBoxStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				c.renderCursor(5, "Back to Main Menu"),
			),
		),
	)
}
