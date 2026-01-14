package calendar

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// WeeklyView handles rendering of the weekly calendar view
type WeeklyView struct {
	calendar *Calendar

	// Styles
	headerStyle             lipgloss.Style
	dateHeaderStyle         lipgloss.Style
	dayNameStyle            lipgloss.Style
	cellStyle               lipgloss.Style
	selectedCellStyle       lipgloss.Style
	todayStyle              lipgloss.Style
	habitLabelStyle         lipgloss.Style
	selectedHabitLabelStyle lipgloss.Style
}

// NewWeeklyView creates a new weekly view renderer
func NewWeeklyView(calendar *Calendar) *WeeklyView {
	return &WeeklyView{
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
			Width(8),
		dayNameStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Align(lipgloss.Center).
			Width(8),
		cellStyle: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(8),
		selectedCellStyle: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(8).
			Background(lipgloss.Color("240")).
			Foreground(lipgloss.Color("15")),
		todayStyle: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(8).
			Foreground(lipgloss.Color("10")).
			Bold(true),
		habitLabelStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Width(20).
			Align(lipgloss.Left),
		selectedHabitLabelStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("240")).
			Width(20).
			Align(lipgloss.Left).
			Bold(true),
	}
}

// RenderHeader renders the fixed header (week info, day names, dates)
func (w *WeeklyView) RenderHeader() string {
	var sb strings.Builder

	// Calculate week number
	_, week := w.calendar.viewStartDate.ISOWeek()
	endDate := w.calendar.viewStartDate.AddDate(0, 0, 6)

	// Calculate day of year and total days in year
	dayOfYear := w.calendar.selectedDate.YearDay()
	year := w.calendar.selectedDate.Year()
	endOfYear := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
	totalDays := endOfYear.YearDay()

	// Header with week range and week number
	header := fmt.Sprintf("Week: %s - %s, %d",
		w.calendar.viewStartDate.Format("Jan 02"),
		endDate.Format("Jan 02"),
		w.calendar.viewStartDate.Year(),
	)
	weekIndicator := fmt.Sprintf("◀ [%d/52] ▶", week)
	dayIndicator := fmt.Sprintf("◀ [%d/%d] ▶", dayOfYear, totalDays)

	// Calculate spacing to fill the width (now with middle element)
	headerLen := len(header)
	indicatorLen := len(weekIndicator)
	dayIndicatorLen := len(dayIndicator)
	availableSpace := w.calendar.width - headerLen - indicatorLen - dayIndicatorLen - 6
	spacingLen := availableSpace / 2
	if spacingLen < 0 {
		spacingLen = 0
	}

	headerLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		w.headerStyle.Render(header),
		strings.Repeat(" ", spacingLen),
		w.headerStyle.Render(dayIndicator),
		strings.Repeat(" ", spacingLen),
		w.headerStyle.Render(weekIndicator),
	)
	sb.WriteString(headerLine + "\n\n")

	// Calculate how many days can fit in the available width
	availableWidth := w.calendar.width - 20
	daysToShow := availableWidth / 8
	if daysToShow < 7 {
		daysToShow = 7
	}

	// Day names row - start from viewStartDate's weekday
	allDayNames := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	startDay := int(w.calendar.viewStartDate.Weekday())
	if startDay == 0 {
		startDay = 7 // Sunday becomes 7
	}

	// Create day names slice starting from the correct day
	dayNamesRow := w.habitLabelStyle.Render("")
	for i := 0; i < daysToShow; i++ {
		dayIndex := (startDay - 1 + i) % 7
		dayName := allDayNames[dayIndex]
		dayNamesRow += w.dayNameStyle.Render(dayName)
	}
	sb.WriteString(dayNamesRow + "\n")

	// Dates row
	datesRow := w.habitLabelStyle.Render("")
	today := time.Now()
	for i := 0; i < daysToShow; i++ {
		date := w.calendar.viewStartDate.AddDate(0, 0, i)
		dateStr := date.Format("02")

		isToday := date.Year() == today.Year() &&
			date.Month() == today.Month() &&
			date.Day() == today.Day()

		if isToday {
			datesRow += w.todayStyle.Render(dateStr)
		} else {
			datesRow += w.dateHeaderStyle.Render(dateStr)
		}
	}
	sb.WriteString(datesRow + "\n")

	// Today indicator
	todayIndicatorRow := w.habitLabelStyle.Render("")
	for i := 0; i < daysToShow; i++ {
		date := w.calendar.viewStartDate.AddDate(0, 0, i)
		isToday := date.Year() == today.Year() &&
			date.Month() == today.Month() &&
			date.Day() == today.Day()

		if isToday {
			todayIndicatorRow += w.cellStyle.Render("▼")
		} else {
			todayIndicatorRow += w.cellStyle.Render("")
		}
	}
	sb.WriteString(todayIndicatorRow)

	return sb.String()
}

// RenderContent renders the scrollable content (habit rows only)
func (w *WeeklyView) RenderContent() string {
	var sb strings.Builder

	// Calculate how many days can fit
	availableWidth := w.calendar.width - 20
	daysToShow := availableWidth / 8
	if daysToShow < 7 {
		daysToShow = 7
	}

	// Habit rows
	for idx, habit := range w.calendar.habits {
		var habitLabel string
		if idx == w.calendar.selectedHabit {
			habitLabel = w.selectedHabitLabelStyle.Render(habit.Name)
		} else {
			habitLabel = w.habitLabelStyle.Render(habit.Name)
		}
		row := habitLabel

		for i := 0; i < daysToShow; i++ {
			date := w.calendar.viewStartDate.AddDate(0, 0, i)
			cellValue := w.calendar.GetCellValue(habit, date)

			isSelected := idx == w.calendar.selectedHabit &&
				date.Year() == w.calendar.selectedDate.Year() &&
				date.Month() == w.calendar.selectedDate.Month() &&
				date.Day() == w.calendar.selectedDate.Day()

			if isSelected {
				row += w.selectedCellStyle.Render(cellValue)
			} else {
				row += w.cellStyle.Render(cellValue)
			}
		}
		sb.WriteString(row + "\n")
	}

	return sb.String()
}

// Render renders the complete weekly view (for backward compatibility)
func (w *WeeklyView) Render() string {
	return w.RenderHeader() + "\n" + w.RenderContent()
}
