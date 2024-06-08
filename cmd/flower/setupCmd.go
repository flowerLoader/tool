package main

import (
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:     "setup",
	Aliases: []string{"init", "configure"},
	Short:   "Setup the flowerLoader environment",
	Long:    "Setup the flowerLoader environment for the first time (or after a reinstall)",
	Example: `flower setup --game-path /path/to/game`,
	Args:    cobra.NoArgs,
	Run:     onSetupCommandRun,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func onSetupCommandRun(cmd *cobra.Command, args []string) {
	// Get the resolved game-path (courtesy of rootCmd PersistentPreRunE)
	gameInstallPath, err := cmd.Flags().GetString("game-path")
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	if err := initFlowerLoader(App.Config, gameInstallPath); err != nil {
		cmd.PrintErrln(err)
		return
	}

	cmd.Println("flowerLoader environment setup complete!")
}
