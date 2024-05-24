package main

import (
	"fmt"
	"time"

	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/spf13/cobra"

	"github.com/flowerLoader/tool/pkg/db/types"
)

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

	fmt.Print("| Name                 | Version | Author               | Installed At             | Last Updated             |\n|----------------------|---------|----------------------|--------------------------|--------------------------|\n")
	for _, installRecord := range records {
		cacheRecord, err := DB.Plugins.CacheGet(installRecord.ID)
		if err != nil || cacheRecord == nil {
			log.Error("Failed to get plugin info from cache", "id", installRecord.ID, "error", err)
			continue
		}

		name := cacheRecord.Name
		if len(name) > 20 {
			name = name[:17] + "..."
		}

		author := cacheRecord.Author
		if len(author) > 20 {
			author = author[:17] + "..."
		}

		fmt.Printf("| %-20s | %-7s | %-20s | %-24s | %-24s |\n",
			name,
			cacheRecord.Version,
			author,
			types.MustParseTime(installRecord.InstalledAt).Local().Format(time.RFC822),
			types.MustParseTime(cacheRecord.UpdatedAt).Local().Format(time.RFC822))
	}
}
