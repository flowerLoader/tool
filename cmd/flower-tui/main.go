package main

import (
	_ "embed"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	zone "github.com/lrstanley/bubblezone"

	"github.com/flowerLoader/tool/pkg/cfg"
)

const (
	APPNAME    = "flower-tui"
	APPVERSION = "0.1.0"
)

var (
	app    *App
	config *cfg.Config

	//go:embed main.json
	MAIN_JSON []byte
)

func main() {
	var err error
	config, err = cfg.LoadFromJSON(MAIN_JSON)
	if err != nil || len(config.Games) == 0 {
		panic("fatal: no games found in config (check main.json and rebuild) error: " + err.Error())
	}

	osInit()

	zone.NewGlobal()
	app = &App{mainRouter: router.New()}
	program := reactea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err = program.Run(); err != nil {
		panic(err)
	}
}
