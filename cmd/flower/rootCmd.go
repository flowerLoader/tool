package main

import (
	"path/filepath"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/cobra"

	"github.com/flowerLoader/tool/pkg/db"
)

var DB *db.DB

func init() {
	rootCmd.PersistentFlags().String("db-path", "",
		"Path to the database file (defaults to game's installation directory)")
	rootCmd.PersistentFlags().String("game-path", "",
		"Path to the game's installation directory")
	rootCmd.PersistentFlags().String("input-path", "dist/src",
		"Path to local plugins to transpile")
	rootCmd.PersistentFlags().String("output-path", "dist/obj",
		"Path to store transpiled plugins")
	rootCmd.PersistentFlags().Bool("debug", false, "Scream and shout")
}

var rootCmd = &cobra.Command{
	Use:     APPNAME,
	Version: APPVERSION,
	Short:   "",
	Long:    "",
	Args:    cobra.NoArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if val, err := cmd.Flags().GetBool("debug"); val && err == nil {
			log.InternalLog.SetLevel(log.LevelDebug)
			log.DefaultLog.SetLevel(log.LevelDebug)
		}

		gameInstallPath, err := cmd.Flags().GetString("game-path")
		if err != nil {
			return err
		}

		if gameInstallPath, err = resolveGamePath(gameInstallPath); err != nil {
			return err
		}

		dbPath, err := cmd.Flags().GetString("db-path")
		if err != nil {
			return err
		}
		if dbPath == "" {
			dbPath = filepath.Join(gameInstallPath, "flower.db")
		}
		if DB, err = db.NewDB(dbPath); err != nil {
			return err
		}

		if err := DB.Migrate(); err != nil {
			return err
		}

		if err := cmd.Root().PersistentFlags().Set("game-path", gameInstallPath); err != nil {
			return err
		}

		if err := cmd.Root().PersistentFlags().Set("db-path", dbPath); err != nil {
			return err
		}

		log.Info("Using Flags (post resolution)",
			"game-path", gameInstallPath,
			"db-path", dbPath)

		// if err := initFlowerLoader(gameInstallPath); err != nil {
		// 	return err
		// }

		return nil
	},
}
