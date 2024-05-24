package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParsePluginName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "github.com/flowerLoader/tool",
			expected: "github.com/flowerLoader/tool",
		},
		{
			name:     "https://www.github.com/flowerLoader/tool",
			expected: "github.com/flowerLoader/tool",
		},
		{
			name:     "/path/to/plugin",
			expected: "local//path/to/plugin", // double-slash is absolute path
		},
		{
			name:     "C:\\path\\to\\plugin",
			expected: "local/C:\\path\\to\\plugin",
		},
		{
			name:     "flowerLoader/limitbreaker",
			expected: "github.com/flowerLoader/limitbreaker",
		},
		{
			name:     "template",
			expected: "github.com/flowerLoader/template",
		},
	}

	Convey("parsePluginName", t, func() {
		for _, test := range tests {
			Convey(test.name, func() {
				So(parsePluginName(test.name), ShouldEqual, test.expected)
			})
		}
	})
}
