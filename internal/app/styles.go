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
}

// NewStyles creates themed styles from a theme
func NewStyles(theme *Theme) *Styles {
	s := &Styles{theme: theme}

	// Header components
	s.Header = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Colors.Header)).
		Bold(true).
		Padding(0, 1)

	// Habit name styling
	s.HabitName = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Colors.HabitName)).
		Align(lipgloss.Left)

	s.HabitSelected = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Colors.HabitSelectedFG)).
		Background(lipgloss.Color(theme.Colors.HabitSelectedBG)).
		Align(lipgloss.Left).
		Bold(true)

	// Habit states
	s.Completed = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Colors.Completed))

	s.Missed = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Colors.Missed))

	s.Untracked = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Colors.Untracked))

	s.Today = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Colors.Today))

	// Special elements
	s.NoteCell = lipgloss.NewStyle().
		Background(lipgloss.Color(theme.Colors.NoteBG))

	s.DBPath = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Colors.DBPath))

	// Command line
	s.CLIPrompt = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Colors.Prompt))

	s.CLIError = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Colors.Error)).
		Bold(true)

	s.CLISuccess = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Colors.Success))

	// View-specific styles
	s.Weekly = NewWeeklyStyles(theme)
	s.Monthly = NewMonthlyStyles(theme)

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
			Foreground(lipgloss.Color(theme.Colors.Header)),

		DateHeader: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Colors.Header)).
			Align(lipgloss.Center).
			Width(8),

		DayName: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Colors.DBPath)).
			Align(lipgloss.Center).
			Width(8),

		Cell: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(8),

		NoteCell: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Colors.NoteFG)).
			Align(lipgloss.Center).
			Width(8).
			Background(lipgloss.Color(theme.Colors.NoteBG)),

		SelectedCell: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(8).
			Background(lipgloss.Color(theme.Colors.HabitSelectedBG)).
			Foreground(lipgloss.Color(theme.Colors.HabitSelectedFG)),

		TodayCell: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(8).
			Foreground(lipgloss.Color(theme.Colors.Today)).
			Bold(true),

		HabitLabel: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Colors.HabitName)).
			Width(20).
			Align(lipgloss.Left),

		SelectedHabitLabel: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Colors.HabitSelectedFG)).
			Background(lipgloss.Color(theme.Colors.HabitSelectedBG)).
			Width(20).
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
			Foreground(lipgloss.Color(theme.Colors.Header)),

		HabitName: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Colors.HabitName)).
			Align(lipgloss.Left),

		SelectedHabit: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Colors.HabitSelectedFG)).
			Background(lipgloss.Color(theme.Colors.HabitSelectedBG)).
			Width(20).
			Align(lipgloss.Left),

		Cell: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(4),

		NoteCell: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Colors.NoteFG)).
			Align(lipgloss.Center).
			Width(4).
			Background(lipgloss.Color(theme.Colors.NoteBG)),

		SelectedCell: lipgloss.NewStyle().
			Background(lipgloss.Color(theme.Colors.HabitSelectedBG)).
			Foreground(lipgloss.Color(theme.Colors.HabitSelectedFG)).
			Align(lipgloss.Center).
			Width(4),

		Empty: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Colors.Empty)).
			Align(lipgloss.Center).
			Width(4),

		Card: lipgloss.NewStyle().
			Padding(1, 1),

		Footer: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Colors.HabitName)),
	}
}
