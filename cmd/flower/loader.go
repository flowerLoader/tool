package main

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"net/url"
// 	"os"
// 	"path/filepath"
// 	"strings"

// 	log "github.com/AlbinoGeek/logxi/v1"
// 	"github.com/evanw/esbuild/pkg/api"
// 	"github.com/go-git/go-git/v5"
// 	"github.com/go-git/go-git/v5/config"
// 	"github.com/go-git/go-git/v5/plumbing"
// )

// // TODO: SUPPORT OTHER GAMES
// func initFlowerLoader(gameInstallPath string) error {
// 	config, err := NewConfig()
// 	if err != nil {
// 		return fmt.Errorf("failed to load config: %w", err)
// 	}

// 	if len(config.Games) == 0 {
// 		return errors.New("no games found in config")
// 	}

// 	for _, game := range config.Games {
// 		if !strings.Contains(gameInstallPath, game.Subsystem.Appname) {
// 			continue
// 		}

// 		if err := installFlowerLoader(game, gameInstallPath); err != nil {
// 			return fmt.Errorf("failed to load loader for game: %w", err)
// 		}
// 	}

// 	return errors.New("no valid loader found for specified (or detected) game(s)")
// }

// func installFlowerLoader(game GameConfig, basePath string) error {
// 	log.Debug("Loading loader", "game", game.Meta.Name, "base", basePath)

// 	parts := strings.Split(game.Loader.ID, "#")
// 	args, err := url.ParseQuery(parts[1])
// 	if err != nil {
// 		return fmt.Errorf("failed to parse loader args: %w", err)
// 	}

// 	// 1) Clone the loader from the specified repository in `game.Loader.ID`
// 	url := parts[0]
// 	if strings.HasPrefix(url, GITHUB_PKG) {
// 		parts = strings.Split(url, "/")
// 		url = fmt.Sprintf("%s/%s.git", GITHUB_URL, strings.Join(parts[1:], "/"))
// 		log.Debug("Resolved as GitHub Repository", "url", url, "args", args)
// 	} else {
// 		log.Info("Unsupported Repository Type", "url", url, "args", args)
// 		return errors.New("unsupported loader repository type (update flowerLoader tool)")
// 	}

// 	clonePath := filepath.Join(basePath, "flower")
// 	repo, err := cloneOrOpen(context.TODO(), url, clonePath, true)
// 	if err != nil {
// 		log.Error("Failed to clone or open repository", "url", url, "path", clonePath, "error", err)
// 		return fmt.Errorf("failed to clone or open repository: %w", err)
// 	}

// 	if branch := args.Get("branch"); branch != "" {
// 		tree, err := repo.Worktree()
// 		if err != nil {
// 			return fmt.Errorf("failed to get worktree: %w", err)
// 		}

// 		// We first need to fetch the branch, as we shallow cloned
// 		log.Debug("Fetching", "branch", branch)
// 		if err := repo.Fetch(&git.FetchOptions{
// 			RefSpecs: []config.RefSpec{
// 				config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/remotes/origin/%s", branch, branch)),
// 			},
// 		}); err != nil {
// 			return fmt.Errorf("failed to fetch branch: %w", err)
// 		}

// 		log.Debug("Checking out", "branch", branch)
// 		if err := tree.Checkout(&git.CheckoutOptions{
// 			Branch: plumbing.ReferenceName(fmt.Sprintf("refs/remotes/origin/%s", branch)),
// 		}); err != nil {
// 			return fmt.Errorf("failed to checkout branch: %w", err)
// 		}
// 	}

// 	// 2) Build the loader
// 	log.Debug("Building loader", "game", game.Meta.Name, "base", basePath)
// 	if err := buildFlowerLoader(game, clonePath); err != nil {
// 		return fmt.Errorf("failed to build loader: %w", err)
// 	}

// 	// 3) Install the loader
// 	log.Debug("Installing loader", "game", game.Meta.Name, "base", basePath)
// 	// ...

// 	return nil
// }

// func buildFlowerLoader(game GameConfig, basePath string) error {
// 	// files, err := os.ReadDir(basePath)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	entryPoints := make([]string, 0)
// 	for _, entry := range game.Loader.Build.Entrypoints {
// 		entryPoints = append(entryPoints, filepath.Join(basePath, entry))
// 	}

// 	if game.Loader.Build.ID != "esbuild" {
// 		return errors.New("unsupported build system")
// 	}

// 	result := api.Build(api.BuildOptions{
// 		Bundle:        true,
// 		EntryPoints:   entryPoints,
// 		Format:        api.FormatESModule,
// 		LegalComments: api.LegalCommentsEndOfFile,
// 		LogLevel:      api.LogLevelDebug,
// 		Platform:      api.PlatformNode,

// 		Banner: map[string]string{
// 			"js": "/* Built with https://github.com/flowerLoader */",
// 		},
// 		Footer: map[string]string{
// 			"js": "/* Built with https://github.com/flowerLoader */",
// 		},
// 	})

// 	if len(result.Errors) > 0 {
// 		fmt.Printf("Failed to transpile %v\n", game.Meta.Name)
// 		for _, err := range result.Errors {
// 			fmt.Printf("%s\n", err.Text)
// 		}

// 		return errors.New("failed to transpile")
// 	}

// 	if len(result.OutputFiles) > 1 {
// 		return errors.New("unexpected >1 output files")
// 	}

// 	resultingJS := result.OutputFiles[0].Contents
// 	outputFilename := filepath.Join(basePath, "dist", "flowerful.js")
// 	if err := os.WriteFile(outputFilename, resultingJS, 0644); err != nil {
// 		return fmt.Errorf("failed to write output file: %w", err)
// 	}

// 	return nil
// }
