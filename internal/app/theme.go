package app

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// ThemeColors represents the color configuration
type ThemeColors struct {
	Header          string `toml:"header"`
	HabitName       string `toml:"habit_name"`
	HabitSelectedFG string `toml:"habit_selected_fg"`
	HabitSelectedBG string `toml:"habit_selected_bg"`
	Completed       string `toml:"completed"`
	Missed          string `toml:"missed"`
	Untracked       string `toml:"untracked"`
	Today           string `toml:"today"`
	NoteBG          string `toml:"note_bg"`
	DBPath          string `toml:"db_path"`
	Prompt          string `toml:"prompt"`
	Error           string `toml:"error"`
	Success         string `toml:"success"`
}

// Theme represents a complete theme configuration
type Theme struct {
	Name   string      `toml:"name"`
	Author string      `toml:"author"`
	Colors ThemeColors `toml:"colors"`
}

// ThemeManager manages theme loading and switching
type ThemeManager struct {
	themes map[string]*Theme
}

// NewThemeManager creates a new theme manager
func NewThemeManager() *ThemeManager {
	return &ThemeManager{
		themes: make(map[string]*Theme),
	}
}

// LoadTheme loads a theme by name
func (tm *ThemeManager) LoadTheme(name string) (*Theme, error) {
	slog.Info("Loading theme ...")
	// Check if already loaded
	if theme, exists := tm.themes[name]; exists {
		return theme, nil
	}

	// Try to load from file
	themePath := filepath.Join(ThemeDir(), name+".toml")
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		// Try built-in themes
		if theme := tm.getBuiltInTheme(name); theme != nil {
			tm.themes[name] = theme
			return theme, nil
		}
		return nil, fmt.Errorf("theme '%s' not found", name)
	}

	// Load from file
	var theme Theme
	if _, err := toml.DecodeFile(themePath, &theme); err != nil {
		return nil, fmt.Errorf("failed to load theme '%s': %w", name, err)
	}

	tm.themes[name] = &theme
	return &theme, nil
}

// ListThemes returns a list of available theme names
func (tm *ThemeManager) ListThemes() []string {
	var themes []string

	// Add built-in themes
	builtIn := []string{"default", "dark", "light", "mono"}
	themes = append(themes, builtIn...)

	// Add user themes
	if entries, err := os.ReadDir(ThemeDir()); err == nil {
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".toml") {
				name := strings.TrimSuffix(entry.Name(), ".toml")
				// Avoid duplicates
				found := false
				for _, existing := range themes {
					if existing == name {
						found = true
						break
					}
				}
				if !found {
					themes = append(themes, name)
				}
			}
		}
	}

	return themes
}

// getBuiltInTheme returns a built-in theme by name
func (tm *ThemeManager) getBuiltInTheme(name string) *Theme {
	switch name {
	case "default":
		return &Theme{
			Name:   "Default",
			Author: "Habits Team",
			Colors: ThemeColors{
				Header:          "12",  // Bright blue
				HabitName:       "14",  // Bright cyan
				HabitSelectedFG: "15",  // Bright white
				HabitSelectedBG: "240", // Gray
				Completed:       "10",  // Bright green
				Missed:          "9",   // Bright red
				Untracked:       "8",   // Gray
				Today:           "11",  // Bright yellow
				NoteBG:          "17",  // Dark blue
				DBPath:          "12",  // Bright blue (changed for testing)
				Prompt:          "12",  // Bright blue
				Error:           "9",   // Bright red
				Success:         "10",  // Bright green
			},
		}
	case "dark":
		return &Theme{
			Name:   "Dark",
			Author: "Habits Team",
			Colors: ThemeColors{
				Header:          "13", // Bright magenta
				HabitName:       "14", // Bright cyan
				HabitSelectedFG: "0",  // Black
				HabitSelectedBG: "15", // Bright white
				Completed:       "12", // Bright blue
				Missed:          "11", // Bright yellow
				Untracked:       "8",  // Gray
				Today:           "12", // Bright blue
				NoteBG:          "17", // Dark blue
				DBPath:          "8",  // Gray
				Prompt:          "13", // Bright magenta
				Error:           "11", // Bright yellow
				Success:         "12", // Bright blue
			},
		}
	case "light":
		return &Theme{
			Name:   "Light",
			Author: "Habits Team",
			Colors: ThemeColors{
				Header:          "4",  // Blue
				HabitName:       "5",  // Magenta
				HabitSelectedFG: "15", // Bright white
				HabitSelectedBG: "4",  // Blue
				Completed:       "3",  // Yellow
				Missed:          "1",  // Red
				Untracked:       "7",  // White/gray
				Today:           "3",  // Yellow
				NoteBG:          "7",  // White/gray
				DBPath:          "7",  // White/gray
				Prompt:          "4",  // Blue
				Error:           "1",  // Red
				Success:         "3",  // Yellow
			},
		}
	case "mono":
		return &Theme{
			Name:   "Monochrome",
			Author: "Habits Team",
			Colors: ThemeColors{
				Header:          "0",  // Black
				HabitName:       "7",  // White
				HabitSelectedFG: "15", // Bright white
				HabitSelectedBG: "0",  // Black
				Completed:       "8",  // Gray
				Missed:          "7",  // White
				Untracked:       "8",  // Gray
				Today:           "8",  // Gray
				NoteBG:          "7",  // White
				DBPath:          "8",  // Gray
				Prompt:          "0",  // Black
				Error:           "7",  // White
				Success:         "8",  // Gray
			},
		}
	default:
		return nil
	}
}
