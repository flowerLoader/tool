package main

import (
	"fmt"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	zone "github.com/lrstanley/bubblezone"
)

type welcomeProps struct {
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
	c.input.Cursor.BlinkSpeed = time.Second / 3
	c.input.Placeholder = "Type to filter games..."
	c.input.Focus()
	c.input.Prompt = ""
	c.input.ShowSuggestions = true
	c.input.Width = 30

	gameNames := make([]string, len(app.config.Games))
	for i, game := range app.config.Games {
		for _, names := range game.Meta.Name {
			gameNames[i] = names
			break
		}
	}
	c.input.SetSuggestions(gameNames)

	return tea.Batch(
		textinput.Blink,
		c.input.Cursor.SetMode(cursor.CursorBlink),
	)
}

func (c *welcomeComponent) Update(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 1)
	c.input, cmds[0] = c.input.Update(msg)

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
			c.input.TextStyle = TextMain
			cmds = append(cmds, c.input.Focus())
		} else {
			c.input.TextStyle = TextDisabled
			c.input.Blur()
		}
	}

	return tea.Batch(cmds...)
}

func (c *welcomeComponent) handleSubmit() tea.Cmd {
	switch c.cursorPos {
	case 0:
		reactea.SetCurrentRoute(fmt.Sprintf("game/%s", c.inputAutocomplete(c.input.Value())))
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

func (c *welcomeComponent) inputAutocomplete(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))

	for _, game := range app.config.Games {
		// sort by locale
		sorted := make([]string, 0, len(game.Meta.Name))
		for locale := range game.Meta.Name {
			sorted = append(sorted, locale)
		}

		slices.Sort(sorted)
		for _, locale := range sorted {
			name := game.Meta.Name[locale]
			if strings.Contains(strings.ToLower(name), text) {
				return name
			}
		}
	}

	return ""
}

const minHeight = 18   // # of terminal lines reserved for header and footer
const spacingRatio = 4 // # of terminal lines per 1 spacing line

func (c *welcomeComponent) renderCursor(pos int, after string) string {
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

func (c *welcomeComponent) renderInput() string {
	elements := []string{
		c.input.View(),
	}

	if guess := c.inputAutocomplete(c.input.Value()); guess == "" {
		elements = append(elements, TextError.Render(
			"Game not found. Please check your spelling and try again.",
		))
	} else {
		elements = append(elements, TextMain.Render(
			fmt.Sprintf("Press Enter to select: %s", guess),
		))
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		elements...,
	)
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
			c.renderCursor(0, c.renderInput()),
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
