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
				Header:          "#87CEEB", // Sky blue
				HabitName:       "#00CED1", // Dark turquoise
				HabitSelectedFG: "#FFFFFF", // White
				HabitSelectedBG: "#696969", // Dim gray
				Completed:       "#32CD32", // Lime green
				Missed:          "#DC143C", // Crimson
				Untracked:       "#708090", // Slate gray
				Today:           "#32CD32", // Lime green
				NoteBG:          "#2F4F4F", // Dark slate gray
				DBPath:          "#708090", // Slate gray
				Prompt:          "#87CEEB", // Sky blue
				Error:           "#DC143C", // Crimson
				Success:         "#32CD32", // Lime green
			},
		}
	case "dark":
		return &Theme{
			Name:   "Dark",
			Author: "Habits Team",
			Colors: ThemeColors{
				Header:          "#FF6B6B", // Coral red
				HabitName:       "#4ECDC4", // Medium turquoise
				HabitSelectedFG: "#000000", // Black
				HabitSelectedBG: "#FFFFFF", // White
				Completed:       "#45B7D1", // Sky blue
				Missed:          "#FFA07A", // Light salmon
				Untracked:       "#778899", // Light slate gray
				Today:           "#45B7D1", // Sky blue
				NoteBG:          "#2F4F4F", // Dark slate gray
				DBPath:          "#778899", // Light slate gray
				Prompt:          "#FF6B6B", // Coral red
				Error:           "#FFA07A", // Light salmon
				Success:         "#45B7D1", // Sky blue
			},
		}
	case "light":
		return &Theme{
			Name:   "Light",
			Author: "Habits Team",
			Colors: ThemeColors{
				Header:          "#2E86AB", // Steel blue
				HabitName:       "#A23B72", // Medium violet red
				HabitSelectedFG: "#FFFFFF", // White
				HabitSelectedBG: "#2E86AB", // Steel blue
				Completed:       "#F18F01", // Orange
				Missed:          "#C73E1D", // Firebrick
				Untracked:       "#A0A0A0", // Dark gray
				Today:           "#F18F01", // Orange
				NoteBG:          "#E8E8E8", // Light gray
				DBPath:          "#A0A0A0", // Dark gray
				Prompt:          "#2E86AB", // Steel blue
				Error:           "#C73E1D", // Firebrick
				Success:         "#F18F01", // Orange
			},
		}
	case "mono":
		return &Theme{
			Name:   "Monochrome",
			Author: "Habits Team",
			Colors: ThemeColors{
				Header:          "#000000", // Black
				HabitName:       "#333333", // Dark gray
				HabitSelectedFG: "#FFFFFF", // White
				HabitSelectedBG: "#000000", // Black
				Completed:       "#666666", // Medium gray
				Missed:          "#999999", // Light gray
				Untracked:       "#CCCCCC", // Very light gray
				Today:           "#666666", // Medium gray
				NoteBG:          "#F0F0F0", // Very light gray
				DBPath:          "#CCCCCC", // Very light gray
				Prompt:          "#000000", // Black
				Error:           "#999999", // Light gray
				Success:         "#666666", // Medium gray
			},
		}
	default:
		return nil
	}
}
