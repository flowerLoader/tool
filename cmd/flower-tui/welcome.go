package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
)

type welcomeProps struct {
	gamePath   string
	sourcePath string
}

type welcomeComponent struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[welcomeProps]

	// Components
	input textinput.Model

	// State
	cursorPos int
	cursorMax int
}

func (c *welcomeComponent) Init(props *welcomeProps) tea.Cmd {
	c.input = textinput.New()
	c.input.CharLimit = 250
	c.input.Cursor.BlinkSpeed = time.Second / 2
	c.input.Placeholder = "Type to filter games..."
	c.input.Prompt = ""
	c.input.Width = 30

	return tea.Batch(
		textinput.Blink,
		c.input.Cursor.SetMode(cursor.CursorBlink),
		c.input.Focus(),
	)
}

func (c *welcomeComponent) Update(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 1)
	c.input, cmds[0] = c.input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			c.cursorPos--
			if c.cursorPos < 0 {
				c.cursorPos = 0
			}

		case "down":
			c.cursorPos++
			if c.cursorPos > c.cursorMax {
				c.cursorPos = c.cursorMax
			}

		case "enter":
			cmds = append(cmds, c.handleSubmit())
		}
	}

	if c.cursorPos == 0 {
		c.input.TextStyle = TextMain
		cmds = append(cmds, c.input.Focus())
	} else {
		c.input.TextStyle = TextDisabled
		c.input.Blur()
	}

	return tea.Batch(cmds...)
}

func (c *welcomeComponent) handleSubmit() tea.Cmd {
	switch c.cursorPos {
	case 0:
		reactea.SetCurrentRoute("gameSelect")
	case 1:
		reactea.SetCurrentRoute("gameAdd")
	case 2:
		reactea.SetCurrentRoute("envManage")
	case 3:
		reactea.SetCurrentRoute("settings")
	case 4:
		return reactea.Destroy
	}

	return nil
}

func (c *welcomeComponent) renderCursor(pos int, after string) string {
	if c.cursorMax < pos {
		c.cursorMax = pos
	}

	if c.cursorPos == pos {
		return TextMain.
			Background(lipgloss.Color("234")).
			Render("  → " + after)
	}

	return TextDisabled.Render("    " + after)
}

const minHeight = 18   // # of terminal lines reserved for header and footer
const spacingRatio = 4 // # of terminal lines per 1 spacing line

func (c *welcomeComponent) Render(width, height int) string {
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
				TextDark.Render(fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)),
			)),

		// Filter Input Entry
		innerBoxStyle.Render(
			c.renderCursor(0, c.input.View()),
		),

		// Help Text
		innerBoxStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				c.renderCursor(1, "Add Unsupported Game (Advanced)"),
				c.renderCursor(2, "Manage Environments"),
				c.renderCursor(3, "TUI Settings"),
				c.renderCursor(4, "Quit"),
			),
		),
	))

	return sb.String()
}
