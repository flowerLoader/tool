package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"

	"github.com/flowerLoader/tool/pkg/db/types"
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

		log.Debug("executing search", "query", sb.String())
		return searchPlugin(cmd.Context(), sb.String())
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func searchPlugin(ctx context.Context, query string) error {
	repos, err := githubSearch(ctx, query)
	if err != nil {
		return err
	}

	if repos.GetTotal() == 0 {
		fmt.Printf("flower > No results found for %s\n", query)
		return nil
	}

	// Update our cache
	records := make([]*types.PluginCacheRecord, 0)
	for _, repo := range repos.Repositories {
		cacheRecord, err := DB.Plugins.CacheGet(*repo.FullName)
		if err != nil || cacheRecord == nil {
			// We have never encountered this repo before, we must interrogate it
			analysis, err := githubRepoAnalyze(ctx, *repo.FullName)
			if err != nil {
				return err
			}

			apiVersion := "n/a"
			if analysis.Dependencies != nil {
				// parse semver from dependencies["@flowerloader/api"] if exists
				ver, ok := analysis.Dependencies["@flowerloader/api"]
				if ok {
					if parsed, err := semver.NewVersion(ver); err == nil {
						apiVersion = parsed.String()
					}
				}
			}

			cacheRecord = &types.PluginCacheRecord{
				ID:         *repo.FullName,
				Version:    analysis.Version,
				Name:       analysis.Name,
				Author:     analysis.Author,
				License:    analysis.License,
				BugsURL:    analysis.Bugs.URL,
				Homepage:   analysis.Homepage,
				APIVersion: apiVersion,
			}
		}

		// Update the cache
		cacheRecord.UpdatedAt = types.FormatTime(repo.UpdatedAt.Time)
		cacheRecord.Author = *repo.Owner.Login
		cacheRecord.Summary = *repo.Description
		cacheRecord.Tags = strings.Join(repo.Topics, ",")

		if err := DB.Plugins.CachePut(cacheRecord); err != nil {
			return err
		}

		records = append(records, cacheRecord)
	}

	// Print the results
	printPluginTable(records)

	fmt.Printf("\n%d/%d results found for %s",
		len(repos.Repositories),
		repos.GetTotal(),
		query)

	return nil
}
