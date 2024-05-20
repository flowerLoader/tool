package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v62/github"
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

		return searchPlugin(cmd.Context(), sb.String())
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func searchPlugin(ctx context.Context, query string) error {
	repos, err := searchGithub(ctx, query)
	if err != nil {
		return err
	}

	if repos.GetTotal() == 0 {
		fmt.Printf("flower > No results found for %s\n", query)
		return nil
	}

	// Update our cache
	for _, repo := range repos.Repositories {
		cacheRecord, err := DB.Plugins.Get(*repo.FullName)
		if err != nil || cacheRecord == nil {
			// We have never encountered this repo before, we must interrogate it
			analysis, err := analyzeRepo(ctx, *repo.FullName)
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
		cacheRecord.Author = *repo.Owner.Login
		cacheRecord.Summary = *repo.Description
		cacheRecord.Tags = strings.Join(repo.Topics, ",")

		if err := DB.Plugins.Upsert(cacheRecord); err != nil {
			return err
		}
	}

	// Print the results
	fmt.Print("| Name                 | Version | Author               | License   | Last Updated               |\n|----------------------|---------|----------------------|-----------|----------------------------|\n")

	for _, repo := range repos.Repositories {
		cacheRecord, err := DB.Plugins.Get(*repo.FullName)
		if err != nil {
			return err
		}

		name := cacheRecord.Name
		if len(name) > 20 {
			name = name[:17] + "..."
		}

		author := cacheRecord.Author
		if len(author) > 20 {
			author = author[:17] + "..."
		}

		fmt.Printf("| %-20s | %-7s | %-20s | %-9s | %-26s |\n",
			name,
			cacheRecord.Version,
			author,
			cacheRecord.License,
			repo.UpdatedAt.Local().Format("January 2, 2006 3:04 PM"))
	}

	fmt.Printf("\n%d/%d results found for %s",
		len(repos.Repositories),
		repos.GetTotal(),
		query)

	return nil
}

func searchGithub(ctx context.Context, query string) (*github.RepositoriesSearchResult, error) {
	client := github.NewClient(nil)

	queryOpts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 10},
		Order:       "desc",
		Sort:        "updated",
		TextMatch:   true,
	}

	log.Debug("Searching GitHub", "query", query)
	repos, _, err := client.Search.Repositories(
		ctx, fmt.Sprintf("%s topic:flower-plugin", query), queryOpts,
	)

	if err != nil {
		return nil, err
	}

	return repos, nil
}

type RepoAnalysis struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Author  string `json:"author,omitempty"`
	License string `json:"license,omitempty"`
	Bugs    struct {
		URL string `json:"url,omitempty"`
	} `json:"bugs,omitempty"`
	Homepage     string            `json:"homepage,omitempty"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
}

func analyzeRepo(ctx context.Context, repo string) (*RepoAnalysis, error) {
	client := github.NewClient(nil)
	owner := strings.Split(repo, "/")[0]
	repo = strings.Split(repo, "/")[1]

	// 1) Which branches are present in the repo?
	log.Debug("Analyzing Repository, Fetching Branches", "owner", owner, "repo", repo)
	branches, _, err := client.Repositories.ListBranches(
		ctx, owner, repo, &github.BranchListOptions{
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		})

	if err != nil {
		return nil, err
	}

	// 2) Which branch is the most recent and contains package.json?
	var packageJSON string
	for _, branch := range branches {
		log.Debug("Analyzing Repository, Fetching package.json ...", "branch", *branch.Name)
		fileContent, _, _, err := client.Repositories.GetContents(
			ctx, owner, repo, "package.json", &github.RepositoryContentGetOptions{
				Ref: *branch.Name,
			})

		if err != nil {
			continue
		}

		data, err := fileContent.GetContent()
		if err != nil {
			return nil, err
		}

		packageJSON = data
		break
	}

	if packageJSON == "" {
		return nil, fmt.Errorf("package.json not found in %s", repo)
	}

	// 3) Parse package.json
	results := new(RepoAnalysis)
	if err := json.Unmarshal([]byte(packageJSON), &results); err != nil {
		return nil, err
	}

	return results, nil
}
