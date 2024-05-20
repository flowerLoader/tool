package db

import "github.com/flowerLoader/tool/pkg/db/types"

type PluginRegistry struct {
	db *DB
}

func (r *PluginRegistry) Get(id string) (*types.PluginCacheRecord, error) {
	stmt, _, err := r.db.conn.Prepare("SELECT id, name, version, author, license, bugs_url, homepage, api_version, tags, summary FROM plugin_cache WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	stmt.BindText(1, id)
	plugin := types.PluginCacheRecord{}

	if stmt.Step() {
		plugin.ID = stmt.ColumnText(0)
		plugin.Name = stmt.ColumnText(1)
		plugin.Version = stmt.ColumnText(2)
		plugin.Author = stmt.ColumnText(3)
		plugin.License = stmt.ColumnText(4)
		plugin.BugsURL = stmt.ColumnText(5)
		plugin.Homepage = stmt.ColumnText(6)
		plugin.APIVersion = stmt.ColumnText(7)
		plugin.Tags = stmt.ColumnText(8)
		plugin.Summary = stmt.ColumnText(9)

		stmt.Reset()
		return &plugin, nil
	}

	stmt.Reset()
	return nil, nil
}

func (reg *PluginRegistry) Upsert(record *types.PluginCacheRecord) error {
	stmt, _, err := reg.db.conn.Prepare("INSERT OR REPLACE INTO plugin_cache (id, name, version, author, license, bugs_url, homepage, api_version, tags, summary) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.BindText(1, record.ID)
	stmt.BindText(2, record.Name)
	stmt.BindText(3, record.Version)
	stmt.BindText(4, record.Author)
	stmt.BindText(5, record.License)
	stmt.BindText(6, record.BugsURL)
	stmt.BindText(7, record.Homepage)
	stmt.BindText(8, record.APIVersion)
	stmt.BindText(9, record.Tags)
	stmt.BindText(10, record.Summary)
	stmt.Step()
	stmt.Reset()

	return nil
}
