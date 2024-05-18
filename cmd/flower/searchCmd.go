package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/v62/github"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"s", "find"},
	Short:   "Search for a plugin",
	Long:    "Search for a plugin by name, author, tags or summary",
	Example: `flower search LimitBreaker`,
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var sb strings.Builder

		// quote-parse using strconv
		for i, arg := range args {
			sb.WriteString(strconv.Quote(arg))
			if i < len(args)-1 {
				sb.WriteString(" ")
			}
		}

		return searchPlugin(cmd.Context(), sb.String())
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func searchPlugin(ctx context.Context, query string) error {
	client := github.NewClient(nil)

	queryOpts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 10},
		Order:       "desc",
		Sort:        "updated",
		TextMatch:   true,
	}

	repos, _, err := client.Search.Repositories(
		ctx, fmt.Sprintf("%s topic:flower-plugin", query), queryOpts,
	)

	if err != nil {
		return err
	}

	if repos.GetTotal() == 0 {
		fmt.Printf("flower > No results found for %s\n", query)
		return nil
	}

	fmt.Print("| Name                                     | Stars  | Forks  | Last Update                |\n|------------------------------------------|--------|--------|----------------------------|\n")

	for _, repo := range repos.Repositories {
		name := *repo.FullName
		if len(name) > 40 {
			name = name[:37] + "..."
		}

		fmt.Printf("| %-40s | %-6d | %-6d | %-26s |\n",
			name,
			*repo.StargazersCount,
			*repo.ForksCount,
			repo.UpdatedAt.Local().Format("January 2, 2006 3:04 PM"))
	}

	fmt.Printf("\n%d/%d results found for %s",
		len(repos.Repositories),
		repos.GetTotal(),
		query)

	return nil
}
