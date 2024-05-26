package db

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flowerLoader/tool/pkg/db/types"
)

func now() string {
	return "2021-07-18 15:04:05+00:00" // rfc3339
}

func TestPluginRegistry(t *testing.T) {
	Convey("Given a new database", t, func() {
		db, err := NewDB(":memory:")
		So(err, ShouldBeNil)
		defer db.Close()

		err = db.Migrate()
		So(err, ShouldBeNil)

		Convey("When a plugin is upserted and retrieved", func() {
			plugin := &types.PluginCacheRecord{
				ID:         "test/repo/tag#commit",
				UpdatedAt:  now(),
				Name:       "Test Plugin",
				Version:    "1.0.0",
				Author:     "Author Name",
				License:    "MIT",
				BugsURL:    "http://example.com/bugs",
				Homepage:   "http://example.com",
				APIVersion: "v1",
				Tags:       "test, plugin",
				Summary:    "A test plugin",
			}

			err = db.Plugins.CachePut(plugin)
			So(err, ShouldBeNil)

			actualPlugin, err := db.Plugins.CacheGet(plugin.ID)
			So(err, ShouldBeNil)
			So(actualPlugin, ShouldResemble, plugin)

			Convey("And when the plugin is updated", func() {
				plugin.Name = "Updated Plugin"
				err = db.Plugins.CachePut(plugin)
				So(err, ShouldBeNil)

				actualPlugin, err := db.Plugins.CacheGet(plugin.ID)
				So(err, ShouldBeNil)
				So(actualPlugin.Name, ShouldEqual, "Updated Plugin")
			})

			Convey("And when a plugin is added", func() {
				installRecord := &types.PluginInstallRecord{
					ID:          plugin.ID,
					Enabled:     true,
					InstalledAt: now(),
					UpdatedAt:   now(),
				}

				err = db.Plugins.Add(installRecord)
				So(err, ShouldBeNil)

				actualRecord, err := db.Plugins.Get(installRecord.ID)
				So(err, ShouldBeNil)
				So(actualRecord, ShouldResemble, installRecord)

				Convey("And when plugins are listed", func() {
					records, err := db.Plugins.List()
					So(err, ShouldBeNil)
					So(records, ShouldHaveLength, 1)
					So(records[0], ShouldResemble, installRecord)
				})

				Convey("And when the plugin is removed", func() {
					err = db.Plugins.Remove(installRecord.ID)
					So(err, ShouldBeNil)

					actualRecord, err := db.Plugins.Get(installRecord.ID)
					So(err, ShouldBeNil)
					So(actualRecord, ShouldBeNil)

					Convey("And when plugins are listed", func() {
						records, err := db.Plugins.List()
						So(err, ShouldBeNil)
						So(records, ShouldBeEmpty)
					})
				})
			})
		})
	})
}
