package ts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/AlbinoGeek/logxi/v1"
)

func detectEntrypoint(sourcePath string) (string, error) {
	fi, err := os.Stat(sourcePath)
	if err != nil {
		return "", fmt.Errorf("failed to stat source path: %w", err)
	}

	// If it's a single file, return it
	if !fi.IsDir() {
		return sourcePath, nil
	}

	// Check for a package.json and use the main field
	pkgFile := filepath.Join(sourcePath, "package.json")
	log.Debug("Checking for package.json", "path", pkgFile)
	if _, err := os.Stat(pkgFile); err == nil {
		type packageJSON struct {
			Main string `json:"main"`
		}

		fileData, err := os.ReadFile(pkgFile)
		if err != nil {
			return "", fmt.Errorf("failed to read package.json: %w", err)
		}

		pkg := packageJSON{}
		if err := json.Unmarshal(fileData, &pkg); err != nil {
			return "", fmt.Errorf("failed to parse package.json: %w", err)
		}

		if pkg.Main != "" {
			mainPath := filepath.Join(sourcePath, pkg.Main)
			log.Debug("Using main field from package.json", "main", mainPath)
			if _, err := os.Stat(mainPath); err == nil {
				return mainPath, nil
			}

			return "", fmt.Errorf("main field in package.json does not exist: %s", mainPath)
		}
	}

	// Check for an index.ts file
	indexTS := filepath.Join(sourcePath, "index.ts")
	if _, err := os.Stat(indexTS); err == nil {
		return indexTS, nil
	}

	// Check if the directory contains only one .ts file, if so, use it
	files, err := os.ReadDir(sourcePath)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	var tsFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) == ".ts" {
			tsFiles = append(tsFiles, file.Name())
		}
	}

	if len(tsFiles) == 1 {
		log.Debug("Using sole .ts file as entrypoint", "file", tsFiles[0])
		return filepath.Join(sourcePath, tsFiles[0]), nil
	}

	// Check for an entrypoint with the same name as the base directory (.ts)
	base := filepath.Base(sourcePath)
	entrypoint := filepath.Join(sourcePath, fmt.Sprintf("%s.ts", base))
	if _, err := os.Stat(entrypoint); err == nil {
		return entrypoint, nil
	}

	// Failed after checking all patterns
	return "", fmt.Errorf("no entrypoint found in %s", sourcePath)
}
