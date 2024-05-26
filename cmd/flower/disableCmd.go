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
	Example: `flower disable FlowerTeam.LimitBreaker`,
	Args:    cobra.MinimumNArgs(1),
	Run:     onDisableCommandRun,
}

func init() {
	rootCmd.AddCommand(disableCmd)
}

func onDisableCommandRun(cmd *cobra.Command, args []string) {
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

	// Mark the plugin as disabled
	plugin.Enabled = false
	if err := DB.Plugins.Update(plugin); err != nil {
		log.Error("Failed to update plugin status in database", "error", err)
		return
	}

	log.Info("Plugin Disabled", "name", fullName)
	fmt.Printf("Plugin %s has been disabled\n", fullName)
}
