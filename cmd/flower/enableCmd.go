package main

import (
	"fmt"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/cobra"
)

var enableCmd = &cobra.Command{
	Use:     "enable",
	Aliases: []string{"en"},
	Short:   "Enable a plugin",
	Long:    "Mark an installed plugin as enabled in the database",
	Example: `flower enable LimitBreaker`,
	Args:    cobra.MinimumNArgs(1),
	Run:     onEnableCommandRun,
}

func init() {
	rootCmd.AddCommand(enableCmd)
}

func onEnableCommandRun(cmd *cobra.Command, args []string) {
	fullName := parsePluginName(args[0])
	log.Debug("Resolved Plugin Name", "input", args[0], "resolved", fullName)

	// Check if the plugin is already installed
	plugin, err := App.DB.Plugins.Get(fullName)
	if err != nil {
		exit(ErrQueryDB, err)
	} else if plugin == nil {
		exit(ErrNotInstalled, fullName)
	}

	// Mark the plugin as enabled
	plugin.Enabled = true
	if err := App.DB.Plugins.Update(plugin); err != nil {
		log.Error("Failed to update plugin status in database", "error", err)
		return
	}

	log.Info("Plugin Enabled", "name", fullName)
	fmt.Printf("Plugin %s has been enabled\n", fullName)
}
