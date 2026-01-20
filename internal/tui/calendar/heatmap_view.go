package calendar

import (
	"fmt"
	"strings"
	"time"

	"example.com/habits/internal/app"
)

// HeatmapView handles rendering of the heatmap calendar view
type HeatmapView struct {
	calendar *Calendar
	styles   *app.HeatmapStyles
}

// NewHeatmapView creates a new heatmap view renderer
func NewHeatmapView(calendar *Calendar, styles *app.HeatmapStyles) *HeatmapView {
	return &HeatmapView{
		calendar: calendar,
		styles:   styles,
	}
}

// RenderHeader renders the header for heatmap view
func (h *HeatmapView) RenderHeader() string {
	now := time.Now()
	year, month, _ := now.Date()

	header := fmt.Sprintf("Heatmap - %s %d", month.String(), year)
	return h.styles.Header.Render(header)
}

// RenderContent renders the heatmap content
func (h *HeatmapView) RenderContent() string {
	if len(h.calendar.habits) == 0 {
		return "No habits to display"
	}

	var sb strings.Builder

	// Get current month
	now := time.Now()
	year, month, _ := now.Date()
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	// Calculate days in month
	daysInMonth := lastOfMonth.Day()

	// For each habit, render a row
	for idx, habit := range h.calendar.habits {
		var habitLabel string
		if idx == h.calendar.selectedHabit {
			habitLabel = h.styles.SelectedHabit.Render(habit.GetDisplayName())
		} else {
			habitLabel = h.styles.HabitName.Render(habit.GetDisplayName())
		}

		// Build the heatmap row
		var cells []string
		for day := 1; day <= daysInMonth; day++ {
			date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
			cell := h.getHeatmapCell(habit, date, idx == h.calendar.selectedHabit && date.Equal(h.calendar.selectedDate))
			cells = append(cells, cell)
		}

		row := habitLabel + strings.Join(cells, "")
		sb.WriteString(row + "\n")
	}

	// Add day labels at the bottom
	var dayLabels []string
	dayLabels = append(dayLabels, h.styles.HabitName.Render("")) // Empty space for habit name column

	for day := 1; day <= daysInMonth; day++ {
		dayStr := fmt.Sprintf("%2d", day)
		dayLabels = append(dayLabels, h.styles.DayLabel.Render(dayStr))
	}

	sb.WriteString(strings.Join(dayLabels, "") + "\n")

	return sb.String()
}

// getHeatmapCell returns the character for a habit on a specific date
func (h *HeatmapView) getHeatmapCell(habit Habit, date time.Time, isSelected bool) string {
	entry, exists := h.calendar.GetEntry(habit.Name, date)
	now := time.Now()

	var symbol string
	if date.After(now) {
		// Future date
		symbol = h.calendar.getUntrackedSymbol(ViewModeHeatmap)
	} else if !exists {
		// No entry - use different symbols based on habit type
		switch habit.Type {
		case HabitTypeBit:
			symbol = h.calendar.getMissedSymbol(ViewModeHeatmap)
		case HabitTypeCount, HabitTypeFloat:
			symbol = h.calendar.getUntrackedSymbol(ViewModeHeatmap)
		default:
			symbol = h.calendar.getMissedSymbol(ViewModeHeatmap)
		}
	} else {
		// Entry exists
		switch habit.Type {
		case HabitTypeBit:
			if entry.Completed {
				symbol = h.calendar.getCompletedSymbol(ViewModeHeatmap)
			} else {
				symbol = h.calendar.getMissedSymbol(ViewModeHeatmap)
			}
		case HabitTypeCount, HabitTypeFloat:
			if entry.Value != "" && entry.Value != "-" {
				// Show actual value, truncate to 1 char for heatmap
				if len(entry.Value) > 1 {
					symbol = entry.Value[:1]
				} else {
					symbol = entry.Value
				}
			} else {
				symbol = h.calendar.getMissedSymbol(ViewModeHeatmap)
			}
		default:
			symbol = "?"
		}
	}

	// Check if this cell has notes
	hasNote := h.calendar.HasEntryNote(habit.Name, date)

	// Apply styling based on symbol type
	var cell string
	if hasNote {
		// Notes override other styling
		cell = h.styles.Note.Render(symbol)
	} else {
		switch symbol {
		case h.calendar.getCompletedSymbol(ViewModeHeatmap):
			cell = h.styles.Completed.Render(symbol)
		case h.calendar.getMissedSymbol(ViewModeHeatmap):
			cell = h.styles.Missed.Render(symbol)
		case h.calendar.getUntrackedSymbol(ViewModeHeatmap):
			cell = h.styles.Future.Render(symbol)
		default:
			// For count/float values, apply completed styling
			if habit.Type == HabitTypeCount || habit.Type == HabitTypeFloat {
				cell = h.styles.Completed.Render(symbol)
			} else {
				cell = symbol
			}
		}
	}

	// Add selection indicator
	if isSelected {
		return "[" + cell + "]"
	}
	return " " + cell + " "
}
