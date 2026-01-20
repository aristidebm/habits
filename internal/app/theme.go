package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/lipgloss"
)

// Theme represents a complete theme configuration
type Theme struct {
	Name   string `toml:"name"`
	Author string `toml:"author"`

	// Core colors
	Header          lipgloss.Color `toml:"header"`
	HabitName       lipgloss.Color `toml:"habit_name"`
	HabitSelectedFG lipgloss.Color `toml:"habit_selected_fg"`
	HabitSelectedBG lipgloss.Color `toml:"habit_selected_bg"`

	// Habit states
	Completed lipgloss.Color `toml:"completed"`
	Missed    lipgloss.Color `toml:"missed"`
	Untracked lipgloss.Color `toml:"untracked"`
	Today     lipgloss.Color `toml:"today"`

	// Special elements
	NoteBG lipgloss.Color `toml:"note_bg"`
	DBPath lipgloss.Color `toml:"db_path"`

	// Command line
	Prompt  lipgloss.Color `toml:"prompt"`
	Error   lipgloss.Color `toml:"error"`
	Success lipgloss.Color `toml:"success"`
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
			Name:            "Default",
			Author:          "Habits Team",
			Header:          "#87CEEB",
			HabitName:       "#00CED1",
			HabitSelectedFG: "#FFFFFF",
			HabitSelectedBG: "#696969",
			Completed:       "#32CD32",
			Missed:          "#DC143C",
			Untracked:       "#708090",
			Today:           "#32CD32",
			NoteBG:          "#2F4F4F",
			DBPath:          "#708090",
			Prompt:          "#87CEEB",
			Error:           "#DC143C",
			Success:         "#32CD32",
		}
	case "dark":
		return &Theme{
			Name:            "Dark",
			Author:          "Habits Team",
			Header:          "#FF6B6B",
			HabitName:       "#4ECDC4",
			HabitSelectedFG: "#000000",
			HabitSelectedBG: "#FFFFFF",
			Completed:       "#45B7D1",
			Missed:          "#FFA07A",
			Untracked:       "#778899",
			Today:           "#45B7D1",
			NoteBG:          "#2F4F4F",
			DBPath:          "#778899",
			Prompt:          "#FF6B6B",
			Error:           "#FFA07A",
			Success:         "#45B7D1",
		}
	case "light":
		return &Theme{
			Name:            "Light",
			Author:          "Habits Team",
			Header:          "#2E86AB",
			HabitName:       "#A23B72",
			HabitSelectedFG: "#FFFFFF",
			HabitSelectedBG: "#2E86AB",
			Completed:       "#F18F01",
			Missed:          "#C73E1D",
			Untracked:       "#A0A0A0",
			Today:           "#F18F01",
			NoteBG:          "#E8E8E8",
			DBPath:          "#A0A0A0",
			Prompt:          "#2E86AB",
			Error:           "#C73E1D",
			Success:         "#F18F01",
		}
	case "mono":
		return &Theme{
			Name:            "Monochrome",
			Author:          "Habits Team",
			Header:          "#000000",
			HabitName:       "#333333",
			HabitSelectedFG: "#FFFFFF",
			HabitSelectedBG: "#000000",
			Completed:       "#666666",
			Missed:          "#999999",
			Untracked:       "#CCCCCC",
			Today:           "#666666",
			NoteBG:          "#F0F0F0",
			DBPath:          "#CCCCCC",
			Prompt:          "#000000",
			Error:           "#999999",
			Success:         "#666666",
		}
	default:
		return nil
	}
}
