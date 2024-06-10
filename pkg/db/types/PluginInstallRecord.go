package types

type PluginInstallRecord struct {
	ID string `db:"id"` // (owner/repo/tag#commit) is the primary key

	// Installation Status
	Enabled     bool   `db:"enabled"`
	Path        string `db:"path"`
	InstalledAt string `db:"installed_at"` // encoded as RFC3339
	UpdatedAt   string `db:"updated_at"`   // encoded as RFC3339
}
