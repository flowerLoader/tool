package main

import (
	"os"
	"path/filepath"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:     "clean",
	Aliases: []string{},
	Short:   "Resets the plugin database",
	Long:    "Deletes the plugin database, removing the installed state and all cached plugin information. This will not remove the plugins themselves, nor the cloned sources.",
	Example: `flower clean`,
	Args:    cobra.NoArgs,
	Run:     onCleanCmdRun,
}

func init() {
	cleanCmd.PersistentFlags().BoolP("force", "f", false,
		"Force the reset without prompting for confirmation")

	rootCmd.AddCommand(cleanCmd)
}

func onCleanCmdRun(cmd *cobra.Command, args []string) {
	stat, err := App.DB.Stat()
	if err != nil {
		exit(ErrQueryDB, err)
	}

	log.Info("Plugin Database Stats",
		"cached", stat.Counts.Cached,
		"cloned", stat.Counts.Cloned,
		"enabled", stat.Counts.Enabled)

	if !confirmReset(cmd) {
		return
	}

	// Close the database connection
	if err := App.DB.Close(); err != nil {
		log.Error("failed to close the plugin database", "error", err)
		return
	}

	// TODO: Make a backup?

	// Reset the plugin database
	// TODO: This should be on App.Config not a flag
	gameInstallPath, err := cmd.Flags().GetString("game-path")
	dbPath := filepath.Join(gameInstallPath, "flower.db")
	if err != nil {
		log.Error("failed to query database path", "error", err)
		return
	}
	log.Warn("Resetting plugin database", "path", dbPath)

	if err := os.Remove(dbPath); err != nil {
		log.Error("failed to delete plugin database", "error", err)
		return
	}

	log.Warn("Plugin database reset")
}

func confirmReset(cmd *cobra.Command) bool {
	forced, err := cmd.Flags().GetBool("force")
	if err != nil {
		log.Error("failed to query force flag", "error", err)
		return false
	}

	if !forced && !promptConfirm("Are you sure you want to reset the plugin database?") {
		return false
	}

	return true
}
