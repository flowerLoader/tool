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
}

func (app *App) Init(reactea.NoProps) tea.Cmd {
	// Components
	app.controls = &controlsFooterComponent{}
	app.controls.Init()

	// Router
	return app.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := &welcomeComponent{}
			return component, component.Init()
		},
		"settings": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := &settingsComponent{}
			return component, component.Init()
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
	return zone.Scan(theme.Gloss(DefaultStyle).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,

			// Main content
			theme.Gloss(BorderStyle).
				Width(innerWidth).
				Height(innerHeight).
				Border(lipgloss.DoubleBorder(), true).
				Render(lipgloss.JoinVertical(
					lipgloss.Left,

					// Header
					theme.Gloss(PrimaryStyle).
						Padding(1, 0).
						Width(innerWidth).
						AlignHorizontal(lipgloss.Center).
						Render(fmt.Sprintf(
							"%s %s\n%s", APPNAME,
							theme.Gloss(SecondaryStyle).Render(fmt.Sprintf("v%s", APPVERSION)),
							theme.Gloss(SecondaryStyle).Render(currentPage),
						)),

					app.mainRouter.Render(innerWidth, innerHeight-4),
				)),

			// Footer
			app.controls.Render(outerWidth, footerHeight),
		)))
}
