package types

type PluginInstallRecord struct {
	ID string // (owner/repo/tag#commit) is the primary key

	// Installation Status
	Enabled     bool
	InstalledAt string // encoded as RFC3339
	UpdatedAt   string // encoded as RFC3339
}
