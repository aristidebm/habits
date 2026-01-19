package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/pressly/goose/v3"
)

type App struct {
	*Store
	config *Config
}

func NewApp(dbPath string) (*App, error) {
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Use config database path if no explicit path provided
	if dbPath == "" {
		dbPath = config.DB.Path
	}

	db, err := OpenDB(dbPath)
	if err != nil {
		return nil, err
	}

	return &App{
		Store:  &Store{db: db},
		config: config,
	}, nil
}

func (a *App) Migrate() error {
	return Migrate(a.Store.db)
}

func (a *App) Export(ctx context.Context) ([]ExportHabit, error) {
	habits, err := a.ListHabits(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list habits: %w", err)
	}

	exportHabits := []ExportHabit{}
	for _, habit := range habits {
		// Get all entries for this habit (no date limit)
		startDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Now().AddDate(10, 0, 0) // Far future
		entries, err := a.ListEntries(ctx, habit.ID, startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("failed to list entries for habit %s: %w", habit.Name, err)
		}

		exportEntries := []ExportEntry{}
		for _, entry := range entries {
			// Get notes for this entry
			notes, err := a.ListNotes(ctx, entry.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to list notes for entry %d: %w", entry.ID, err)
			}

			exportNotes := []ExportNote{}
			for _, note := range notes {
				exportNotes = append(exportNotes, ExportNote{
					Note: note.Note,
				})
			}

			exportEntries = append(exportEntries, ExportEntry{
				ID:    entry.ID,
				Date:  entry.Date.Format("2006-01-02"),
				Value: entry.Value,
				Notes: exportNotes,
			})
		}

		// Ensure entries is never nil
		if exportEntries == nil {
			exportEntries = []ExportEntry{}
		}

		exportHabits = append(exportHabits, ExportHabit{
			ID:      habit.ID,
			Name:    habit.Name,
			Goal:    habit.Goal,
			Kind:    habit.Type.String(),
			Entries: exportEntries,
		})
	}

	return exportHabits, nil
}

// GetConfig returns the application configuration
func (a *App) GetConfig() *Config {
	return a.config
}

// LoadConfig loads the configuration from file
func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	configPath := ConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Config file doesn't exist, create with defaults
		if err := EnsureConfigDir(); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}
		if err := SaveConfig(config); err != nil {
			return nil, fmt.Errorf("failed to save default config: %w", err)
		}
		return config, nil
	}

	// Load existing config
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return config, nil
}

// SaveConfig saves the configuration to file
func SaveConfig(config *Config) error {
	if err := EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := ConfigPath()
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

func (a *App) Close() error {
	return a.Store.Close()
}
