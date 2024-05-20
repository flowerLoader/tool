package db

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"github.com/flowerLoader/tool/pkg/db/types"
)

func now() string {
	return "2021-07-18 15:04:05"
}

func TestPluginRegistry(t *testing.T) {
	convey.Convey("Given a new database", t, func() {
		db, err := NewDB(":memory:")
		convey.So(err, convey.ShouldBeNil)
		defer db.Close()

		err = db.Migrate()
		convey.So(err, convey.ShouldBeNil)

		convey.Convey("When a plugin is upserted and retrieved", func() {
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
			convey.So(err, convey.ShouldBeNil)

			actualPlugin, err := db.Plugins.CacheGet(plugin.ID)
			convey.So(err, convey.ShouldBeNil)
			convey.So(actualPlugin, convey.ShouldResemble, plugin)

			convey.Convey("And when the plugin is updated", func() {
				plugin.Name = "Updated Plugin"
				err = db.Plugins.CachePut(plugin)
				convey.So(err, convey.ShouldBeNil)

				actualPlugin, err := db.Plugins.CacheGet(plugin.ID)
				convey.So(err, convey.ShouldBeNil)
				convey.So(actualPlugin.Name, convey.ShouldEqual, "Updated Plugin")
			})
		})
	})
}
