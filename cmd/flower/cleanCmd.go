package main

import (
	"os"

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
	rootCmd.AddCommand(cleanCmd)
}

func onCleanCmdRun(cmd *cobra.Command, args []string) {
	stat, err := DB.Stat()
	if err != nil {
		log.Error("failed to query plugin database stats", "error", err)
		return
	}

	log.Info("Plugin Database Stats",
		"cached", stat.Counts.Cached,
		"cloned", stat.Counts.Cloned,
		"enabled", stat.Counts.Enabled)

	// Prompt the user to confirm the reset
	if !promptConfirm("Are you sure you want to reset the plugin database?") {
		return
	}

	// Close the database connection
	if err := DB.Close(); err != nil {
		log.Error("failed to close the plugin database", "error", err)
		return
	}

	// Reset the plugin database
	dbPath, err := cmd.Flags().GetString("db-path")
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
