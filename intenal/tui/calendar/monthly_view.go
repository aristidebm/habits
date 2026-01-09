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

// Render renders the monthly view
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

	// Build date headers (showing all days 1-31)
	dateHeaderRow := m.habitLabelStyle.Render("")
	today := time.Now()

	for day := 1; day <= daysInMonth; day++ {
		date := time.Date(m.calendar.viewMonth.Year(), m.calendar.viewMonth.Month(), day, 0, 0, 0, 0, time.UTC)
		dateStr := fmt.Sprintf("%d", day)

		isToday := date.Year() == today.Year() &&
			date.Month() == today.Month() &&
			date.Day() == today.Day()

		if isToday {
			dateHeaderRow += m.todayStyle.Render(dateStr)
		} else {
			dateHeaderRow += m.dateHeaderStyle.Render(dateStr)
		}
	}
	sb.WriteString(dateHeaderRow + "\n")

	// Day of week abbreviations (M T W T F S S pattern)
	dayOfWeekRow := m.habitLabelStyle.Render("")
	dayAbbrevs := []string{"S", "M", "T", "W", "T", "F", "S"}

	for day := 1; day <= daysInMonth; day++ {
		date := time.Date(m.calendar.viewMonth.Year(), m.calendar.viewMonth.Month(), day, 0, 0, 0, 0, time.UTC)
		weekday := int(date.Weekday())
		dayOfWeekRow += m.dayNameStyle.Render(dayAbbrevs[weekday])
	}
	sb.WriteString(dayOfWeekRow + "\n\n")

	// Habit rows - show all habits
	for idx, habit := range m.calendar.habits {
		var habitLabel string
		if idx == m.calendar.selectedHabit {
			habitLabel = m.selectedHabitStyle.Render(habit.Name)
		} else {
			habitLabel = m.habitLabelStyle.Render(habit.Name)
		}
		row := habitLabel

		for day := 1; day <= daysInMonth; day++ {
			date := time.Date(m.calendar.viewMonth.Year(), m.calendar.viewMonth.Month(), day, 0, 0, 0, 0, time.UTC)
			cellValue := m.getCompactCellValue(habit, date)

			// Check if this cell is selected
			isSelected := idx == m.calendar.selectedHabit &&
				date.Year() == m.calendar.selectedDate.Year() &&
				date.Month() == m.calendar.selectedDate.Month() &&
				date.Day() == m.calendar.selectedDate.Day()

			if isSelected {
				row += m.selectedCellStyle.Render(cellValue)
			} else {
				row += m.cellStyle.Render(cellValue)
			}
		}
		sb.WriteString(row + "\n")
	}

	sb.WriteString("\n")

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
