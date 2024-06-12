package cfg

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/viper"
)

var BaseDirectory = "flower"
var DefaultFilename = "flower_env.json"

func LoadFromJSON(jsonBytes []byte) (*Config, error) {
	toolConfig := new(Config)
	if err := json.Unmarshal(jsonBytes, &toolConfig); err != nil {
		return nil, err
	}

	{
		parts := strings.SplitN(DefaultFilename, ".", 2)
		viper.SetConfigName(parts[0])
		viper.SetConfigType(parts[1])
	}

	// on *nix this is $HOME/.config/APPNAME/config.json
	// on Windows this is %APPDATA%\APPNAME\config.json
	osConfigPath := filepath.Join(os.Getenv("HOME"), ".config", BaseDirectory)
	if runtime.GOOS == "windows" {
		osConfigPath = filepath.Join(os.Getenv("APPDATA"), BaseDirectory)
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
			filepath.Join(osConfigPath, DefaultFilename),
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
