package main

import (
	"fmt"
	"strings"

	log "github.com/AlbinoGeek/logxi/v1"
)

func promptConfirm(prompt string) bool {
	log.Info("Prompting for confirmation", "prompt", prompt)
	print(prompt, "[y/N]: ")

	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		if err.Error() != "unexpected newline" {
			log.Error("Unexpected error reading confirmation response", "error", err)
		}

		return false
	}

	return strings.ToLower(response) == "y"
}
