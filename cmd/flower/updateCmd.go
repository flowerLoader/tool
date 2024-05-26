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
	Example: `flower update FlowerTeam.LimitBreaker`,
	Args:    cobra.MinimumNArgs(1),
	Run:     onUpdateCommandRun,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func onUpdateCommandRun(cmd *cobra.Command, args []string) {
	inputPath, err := cmd.Flags().GetString("input-path")
	if err != nil {
		log.Error("Failed to get input-path flag", "error", err)
		return
	}

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

	setupProgress()

	if strings.HasPrefix(fullName, GITHUB_PKG) {
		if err := updatePluginGithub(cmd.Context(), inputPath, plugin); err != nil {
			log.Error("Failed to update GitHub Plugin", "error", err)
		}
		return
	}

	if err := updatePluginLocal(inputPath, plugin); err != nil {
		log.Error("Failed to update Local Plugin", "error", err)
		return
	}
}

func updatePluginGithub(ctx context.Context, inputPath string, plugin *types.PluginInstallRecord) error {
	_, done := newTracker("Updating " + plugin.ID)
	log.Debug("Updating GitHub Plugin", "name", plugin.ID)
	t := time.Now()
	fullPath := fmt.Sprintf("%s/%s", inputPath, plugin.ID)
	if err := cloneGitPlugin(ctx, "https://github.com", fullPath, plugin.ID); err != nil {
		return err
	}
	done()
	log.Debug("Updating GitHub Plugin", "name", plugin.ID, "took", time.Since(t).String())

	// Update the plugin in the database
	return DB.Plugins.Update(&types.PluginInstallRecord{
		ID:          plugin.ID,
		Enabled:     true,
		InstalledAt: plugin.InstalledAt,
		UpdatedAt:   types.FormatTime(time.Now()),
	})
}

func updatePluginLocal(inputPath string, plugin *types.PluginInstallRecord) error {
	log.Debug("Updating Local Plugin", "name", plugin.ID)

	// Check if the plugin exists
	fullPath := fmt.Sprintf("%s/%s", inputPath, plugin.ID)
	if _, err := os.Stat(fullPath); err != nil {
		return err
	}

	// Update the plugin in the database
	return DB.Plugins.Update(&types.PluginInstallRecord{
		ID:          plugin.ID,
		Enabled:     true,
		InstalledAt: plugin.InstalledAt,
		UpdatedAt:   types.FormatTime(time.Now()),
	})
}
