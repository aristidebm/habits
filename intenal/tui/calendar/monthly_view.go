package calendar

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// MonthlyView handles rendering of the monthly calendar view
type MonthlyView struct {
	calendar *Calendar

	// Styles
	headerStyle        lipgloss.Style
	dateHeaderStyle    lipgloss.Style
	dayNameStyle       lipgloss.Style
	habitLabelStyle    lipgloss.Style
	selectedHabitStyle lipgloss.Style
	cellStyle          lipgloss.Style
	selectedCellStyle  lipgloss.Style
	todayStyle         lipgloss.Style
	emptyStyle         lipgloss.Style
}

// NewMonthlyView creates a new monthly view renderer
func NewMonthlyView(calendar *Calendar) *MonthlyView {
	return &MonthlyView{
		calendar: calendar,

		// Initialize styles
		headerStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12")).
			Padding(0, 1),
		dateHeaderStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Bold(true).
			Align(lipgloss.Center).
			Width(5),
		dayNameStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Align(lipgloss.Center).
			Width(5),
		habitLabelStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Width(20).
			Align(lipgloss.Left),
		selectedHabitStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("240")).
			Width(20).
			Align(lipgloss.Left).
			Bold(true),
		cellStyle: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(5),
		selectedCellStyle: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(5).
			Background(lipgloss.Color("240")).
			Foreground(lipgloss.Color("15")),
		todayStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true).
			Align(lipgloss.Center).
			Width(5),
		emptyStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("235")).
			Align(lipgloss.Center).
			Width(5),
	}
}

// Render renders the monthly view in Dijo-style grid format
func (m *MonthlyView) Render() string {
	if len(m.calendar.habits) == 0 {
		return "No habits to display"
	}

	var sb strings.Builder

	// Header with month/year
	header := fmt.Sprintf("%s %d",
		m.calendar.viewMonth.Format("January"),
		m.calendar.viewMonth.Year(),
	)
	monthIndicator := fmt.Sprintf("◀ [Month %d/12] ▶", int(m.calendar.viewMonth.Month()))

	// Calculate spacing
	headerLen := len(header)
	indicatorLen := len(monthIndicator)
	spacingLen := m.calendar.width - headerLen - indicatorLen - 4
	if spacingLen < 0 {
		spacingLen = 0
	}

	headerLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.headerStyle.Render(header),
		strings.Repeat(" ", spacingLen),
		m.headerStyle.Render(monthIndicator),
	)
	sb.WriteString(headerLine + "\n\n")

	// Get first and last day of month
	nextMonth := m.calendar.viewMonth.AddDate(0, 1, 0)
	lastDay := nextMonth.AddDate(0, 0, -1)
	daysInMonth := lastDay.Day()

	// Calculate week layout - find first day of month and its weekday
	firstDay := time.Date(m.calendar.viewMonth.Year(), m.calendar.viewMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	firstWeekday := int(firstDay.Weekday())
	if firstWeekday == 0 {
		firstWeekday = 7 // Sunday becomes 7
	}

	// Calculate number of weeks needed (including partial weeks)
	numWeeks := (daysInMonth+firstWeekday-2)/7 + 1

	// Build day of week headers (M T W T F S S)
	dayOfWeekRow := m.habitLabelStyle.Render("  ")
	dayAbbrevs := []string{"M", "T", "W", "T", "F", "S", "S"}

	for _, dayAbbrev := range dayAbbrevs {
		dayOfWeekRow += m.dayNameStyle.Render(dayAbbrev)
	}
	sb.WriteString(dayOfWeekRow + "\n")

	today := time.Now()

	// Render each habit with its own calendar grid
	for idx, habit := range m.calendar.habits {
		var habitLabel string
		if idx == m.calendar.selectedHabit {
			habitLabel = m.selectedHabitStyle.Render(habit.Name)
		} else {
			habitLabel = m.habitLabelStyle.Render(habit.Name)
		}
		sb.WriteString(habitLabel + "\n")

		// Render calendar grid for this habit
		for week := 0; week < numWeeks; week++ {
			weekRow := m.habitLabelStyle.Render("  ") // Indent for habit label

			for dayOfWeek := 0; dayOfWeek < 7; dayOfWeek++ {
				// Calculate the actual day number
				dayNumber := week*7 + dayOfWeek - (firstWeekday - 2)

				if dayNumber < 1 || dayNumber > daysInMonth {
					// Empty cell for days outside current month
					weekRow += m.emptyStyle.Render(" ")
					continue
				}

				date := time.Date(m.calendar.viewMonth.Year(), m.calendar.viewMonth.Month(), dayNumber, 0, 0, 0, 0, time.UTC)
				cellValue := m.getCompactCellValue(habit, date)

				// Check if this cell is selected
				isSelected := idx == m.calendar.selectedHabit &&
					date.Year() == m.calendar.selectedDate.Year() &&
					date.Month() == m.calendar.selectedDate.Month() &&
					date.Day() == m.calendar.selectedDate.Day()

				// Check if this cell is today
				isToday := date.Year() == today.Year() &&
					date.Month() == today.Month() &&
					date.Day() == today.Day()

				if isSelected {
					weekRow += m.selectedCellStyle.Render(cellValue)
				} else if isToday {
					weekRow += m.todayStyle.Render(cellValue)
				} else {
					weekRow += m.cellStyle.Render(cellValue)
				}
			}

			sb.WriteString(weekRow + "\n")
		}

		sb.WriteString("\n") // Spacing between habits
	}

	// Stats for selected habit
	if m.calendar.selectedHabit < len(m.calendar.habits) {
		currentHabit := m.calendar.habits[m.calendar.selectedHabit]
		stats := m.calculateStats(currentHabit)
		statsLine := fmt.Sprintf("Selected: %s | Streak: %d days | Completion: %d%%",
			currentHabit.Name, stats.Streak, stats.CompletionRate)
		sb.WriteString(m.habitLabelStyle.Render(statsLine) + "\n")
	}

	// Fill remaining vertical space
	lines := strings.Count(sb.String(), "\n")
	for i := lines; i < m.calendar.height-2; i++ {
		sb.WriteString("\n")
	}

	return sb.String()
}

