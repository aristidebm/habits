package app

import (
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	DB    DBConfig    `toml:"db"`
	Views ViewsConfig `toml:"views"`
}

// DBConfig represents database configuration
type DBConfig struct {
	Path string `toml:"path"`
}

// ViewsConfig represents view-specific configurations
type ViewsConfig struct {
	Weekly  ViewConfig `toml:"weekly"`
	Monthly ViewConfig `toml:"monthly"`
	Heatmap ViewConfig `toml:"heatmap"`
}

// ViewConfig represents configuration for a specific view
type ViewConfig struct {
	Missed    string `toml:"missed"`
	Completed string `toml:"completed"`
	Untracked string `toml:"untracked"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()

	return &Config{
		DB: DBConfig{
			Path: filepath.Join(homeDir, ".local", "state", "habits", "habits.db"),
		},
		Views: ViewsConfig{
			Weekly: ViewConfig{
				Missed:    "○",
				Completed: "●",
				Untracked: "✗",
			},
			Monthly: ViewConfig{
				Missed:    "○",
				Completed: "●",
				Untracked: "◦",
			},
			Heatmap: ViewConfig{
				Missed:    "○",
				Completed: "●",
				Untracked: "✗",
			},
		},
	}
}

// ConfigPath returns the default configuration file path
func ConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "habits", "config.toml")
}

// EnsureConfigDir ensures the configuration directory exists
func EnsureConfigDir() error {
	configPath := ConfigPath()
	configDir := filepath.Dir(configPath)
	return os.MkdirAll(configDir, 0755)
}
