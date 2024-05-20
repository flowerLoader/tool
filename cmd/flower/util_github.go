package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/google/go-github/v62/github"
)

var ghClient = github.NewClient(nil)

func githubSearch(ctx context.Context, query string) (*github.RepositoriesSearchResult, error) {
	log.Debug("Searching GitHub for repositories", "query", query)
	repos, _, err := ghClient.Search.Repositories(
		ctx, fmt.Sprintf("%s topic:flower-plugin", query), &github.SearchOptions{
			ListOptions: github.ListOptions{PerPage: 10},
			Order:       "desc",
			Sort:        "updated",
			TextMatch:   true,
		})

	if err != nil {
		return nil, err
	}

	return repos, nil
}

type PackageJSON struct {
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

func githubRepoAnalyze(ctx context.Context, repo string) (*PackageJSON, error) {
	owner := strings.Split(repo, "/")[0]
	repo = strings.Split(repo, "/")[1]

	// 1) Which branches are present in the repo?
	log.Debug("Analyzing Repository, Fetching Branches", "owner", owner, "repo", repo)
	branches, _, err := ghClient.Repositories.ListBranches(
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
		fileContent, _, _, err := ghClient.Repositories.GetContents(
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
	results := new(PackageJSON)
	if err := json.Unmarshal([]byte(packageJSON), &results); err != nil {
		return nil, err
	}

	return results, nil
}
