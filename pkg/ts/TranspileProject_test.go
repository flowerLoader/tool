package ts

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

		// remove banner/footer for testing
		JS_BANNER = ""
		JS_FOOTER = ""

		dstFile := filepath.Join(tmpSrc, "test.js")
		Convey("should transpile a single file", func() {
			So(os.WriteFile(filepath.Join(tmpSrc, "test.ts"),
				[]byte("console.log('hello world')"), 0777), ShouldBeNil)
			So(TranspileProject(tmpSrc, dstFile), ShouldBeNil)

			_, err := os.Stat(dstFile)
			So(err, ShouldBeNil)
			defer os.Remove(dstFile)

			contents, err := os.ReadFile(dstFile)
			So(err, ShouldBeNil)
			So(string(contents), ShouldResemble, fmt.Sprintf("// %s/test.ts\nconsole.log(\"hello world\");\n", strings.Replace(tmpSrc, "\\", "/", -1)))
		})

		Convey("should transpile multiple files", func() {
			So(os.WriteFile(filepath.Join(tmpSrc, "index.ts"),
				[]byte("import { fn } from './fnTest'\nfn()"), 0777), ShouldBeNil)
			So(os.WriteFile(filepath.Join(tmpSrc, "fnTest.ts"),
				[]byte("export function fn() { console.log('hello world') }"), 0777), ShouldBeNil)

			So(TranspileProject(tmpSrc, dstFile), ShouldBeNil)

			_, err := os.Stat(dstFile)
			So(err, ShouldBeNil)
			defer os.Remove(dstFile)

			contents, err := os.ReadFile(dstFile)
			So(err, ShouldBeNil)
			So(string(contents), ShouldResemble, fmt.Sprintf("// %s/fnTest.ts\nfunction fn() {\n  console.log(\"hello world\");\n}\n\n// %s/index.ts\nfn();\n", strings.Replace(tmpSrc, "\\", "/", -1), strings.Replace(tmpSrc, "\\", "/", -1)))
		})

		Convey("should transpile with debug", func() {
			So(os.WriteFile(filepath.Join(tmpSrc, "debug.ts"),
				[]byte("const debuglogging = false;\nconsole.log(debuglogging);"), 0777), ShouldBeNil)
			So(TranspileProject(tmpSrc, dstFile, WithEntrypoints(
				filepath.Join(tmpSrc, "debug.ts"),
			), WithDebugMode(true)), ShouldBeNil)

			_, err := os.Stat(dstFile)
			So(err, ShouldBeNil)
			defer os.Remove(dstFile)

			contents, err := os.ReadFile(dstFile)
			So(err, ShouldBeNil)
			So(string(contents), ShouldResemble, fmt.Sprintf("// %s/debug.ts\nvar debuglogging = true;\nconsole.log(debuglogging);\n", strings.Replace(tmpSrc, "\\", "/", -1)))
		})

		Convey("can transpile from package.json", func() {
			So(os.WriteFile(filepath.Join(tmpSrc, "package.json"),
				[]byte(`{"main":"whatever.ts"}`), 0777), ShouldBeNil)
			So(os.WriteFile(filepath.Join(tmpSrc, "whatever.ts"),
				[]byte("console.log('hello world')"), 0777), ShouldBeNil)
			So(TranspileProject(tmpSrc, dstFile), ShouldBeNil)

			_, err := os.Stat(dstFile)
			So(err, ShouldBeNil)
			defer os.Remove(dstFile)

			contents, err := os.ReadFile(dstFile)
			So(err, ShouldBeNil)
			So(string(contents), ShouldResemble, fmt.Sprintf("// %s/whatever.ts\nconsole.log(\"hello world\");\n", strings.Replace(tmpSrc, "\\", "/", -1)))
		})
	})
}
