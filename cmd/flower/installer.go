package main

import (
	"errors"
	"os"
	"path/filepath"

	log "github.com/AlbinoGeek/logxi/v1"

	"github.com/flowerLoader/tool/pkg/steam"
)

var (
	errGameNotFound = errors.New("game directory not found")

	pathParts       = filepath.Join("steamapps", "common", "isekainosouzousha")
	pathPartLinux   = filepath.Join("gamedata", "game")
	pathPartWindows = "game"
)

// findGameInstallationPath attempts to find the game's installation path by
// looking at the Steam library folders. If the game is not found in any of the
// library folders, an error is returned.
func findGameInstallationPath() (string, error) {
	st := steam.NewSteam()
	if err := st.Find(); err != nil {
		log.Error("failed to locate Steam installation", "error", err)
		return "", err
	}

	log.Debug("findGameInstallationPath", "libraryFolders", st.LibraryFolders)
	for _, libraryFolder := range st.LibraryFolders {
		path := filepath.Join(libraryFolder, pathParts)
		if stat, err := os.Stat(path); err == nil && stat.IsDir() {
			log.Debug("findGameInstallationPath: found", "path", path)
			return path, nil
		}
	}

	return "", errGameNotFound
}

// firstExisting returns the first existing path from the given paths. If none
// of the paths exist, an empty string is returned instead. No other property
// of a given path is checked, simply the existence as-is.
func firstExisting(paths ...string) string {
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// resolveGamePath checks if the given path is valid, resolving it to the game's
// nw.js HTML directory. If the path is empty, the game's installation path will
// be detected using Steam's library folders. If the game is not found in any of
// the library folders, or the given path is invalid, an error is returned.
func resolveGamePath(path string) (resolvedPath string, err error) {
	log.Debug("resolveGamePath: start", "path", path)
	if path == "" {
		if path, err = findGameInstallationPath(); err != nil {
			return "", err // error already logged
		}
	}

	linuxPath := filepath.Join(path, pathPartLinux)
	windowsPath := filepath.Join(path, pathPartWindows)
	if path = firstExisting(linuxPath, windowsPath); path == "" {
		log.Error("failed to resolve game installation path, guessed paths:",
			"linux", linuxPath, "windows", windowsPath)
		return "", errGameNotFound
	}

	log.Debug("resolveGamePath: resolved", "path", path)
	return path, nil
}
