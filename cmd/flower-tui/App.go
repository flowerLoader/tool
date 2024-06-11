package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
)

type App struct {
	reactea.BasicComponent                         // AfterUpdate()
	reactea.BasicPropfulComponent[reactea.NoProps] // UpdateProps() and Props()

	mainRouter reactea.Component[router.Props]
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			// component := input.New()

			// return component, component.Init(input.Props{
			// 	SetText: c.setText, // Can also use "lambdas" (function can be created here)
			// })

			return nil, nil
		},
	})
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}

	return c.mainRouter.Update(msg)
}

func (c *App) Render(width, height int) string {
	return c.mainRouter.Render(width, height)
}
