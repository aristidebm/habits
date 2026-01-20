package calendar

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"example.com/habits/internal/app"
)

// WeeklyView handles rendering of the weekly calendar view
type WeeklyView struct {
	calendar *Calendar
	styles   *app.WeeklyStyles
}

// NewWeeklyView creates a new weekly view renderer
func NewWeeklyView(calendar *Calendar, styles *app.WeeklyStyles) *WeeklyView {
	return &WeeklyView{
		calendar: calendar,
		styles:   styles,
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
		w.styles.Header.Render(header),
		strings.Repeat(" ", spacingLen),
		w.styles.Header.Render(dayIndicator),
		strings.Repeat(" ", spacingLen),
		w.styles.Header.Render(weekIndicator),
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
	dayNames := allDayNames[startDay-1:]
	dayNames = append(dayNames, allDayNames[:startDay-1]...)
	dayNamesRow := w.styles.HabitLabel.Render("")
	for _, dayName := range dayNames {
		dayNamesRow += w.styles.DayName.Render(dayName)
	}
	datesRow := w.styles.HabitLabel.Render("")
	today := time.Now()
	for i := 0; i < daysToShow; i++ {
		date := w.calendar.viewStartDate.AddDate(0, 0, i)
		dateStr := date.Format("02")

		isToday := date.Year() == today.Year() &&
			date.Month() == today.Month() &&
			date.Day() == today.Day()

		if isToday {
			datesRow += w.styles.TodayCell.Render(dateStr)
		} else {
			datesRow += w.styles.DateHeader.Render(dateStr)
		}
	}
	todayIndicatorRow := w.styles.HabitLabel.Render("")
	for i := 0; i < daysToShow; i++ {
		date := w.calendar.viewStartDate.AddDate(0, 0, i)
		isToday := date.Year() == today.Year() &&
			date.Month() == today.Month() &&
			date.Day() == today.Day()
		if isToday {
			todayIndicatorRow += w.styles.Cell.Render("▼")
		} else {
			todayIndicatorRow += w.styles.Cell.Render("")
		}
	}
	sb.WriteString(datesRow + "\n")
	for i := 0; i < daysToShow; i++ {
		date := w.calendar.viewStartDate.AddDate(0, 0, i)
		isToday := date.Year() == today.Year() &&
			date.Month() == today.Month() &&
			date.Day() == today.Day()

		if isToday {
			todayIndicatorRow += w.styles.Cell.Render("▼")
		} else {
			todayIndicatorRow += w.styles.Cell.Render("")
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
			habitLabel = w.styles.SelectedHabitLabel.Render(habit.GetDisplayName())
		} else {
			habitLabel = w.styles.HabitLabel.Render(habit.GetDisplayName())
		}
		row := habitLabel

		for i := 0; i < daysToShow; i++ {
			date := w.calendar.viewStartDate.AddDate(0, 0, i)
			cellValue := w.calendar.GetCellValue(habit, date, ViewModeWeekly)

			isSelected := idx == w.calendar.selectedHabit &&
				date.Year() == w.calendar.selectedDate.Year() &&
				date.Month() == w.calendar.selectedDate.Month() &&
				date.Day() == w.calendar.selectedDate.Day()

			hasNote := w.calendar.HasEntryNote(habit.Name, date)

			if isSelected {
				row += w.styles.SelectedCell.Render(cellValue)
			} else if hasNote {
				row += w.styles.NoteCell.Render(cellValue)
			} else {
				row += w.styles.Cell.Render(cellValue)
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
