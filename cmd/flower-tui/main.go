package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
)

const (
	APPNAME    = "flower-tui"
	APPVERSION = "0.1.0"
)

func main() {
	app := &App{mainRouter: router.New()}
	program := reactea.NewProgram(app, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
