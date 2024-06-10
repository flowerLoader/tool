package main

import (
	"fmt"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/cobra"
)

var disableCmd = &cobra.Command{
	Use:     "disable",
	Aliases: []string{"dis"},
	Short:   "Disable a plugin",
	Long:    "Mark an installed plugin as disabled in the database",
	Example: `flower disable LimitBreaker`,
	Args:    cobra.MinimumNArgs(1),
	Run:     onDisableCommandRun,
}

func init() {
	rootCmd.AddCommand(disableCmd)
}

func onDisableCommandRun(cmd *cobra.Command, args []string) {
	fullName := parsePluginName(args[0])
	log.Debug("Resolved Plugin Name", "input", args[0], "resolved", fullName)

	// Check if the plugin is already installed
	plugin, err := App.DB.Plugins.Get(fullName)
	if err != nil {
		exit(ErrQueryDB, err)
	} else if plugin == nil {
		exit(ErrNotInstalled, fullName)
	}

	// Mark the plugin as disabled
	plugin.Enabled = false
	if err := App.DB.Plugins.Update(plugin); err != nil {
		log.Error("Failed to update plugin status in database", "error", err)
		return
	}

	log.Info("Plugin Disabled", "name", fullName)
	fmt.Printf("Plugin %s has been disabled\n", fullName)
}
