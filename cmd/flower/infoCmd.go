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
	query := strings.Join(args, " ")
	cacheRecord, err := DB.Plugins.CacheGet(query)
	if err != nil {
		log.Error("Failed to get plugin info from cache", "error", err)
		return
	}

	if cacheRecord == nil {
		fmt.Printf("flower > Try running 'flower search %s' to find it\n", query)
		return
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
