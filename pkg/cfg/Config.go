package cfg

type Config struct {
	App       AppConfig        `json:"app"`
	Games     []GameConfig     `json:"games"`
	Providers []ProviderConfig `json:"providers"`
}

type GameMeta struct {
	Name map[string]string `json:"name"` // [en] => "Game Name"
	OS   []string          `json:"os"`
}

type GameSubsystem struct {
	ID      string `json:"id"`
	AppID   int    `json:"app_id"`
	AppName string `json:"app_name"`
}

type GameConfig struct {
	Loader    LoaderConfig  `json:"loader"`
	Meta      GameMeta      `json:"meta"`
	Subsystem GameSubsystem `json:"subsystem"`
}

type LoaderBuildConfig struct {
	ID          string   `json:"id"`
	Entrypoints []string `json:"entrypoints"`
}

type LoaderConfig struct {
	ID    string            `json:"id"`
	Base  []string          `json:"base"`
	Build LoaderBuildConfig `json:"build"`
}

type AppConfig struct {
	Environments []EnvironmentConfig `json:"environments"`
}

type ProviderConfig struct {
	Name      string   `json:"name"`
	Hosts     []string `json:"hosts"`
	Schemas   []string `json:"schemas"`
	Subsystem string   `json:"subsystem"`
}

type EnvironmentConfig struct {
	DB   string `json:"db"`   // Path to the database file
	Game string `json:"game"` // Game ID (format: "subsystem:app_id")
	Path string `json:"path"` // Path to the game's installation directory
}
