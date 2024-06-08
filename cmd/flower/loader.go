package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/codeclysm/extract/v3"
	"github.com/evanw/esbuild/pkg/api"
)

// TODO: SUPPORT OTHER GAMES
func initFlowerLoader(config *Config, gameInstallPath string) error {
	for _, game := range config.Games {
		if !strings.Contains(gameInstallPath, game.Subsystem.AppName) {
			continue
		}

		if err := installFlowerLoader(game, gameInstallPath); err != nil {
			return fmt.Errorf("failed to install loader for game: %w", err)
		}

		return nil
	}

	return errors.New("no valid loader found for specified (or detected) game(s)")
}

func installFlowerLoader(game GameConfig, gameInstallPath string) error {
	log.Debug("Loading loader",
		"game", game.Meta.Name,
		"gamePath", gameInstallPath)

	parts := strings.Split(game.Loader.ID, "#")

	// 1) Clone the loader from the specified repository in `game.Loader.ID`
	url := parts[0]
	if strings.HasPrefix(url, GITHUB_PKG) {
		parts = strings.Split(url, "/")
		url = fmt.Sprintf("%s/%s.git", GITHUB_URL, strings.Join(parts[1:], "/"))
		log.Debug("Resolved as GitHub Repository", "url", url)
	} else {
		log.Info("Unsupported Repository Type", "url", url)
		return errors.New("unsupported loader repository type (update flowerLoader tool)")
	}

	sourcePath, err := rootCmd.Flags().GetString("source-path")
	if err != nil {
		return fmt.Errorf("failed to get source-path flag: %w", err)
	}

	clonePath := filepath.Join(sourcePath, "_loader")
	if _, err = cloneOrOpen(context.TODO(), url, clonePath, true); err != nil {
		log.Error("Failed to clone or open repository", "url", url, "path", clonePath, "error", err)
		return fmt.Errorf("failed to clone or open repository: %w", err)
	}

	// 2) Build the loader
	buildPath := filepath.Join(clonePath, "build")
	installPath := filepath.Join(gameInstallPath, "flower")

	// 2.5) Download dependencies for loader
	jsonPath := filepath.Join(clonePath, "package-lock.json")
	if err := installDependencies(jsonPath); err != nil {
		return fmt.Errorf("failed to install loader dependencies: %w", err)
	}

	log.Debug("Building loader",
		"game", game.Meta.Name,
		"sourcePath", clonePath,
		"buildPath", buildPath)
	if err := buildFlowerLoader(game, clonePath, buildPath); err != nil {
		return fmt.Errorf("failed to build loader: %w", err)
	}

	// 3) Install the loader
	log.Debug("Installing loader",
		"game", game.Meta.Name,
		"buildPath", buildPath,
		"installPath", installPath)
	// ...

	if err = os.MkdirAll(filepath.Join(installPath, "flower-plugins"), 0755); err != nil {
		return fmt.Errorf("failed to create flower-plugins directory: %w", err)
	}

	if err = copyAll(
		filepath.Join(clonePath),
		filepath.Join(installPath),
		[]string{"logger.css", "logger.html"},
	); err != nil {
		return fmt.Errorf("failed to copy files: %w", err)
	}

	if err = copyAll(
		filepath.Join(buildPath),
		filepath.Join(installPath),
		[]string{"config.js", "flowerful.js"},
	); err != nil {
		return fmt.Errorf("failed to copy files: %w", err)
	}

	return nil
}

func copyAll(src, dst string, includes []string) error {
	for _, include := range includes {
		srcPath := filepath.Join(src, include)
		dstPath := filepath.Join(dst, include)

		fileData, err := os.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		if err = os.WriteFile(dstPath, fileData, 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		log.Debug("Installed file", "dst", dstPath)
	}

	return nil
}

type Dependency struct {
	Version   string            `json:"version"`
	Resolved  string            `json:"resolved"`
	Integrity string            `json:"integrity"`
	Dev       bool              `json:"dev,omitempty"`
	Requires  map[string]string `json:"requires,omitempty"`
}
type PackageLog struct {
	Dependencies map[string]Dependency `json:"dependencies"`
}

func getDepsForPackage(jsonPath string) (map[string]Dependency, error) {
	js, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read package-lock.json: %w", err)
	}

	var val PackageLog = PackageLog{}
	if err = json.Unmarshal(js, &val); err != nil {
		return nil, fmt.Errorf("failed to unmarshal package-lock.json: %w", err)
	}

	return val.Dependencies, nil
}

