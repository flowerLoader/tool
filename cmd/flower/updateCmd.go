package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/cobra"

	"github.com/flowerLoader/tool/pkg/db/types"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"up", "upgrade"},
	Short:   "Update a plugin",
	Long:    "Update an installed plugin to the latest version",
	Example: `flower update LimitBreaker`,
	Args:    cobra.MinimumNArgs(1),
	Run:     onUpdateCommandRun,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func onUpdateCommandRun(cmd *cobra.Command, args []string) {
	// TODO: This should be via App.Config instead of a flag.
	sourcePath, err := cmd.Flags().GetString("source-path")
	if err != nil {
		log.Error("Failed to get source-path flag", "error", err)
		return
	}

	name := args[0]
	if strings.ToLower(name) == "all" {
		updateAllPlugins(cmd, sourcePath)
		return
	}

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

	setupProgress()

	if strings.HasPrefix(fullName, GITHUB_PKG) {
		if err := updatePluginGithub(cmd.Context(), sourcePath, plugin); err != nil {
			log.Error("Failed to update GitHub Plugin", "error", err)
		}
		return
	}

	if err := updatePluginLocal(sourcePath, plugin); err != nil {
		log.Error("Failed to update Local Plugin", "error", err)
		return
	}
}

func updateAllPlugins(cmd *cobra.Command, sourcePath string) {
	plugins, err := App.DB.Plugins.List()
	if err != nil {
		log.Error("Failed to list installed plugins", "error", err)
		return
	}

	setupProgress()

	for _, plugin := range plugins {
		fullName := plugin.ID
		log.Debug("Updating Plugin", "name", fullName)

		if strings.HasPrefix(fullName, GITHUB_PKG) {
			if err := updatePluginGithub(cmd.Context(), sourcePath, plugin); err != nil {
				log.Error("Failed to update GitHub Plugin", "error", err)
			}
		} else {
			if err := updatePluginLocal(sourcePath, plugin); err != nil {
				log.Error("Failed to update Local Plugin", "error", err)
			}
		}
	}
}

func updatePluginGithub(ctx context.Context, sourcePath string, plugin *types.PluginInstallRecord) error {
	_, done := newTracker("Updating " + plugin.ID)
	log.Debug("Updating GitHub Plugin", "name", plugin.ID)
	t := time.Now()
	fullPath := fmt.Sprintf("%s/%s", sourcePath, plugin.ID)
	if err := cloneGitPlugin(ctx, GITHUB_URL, fullPath, plugin.ID); err != nil {
		return err
	}
	done()
	log.Debug("Updating GitHub Plugin", "name", plugin.ID, "took", time.Since(t).String())

	// Update the plugin in the database
	return App.DB.Plugins.Update(&types.PluginInstallRecord{
		ID:          plugin.ID,
		Enabled:     true,
		InstalledAt: plugin.InstalledAt,
		UpdatedAt:   types.FormatTime(time.Now()),
	})
}

func updatePluginLocal(sourcePath string, plugin *types.PluginInstallRecord) error {
	log.Debug("Updating Local Plugin", "name", plugin.ID)

	// Check if the plugin exists
	fullPath := fmt.Sprintf("%s/%s", sourcePath, plugin.ID)
	if _, err := os.Stat(fullPath); err != nil {
		return err
	}

	// Update the plugin in the database
	return App.DB.Plugins.Update(&types.PluginInstallRecord{
		ID:          plugin.ID,
		Enabled:     true,
		InstalledAt: plugin.InstalledAt,
		UpdatedAt:   types.FormatTime(time.Now()),
	})
}
