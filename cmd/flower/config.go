package main

import (
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/viper"
)

type LoaderConfig struct {
	ID    string
	Base  []string
	Build struct {
		ID          string
		Entrypoints []string
	}
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

type EnvironmentConfig struct {
	DB   string `json:"db"`   // Path to the database file
	Game string `json:"game"` // Game ID (format: "subsystem:app_id")
	Path string `json:"path"` // Path to the game's installation directory
}

type AppConfig struct {
	Environments []EnvironmentConfig
}

type Config struct {
	App       AppConfig
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

	// use Viper to load AppConfig, appending to toolConfig
	viper.SetConfigName("flower_env")
	viper.SetConfigType("json")

	// on *nix this is $HOME/.config/APPNAME/config.json
	// on Windows this is %APPDATA%\APPNAME\config.json
	osConfigPath := filepath.Join(os.Getenv("HOME"), ".config", APPNAME)
	if runtime.GOOS == "windows" {
		osConfigPath = filepath.Join(os.Getenv("APPDATA"), APPNAME)
	}
	if err := os.Mkdir(osConfigPath, 0755); err != nil && !os.IsExist(err) {
		return nil, err
	}

	log.Debug("Loading Application Configuration",
		"osConfigPath", osConfigPath)
	viper.AddConfigPath(osConfigPath)

	if err := viper.ReadInConfig(); err == nil {
		if err := viper.Unmarshal(&toolConfig.App); err != nil {
			return toolConfig, err
		}

		log.Debug("Loaded Application Configuration",
			"filepath", viper.ConfigFileUsed(),
			"config.environments", toolConfig.App.Environments)
	} else if strings.Contains(err.Error(), "Not Found in ") {
		if err := viper.WriteConfigAs(
			filepath.Join(osConfigPath, "flower_env.json"),
		); err != nil {
			return toolConfig, err
		}

		log.Warn("Created Default Application Configuration",
			"filepath", viper.ConfigFileUsed(),
			"environments", toolConfig.App.Environments)
	} else {
		log.Warn("Failed to Load Application Configuration",
			"error", err)
	}

	return toolConfig, nil
}
