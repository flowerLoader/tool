package main

import (
	"github.com/spf13/cobra"
)

const (
	APPNAME    = "flower"
	APPVERSION = "0.1.0"
)

var (
	rootCmd = &cobra.Command{
		Use:     APPNAME,
		Version: APPVERSION,
		Short:   "",
		Long:    "",
		Run:     rootCmdRun,
	}
)

func init() {

}

func main() {
	rootCmd.Execute()
}

func rootCmdRun(cmd *cobra.Command, args []string) {
	cmd.Help()
}
