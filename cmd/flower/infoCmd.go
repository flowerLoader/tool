package main

import (
	"fmt"
	"strings"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"i", "details", "lookup"},
	Short:   "Get information about a plugin",
	Long:    "Get detailed information about a plugin by name",
	Example: `flower info LimitBreaker`,
	Args:    cobra.MinimumNArgs(1),
	Run:     onInfoCommandRun,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func onInfoCommandRun(cmd *cobra.Command, args []string) {
	queryArgs := make([]string, len(args))
	for i, arg := range args {
		if parsed := parsePluginName(arg); parsed != "" {
			queryArgs[i] = parsed
		} else {
			queryArgs[i] = arg
		}
	}

	query := strings.Join(queryArgs, " ")
	cacheRecord, err := DB.Plugins.CacheGet(query)
	if err != nil {
		log.Error("Failed to get plugin info from cache", "error", err)
		return
	}

	if cacheRecord == nil {
		query = args[0]

		log.Warn("Plugin not found in cache. Attempting to search...", "query", query)
		if err := searchPlugin(cmd.Context(), query); err != nil {
			log.Error("Failed to search for plugin", "error", err)
			return
		}

		cacheRecord, err = DB.Plugins.CacheGet(parsePluginName(query))
		if err != nil {
			log.Error("Failed to get plugin info from cache", "error", err)
			return
		}

		if cacheRecord == nil {
			log.Error("Plugin not found in cache or search results")
			return
		}

		fmt.Printf("\n\n") // Add some space between search results and plugin info
	}

	fmt.Printf(strings.TrimSpace(`
ID             : %s
Name           : %s
Author         : %s
Version        : %s
API Version    : %s
License        : %s
Tags           : %s
Last Updated   : %s
Summary        : %s
Homepage       : %s
Report Bugs At : %s
`), cacheRecord.ID, cacheRecord.Name, cacheRecord.Author, cacheRecord.Version, cacheRecord.APIVersion, cacheRecord.License, cacheRecord.Tags, cacheRecord.UpdatedAt, cacheRecord.Summary, cacheRecord.Homepage, cacheRecord.BugsURL)
}
