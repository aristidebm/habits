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
	Header          lipgloss.AdaptiveColor `toml:"header"`
	HabitName       lipgloss.AdaptiveColor `toml:"habit_name"`
	HabitSelectedFG lipgloss.AdaptiveColor `toml:"habit_selected_fg"`
	HabitSelectedBG lipgloss.AdaptiveColor `toml:"habit_selected_bg"`

	// Habit states
	Completed lipgloss.AdaptiveColor `toml:"completed"`
	Missed    lipgloss.AdaptiveColor `toml:"missed"`
	Untracked lipgloss.AdaptiveColor `toml:"untracked"`
	Today     lipgloss.AdaptiveColor `toml:"today"`

	// Special elements
	NoteBG lipgloss.AdaptiveColor `toml:"note_bg"`
	DBPath lipgloss.AdaptiveColor `toml:"db_path"`

	// Command line
	Prompt  lipgloss.AdaptiveColor `toml:"prompt"`
	Error   lipgloss.AdaptiveColor `toml:"error"`
	Success lipgloss.AdaptiveColor `toml:"success"`
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
			Header:          lipgloss.AdaptiveColor{Light: "#87CEEB", Dark: "#87CEEB"},
			HabitName:       lipgloss.AdaptiveColor{Light: "#00CED1", Dark: "#00CED1"},
			HabitSelectedFG: lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"},
			HabitSelectedBG: lipgloss.AdaptiveColor{Light: "#696969", Dark: "#696969"},
			Completed:       lipgloss.AdaptiveColor{Light: "#32CD32", Dark: "#32CD32"},
			Missed:          lipgloss.AdaptiveColor{Light: "#DC143C", Dark: "#DC143C"},
			Untracked:       lipgloss.AdaptiveColor{Light: "#708090", Dark: "#708090"},
			Today:           lipgloss.AdaptiveColor{Light: "#32CD32", Dark: "#32CD32"},
			NoteBG:          lipgloss.AdaptiveColor{Light: "#2F4F4F", Dark: "#2F4F4F"},
			DBPath:          lipgloss.AdaptiveColor{Light: "#708090", Dark: "#708090"},
			Prompt:          lipgloss.AdaptiveColor{Light: "#87CEEB", Dark: "#87CEEB"},
			Error:           lipgloss.AdaptiveColor{Light: "#DC143C", Dark: "#DC143C"},
			Success:         lipgloss.AdaptiveColor{Light: "#32CD32", Dark: "#32CD32"},
		}
	case "dark":
		return &Theme{
			Name:            "Dark",
			Author:          "Habits Team",
			Header:          lipgloss.AdaptiveColor{Light: "#FF6B6B", Dark: "#FF6B6B"},
			HabitName:       lipgloss.AdaptiveColor{Light: "#4ECDC4", Dark: "#4ECDC4"},
			HabitSelectedFG: lipgloss.AdaptiveColor{Light: "#000000", Dark: "#000000"},
			HabitSelectedBG: lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"},
			Completed:       lipgloss.AdaptiveColor{Light: "#45B7D1", Dark: "#45B7D1"},
			Missed:          lipgloss.AdaptiveColor{Light: "#FFA07A", Dark: "#FFA07A"},
			Untracked:       lipgloss.AdaptiveColor{Light: "#778899", Dark: "#778899"},
			Today:           lipgloss.AdaptiveColor{Light: "#45B7D1", Dark: "#45B7D1"},
			NoteBG:          lipgloss.AdaptiveColor{Light: "#2F4F4F", Dark: "#2F4F4F"},
			DBPath:          lipgloss.AdaptiveColor{Light: "#778899", Dark: "#778899"},
			Prompt:          lipgloss.AdaptiveColor{Light: "#FF6B6B", Dark: "#FF6B6B"},
			Error:           lipgloss.AdaptiveColor{Light: "#FFA07A", Dark: "#FFA07A"},
			Success:         lipgloss.AdaptiveColor{Light: "#45B7D1", Dark: "#45B7D1"},
		}
	case "light":
		return &Theme{
			Name:            "Light",
			Author:          "Habits Team",
			Header:          lipgloss.AdaptiveColor{Light: "#2E86AB", Dark: "#2E86AB"},
			HabitName:       lipgloss.AdaptiveColor{Light: "#A23B72", Dark: "#A23B72"},
			HabitSelectedFG: lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"},
			HabitSelectedBG: lipgloss.AdaptiveColor{Light: "#2E86AB", Dark: "#2E86AB"},
			Completed:       lipgloss.AdaptiveColor{Light: "#F18F01", Dark: "#F18F01"},
			Missed:          lipgloss.AdaptiveColor{Light: "#C73E1D", Dark: "#C73E1D"},
			Untracked:       lipgloss.AdaptiveColor{Light: "#A0A0A0", Dark: "#A0A0A0"},
			Today:           lipgloss.AdaptiveColor{Light: "#F18F01", Dark: "#F18F01"},
			NoteBG:          lipgloss.AdaptiveColor{Light: "#E8E8E8", Dark: "#E8E8E8"},
			DBPath:          lipgloss.AdaptiveColor{Light: "#A0A0A0", Dark: "#A0A0A0"},
			Prompt:          lipgloss.AdaptiveColor{Light: "#2E86AB", Dark: "#2E86AB"},
			Error:           lipgloss.AdaptiveColor{Light: "#C73E1D", Dark: "#C73E1D"},
			Success:         lipgloss.AdaptiveColor{Light: "#F18F01", Dark: "#F18F01"},
		}
	case "mono":
		return &Theme{
			Name:            "Monochrome",
			Author:          "Habits Team",
			Header:          lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"},
			HabitName:       lipgloss.AdaptiveColor{Light: "#333333", Dark: "#CCCCCC"},
			HabitSelectedFG: lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"},
			HabitSelectedBG: lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"},
			Completed:       lipgloss.AdaptiveColor{Light: "#666666", Dark: "#999999"},
			Missed:          lipgloss.AdaptiveColor{Light: "#999999", Dark: "#666666"},
			Untracked:       lipgloss.AdaptiveColor{Light: "#CCCCCC", Dark: "#333333"},
			Today:           lipgloss.AdaptiveColor{Light: "#666666", Dark: "#999999"},
			NoteBG:          lipgloss.AdaptiveColor{Light: "#F0F0F0", Dark: "#1A1A1A"},
			DBPath:          lipgloss.AdaptiveColor{Light: "#CCCCCC", Dark: "#333333"},
			Prompt:          lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"},
			Error:           lipgloss.AdaptiveColor{Light: "#999999", Dark: "#666666"},
			Success:         lipgloss.AdaptiveColor{Light: "#666666", Dark: "#999999"},
		}
	default:
		return nil
	}
}
