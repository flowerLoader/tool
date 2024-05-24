package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/cobra"

	"github.com/flowerLoader/tool/pkg/db/types"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"i", "add", "get", "fetch"},
	Short:   "Add a plugin",
	Long:    "Add a plugin from a git repository or local file",
	Example: `flower add FlowerTeam.LimitBreaker`,
	Args:    cobra.MinimumNArgs(1),
	Run:     onAddCommandRun,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func onAddCommandRun(cmd *cobra.Command, args []string) {
	// `plugin` is either a git repository (full URL, or org/repo), or local path
	// If it's a git repository, clone it into the plugins directory
	// If it's a local path, add a reference to it in the plugins directory

	// Parse the plugin name
	name := args[0]
	fullName := parsePluginName(name)
	log.Debug("Resolved Plugin Name", "input", name, "resolved", fullName)

	// Check if the plugin is already installed
	plugin, err := DB.Plugins.Get(fullName)
	if err != nil {
		log.Error("Failed to query plugin database", "error", err)
		return
	}

	if plugin != nil {
		log.Warn("Plugin Already Installed", "name", fullName)
		fmt.Printf("Plugin %s is already installed\n", fullName)

		withoutNS := strings.SplitN(fullName, "/", 2)[1]
		fmt.Printf("To update the plugin, use `flower update %s`\n", withoutNS)
		return
	}

	inputPath, err := rootCmd.Flags().GetString("input-path")
	if err != nil {
		log.Error("Failed to query input path", "error", err)
		return
	}

	if strings.HasPrefix(fullName, "github.com/") {
		if err := installPluginGithub(cmd.Context(), inputPath, fullName); err != nil {
			log.Error("Failed to install GitHub Plugin", "error", err)
		}
		return
	}

	if err := installPluginLocal(cmd, fullName); err != nil {
		log.Error("Failed to install Local Plugin", "error", err)
		return
	}
}

func installPluginGithub(ctx context.Context, pluginRoot, fullName string) error {
	log.Debug("Installing GitHub Plugin", "name", fullName)
	t := time.Now()
	if err := cloneGitPlugin(ctx, "https://github.com", fmt.Sprintf(
		"%s/%s", pluginRoot, fullName), fullName); err != nil {
		return err
	}
	log.Debug("Installing GitHub Plugin", "name", fullName, "took", time.Since(t).String())

	// Add the plugin to the database
	return DB.Plugins.Add(&types.PluginInstallRecord{
		ID:          fullName,
		Enabled:     true,
		InstalledAt: types.FormatTime(time.Now()),
		Path:        fmt.Sprintf("{INPUT}%s", fullName),
	})
}

func installPluginLocal(cmd *cobra.Command, fullName string) error {
	log.Debug("Installing Local Plugin", "name", fullName)

	// Check if the plugin exists
	if _, err := os.Stat(fullName); err != nil {
		// expand using inputPath
		inputPath, err := cmd.Flags().GetString("input-path")
		if err != nil {
			return err
		}

		fullName = fmt.Sprintf("%s/%s", inputPath, fullName)
		if _, err := os.Stat(fullName); err != nil {
			return err
		}
	}

	// Add the plugin to the database
	return DB.Plugins.Add(&types.PluginInstallRecord{
		ID:          fullName,
		Enabled:     true,
		InstalledAt: types.FormatTime(time.Now()),
		Path:        fullName,
	})
}

var (
	ErrPluginAlreadyInstalled = errors.New("plugin is already installed")
)