func installDependencies(jsonPath string) error {
	allDeps, err := getDepsForPackage(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to get dependencies: %w", err)
	}

	// Create node_modules directory
	if err = os.MkdirAll("node_modules", 0755); err != nil {
		return fmt.Errorf("failed to create node_modules directory: %w", err)
	}

	{
		// Get all dependency names
		keys := make([]string, 0, len(allDeps))
		for key := range allDeps {
			keys = append(keys, key)
		}

		// Remove already installed dependencies
		var banned = []string{"@types/node", "typescript"}
		for _, key := range keys {
			if slices.Contains(banned, key) {
				delete(allDeps, key)
			} else if _, err = os.Stat(filepath.Join("node_modules", key)); err == nil {
				delete(allDeps, key)
			} else {
				log.Debug("Dependency not yet installed", "dep", key)
			}
		}
	}

	if len(allDeps) == 0 {
		fmt.Println("All dependencies are already installed")
		return nil
	}

	setupProgress()
	n := int64(len(allDeps))
	work, done := newTrackerOf(fmt.Sprintf("Installing %d dependencies ...", n), n)
	defer done()

	// Process all dependencies in parallel
	wg := new(sync.WaitGroup)
	log.Warn("Installing dependencies", "count", len(allDeps))
	ch := make(chan string, len(allDeps))
	for dep := range allDeps {
		ch <- dep
	}
	close(ch)

	consumer := func(ch <-chan string) {
		defer wg.Done()

		for {
			select {
			case depName := <-ch:
				if depName == "" { // Skip empty strings
					continue
				}

				if err := installNPMDependency(depName, allDeps[depName]); err != nil {
					log.Error("Failed to install dependency", "dep", depName, "error", err)
				}

				work(1)
			default:
				return
			}
		}
	}

	wg.Add(1)
	go consumer(ch)

	wg.Add(1)
	go consumer(ch)

	wg.Wait()
	return nil
}

func installNPMDependency(name string, info Dependency) error {
	url := info.Resolved
	if strings.HasPrefix(url, "file:") {
		return errors.New("unsupported dependency type: file")
	}

	if strings.HasPrefix(url, "git+") {
		return errors.New("unsupported dependency type: git+")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request for %s: %w", name, err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", name, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download %s: %s", name, res.Status)
	}

	// Extract the tarball
	if err = os.MkdirAll(filepath.Join("node_modules", name), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", name, err)
	}

	ctx := context.TODO()
	if err = extract.Gz(ctx, res.Body, filepath.Join("node_modules", name), func(s string) string {
		// NPM tarballs have a `package/` prefix
		return strings.Replace(s, "package/", "", 1)
	}); err != nil {
		return fmt.Errorf("failed to extract tarball for %s: %w", name, err)
	}

	return nil
}

func buildFlowerLoader(game GameConfig, sourcePath string, outputPath string) error {
	// files, err := os.ReadDir(basePath)
	// if err != nil {
	// 	return err
	// }

	entryPoints := make([]string, 0)
	for _, entry := range game.Loader.Build.Entrypoints {
		entryPoints = append(entryPoints, filepath.Join(sourcePath, entry))
	}

	if game.Loader.Build.ID != "esbuild" {
		return errors.New("unsupported build system")
	}

	logLevel := api.LogLevelInfo
	if rootCmd.Flags().Changed("debug") {
		logLevel = api.LogLevelDebug
	}
	result := api.Build(api.BuildOptions{
		Bundle:        true,
		EntryPoints:   entryPoints,
		Format:        api.FormatESModule,
		LegalComments: api.LegalCommentsEndOfFile,
		LogLevel:      logLevel,
		Platform:      api.PlatformNode,

		Banner: map[string]string{
			"js": "/* Built with https://github.com/flowerLoader */",
		},
		Footer: map[string]string{
			"js": "/* Built with https://github.com/flowerLoader */",
		},
	})

	if len(result.Errors) > 0 {
		fmt.Printf("Failed to transpile %v\n", game.Meta.Name)
		for _, err := range result.Errors {
			fmt.Printf("%s\n", err.Text)
		}

		return errors.New("failed to transpile")
	}

	if len(result.OutputFiles) > 1 {
		return errors.New("unexpected >1 output files")
	}

	resultingJS := result.OutputFiles[0].Contents
	outputFilename := filepath.Join(outputPath, "flowerful.js")
	if err := os.WriteFile(outputFilename, resultingJS, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}
