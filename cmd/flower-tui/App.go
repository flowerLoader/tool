package main

import (
	_ "embed"
	"fmt"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	zone "github.com/lrstanley/bubblezone"
)

type App struct {
	reactea.BasicComponent                         // AfterUpdate()
	reactea.BasicPropfulComponent[reactea.NoProps] // UpdateProps() and Props()

	// Components
	mainRouter reactea.Component[router.Props]
	controls   *controlsFooterComponent

	// State
	theme map[string]string
}

func (app *App) Init(reactea.NoProps) tea.Cmd {
	// Components
	app.controls = &controlsFooterComponent{}
	app.controls.Init()

	app.theme = make(map[string]string)
	app.theme["Border"] = ANSIBorder
	app.theme["Primary"] = ANSIPrimary
	app.theme["Secondary"] = ANSISecondary
	app.theme["Disabled"] = ANSIDisabled
	app.theme["Error"] = ANSIError

	// Router
	return app.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := &welcomeComponent{}
			return component, component.Init()
		},
		"settings": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := &settingsComponent{}
			return component, component.Init(&settingsProps{
				getThemeColor: func(key string) string {
					return app.theme[key]
				},
				setThemeColor: func(key, value string) {
					app.theme[key] = value

					switch key {
					case "Border":
						ANSIBorder = value
						ColorBorder = lipgloss.NewStyle().BorderForeground(lipgloss.Color(ANSIBorder))
					case "Primary":
						ANSIPrimary = value
						ColorPrimary = lipgloss.NewStyle().Foreground(lipgloss.Color(ANSIPrimary))
					case "Secondary":
						ANSISecondary = value
						ColorSecondary = lipgloss.NewStyle().Foreground(lipgloss.Color(ANSISecondary))
					case "Disabled":
						ANSIDisabled = value
						ColorDisabled = lipgloss.NewStyle().Foreground(lipgloss.Color(ANSIDisabled))
					case "Error":
						ANSIError = value
						ColorError = lipgloss.NewStyle().Foreground(lipgloss.Color(ANSIError))
					}
				},
			})
		},
	})
}

func (app *App) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}

	return app.mainRouter.Update(msg)
}

func (app *App) Render(outerWidth, outerHeight int) string {
	footerHeight := 1
	outerHeight -= footerHeight

	innerHeight := outerHeight - 2 // Subtract 2 for the border
	innerWidth := outerWidth - 2   // Subtract 2 for the border

	currentPage := reactea.CurrentRoute()
	if currentPage == "default" || currentPage == "" {
		currentPage = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	} else {
		currentPage = strings.ToUpper(currentPage[:1]) + strings.ToLower(currentPage[1:])
	}

	// Render the main components
	return zone.Scan(lipgloss.NewStyle().
		Background(lipgloss.Color(ANSIBackground)).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,

			// Main content
			ColorBorder.
				Width(innerWidth).
				Height(innerHeight).
				Border(lipgloss.DoubleBorder(), true).
				Render(lipgloss.JoinVertical(
					lipgloss.Left,

					// Header
					lipgloss.NewStyle().
						Padding(1, 0).
						Width(innerWidth).
						AlignHorizontal(lipgloss.Center).
						Background(lipgloss.Color(ANSIBackground)).
						MarginBackground(lipgloss.Color(ANSIBackground)).
						Render(fmt.Sprintf(
							"%s%s\n%s",
							ColorPrimary.Bold(true).Render(fmt.Sprintf("%s ", APPNAME)),
							ColorSecondary.Render(fmt.Sprintf("v%s", APPVERSION)),
							ColorSecondary.Render(currentPage),
						)),

					app.mainRouter.Render(innerWidth, innerHeight-4),
				)),

			// Footer
			app.controls.Render(outerWidth, footerHeight),
		)))
}
