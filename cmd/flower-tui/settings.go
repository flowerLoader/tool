package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	zone "github.com/lrstanley/bubblezone"
)

type settingsProps struct {
}

type settingsComponent struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[settingsProps]

	// Components
	primaryMainColor FormField
	primaryDarkColor FormField
	disabledColor    FormField
	errorColor       FormField

	// State
	cursorPos int
	cursorMax int
}

func (c *settingsComponent) Init(props *settingsProps) tea.Cmd {
	c.primaryMainColor = NewFormField("Primary Main Color", "#", "(default: 176)", "")
	c.primaryDarkColor = NewFormField("Primary Dark Color", "#", "(default: 96)", "")
	c.disabledColor = NewFormField("Disabled Color", "#", "(default: 243)", "")
	c.errorColor = NewFormField("Error Color", "#", "(default: 9)", "")

	return tea.Batch(
		textinput.Blink,
		c.primaryMainColor.Focus(),
		c.primaryMainColor.Cursor.SetMode(cursor.CursorBlink),
		c.primaryDarkColor.Cursor.SetMode(cursor.CursorBlink),
		c.disabledColor.Cursor.SetMode(cursor.CursorBlink),
		c.errorColor.Cursor.SetMode(cursor.CursorBlink),
	)
}

func (c *settingsComponent) Update(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	var cmd tea.Cmd

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

	switch c.cursorPos {
	case 0:
		c.primaryMainColor, cmd = c.primaryMainColor.Update(msg)
	case 1:
		c.primaryDarkColor, cmd = c.primaryDarkColor.Update(msg)
	case 2:
		c.disabledColor, cmd = c.disabledColor.Update(msg)
	case 3:
		c.errorColor, cmd = c.errorColor.Update(msg)
	}
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	if prevPos != c.cursorPos {
		if c.cursorPos == 0 {
			c.primaryMainColor.TextStyle = TextMain
			cmds = append(cmds, c.primaryMainColor.Focus())
		} else {
			c.primaryMainColor.TextStyle = TextDisabled
			c.primaryMainColor.Blur()
		}

		if c.cursorPos == 1 {
			c.primaryDarkColor.TextStyle = TextMain
			cmds = append(cmds, c.primaryDarkColor.Focus())
		} else {
			c.primaryDarkColor.TextStyle = TextDisabled
			c.primaryDarkColor.Blur()
		}

		if c.cursorPos == 2 {
			c.disabledColor.TextStyle = TextMain
			cmds = append(cmds, c.disabledColor.Focus())
		} else {
			c.disabledColor.TextStyle = TextDisabled
			c.disabledColor.Blur()
		}

		if c.cursorPos == 3 {
			c.errorColor.TextStyle = TextMain
			cmds = append(cmds, c.errorColor.Focus())
		} else {
			c.errorColor.TextStyle = TextDisabled
			c.errorColor.Blur()
		}
	}

	return tea.Batch(cmds...)
}

func (c *settingsComponent) handleSubmit() tea.Cmd {
	switch c.cursorPos {
	case 0: // Primary Main Color
	case 1: // Primary Dark Color
	case 2: // Disabled Color
	case 3: // Error Color
	case 4: // Save Changes
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
		elements[0] = TextMain.
			Background(lipgloss.Color("234")).
			Render(" → ")
	} else {
		elements[0] = TextDisabled.
			Render("   ")
	}

	elements = append(elements, after)
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
}

func (c *settingsComponent) Render(width, height int) string {
	usableHeight := height - 2
	usableWidth := width - 2

	spacing := 0
	if usableHeight > minHeight {
		spacing = (usableHeight - minHeight) / spacingRatio
	}

	var innerBoxStyle = lipgloss.NewStyle().
		Width(usableWidth-8).
		Margin(spacing, 4, 0).
		Padding(spacing/2, 0).
		BorderForeground(ColorPrimaryMain).
		Border(lipgloss.RoundedBorder(), true)

	// c.filterInput.Width = usableWidth - 36

	var sb strings.Builder
	sb.WriteString(lipgloss.JoinVertical(
		lipgloss.Top,

		// Header
		lipgloss.NewStyle().
			Padding(1, 0).
			Width(usableWidth).
			AlignHorizontal(lipgloss.Center).
			Render(fmt.Sprintf(
				"%s %s\n%s",
				TextMain.Bold(true).Render(APPNAME),
				TextDark.Render(fmt.Sprintf("v%s", APPVERSION)),
				TextDisabled.Render("Settings"),
			)),

		// Input Entries
		innerBoxStyle.Render(lipgloss.JoinVertical(
			lipgloss.Left,
			c.renderCursor(0, c.primaryMainColor.View()),
			c.renderCursor(1, c.primaryDarkColor.View()),
			c.renderCursor(2, c.disabledColor.View()),
			c.renderCursor(3, c.errorColor.View()),
		)),

		// Help Text
		innerBoxStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				c.renderCursor(4, TextDisabled.Render("No Changes to Save")),
				c.renderCursor(5, "Back to Main Menu"),
			),
		),
	))

	return sb.String()
}
