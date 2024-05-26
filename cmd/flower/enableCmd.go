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
	Example: `flower enable FlowerTeam.LimitBreaker`,
	Args:    cobra.MinimumNArgs(1),
	Run:     onEnableCommandRun,
}

func init() {
	rootCmd.AddCommand(enableCmd)
}

func onEnableCommandRun(cmd *cobra.Command, args []string) {
	name := args[0]
	fullName := parsePluginName(name)
	log.Debug("Resolved Plugin Name", "input", name, "resolved", fullName)

	// Check if the plugin is installed
	plugin, err := DB.Plugins.Get(fullName)
	if err != nil {
		log.Error("Failed to query plugin database", "error", err)
		return
	}

	if plugin == nil {
		log.Warn("Plugin Not Installed", "name", fullName)
		fmt.Printf("Plugin %s is not installed\n", fullName)
		return
	}

	// Mark the plugin as enabled
	plugin.Enabled = true
	if err := DB.Plugins.Update(plugin); err != nil {
		log.Error("Failed to update plugin status in database", "error", err)
		return
	}

	log.Info("Plugin Enabled", "name", fullName)
	fmt.Printf("Plugin %s has been enabled\n", fullName)
}
