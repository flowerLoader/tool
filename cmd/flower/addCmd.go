package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath" // Added for cross-platform path handling
	"strings"
	"time"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/cobra"

	"github.com/flowerLoader/tool/pkg/db/types"
)

const GITHUB_URL = "https://github.com"
const GITHUB_PKG = "github.com"

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a", "install", "get", "fetch"},
	Short:   "Add a plugin",
	Long:    "Add a plugin from a git repository or local file",
	Example: `flower add LimitBreaker`,
	Args:    cobra.MinimumNArgs(1),
	Run:     onAddCommandRun,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func onAddCommandRun(cmd *cobra.Command, args []string) {
	fullName := parsePluginName(args[0])
	log.Debug("Resolved Plugin Name", "input", args[0], "resolved", fullName)

	// Check if the plugin is already installed
	if plugin, err := App.DB.Plugins.Get(fullName); err != nil {
		exit(ErrQueryDB, err)
	} else if plugin != nil {
		exit(ErrAlreadyInstalled, fullName)
	}

	sourcePath, err := rootCmd.Flags().GetString("source-path")
	if err != nil {
		log.Error("Failed to query source-path", "error", err)
		return
	}

	setupProgress()

	if strings.HasPrefix(fullName, GITHUB_PKG) {
		if err := installPluginGithub(cmd.Context(), sourcePath, fullName); err != nil {
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
	_, done := newTracker("Installing " + fullName)
	log.Debug("Installing GitHub Plugin", "name", fullName)
	t := time.Now()
	clonePath := filepath.Join(pluginRoot, fullName)
	if err := cloneGitPlugin(ctx, GITHUB_URL, clonePath, fullName); err != nil {
		return err
	}
	done()
	log.Debug("Installing GitHub Plugin", "name", fullName, "took", time.Since(t).String())

	// Add the plugin to the database
	return App.DB.Plugins.Add(&types.PluginInstallRecord{
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
		// expand using sourcePath
		sourcePath, err := cmd.Flags().GetString("source-path")
		if err != nil {
			return err
		}

		fullName = filepath.Join(sourcePath, fullName)
		if _, err := os.Stat(fullName); err != nil {
			return err
		}
	}

	// Add the plugin to the database
	return App.DB.Plugins.Add(&types.PluginInstallRecord{
		ID:          fullName,
		Enabled:     true,
		InstalledAt: types.FormatTime(time.Now()),
		Path:        fullName,
	})
}
