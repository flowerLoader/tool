package main

import (
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
)

const (
	APPNAME    = "flower-tui"
	APPVERSION = "0.1.0"
)

func main() {
	app := &App{mainRouter: router.New()}
	program := reactea.NewProgram(app)
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
