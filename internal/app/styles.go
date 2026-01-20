package app

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles contains all themed lipgloss styles for the application
type Styles struct {
	theme *Theme

	// Header components
	Header lipgloss.Style

	// Habit name styling
	HabitName     lipgloss.Style
	HabitSelected lipgloss.Style

	// Habit states
	Completed lipgloss.Style
	Missed    lipgloss.Style
	Untracked lipgloss.Style
	Today     lipgloss.Style

	// Special elements
	NoteCell lipgloss.Style
	DBPath   lipgloss.Style

	// Command line
	CLIPrompt  lipgloss.Style
	CLIError   lipgloss.Style
	CLISuccess lipgloss.Style

	// View-specific styles
	Weekly  *WeeklyStyles
	Monthly *MonthlyStyles
	Heatmap *HeatmapStyles
}

// NewStyles creates themed styles from a theme
func NewStyles(theme *Theme) *Styles {
	s := &Styles{theme: theme}

	// Header components
	s.Header = lipgloss.NewStyle().
		Foreground(theme.Header)

	// Habit name styling
	s.HabitName = lipgloss.NewStyle().
		Foreground(theme.HabitName).
		Align(lipgloss.Left)

	s.HabitSelected = lipgloss.NewStyle().
		Foreground(theme.HabitSelectedFG).
		Background(theme.HabitSelectedBG).
		Align(lipgloss.Left)

	// Habit states
	s.Completed = lipgloss.NewStyle().
		Foreground(theme.Completed)

	s.Missed = lipgloss.NewStyle().
		Foreground(theme.Missed)

	s.Untracked = lipgloss.NewStyle().
		Foreground(theme.Untracked)

	s.Today = lipgloss.NewStyle().
		Foreground(theme.Today)

	// Special elements
	s.NoteCell = lipgloss.NewStyle().
		Background(theme.NoteBG)

	s.DBPath = lipgloss.NewStyle().
		Foreground(theme.DBPath)

	// Command line
	s.CLIPrompt = lipgloss.NewStyle().
		Foreground(theme.Prompt)

	s.CLIError = lipgloss.NewStyle().
		Foreground(theme.Error).
		Bold(true)

	s.CLISuccess = lipgloss.NewStyle().
		Foreground(theme.Success)

	// View-specific styles
	s.Weekly = NewWeeklyStyles(theme)
	s.Monthly = NewMonthlyStyles(theme)
	s.Heatmap = NewHeatmapStyles(theme)

	return s
}

// WeeklyStyles contains styles specific to the weekly view
type WeeklyStyles struct {
	Header             lipgloss.Style
	DateHeader         lipgloss.Style
	DayName            lipgloss.Style
	Cell               lipgloss.Style
	NoteCell           lipgloss.Style
	SelectedCell       lipgloss.Style
	TodayCell          lipgloss.Style
	HabitLabel         lipgloss.Style
	SelectedHabitLabel lipgloss.Style
}

// NewWeeklyStyles creates weekly view styles
func NewWeeklyStyles(theme *Theme) *WeeklyStyles {
	return &WeeklyStyles{
		Header: lipgloss.NewStyle().
			Foreground(theme.Header),

		DateHeader: lipgloss.NewStyle().
			Foreground(theme.HabitSelectedFG).
			Align(lipgloss.Center),

		DayName: lipgloss.NewStyle().
			Foreground(theme.DBPath).
			Align(lipgloss.Center),

		Cell: lipgloss.NewStyle().
			Align(lipgloss.Center),

		NoteCell: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Background(theme.NoteBG),

		SelectedCell: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Background(theme.HabitSelectedBG).
			Foreground(theme.HabitSelectedFG),

		TodayCell: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Foreground(theme.Today).
			Bold(true),

		HabitLabel: lipgloss.NewStyle().
			Foreground(theme.HabitName).
			Align(lipgloss.Left),

		SelectedHabitLabel: lipgloss.NewStyle().
			Foreground(theme.HabitSelectedFG).
			Background(theme.HabitSelectedBG).
			Align(lipgloss.Left),
	}
}

// MonthlyStyles contains styles specific to the monthly view
type MonthlyStyles struct {
	Header        lipgloss.Style
	HabitName     lipgloss.Style
	SelectedHabit lipgloss.Style
	Cell          lipgloss.Style
	NoteCell      lipgloss.Style
	SelectedCell  lipgloss.Style
	Empty         lipgloss.Style
	Card          lipgloss.Style
	Footer        lipgloss.Style
}

// NewMonthlyStyles creates monthly view styles
func NewMonthlyStyles(theme *Theme) *MonthlyStyles {
	return &MonthlyStyles{
		Header: lipgloss.NewStyle().
			Foreground(theme.Header),

		HabitName: lipgloss.NewStyle().
			Foreground(theme.HabitName).
			Align(lipgloss.Left),

		SelectedHabit: lipgloss.NewStyle().
			Foreground(theme.HabitSelectedFG).
			Align(lipgloss.Left),

		Cell: lipgloss.NewStyle().
			Foreground(theme.HabitSelectedFG).
			Align(lipgloss.Center),

		NoteCell: lipgloss.NewStyle().
			Foreground(theme.HabitSelectedFG).
			Align(lipgloss.Center).
			Background(theme.NoteBG),

		SelectedCell: lipgloss.NewStyle().
			Background(theme.HabitSelectedBG).
			Foreground(theme.HabitSelectedFG).
			Align(lipgloss.Center),

		Empty: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#235", Dark: "#888"}).
			Align(lipgloss.Center),

		Card: lipgloss.NewStyle(),

		Footer: lipgloss.NewStyle().
			Foreground(theme.HabitName),
	}
}

// HeatmapStyles contains styles specific to the heatmap view
type HeatmapStyles struct {
	Header        lipgloss.Style
	HabitName     lipgloss.Style
	SelectedHabit lipgloss.Style
	Completed     lipgloss.Style
	Missed        lipgloss.Style
	Note          lipgloss.Style
	Future        lipgloss.Style
	DayLabel      lipgloss.Style
}

// NewHeatmapStyles creates heatmap view styles
func NewHeatmapStyles(theme *Theme) *HeatmapStyles {
	return &HeatmapStyles{
		Header: lipgloss.NewStyle().
			Foreground(theme.Header),

		HabitName: lipgloss.NewStyle().
			Foreground(theme.HabitName).
			Align(lipgloss.Left),

		SelectedHabit: lipgloss.NewStyle().
			Foreground(theme.HabitSelectedFG).
			Background(theme.HabitSelectedBG).
			Align(lipgloss.Left),

		Completed: lipgloss.NewStyle().
			Foreground(theme.Completed),

		Missed: lipgloss.NewStyle().
			Foreground(theme.Missed),

		Note: lipgloss.NewStyle().
			Background(theme.NoteBG),

		Future: lipgloss.NewStyle().
			Foreground(theme.Untracked),

		DayLabel: lipgloss.NewStyle().
			Foreground(theme.DBPath).
			Align(lipgloss.Center),
	}
}
