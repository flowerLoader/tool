package main

import (
	"fmt"
	"net/url"
	"strings"
)

// parsePluginName takes a partial plugin name (full URL, org/repo, or local
// path) and returns the full name of the plugin ({github.com|local}/org/repo)
func parsePluginName(name string) string {
	u, err := url.Parse(name)
	if err == nil && u.Scheme != "" && u.Hostname() != "" {
		// https://www.github.com/flowerLoader/tool
		// -> github.com/flowerLoader/tool
		return strings.TrimPrefix(u.Hostname(), "www.") + u.Path
	}

	parts := strings.Split(name, "/")
	if len(parts) == 1 && !strings.Contains(name, "\\") {
		// template -> github.com/flowerLoader/template
		return fmt.Sprintf("%s/flowerLoader/%s", GITHUB_PKG, name)
	}

	if parts[0] != "local" && len(parts) == 2 {
		// flowerLoader/tool
		// -> github.com/flowerLoader/tool
		return fmt.Sprintf("%s/%s", GITHUB_PKG, name)
	}

	// github.com/flowerLoader/tool
	// -> github.com/flowerLoader/tool
	if parts[0] == GITHUB_PKG && len(parts) == 3 {
		return name
	}

	// local/some/path/...
	// -> local/some/path/...
	// C:\path\to\...
	// -> local/C:\path\to\...
	return "local/" + strings.Join(parts, "/")
}
