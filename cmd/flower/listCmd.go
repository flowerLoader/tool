package main

import (
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/cobra"

	"github.com/flowerLoader/tool/pkg/db/types"
)

var listOutputTemplate = template.Must(template.New("pluginList").Parse(strings.TrimSpace(`
| Name                 | Version | Author               | Installed At             | Last Updated             |
|----------------------|---------|----------------------|--------------------------|--------------------------|{{range .}}
| {{.Name | printf "%-20s"}} | {{.Version | printf "%-7s"}} | {{.Author | printf "%-20s"}} | {{.InstalledAt | printf "%-24s"}} | {{.UpdatedAt | printf "%-24s"}} |{{end}}

Total Installed: {{len .}}`)))

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l", "installed"},
	Short:   "List installed plugins",
	Long:    "List installed plugins by name, author, tags or summary",
	Example: `flower list`,
	Args:    cobra.NoArgs,
	Run:     onListCommandRun,
}

func init() {
	listCmd.PersistentFlags().String("author", "", "Only show plugins by this author (Regular Expression Match, such as 'flower.*')")
	rootCmd.AddCommand(listCmd)
}

func onListCommandRun(cmd *cobra.Command, args []string) {
	authorFilter, err := cmd.Flags().GetString("author")
	if err != nil {
		log.Error("Failed to get author flag", "error", err)
		return
	}

	records, err := App.DB.Plugins.List()
	if err != nil {
		log.Error("Failed to list plugins", "error", err)
		return
	}

	if len(records) == 0 {
		log.Info("No plugins installed")
		return
	}
	log.Info("Installed Plugins", "count", len(records))
	toRender := filterPlugins(records, authorFilter)

	if len(toRender) == 0 {
		log.Warn("No plugins matched the filter",
			"author", authorFilter,
			"ignored", len(records))
		return
	}

	if err := listOutputTemplate.Execute(os.Stdout, toRender); err != nil {
		log.Error("Failed to execute template (can't print?)", "error", err)
	}
}

func filterPlugins(
	records []*types.PluginInstallRecord,
	authorFilter string,
) map[string]interface{} {
	toRender := make(map[string]interface{})
	filters := compileFilters(authorFilter)

	for _, record := range records {
		cacheRecord, err := App.DB.Plugins.CacheGet(record.ID)
		if err != nil || cacheRecord == nil {
			exit(ErrQueryDB, err)
		}

		if applyFilters(cacheRecord, filters) {
			toRender[record.ID] = map[string]interface{}{
				"Name":        cacheRecord.Name,
				"Version":     cacheRecord.Version,
				"Author":      cacheRecord.Author,
				"InstalledAt": types.MustParseTime(record.InstalledAt).Local().Format(time.RFC822),
				"UpdatedAt":   types.MustParseTime(cacheRecord.UpdatedAt).Local().Format(time.RFC822),
			}
		}
	}

	return toRender
}

func compileFilters(authorFilter string) []func(*types.PluginCacheRecord) bool {
	filters := make([]func(*types.PluginCacheRecord) bool, 0)
	if authorFilter != "" {
		pattern, err := regexp.Compile(authorFilter)
		if err != nil {
			log.Error("Failed to compile author filter", "error", err)
			return nil
		}

		filters = append(filters, func(record *types.PluginCacheRecord) bool {
			return pattern.MatchString(record.Author)
		})
	}

	return filters
}

func applyFilters(
	record *types.PluginCacheRecord,
	filters []func(*types.PluginCacheRecord) bool,
) bool {
	for _, filter := range filters {
		if !filter(record) {
			return false
		}
	}

	return true
}
