package main

import (
	"errors"
	"os"

	log "github.com/AlbinoGeek/logxi/v1"
)

var (
	ErrAlreadyInstalled = errors.New("plugin is already installed, use `flower update` instead")
	ErrNameTaken        = errors.New("plugin name is already taken, try a different name")
	ErrNotInstalled     = errors.New("plugin is not installed, use `flower add` instead")
	ErrQueryDB          = errors.New("failed to query plugin database, check logs for details")
)

func exit(err error, args ...interface{}) {
	log.Error("Exiting",
		"error", err,
		"args", args)

	os.Exit(1)
}
