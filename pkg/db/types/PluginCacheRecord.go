package types

type PluginCacheRecord struct {
	ID string // (owner/repo/tag#commit) is the primary key

	// META
	GUID    string
	Version string
	Name    string
	Author  string

	// GITHUB
	Tags    string
	Summary string
}
