package main

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/flowerLoader/tool/pkg/db"
)

//const PLUGINS_ROOT = "flowerful-plugins"
////go:embed all:installer
// var assets embed.FS

var (
	gameInstallPath  string
	pluginInputPath  string
	pluginOutputPath string

	DB *db.DB
)

func init() {
	rootCmd.PersistentFlags().StringVar(&gameInstallPath, "game-path", "",
		"Path to the game's installation directory")
	rootCmd.PersistentFlags().StringVar(&pluginInputPath, "input-path", "dist/src",
		"Path to local plugins to transpile")
	rootCmd.PersistentFlags().StringVar(&pluginOutputPath, "output-path", "dist/obj",
		"Path to store transpiled plugins")
	rootCmd.PersistentFlags().Bool("debug", false, "Scream and shout")

	must := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	must(viper.BindPFlag("game-path", rootCmd.PersistentFlags().Lookup("game-path")))
	must(viper.BindPFlag("input-path", rootCmd.PersistentFlags().Lookup("input-path")))
	must(viper.BindPFlag("output-path", rootCmd.PersistentFlags().Lookup("output-path")))
	must(viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")))
}

var rootCmd = &cobra.Command{
	Use:     APPNAME,
	Version: APPVERSION,
	Short:   "",
	Long:    "",
	Args:    cobra.NoArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if gameInstallPath, err = resolveGamePath(gameInstallPath); err != nil {
			return err
		}

		viper.Set("game-path", gameInstallPath)

		dbPath := filepath.Join(gameInstallPath, "flower.db")
		if DB, err = db.NewDB(dbPath); err != nil {
			return err
		}

		if err := DB.Migrate(); err != nil {
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
	},
}
