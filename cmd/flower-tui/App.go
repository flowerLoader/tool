package main

import (
	_ "embed"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	zone "github.com/lrstanley/bubblezone"

	"github.com/flowerLoader/tool/pkg/cfg"
)

//go:embed main.json
var MAIN_JSON []byte

type App struct {
	reactea.BasicComponent                         // AfterUpdate()
	reactea.BasicPropfulComponent[reactea.NoProps] // UpdateProps() and Props()

	// Components
	mainRouter reactea.Component[router.Props]
	controls   *controlsFooterComponent

	// State
	config     *cfg.Config
	gamePath   string
	sourcePath string
}

func (app *App) Init(reactea.NoProps) tea.Cmd {
	// State
	var err error
	app.config, err = cfg.LoadFromJSON(MAIN_JSON)
	if err != nil || len(app.config.Games) == 0 {
		panic("fatal: no games found in config (check main.json and rebuild) error: " + err.Error())
	}

	osInit()

	// Components
	app.controls = &controlsFooterComponent{}
	app.controls.Init(&controlsFooterProps{})

	// Router
	return app.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := &welcomeComponent{}
			return component, component.Init(&welcomeProps{})
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

	// Render the main components
	return zone.Scan(lipgloss.JoinVertical(
		lipgloss.Left,

		// Main content
		lipgloss.NewStyle().
			Width(innerWidth).
			Height(innerHeight).
			BorderForeground(ColorPrimaryMain).
			Border(lipgloss.DoubleBorder(), true).
			Render(app.mainRouter.Render(innerWidth, innerHeight)),

		// Footer
		app.controls.Render(outerWidth, footerHeight),
	))
}