// getCompactCellValue returns a compact display value for a habit on a specific date
func (m *MonthlyView) getCompactCellValue(habit Habit, date time.Time) string {
	entry, exists := m.calendar.GetEntry(habit.Name, date)

	today := time.Now()
	if date.After(today) {
		return "·"
	}

	if !exists {
		return "-"
	}

	switch habit.Type {
	case HabitTypeBit:
		if entry.Completed {
			return "●"
		}
		return "-"
	case HabitTypeCount:
		if entry.Value == "-" || entry.Value == "" {
			return "-"
		}
		// Show a filled circle for any count value
		return "●"
	case HabitTypeFloat:
		if entry.Value == "-" || entry.Value == "" {
			return "-"
		}
		// Show a filled circle for any float value
		return "●"
	}

	return "-"
}

// HabitStats holds statistics for a habit
type HabitStats struct {
	Streak         int
	CompletionRate int
}

// calculateStats calculates statistics for a habit in the current month
func (m *MonthlyView) calculateStats(habit Habit) HabitStats {
	stats := HabitStats{}

	// Calculate completion rate for the month
	firstDay := m.calendar.viewMonth
	nextMonth := m.calendar.viewMonth.AddDate(0, 1, 0)
	lastDay := nextMonth.AddDate(0, 0, -1)
	today := time.Now()

	totalDays := 0
	completedDays := 0

	for d := firstDay; !d.After(lastDay) && !d.After(today); d = d.AddDate(0, 0, 1) {
		totalDays++
		entry, exists := m.calendar.GetEntry(habit.Name, d)
		if exists && (entry.Completed || (entry.Value != "" && entry.Value != "-")) {
			completedDays++
		}
	}

	if totalDays > 0 {
		stats.CompletionRate = (completedDays * 100) / totalDays
	}

	// Calculate current streak
	currentDate := today
	for {
		entry, exists := m.calendar.GetEntry(habit.Name, currentDate)
		if !exists || (!entry.Completed && (entry.Value == "" || entry.Value == "-")) {
			break
		}
		stats.Streak++
		currentDate = currentDate.AddDate(0, 0, -1)

		// Stop if we go back more than a year
		if stats.Streak > 365 {
			break
		}
	}

	return stats
}
