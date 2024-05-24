package main

import (
	"fmt"
	"strings"

	log "github.com/AlbinoGeek/logxi/v1"
)

func promptConfirm(prompt string) bool {
	log.Warn(prompt, " [y/N]")

	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		log.Error("Failed to read response", "error", err)
		return false
	}

	return strings.ToLower(response) == "y"
}
