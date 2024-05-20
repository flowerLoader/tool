package types

type PluginCacheRecord struct {
	ID string // (owner/repo/tag#commit) is the primary key

	// These fields are read from `package.json` of the plugin
	Name     string
	Version  string
	Author   string
	License  string
	BugsURL  string
	Homepage string

	// This is read from `package.json` `dependencies` `@flowerloader/api` version
	APIVersion string

	// These fiels are read from the GitHub API
	Tags    string
	Summary string
}
