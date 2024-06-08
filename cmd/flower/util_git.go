package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/go-git/go-git/v5"
)

func cloneGitPlugin(ctx context.Context, baseURL, clonePath, fullName string) error {
	parts := strings.Split(fullName, "/")
	if len(parts) < 2 {
		return errors.New("invalid plugin name")
	}

	// Get the owner and repo name
	repoName, owner := parts[len(parts)-1], parts[len(parts)-2]
	if i := strings.Index(owner, "."); i != -1 {
		owner = owner[:i] // Remove the ".git" suffix if it exists
	}

	// Check if the plugin is already cloned
	gitURL := fmt.Sprintf("%s/%s/%s.git", baseURL, owner, repoName)
	log.Debug("Plugin from Git", "name", fullName, "url", gitURL, "path", clonePath)

	t := time.Now()
	repo, err := cloneOrOpen(ctx, gitURL, clonePath, true)
	if err != nil {
		return err
	}
	log.Debug("Plugin from Git", "name", fullName, "took", time.Since(t).String())

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

	log.Info("Plugin Installed from Git",
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
		if err := wt.Checkout(&git.CheckoutOptions{
			Branch: "refs/heads/main",
			Create: false,
			Force:  true,
			Keep:   false,
		}); err != nil && err != git.NoErrAlreadyUpToDate {
			return repo, fmt.Errorf("failed to checkout: %w", err)
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
			Progress:          os.Stdout,
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
