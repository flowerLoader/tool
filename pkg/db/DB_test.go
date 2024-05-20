package db

import (
	"testing"

	"github.com/flowerLoader/tool/pkg/db/types"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func TestDB(t *testing.T) {
	t.Run("NewDB", func(t *testing.T) {
		t.Run("should return a new DB instance", func(t *testing.T) {
			db, err := NewDB(":memory:")
			if err != nil {
				t.Errorf("NewDB() error = %v", err)
				return
			}
			if db == nil {
				t.Error("NewDB() db is nil")
				return
			}
		})
	})

	t.Run("Migrate", func(t *testing.T) {
		t.Run("should create the plugin_cache and plugin_install tables", func(t *testing.T) {
			db, err := NewDB(":memory:")
			if err != nil {
				t.Errorf("NewDB() error = %v", err)
				return
			}

			if err := db.Migrate(); err != nil {
				t.Errorf("Migrate() error = %v", err)
				return
			}

			stmt, _, err := db.conn.Prepare("SELECT name FROM sqlite_master WHERE type='table' AND name IN ('plugin_cache', 'plugin_install')")
			if err != nil {
				t.Errorf("Prepare() error = %v", err)
				return
			}

			if haveData := stmt.Step(); !haveData {
				t.Errorf("Step() error = %v", err)
				return
			}

			if stmt.ColumnText(0) != "plugin_cache" {
				t.Errorf("Step() got = %v, want plugin_cache", stmt.ColumnText(0))
				return
			}

			if haveData := stmt.Step(); !haveData {
				t.Errorf("Step() error = %v", err)
				return
			}

			if stmt.ColumnText(0) != "plugin_install" {
				t.Errorf("Step() got = %v, want plugin_install", stmt.ColumnText(0))
				return
			}

			if haveData := stmt.Step(); haveData {
				t.Error("Step() expected error, got nil")
				return
			}

			if err := stmt.Close(); err != nil {
				t.Errorf("Finalize() error = %v", err)
				return
			}

			if err := db.Close(); err != nil {
				t.Errorf("Close() error = %v", err)
				return
			}
		})
	})

	t.Run("InsertPluginCache", func(t *testing.T) {
		t.Run("should insert a new record into the plugin_cache table", func(t *testing.T) {
			db, err := NewDB(":memory:")
			if err != nil {
				t.Errorf("NewDB() error = %v", err)
				return
			}

			if err := db.Migrate(); err != nil {
				t.Errorf("Migrate() error = %v", err)
				return
			}

			record := types.PluginCacheRecord{
				ID:      "id",
				GUID:    "guid",
				Version: "version",
				Name:    "name",
				Author:  "author",
				Tags:    "tags,tag2,tag3",
				Summary: "summary",
			}

			if err := db.Plugins.CacheUpdate(record); err != nil {
				t.Errorf("InsertPluginCache() error = %v", err)
				return
			}

			res, err := db.Plugins.CacheGet("id")
			if err != nil {
				t.Errorf("CacheGet() error = %v", err)
				return
			}

			if res.ID != record.ID {
				t.Errorf("CacheGet() got = %v, want %v", res.ID, record.ID)
				return
			}

			if res.GUID != record.GUID {
				t.Errorf("CacheGet() got = %v, want %v", res.GUID, record.GUID)
				return
			}

			if res.Version != record.Version {
				t.Errorf("CacheGet() got = %v, want %v", res.Version, record.Version)
				return
			}

			if res.Name != record.Name {
				t.Errorf("CacheGet() got = %v, want %v", res.Name, record.Name)
				return
			}

			if res.Author != record.Author {
				t.Errorf("CacheGet() got = %v, want %v", res.Author, record.Author)
				return
			}

			if res.Tags != record.Tags {
				t.Errorf("CacheGet() got = %v, want %v", res.Tags, record.Tags)
				return
			}

			if res.Summary != record.Summary {
				t.Errorf("CacheGet() got = %v, want %v", res.Summary, record.Summary)
				return
			}

			if err := db.Close(); err != nil {
				t.Errorf("Close() error = %v", err)
				return
			}
		})
	})
}
