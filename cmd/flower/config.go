package main

import (
	_ "embed"
	"encoding/json"
)

type LoaderConfig struct {
	ID    string
	Base  []string
	Build struct {
		ID          string
		Entrypoints []string
	}
	Map map[string]string
}

type GameConfig struct {
	Loader LoaderConfig
	Meta   struct {
		Name map[string]string // [en] => "Game Name"
		OS   []string          `json:"os"`
	}
	Subsystem struct {
		ID      string `json:"id"`
		AppID   int    `json:"app_id"`
		AppName string `json:"app_name"`
	}
}

type ProviderConfig struct {
	Name      string
	Hosts     []string
	Schemas   []string
	Subsystem string
}

type Config struct {
	Games     []GameConfig
	Providers []ProviderConfig
}

//go:embed main.json
var MAIN_JSON []byte

func NewConfig() (*Config, error) {
	toolConfig := new(Config)

	if err := json.Unmarshal(MAIN_JSON, &toolConfig); err != nil {
		return nil, err
	}

	return toolConfig, nil
}
