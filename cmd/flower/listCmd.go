package main

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l", "installed"},
	Short:   "List installed plugins",
	Long:    "List installed plugins by name, author, tags or summary",
	Example: `flower list`,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var sb strings.Builder

		// quote-parse using strconv
		for i, arg := range args {
			sb.WriteString(strconv.Quote(arg))
			if i < len(args)-1 {
				sb.WriteString(" ")
			}
		}

		return listPlugins(cmd.Context(), sb.String())
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listPlugins(ctx context.Context, query string) error {
	return errors.New("not implemented")
}
