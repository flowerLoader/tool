package main

import (
	"context"
	"errors"
	"strconv"
	"strings"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l", "installed"},
	Short:   "List installed plugins",
	Long:    "List installed plugins by name, author, tags or summary",
	Example: `flower list`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var sb strings.Builder

		// quote-parse using strconv
		for i, arg := range args {
			sb.WriteString(strconv.Quote(arg))
			if i < len(args)-1 {
				sb.WriteString(" ")
			}
		}

		if err := listPlugins(cmd.Context(), sb.String()); err != nil {
			log.Fatal("failed to list plugins", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listPlugins(ctx context.Context, query string) error {
	return errors.New("not implemented")
}
