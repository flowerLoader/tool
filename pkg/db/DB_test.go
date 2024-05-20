package db

import (
	"testing"

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
}
