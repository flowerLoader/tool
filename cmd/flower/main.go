package main

import (
	"os"
	"path/filepath"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	APPNAME    = "flower"
	APPVERSION = "0.1.0"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AutomaticEnv()
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if dir, err := os.UserConfigDir(); err == nil {
		viper.AddConfigPath(filepath.Join(dir, "flower"))
	}

	if dir, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(filepath.Join(dir, ".flower"))
	}

	if err := viper.ReadInConfig(); err == nil {
		log.Debug("Loaded Configuration", "file", viper.ConfigFileUsed())
	}

	if viper.GetBool("debug") {
		log.DefaultLog.SetLevel(log.LevelDebug)
		log.InternalLog.SetLevel(log.LevelDebug)
	} else {
		log.DefaultLog.SetLevel(log.LevelInfo)
		log.InternalLog.SetLevel(log.LevelInfo)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		println(os.Stderr, err)
		os.Exit(1)
	}
}
