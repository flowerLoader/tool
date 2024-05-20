package db

import (
	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/embed"
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
