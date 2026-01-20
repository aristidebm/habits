package app

import (
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	DB    DBConfig    `toml:"db"`
	Views ViewsConfig `toml:"views"`
	Theme ThemeConfig `toml:"theme"`
}

// DBConfig represents database configuration
type DBConfig struct {
	Path string `toml:"path"`
}

// ViewsConfig represents view-specific configurations
type ViewsConfig struct {
	Weekly  ViewConfig `toml:"weekly"`
	Monthly ViewConfig `toml:"monthly"`
}

// ViewConfig represents configuration for a specific view
type ViewConfig struct {
	Missed    string `toml:"missed"`
	Completed string `toml:"completed"`
	Untracked string `toml:"untracked"`
}

// ThemeConfig represents theme configuration
type ThemeConfig struct {
	Name string `toml:"name"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		DB: DBConfig{
			Path: StatePath(),
		},
		Views: ViewsConfig{
			Weekly: ViewConfig{
				Missed:    "⛌",
				Completed: "🗸",
				Untracked: "○",
			},
			Monthly: ViewConfig{
				Missed:    "⛌",
				Completed: "🗸",
				Untracked: "○",
			},
		},
		Theme: ThemeConfig{
			Name: "default",
		},
	}
}

// ConfigPath returns the default configuration file path
func ConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "habits", "config.toml")
}

// StatePath returns the default database file path
func StatePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".local", "state", "habits", "habits.db")
}

// EnsureConfigDir ensures the configuration directory exists
func EnsureConfigDir() error {
	configPath := ConfigPath()
	return EnsureDir(configPath)
}

// EnsureConfigDir ensures the configuration directory exists
func EnsureDir(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0755)
}

// ThemeDir returns the default theme directory path
func ThemeDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "habits", "themes")
}

// EnsureThemeDir ensures the theme directory exists
func EnsureThemeDir() error {
	return os.MkdirAll(ThemeDir(), 0755)
}
