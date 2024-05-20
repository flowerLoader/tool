package db

import (
	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/embed"

	"github.com/flowerLoader/tool/pkg/db/types"
)

type DB struct {
	conn *sqlite3.Conn

	Plugins *PluginRegistry
}

func NewDB(filename string) (*DB, error) {
	conn, err := sqlite3.Open(filename)
	if err != nil {
		return nil, err
	}

	db := &DB{
		conn: conn,
	}

	db.Plugins = &PluginRegistry{
		db: db,
	}

	return db, nil
}

func (db *DB) Migrate() error {
	if err := db.conn.Exec(`CREATE TABLE IF NOT EXISTS plugin_cache (
		-- (owner/repo/tag#commit) is the primary key
		id TEXT PRIMARY KEY,

		-- META
		guid TEXT NOT NULL,
		version TEXT NOT NULL,
		name TEXT NOT NULL,
		author TEXT NOT NULL,

		-- GITHUB
		tags TEXT NOT NULL,
		summary TEXT NOT NULL
	)`); err != nil {
		return err
	}

	if err := db.conn.Exec(`CREATE TABLE IF NOT EXISTS plugin_install (
		-- (owner/repo/tag#commit) is the primary key
		-- refers to plugin_cache.id
		id TEXT PRIMARY KEY,

		-- Installation Status
		enabled TEXT NOT NULL,
		installed TEXT NOT NULL,
		updated TEXT NOT NULL
	)`); err != nil {
		return err
	}

	return nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

// ---

type PluginRegistry struct {
	db *DB
}

func (p *PluginRegistry) CacheGet(id string) (*types.PluginCacheRecord, error) {
	prepared := `SELECT
		id, guid, version, name, author, tags, summary
	FROM plugin_cache WHERE id = ?`

	stmt, _, err := p.db.conn.Prepare(prepared)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	stmt.BindText(1, id)

	if hasRow := stmt.Step(); !hasRow {
		return nil, err
	}

	return &types.PluginCacheRecord{
		ID:      stmt.ColumnText(0),
		GUID:    stmt.ColumnText(1),
		Version: stmt.ColumnText(2),
		Name:    stmt.ColumnText(3),
		Author:  stmt.ColumnText(4),
		Tags:    stmt.ColumnText(5),
		Summary: stmt.ColumnText(6),
	}, nil
}

func (p *PluginRegistry) CacheUpdate(plugin types.PluginCacheRecord) error {
	prepared := `INSERT INTO plugin_cache (
		id, guid, version, name, author, tags, summary
	) VALUES (
		?, ?, ?, ?, ?, ?, ?
	)`

	stmt, _, err := p.db.conn.Prepare(prepared)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if err := stmt.BindText(1, plugin.ID); err != nil {
		return err
	}
	if err := stmt.BindText(2, plugin.GUID); err != nil {
		return err
	}
	if err := stmt.BindText(3, plugin.Version); err != nil {
		return err
	}
	if err := stmt.BindText(4, plugin.Name); err != nil {
		return err
	}
	if err := stmt.BindText(5, plugin.Author); err != nil {
		return err
	}
	if err := stmt.BindText(6, plugin.Tags); err != nil {
		return err
	}
	if err := stmt.BindText(7, plugin.Summary); err != nil {
		return err
	}

	return stmt.Exec()
}
