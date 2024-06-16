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

type settingsProps struct {
	getThemeColor func(string) string
	setThemeColor func(string, string)
}

type settingsComponent struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[*settingsProps]

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

func (c *settingsComponent) Init(props *settingsProps) tea.Cmd {
	c.UpdateProps(props)
	getThemeColor := props.getThemeColor

	val := getThemeColor("Border")
	c.borderColor = NewFormField("Border", "#", "(default: 103)", val)

	val = getThemeColor("Primary")
	c.primaryColor = NewFormField("Primary", "#", "(default: 176)", val)

	val = getThemeColor("Secondary")
	c.secondaryColor = NewFormField("Secondary", "#", "(default: 96)", val)

	val = getThemeColor("Disabled")
	c.disabledColor = NewFormField("Disabled", "#", "(default: 243)", val)

	val = getThemeColor("Error")
	c.errorColor = NewFormField("Error", "#", "(default: 9)", val)

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
		color := c.borderColor.Value()
		c.borderColor, cmd = c.borderColor.Update(msg)
		if newColor := c.borderColor.Value(); color != newColor {
			cmds = append(cmds, c.setThemeColor("Border", newColor))
		}
	case 1:
		color := c.primaryColor.Value()
		c.primaryColor, cmd = c.primaryColor.Update(msg)
		if newColor := c.primaryColor.Value(); color != newColor {
			cmds = append(cmds, c.setThemeColor("Primary", newColor))
		}
	case 2:
		color := c.secondaryColor.Value()
		c.secondaryColor, cmd = c.secondaryColor.Update(msg)
		if newColor := c.secondaryColor.Value(); color != newColor {
			cmds = append(cmds, c.setThemeColor("Secondary", newColor))
		}
	case 3:
		color := c.disabledColor.Value()
		c.disabledColor, cmd = c.disabledColor.Update(msg)
		if newColor := c.disabledColor.Value(); color != newColor {
			cmds = append(cmds, c.setThemeColor("Disabled", newColor))
		}
	case 4:
		color := c.errorColor.Value()
		c.errorColor, cmd = c.errorColor.Update(msg)
		if newColor := c.errorColor.Value(); color != newColor {
			cmds = append(cmds, c.setThemeColor("Error", newColor))
		}
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
		elements[0] = ColorPrimary.Render(" → ")
		elements = append(elements, ColorPrimary.Render(" "))
		elements = append(elements, ColorPrimary.Render(after))
	} else {
		elements[0] = ColorDisabled.Render("   ")
		elements = append(elements, ColorDisabled.Render(" "))
		elements = append(elements, ColorDisabled.Render(after))
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
		c.borderColor.TextStyle = ColorDisabled
		c.primaryColor.TextStyle = ColorDisabled
		c.secondaryColor.TextStyle = ColorDisabled
		c.disabledColor.TextStyle = ColorDisabled
		c.errorColor.TextStyle = ColorDisabled
	}

	if c.cursorPos == 0 {
		c.borderColor.TextStyle = ColorPrimary
		c.borderColor.Focus()
	} else {
		c.borderColor.TextStyle = ColorDisabled
		c.borderColor.Blur()
	}

	if c.cursorPos == 1 {
		c.primaryColor.TextStyle = ColorPrimary
		c.primaryColor.Focus()
	} else {
		c.primaryColor.TextStyle = ColorDisabled
		c.primaryColor.Blur()
	}

	if c.cursorPos == 2 {
		c.secondaryColor.TextStyle = ColorPrimary
		c.secondaryColor.Focus()
	} else {
		c.secondaryColor.TextStyle = ColorDisabled
		c.secondaryColor.Blur()
	}

	if c.cursorPos == 3 {
		c.disabledColor.TextStyle = ColorPrimary
		c.disabledColor.Focus()
	} else {
		c.disabledColor.TextStyle = ColorDisabled
		c.disabledColor.Blur()
	}

	if c.cursorPos == 4 {
		c.errorColor.TextStyle = ColorPrimary
		c.errorColor.Focus()
	} else {
		c.errorColor.TextStyle = ColorDisabled
		c.errorColor.Blur()
	}
}

func (c *settingsComponent) setThemeColor(key, value string) tea.Cmd {
	return func() tea.Msg {
		c.Props().setThemeColor(key, value)
		return nil
	}
}

func (c *settingsComponent) Render(width, height int) string {
	spacing := 0
	if height > minHeight {
		spacing = (height - minHeight) / spacingRatio
	}

	var innerBoxStyle = ColorBorder.
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
