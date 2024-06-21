package main

import (
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"start", "launch"},
	Short:   "Runs the game with your enabled mods",
	Long:    "Run the current game with your configured mods and settings",
	Example: `flower run`,
	Args:    cobra.NoArgs,
	Run:     onRunCommandRun,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func onRunCommandRun(cmd *cobra.Command, args []string) {
	// Get the resolved game-path (courtesy of rootCmd PersistentPreRunE)
	gameInstallPath, err := cmd.Flags().GetString("game-path")
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	// Ensure flower is installed
	if !isFlowerInstalled(gameInstallPath) {
		cmd.PrintErrln("Flower is not installed in your game directory")
		return
	}

	// Run the game with the correct environment arguments
	runGame(gameInstallPath)
}

func isFlowerInstalled(gameInstallPath string) bool {
	_, err := os.Stat(filepath.Join(gameInstallPath, "flower"))
	return err == nil
}

func runGame(gameInstallPath string) {
	gameInstallPath, err := filepath.Abs(filepath.Join(gameInstallPath, ".."))
	if err != nil {
		panic(err)
	}

	if filepath.Base(gameInstallPath) == "gamedata" {
		gameInstallPath, err = filepath.Abs(filepath.Join(gameInstallPath, ".."))
		if err != nil {
			panic(err)
		}
	}

	cmd := exec.Command(filepath.Join(gameInstallPath, "Game.exe"))
	cmd.Dir = gameInstallPath
	cmd.Env = append(os.Environ(), "WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS=--allow-file-access-from-files")

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	log.Warn("Game started", "pid", cmd.Process.Pid)
	os.Exit(0)
}
