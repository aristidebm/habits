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
	headerStyle    lipgloss.Style
	subHeaderStyle lipgloss.Style
	dayNameStyle   lipgloss.Style
	dateStyle      lipgloss.Style
	todayDateStyle lipgloss.Style
	valueStyle     lipgloss.Style
	selectedStyle  lipgloss.Style
	emptyDateStyle lipgloss.Style
	statsStyle     lipgloss.Style
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
		subHeaderStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Padding(0, 1),
		dayNameStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Align(lipgloss.Center).
			Width(10),
		dateStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Align(lipgloss.Right).
			Width(10),
		todayDateStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true).
			Align(lipgloss.Right).
			Width(10),
		valueStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Align(lipgloss.Center).
			Width(10),
		selectedStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("240")).
			Foreground(lipgloss.Color("15")).
			Align(lipgloss.Center).
			Width(10),
		emptyDateStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("235")).
			Align(lipgloss.Center).
			Width(10),
		statsStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Padding(0, 1),
	}
}

// Render renders the monthly view
func (m *MonthlyView) Render() string {
	if len(m.calendar.habits) == 0 {
		return "No habits to display"
	}

	var sb strings.Builder

	currentHabit := m.calendar.habits[m.calendar.selectedHabit]

	// Header with month/year and habit info
	header := fmt.Sprintf("%s %d - %s",
		m.calendar.viewMonth.Format("January"),
		m.calendar.viewMonth.Year(),
		currentHabit.Name,
	)
	habitIndicator := fmt.Sprintf("[Habit %d/%d]", m.calendar.selectedHabit+1, len(m.calendar.habits))

	// Calculate spacing
	headerLen := len(header)
	indicatorLen := len(habitIndicator)
	spacingLen := m.calendar.width - headerLen - indicatorLen - 4
	if spacingLen < 0 {
		spacingLen = 0
	}

	headerLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.headerStyle.Render(header),
		strings.Repeat(" ", spacingLen),
		m.subHeaderStyle.Render(habitIndicator),
	)
	sb.WriteString(headerLine + "\n")
	sb.WriteString(strings.Repeat("─", m.calendar.width) + "\n")

	// Day names header
	dayNames := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	dayNamesRow := ""
	for _, day := range dayNames {
		dayNamesRow += m.dayNameStyle.Render(day)
	}
	sb.WriteString(dayNamesRow + "\n")

	// Get first day of month and calculate starting position
	firstDay := m.calendar.viewMonth
	weekday := int(firstDay.Weekday()) // 0 = Sunday

	// Get last day of month
	nextMonth := m.calendar.viewMonth.AddDate(0, 1, 0)
	lastDay := nextMonth.AddDate(0, 0, -1)
	daysInMonth := lastDay.Day()

	today := time.Now()
	currentDay := 1
	done := false

	// Generate calendar rows
	for week := 0; week < 6 && !done; week++ {
		sb.WriteString(strings.Repeat("─", m.calendar.width) + "\n")

		// Date row
		dateRow := ""
		for day := 0; day < 7; day++ {
			if (week == 0 && day < weekday) || currentDay > daysInMonth {
				dateRow += m.emptyDateStyle.Render("")
			} else {
				date := time.Date(m.calendar.viewMonth.Year(), m.calendar.viewMonth.Month(), currentDay, 0, 0, 0, 0, time.UTC)
				dateStr := fmt.Sprintf("%d", currentDay)

				isToday := date.Year() == today.Year() &&
					date.Month() == today.Month() &&
					date.Day() == today.Day()

				if isToday {
					dateRow += m.todayDateStyle.Render(dateStr)
				} else {
					dateRow += m.dateStyle.Render(dateStr)
				}

				if currentDay == daysInMonth {
					done = true
				}
				currentDay++
			}
		}
		sb.WriteString(dateRow + "\n")

		// Value row
		valueRow := ""
		currentDay -= 7
		if currentDay < 1 {
			currentDay = 1
		}

		for day := 0; day < 7; day++ {
			if (week == 0 && day < weekday) || currentDay > daysInMonth {
				valueRow += m.emptyDateStyle.Render("")
			} else {
				date := time.Date(m.calendar.viewMonth.Year(), m.calendar.viewMonth.Month(), currentDay, 0, 0, 0, 0, time.UTC)
				cellValue := m.calendar.GetCellValue(currentHabit, date)

				isSelected := date.Year() == m.calendar.selectedDate.Year() &&
					date.Month() == m.calendar.selectedDate.Month() &&
					date.Day() == m.calendar.selectedDate.Day()

				if isSelected {
					valueRow += m.selectedStyle.Render(cellValue)
				} else {
					valueRow += m.valueStyle.Render(cellValue)
				}
				currentDay++
			}
		}
		sb.WriteString(valueRow + "\n")
	}

	sb.WriteString(strings.Repeat("─", m.calendar.width) + "\n")

	// Statistics
	stats := m.calculateStats(currentHabit)
	statsLine := fmt.Sprintf("[TAB] Switch view  [Stats] Streak: %d  Rate: %d%%",
		stats.Streak, stats.CompletionRate)
	sb.WriteString(m.statsStyle.Render(statsLine) + "\n")

	// Fill remaining vertical space
	lines := strings.Count(sb.String(), "\n")
	for i := lines; i < m.calendar.height-2; i++ {
		sb.WriteString("\n")
	}

	return sb.String()
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
