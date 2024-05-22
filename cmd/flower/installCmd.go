package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"i", "add", "get", "fetch"},
	Short:   "Install a plugin",
	Long:    "Install a plugin from a git repository or local file",
	Example: `flower install FlowerTeam.LimitBreaker`,
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		plugin := args[0]
		return installPlugin(cmd.Context(), plugin)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func installPlugin(ctx context.Context, name string) error {
	// `plugin` is either a git repository (full URL, or org/repo), or local path
	// If it's a git repository, clone it into the plugins directory
	// If it's a local path, add a reference to it in the plugins directory

	// Parse the plugin name
	fullName := parsePluginName(name)
	log.Debug("Resolved Plugin Name", "input", name, "resolved", fullName)

	// Check if the plugin is already installed
	// plugin, err := DB.Plugins.Get(fullName)
	// if err != nil {
	// 	return err
	// }

	// if plugin != nil {
	// 	return ErrPluginAlreadyInstalled
	// }

	if strings.HasPrefix(fullName, "github.com/") {
		return installPluginGithub(ctx, fullName)
	}

	// return installPluginLocal(ctx, fullName)
	return nil
}

func installPluginGithub(ctx context.Context, fullName string) error {
	log.Debug("Installing GitHub Plugin", "name", fullName)
	t := time.Now()
	if err := cloneGitHubPlugin(ctx, fullName); err != nil {
		return err
	}
	log.Debug("Installing GitHub Plugin", "name", fullName, "took", time.Since(t).String())

	// Add the plugin to the database
	// return DB.Plugins.Add(&types.PluginCacheRecord{
	// 	ID: fullName,
	// })
	return nil
}

// func installPluginLocal(ctx context.Context, fullName string) error {
// 	// Add the plugin to the database
// 	return DB.Plugins.Add(&types.PluginCacheRecord{
// 		ID: fullName,
// 	})
// }

func cloneGitHubPlugin(ctx context.Context, fullName string) error {
	var (
		clonePath = fmt.Sprintf("%s/%s", viper.GetString("input-path"), fullName)
		parts     = strings.Split(fullName, "/")
	)

	if len(parts) < 2 {
		return errors.New("invalid plugin name")
	}

	// Get the owner and repo name
	repoName, owner := parts[len(parts)-1], parts[len(parts)-2]
	if i := strings.Index(owner, "."); i != -1 {
		owner = owner[:i] // Remove the ".git" suffix if it exists
	}

	// Check if the plugin is already cloned
	gitURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, repoName)
	log.Debug("Plugin from GitHub", "name", fullName, "url", gitURL, "path", clonePath)

	t := time.Now()
	repo, err := cloneOrOpen(ctx, gitURL, clonePath, true)
	if err != nil {
		return err
	}
	log.Debug("Plugin from GitHub", "name", fullName, "took", time.Since(t).String())

	// Get the latest commit
	ref, err := repo.Head()
	if err != nil {
		return err
	}

	// Get the Branch
	branch := strings.TrimPrefix(ref.Name().String(), "refs/heads/")
	if branch == "" {
		branch = "main"
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return err
	}

	log.Info("GitHub Plugin Installed",
		"name", fullName, "branch", branch,
		"commit", commit.Hash.String(),
		"message", commit.Message,
		"author", commit.Author.String(),
		"date", commit.Author.When)

	return nil
}

// cloneOrOpen clones a git repository if it doesn't exist, or opens it if it
// does. If `pull` is true, it will attempt to pull the latest changes from the
// remote repository, otherwise it will be left as-is. The repository is
// returned regardless of whether it was cloned or opened. If an error occurs
// during the process, the repository will be returned along with the error.
func cloneOrOpen(ctx context.Context, url, path string, pull bool) (*git.Repository, error) {
	if repo, err := git.PlainOpen(path); err == nil {
		log.Debug("Opening Existing Repository", "path", path, "url", url, "pull", pull)
		if !pull {
			return repo, nil
		}

		// Check if the plugin is up-to-date
		wt, err := repo.Worktree()
		if err != nil {
			return repo, err
		}

		log.Debug("Pulling Repository", "path", path, "url", url)
		t := time.Now()
		if err := wt.PullContext(ctx, &git.PullOptions{
			RemoteName: "origin",
		}); err != nil && err != git.NoErrAlreadyUpToDate {
			return repo, err
		}

		log.Debug("Pulling Repository", "path", path, "took", time.Since(t).String())
		return repo, nil
	} else if err != git.ErrRepositoryNotExists {
		return nil, err
	}

	log.Debug("Cloning Repository", "path", path, "url", url)
	t := time.Now()
	repo, err := git.PlainCloneContext(
		ctx, path, false, &git.CloneOptions{
			Depth:             1,
			Progress:          nil,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			RemoteName:        "origin",
			ShallowSubmodules: true,
			Tags:              git.NoTags,
			URL:               url,
			// ReferenceName:     "",
			// SingleBranch:      true,
		})
	if err != nil {
		return repo, err
	}

	log.Debug("Cloning Repository", "path", path, "took", time.Since(t).String())
	return repo, nil
}

// parsePluginName takes a partial plugin name (full URL, org/repo, or local
// path) and returns the full name of the plugin ({github.com|local}/org/repo)
func parsePluginName(name string) string {
	u, err := url.Parse(name)
	if err == nil && u.Scheme != "" {
		// https://www.github.com/flowerLoader/tool
		// -> github.com/flowerLoader/tool
		return strings.TrimPrefix(u.Hostname(), "www.") + u.Path
	}

	// /path/to/plugin (or C:\path\to\plugin)
	// -> local/pluginAuthor/pluginName
	if strings.HasPrefix(name, "/") || strings.Contains(name, ":\\") {
		return "local/" + strings.TrimPrefix(name, "/")
	}

	// flowerLoader/tool
	// -> github.com/flowerLoader/tool
	return "github.com/" + name
}

var (
	ErrPluginAlreadyInstalled = errors.New("plugin is already installed")
)
