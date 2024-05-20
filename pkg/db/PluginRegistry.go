package db

import (
	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/flowerLoader/tool/pkg/db/types"
)

type IPluginRegistry interface {
	CacheGet(id string) (*types.PluginCacheRecord, error)
	CachePut(record *types.PluginCacheRecord) error
}

// Ensure PluginRegistry implements IPluginRegistry
var _ IPluginRegistry = (*PluginRegistry)(nil)

type PluginRegistry struct {
	db  *DB
	log log.Logger
}

const SELECT_PLUGIN_CACHE = `SELECT * FROM plugin_cache WHERE id = ?`
const UPSERT_PLUGIN_CACHE = `INSERT OR REPLACE INTO plugin_cache (
	id, updated_at, name, version, author, license, bugs_url, homepage, api_version, tags, summary
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

func (r *PluginRegistry) CacheGet(id string) (*types.PluginCacheRecord, error) {
	r.log.Debug("searching for plugin", "id", id)

	stmt, err := r.db.conn.Prepare(SELECT_PLUGIN_CACHE)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	record := new(types.PluginCacheRecord)
	return record, stmt.QueryRow(id).Scan(&record.ID, &record.UpdatedAt, &record.Name, &record.Version, &record.Author, &record.License, &record.BugsURL, &record.Homepage, &record.APIVersion, &record.Tags, &record.Summary)
}

func (reg *PluginRegistry) CachePut(record *types.PluginCacheRecord) error {
	reg.log.Debug("upserting plugin", "id", record.ID)

	stmt, err := reg.db.conn.Prepare(UPSERT_PLUGIN_CACHE)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(record.ID, record.UpdatedAt, record.Name, record.Version, record.Author, record.License, record.BugsURL, record.Homepage, record.APIVersion, record.Tags, record.Summary)

	return err
}
