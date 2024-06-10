package ts

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	log "github.com/AlbinoGeek/logxi/v1"
)

var (
	// Regex to match a variable declaration line in TypeScript (const/var/let variable = value;)
	variableRegex = regexp.MustCompile(`^\s*(const|var|let)\s+(\w+)\s*=\s*[^;]*`)
)

func setVariable(filename, variable, value string) error {
	// Rename the original file to a backup file
	backupFilename := filename + ".bak"
	if err := os.Rename(filename, backupFilename); err != nil {
		return fmt.Errorf("failed to rename original file to backup: %w", err)
	}

	// Open the backup file for reading
	bakFile, err := os.Open(backupFilename)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer bakFile.Close()

	// Create a new file with the original filename for writing the modified content
	newFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create new file: %w", err)
	}
	defer newFile.Close()

	scanner := bufio.NewScanner(bakFile)
	writer := bufio.NewWriter(newFile)
	for scanner.Scan() {
		line := scanner.Text()
		if variableRegex.MatchString(line) {
			// Extract the parts of the matched line
			parts := variableRegex.FindStringSubmatch(line)
			if len(parts) > 2 && parts[2] == variable {
				// Replace only if the variable name matches
				line = fmt.Sprintf("%s %s = %s;", parts[1], variable, value)

				log.Debug("Updated variable",
					"filename", filename,
					"variable", variable,
					"value", value)
			}
		}
		// Write the (possibly modified) line to the new file
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write to new file: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading backup file: %w", err)
	}

	// Flush the writer to ensure all data is written to the new file
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	// Close both files
	newFile.Close()
	bakFile.Close()

	// Remove the backup file as everything went well
	if err := os.Remove(backupFilename); err != nil {
		return fmt.Errorf("failed to remove backup file: %w", err)
	}

	return nil
}
