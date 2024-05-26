package main

import (
	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{},
	Short:   "Create a new plugin",
	Long:    "Create a new plugin from a template repository. Queries the user for the plugin name, description, etc. and then clones the template repository into the plugins directory.",
	Example: `flower new`,
	Args:    cobra.ExactArgs(1),
	Run:     onNewPluginRun,
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func onNewPluginRun(cmd *cobra.Command, args []string) {
	// Parse the plugin name
	name := args[0]
	fullName := parsePluginName(name)
	log.Debug("Resolved Plugin Name", "input", name, "resolved", fullName)

	// Check if the plugin is already installed
	plugin, err := DB.Plugins.Get(fullName)
	if err != nil {
		log.Fatal("failed to query local plugin database", "name", fullName, "error", err)
	}

	if plugin != nil {
		log.Error("A plugin with the same name already exists", "name", fullName)
		return
	}

	// Clone the template repository and install it, change the git remote, etc.

	setupProgress()
	if err := installPluginLocal(cmd, fullName); err != nil {
		log.Fatal("failed to install new plugin", "plugin", fullName, "error", err)
	}
}
