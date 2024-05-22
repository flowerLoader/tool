package db

import (
	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sqlx.DB

	Plugins IPluginRegistry
}

type Stats struct {
	Counts struct {
		Cached  int
		Cloned  int
		Enabled int
	}
}

func NewDB(filename string) (db *DB, err error) {
	db = new(DB)

	db.Plugins = &PluginRegistry{
		db:  db,
		log: log.New("db.plugins"),
	}

	if db.conn, err = sqlx.Connect("sqlite3", filename); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Migrate() error {
	if _, err := db.conn.Exec(`CREATE TABLE IF NOT EXISTS plugin_cache (
		-- (owner/repo/tag#commit) is the primary key
		id TEXT PRIMARY KEY,
		updated_at TEXT NOT NULL,

		-- META
		name TEXT NOT NULL,
		version TEXT NOT NULL,
		author TEXT NOT NULL,
		license TEXT NOT NULL,
		bugs_url TEXT NOT NULL,
		homepage TEXT NOT NULL,
		api_version TEXT NOT NULL,

		-- GITHUB
		tags TEXT NOT NULL,
		summary TEXT NOT NULL
	)`); err != nil {
		return err
	}

	if _, err := db.conn.Exec(`CREATE TABLE IF NOT EXISTS plugin_install (
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

func (db *DB) Stat() (stat Stats, err error) {
	return stat, db.conn.Get(&stat, `SELECT
		(SELECT COUNT(*) FROM plugin_cache) AS cached,
		(SELECT COUNT(*) FROM plugin_install) AS cloned,
		(SELECT COUNT(*) FROM plugin_install WHERE enabled = 'true') AS enabled`)
}
