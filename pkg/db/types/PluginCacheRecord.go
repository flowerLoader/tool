package types

type PluginCacheRecord struct {
	ID        string `db:"id"`         // (owner/repo/tag#commit) is the primary key
	UpdatedAt string `db:"updated_at"` // Last time the SOURCE material was updated

	// These fields are read from `package.json` of the plugin
	Name     string `db:"name"`
	Version  string `db:"version"`
	Author   string `db:"author"`
	License  string `db:"license"`
	BugsURL  string `db:"bugs_url"`
	Homepage string `db:"homepage"`

	// This is read from `package.json` `dependencies` `@flowerloader/api` version
	APIVersion string `db:"api_version"`

	// These fiels are read from the GitHub API
	Tags    string `db:"tags"` // Comma separated list of tags
	Summary string `db:"summary"`
}
