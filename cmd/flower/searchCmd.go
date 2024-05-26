package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	setupProgress()

	_, done := newTracker("Retrieving search results from GitHub")
	repos, err := githubSearch(ctx, query)
	done()
	if err != nil {
		return err
	}

	if repos.GetTotal() == 0 {
		fmt.Printf("flower > No results found for %s\n", query)
		return nil
	}
	log.Debug("search results", "total", repos.GetTotal())

	// Update our cache
	records := make([]*types.PluginCacheRecord, 0)
	for _, repo := range repos.Repositories {
		cacheRecord, err := DB.Plugins.CacheGet(
			fmt.Sprintf("%s/%s", GITHUB_PKG, *repo.FullName))
		if err != nil || cacheRecord == nil {
			_, done = newTracker(fmt.Sprintf("Analyzing %s", *repo.FullName))

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
				ID:         fmt.Sprintf("%s/%s", GITHUB_PKG, *repo.FullName),
				Version:    analysis.Version,
				Name:       analysis.Name,
				Author:     analysis.Author,
				License:    analysis.License,
				BugsURL:    analysis.Bugs.URL,
				Homepage:   analysis.Homepage,
				APIVersion: apiVersion,
			}

			done()
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
	pw.Stop()
	time.Sleep(time.Millisecond * 10)

	// Print the results
	fmt.Print("\n| Name                 | Version | Author               | License      | Last Updated             |\n|----------------------|---------|----------------------|--------------|--------------------------|\n")
	for _, cacheRecord := range records {
		name := cacheRecord.Name
		if len(name) > 20 {
			name = name[:17] + "..."
		}

		author := cacheRecord.Author
		if len(author) > 20 {
			author = author[:17] + "..."
		}

		fmt.Printf("| %-20s | %-7s | %-20s | %-12s | %-24s |\n",
			name,
			cacheRecord.Version,
			author,
			cacheRecord.License,
			types.MustParseTime(cacheRecord.UpdatedAt).Local().Format(time.RFC822))
	}

	fmt.Printf("\n%d/%d results found for %s",
		len(repos.Repositories),
		repos.GetTotal(),
		query)

	return nil
}
