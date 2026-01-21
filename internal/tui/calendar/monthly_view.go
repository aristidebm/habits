package calendar

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"example.com/habits/internal/app"
)

// MonthlyView handles rendering of the monthly calendar view
type MonthlyView struct {
	calendar *Calendar
	styles   *app.MonthlyStyles
}

// NewMonthlyView creates a new monthly view renderer
func NewMonthlyView(calendar *Calendar, styles *app.MonthlyStyles) *MonthlyView {
	return &MonthlyView{
		calendar: calendar,
		styles:   styles,
	}
}

// RenderHeader renders fixed header (month/year)
func (m *MonthlyView) RenderHeader() string {
	header := fmt.Sprintf("%s %d",
		m.calendar.viewMonth.Format("January"),
		m.calendar.viewMonth.Year(),
	)
	return m.styles.Header.Render(header)
}

// RenderContent renders scrollable content (habit cards)
func (m *MonthlyView) RenderContent() string {
	if len(m.calendar.habits) == 0 {
		return "No habits to display"
	}

	var sb strings.Builder

	// Calculate card width (7 days * 4 chars per cell + padding)
	cardWidth := 7*4 + 4
	cardsPerRow := m.calendar.width / cardWidth
	if cardsPerRow < 1 {
		cardsPerRow = 1
	}

	// Create habit cards
	habitCards := make([]string, len(m.calendar.habits))
	for i, habit := range m.calendar.habits {
		habitCards[i] = m.renderHabitCard(habit, i == m.calendar.selectedHabit)
	}

	// Arrange cards in rows
	for i := 0; i < len(habitCards); i += cardsPerRow {
		end := i + cardsPerRow
		if end > len(habitCards) {
			end = len(habitCards)
		}

		rowCards := habitCards[i:end]
		// Join cards horizontally
		row := lipgloss.JoinHorizontal(lipgloss.Top, rowCards...)
		sb.WriteString(row + "\n\n")
	}

	return sb.String()
}

// Render renders complete monthly view (for backward compatibility)
func (m *MonthlyView) Render() string {
	return m.RenderHeader() + "\n\n" + m.RenderContent()
}

// renderHabitCard renders a single habit as a card with its monthly calendar
func (m *MonthlyView) renderHabitCard(habit Habit, isSelected bool) string {
	var card strings.Builder

	// Habit name
	var habitName string
	if isSelected {
		habitName = m.styles.SelectedHabit.Render(habit.GetDisplayName())
	} else {
		habitName = m.styles.HabitName.Render(habit.GetDisplayName())
	}
	card.WriteString(habitName + "\n\n")

	// Get first day of month and calculate starting position
	firstDay := m.calendar.viewMonth
	weekday := int(firstDay.Weekday()) // 0 = Sunday

	// Get last day of month
	nextMonth := m.calendar.viewMonth.AddDate(0, 1, 0)
	lastDay := nextMonth.AddDate(0, 0, -1)
	daysInMonth := lastDay.Day()

	currentDay := 1

	// Generate calendar grid (up to 6 rows for weeks)
	for week := 0; week < 6; week++ {
		var weekRow string

		for day := 0; day < 7; day++ {
			if (week == 0 && day < weekday) || currentDay > daysInMonth {
				// Empty cell
				weekRow += m.styles.Empty.Render(m.calendar.getUntrackedSymbol(ViewModeMonthly))
			} else {
				date := time.Date(m.calendar.viewMonth.Year(), m.calendar.viewMonth.Month(), currentDay, 0, 0, 0, 0, time.UTC)

				// Check if this is the selected date
				isSelectedDate := date.Year() == m.calendar.selectedDate.Year() &&
					date.Month() == m.calendar.selectedDate.Month() &&
					date.Day() == m.calendar.selectedDate.Day() &&
					isSelected

				cellValue := m.getCompactCellValue(habit, date)
				hasNote := m.calendar.HasEntryNote(habit.Name, date)

				if isSelectedDate {
					weekRow += m.styles.SelectedCell.Render(cellValue)
				} else if hasNote {
					weekRow += m.styles.NoteCell.Render(cellValue)
				} else {
					weekRow += m.styles.Cell.Render(cellValue)
				}

				currentDay++
			}
		}

		card.WriteString(weekRow + "\n")

		// Stop if we've shown all days
		if currentDay > daysInMonth {
			break
		}
	}

	return m.styles.Card.Render(card.String())
}

// getCompactCellValue returns a compact display value for a habit on a specific date
func (m *MonthlyView) getCompactCellValue(habit Habit, date time.Time) string {
	entry, exists := m.calendar.GetEntry(habit.Name, date)

	today := time.Now()
	completedSymbol := m.calendar.getCompletedSymbol(ViewModeMonthly)
	missedSymbol := m.calendar.getMissedSymbol(ViewModeMonthly)
	untrackedSymbol := m.calendar.getUntrackedSymbol(ViewModeMonthly)

	if date.After(today) {
		return untrackedSymbol
	}

	if !exists {
		// For past dates with no entry, use different symbols based on habit type
		if date.After(today) {
			return untrackedSymbol
		}
		switch habit.Type {
		case HabitTypeBit:
			return missedSymbol
		case HabitTypeCount, HabitTypeFloat:
			return untrackedSymbol
		default:
			return missedSymbol
		}
	}

	switch habit.Type {
	case HabitTypeBit:
		if entry.Completed {
			return completedSymbol
		}
		return missedSymbol
	case HabitTypeCount, HabitTypeFloat:
		if entry.Value < 0 {
			return missedSymbol
		}
		// Show actual value, truncate if too long
		valStr := fmt.Sprintf("%.0f", entry.Value)
		if len(valStr) > 3 {
			return valStr[:3]
		}
		return valStr
	}

	return missedSymbol
}

// getSelectedDateStats returns number of completed and remaining habits for selected date
func (m *MonthlyView) getSelectedDateStats() (completed, remaining int) {
	selectedDate := m.calendar.selectedDate

	for _, habit := range m.calendar.habits {
		entry, exists := m.calendar.GetEntry(habit.Name, selectedDate)

		if exists {
			switch habit.Type {
			case HabitTypeBit:
				if entry.Completed {
					completed++
				} else {
					remaining++
				}
			case HabitTypeCount, HabitTypeFloat:
				if entry.Value >= 0 {
					completed++
				} else {
					remaining++
				}
			}
		} else {
			remaining++
		}
	}

	return completed, remaining
}
