package main

import (
	"fmt"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "uninstall", "delete"},
	Short:   "Remove a plugin",
	Long:    "Remove a plugin by name from the local database. This will not delete the plugin files.",
	Example: `flower remove FlowerTeam.LimitBreaker`,
	Args:    cobra.MinimumNArgs(1),
	Run:     onRemoveCommandRun,
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.PersistentFlags().BoolP("force", "f", false, "Force the removal without prompting for confirmation")
}

func onRemoveCommandRun(cmd *cobra.Command, args []string) {
	name := args[0]
	fullName := parsePluginName(name)
	log.Debug("Resolved Plugin Name", "input", name, "resolved", fullName)

	// Check if the plugin is installed
	plugin, err := App.DB.Plugins.Get(fullName)
	if err != nil {
		log.Error("Failed to query plugin database", "error", err)
		return
	}

	if plugin == nil {
		log.Warn("Plugin Not Installed", "name", fullName)
		fmt.Printf("Plugin %s is not installed\n", fullName)
		return
	}

	// Prompt the user to confirm the removal
	forced, err := cmd.Flags().GetBool("force")
	if err != nil {
		log.Error("Failed to query force flag", "error", err)
		return
	}

	if !forced {
		if !promptConfirm(fmt.Sprintf("Are you sure you want to remove the plugin %s?", fullName)) {
			return
		}
	}

	// Remove the plugin from the database
	if err := App.DB.Plugins.Remove(fullName); err != nil {
		log.Error("Failed to remove plugin from database", "error", err)
		return
	}

	log.Info("Plugin Removed", "name", fullName)
	fmt.Printf("Plugin %s has been removed\n", fullName)
}
