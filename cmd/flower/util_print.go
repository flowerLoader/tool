package main

import (
	"fmt"

	"github.com/flowerLoader/tool/pkg/db/types"
)

func printPluginTable(records []*types.PluginCacheRecord) {
	fmt.Print("| Name                 | Version | Author               | License   | Last Updated               |\n|----------------------|---------|----------------------|-----------|----------------------------|\n")
	for _, cacheRecord := range records {
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
			cacheRecord.UpdatedAt)
	}
}
