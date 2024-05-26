package main

import (
	"os"
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
	rootCmd.AddCommand(listCmd)
}

func onListCommandRun(cmd *cobra.Command, args []string) {
	records, err := DB.Plugins.List()
	if err != nil {
		log.Error("Failed to list plugins", "error", err)
		return
	}

	if len(records) == 0 {
		log.Info("No plugins installed")
		return
	}
	log.Info("Installed Plugins", "count", len(records))

	data := make(map[string]interface{})
	for _, record := range records {
		cacheRecord, err := DB.Plugins.CacheGet(record.ID)
		if err != nil || cacheRecord == nil {
			log.Error("Failed to get plugin info from cache", "id", record.ID, "error", err)
			continue
		}

		// Add the cache record to the data map
		data[record.ID] = map[string]interface{}{
			"Name":        cacheRecord.Name,
			"Version":     cacheRecord.Version,
			"Author":      cacheRecord.Author,
			"InstalledAt": types.MustParseTime(record.InstalledAt).Local().Format(time.RFC822),
			"UpdatedAt":   types.MustParseTime(cacheRecord.UpdatedAt).Local().Format(time.RFC822),
		}
	}

	if err := listOutputTemplate.Execute(os.Stdout, data); err != nil {
		log.Error("Failed to execute template (can't print?)", "error", err)
	}
}
