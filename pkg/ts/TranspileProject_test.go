package ts

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTranspileProject(t *testing.T) {
	Convey("TranspileProject", t, func() {
		// make temp directory for testing sources
		tmpSrc := filepath.Join(os.TempDir(), "_flowerLoader_test_")
		So(os.MkdirAll(tmpSrc, 0777), ShouldBeNil)
		defer os.RemoveAll(tmpSrc)

		Convey("should error on invalid sources", func() {
			So(TranspileProject("/does/not/exist", ""), ShouldNotBeNil)
			So(TranspileProject("go.mod", ""), ShouldNotBeNil)
		})

		Convey("should error on no entrypoint", func() {
			err := TranspileProject(tmpSrc, "")
			So(errors.Is(err, ErrNoEntrypoint), ShouldBeTrue)
		})

		Convey("should error on invalid entrypoint", func() {
			err := TranspileProject(tmpSrc, "", WithEntrypoints("invalid"))
			So(errors.Is(err, ErrNoEntrypoint), ShouldBeTrue)
		})

		Convey("given valid sources", func() {
			// create a test source file
			So(os.WriteFile(filepath.Join(tmpSrc, "test.ts"),
				[]byte("console.log('hello world')"), 0777), ShouldBeNil)

			dstFile := filepath.Join(tmpSrc, "test.js")
			Convey("should transpile a single file", func() {
				So(TranspileProject(tmpSrc, dstFile), ShouldBeNil)

				_, err := os.Stat(dstFile)
				So(err, ShouldBeNil)

				// check contents of
				contents, err := os.ReadFile(dstFile)
				So(err, ShouldBeNil)
				So(string(contents), ShouldEqual, "console.log('hello world')")
			})
		})
	})
}
