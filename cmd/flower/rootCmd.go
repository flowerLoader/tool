package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//const PLUGINS_ROOT = "flowerful-plugins"
////go:embed all:installer
// var assets embed.FS

var (
	gameInstallPath  string
	pluginInputPath  string
	pluginOutputPath string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&gameInstallPath, "game-path", "",
		"Path to the game's installation directory")
	rootCmd.PersistentFlags().StringVar(&pluginInputPath, "input-path", "source",
		"Path to local plugins to transpile")
	rootCmd.PersistentFlags().StringVar(&pluginOutputPath, "output-path", "obj",
		"Path to store transpiled plugins")
	rootCmd.PersistentFlags().Bool("debug", false, "Scream and shout")

	viper.BindPFlag("game-path", rootCmd.PersistentFlags().Lookup("game-path"))
	viper.BindPFlag("input-path", rootCmd.PersistentFlags().Lookup("input-path"))
	viper.BindPFlag("output-path", rootCmd.PersistentFlags().Lookup("output-path"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

var rootCmd = &cobra.Command{
	Use:     APPNAME,
	Version: APPVERSION,
	Short:   "",
	Long:    "",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if gameInstallPath, err = resolveGamePath(gameInstallPath); err != nil {
			return err
		}

		viper.Set("game-path", gameInstallPath)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}