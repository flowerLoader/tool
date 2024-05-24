package db

import (
	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/flowerLoader/tool/pkg/db/types"
)

type IPluginRegistry interface {
	CacheGet(id string) (*types.PluginCacheRecord, error)
	CachePut(record *types.PluginCacheRecord) error

	Add(record *types.PluginInstallRecord) error
	Get(id string) (*types.PluginInstallRecord, error)
	List() ([]*types.PluginInstallRecord, error)
}

// Ensure PluginRegistry implements IPluginRegistry
var _ IPluginRegistry = (*PluginRegistry)(nil)

type PluginRegistry struct {
	db  *DB
	log log.Logger
}

const SELECT_ALL_PLUGIN = `SELECT * FROM plugin_install`
const SELECT_PLUGIN = `SELECT * FROM plugin_install WHERE id = ?`
const SELECT_PLUGIN_CACHE = `SELECT * FROM plugin_cache WHERE id = ?`
const UPSERT_PLUGIN_CACHE = `INSERT OR REPLACE INTO plugin_cache (
	id, updated_at, name, version, author, license, bugs_url, homepage, api_version, tags, summary
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
const INSERT_PLUGIN = `INSERT INTO plugin_install (
	id, enabled, installed, updated
) VALUES (?, ?, ?, ?)`

func (r *PluginRegistry) CacheGet(id string) (*types.PluginCacheRecord, error) {
	r.log.Debug("searching for plugin", "id", id)

	stmt, err := r.db.conn.Prepare(SELECT_PLUGIN_CACHE)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	record := new(types.PluginCacheRecord)
	err = stmt.QueryRow(id).Scan(&record.ID, &record.UpdatedAt, &record.Name, &record.Version, &record.Author, &record.License, &record.BugsURL, &record.Homepage, &record.APIVersion, &record.Tags, &record.Summary)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil, nil // Not found
	}

	return record, err
}

func (r *PluginRegistry) CachePut(record *types.PluginCacheRecord) error {
	r.log.Debug("upserting plugin", "id", record.ID)

	stmt, err := r.db.conn.Prepare(UPSERT_PLUGIN_CACHE)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(record.ID, record.UpdatedAt, record.Name, record.Version, record.Author, record.License, record.BugsURL, record.Homepage, record.APIVersion, record.Tags, record.Summary)

	return err
}

func (r *PluginRegistry) Add(record *types.PluginInstallRecord) error {
	r.log.Debug("adding plugin", "id", record.ID)

	stmt, err := r.db.conn.Prepare(INSERT_PLUGIN)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(record.ID, record.Enabled, record.InstalledAt, record.UpdatedAt)
	return err
}

func (r *PluginRegistry) Get(id string) (*types.PluginInstallRecord, error) {
	r.log.Debug("searching for plugin", "id", id)

	stmt, err := r.db.conn.Prepare(SELECT_PLUGIN)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	record := new(types.PluginInstallRecord)
	err = stmt.QueryRow(id).Scan(&record.ID, &record.Enabled, &record.InstalledAt, &record.UpdatedAt)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil, nil // Not found
	}

	return record, err
}

func (r *PluginRegistry) List() ([]*types.PluginInstallRecord, error) {
	r.log.Debug("listing all plugins")

	rows, err := r.db.conn.Query(SELECT_ALL_PLUGIN)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*types.PluginInstallRecord
	for rows.Next() {
		record := new(types.PluginInstallRecord)
		err = rows.Scan(&record.ID, &record.Enabled, &record.InstalledAt, &record.UpdatedAt)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}
