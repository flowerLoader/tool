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
	fullName := parsePluginName(args[0])
	log.Debug("Resolved Plugin Name", "input", args[0], "resolved", fullName)

	// Check if the plugin is already installed
	plugin, err := App.DB.Plugins.Get(fullName)
	if err != nil {
		exit(ErrQueryDB, err)
	} else if plugin != nil {
		exit(ErrNameTaken, fullName)
	}

	// Clone the template repository and install it, change the git remote, etc.

	setupProgress()
	if err := installPluginLocal(cmd, fullName); err != nil {
		log.Error("failed to install new plugin", "plugin", fullName, "error", err)
		return
	}
}
