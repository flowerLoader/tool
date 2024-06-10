package main

import (
	"path/filepath"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/cobra"

	"github.com/flowerLoader/tool/pkg/db"
)

type Application struct {
	Config *Config
	DB     *db.DB
}

var App Application

func init() {
	rootCmd.PersistentFlags().String("game-path", "", "Path to the game's installation directory")
	rootCmd.PersistentFlags().String("source-path", "dist/src", "Path to local plugin source code")
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
			log.Info("Debugging Enabled (by flag)")
		}

		App.Config, err = NewConfig()
		if err != nil || len(App.Config.Games) == 0 {
			panic("fatal: no games found in config (check main.json and rebuild) error: " + err.Error())
		}

		gameInstallPath, err := cmd.Flags().GetString("game-path")
		if err != nil {
			return err
		}

		if gameInstallPath, err = resolveGamePath(gameInstallPath); err != nil {
			return err
		}

		// TODO: dbPath should be derived from App.Config.Environment[i] where
		// TODO   i == id of resolved gameInstallPath, or new if not found
		dbPath := filepath.Join(gameInstallPath, "flower.db")
		if App.DB, err = db.NewDB(dbPath); err != nil {
			return err
		}

		if err := App.DB.Migrate(); err != nil {
			return err
		}

		if err := cmd.Root().PersistentFlags().Set("game-path", gameInstallPath); err != nil {
			return err
		}

		// TODO: This should be stored on App.Config somewhere instead of a flag (see above)
		log.Info("Using Environment Configuration",
			"db-path", dbPath,
			"game-path", gameInstallPath,
			"source-path", cmd.Flag("source-path").Value.String())

		return nil
	},
}
