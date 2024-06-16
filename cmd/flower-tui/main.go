package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"

	zone "github.com/lrstanley/bubblezone"
)

const (
	APPNAME    = "flower-tui"
	APPVERSION = "0.1.0"
)

var app *App

func main() {
	zone.NewGlobal()
	app = &App{mainRouter: router.New()}
	program := reactea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
