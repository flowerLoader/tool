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
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("executing info", "query", args[0])
		if err := infoPlugin(args[0]); err != nil {
			log.Fatal("failed to get plugin info", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func infoPlugin(query string) error {
	cacheRecord, err := DB.Plugins.CacheGet(query)
	if err != nil {
		return err
	}

	if cacheRecord == nil {
		fmt.Printf("flower > Try running 'flower search %s' to find it\n", query)
		return nil
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

	return nil
}
