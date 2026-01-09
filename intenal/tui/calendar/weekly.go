package calendar

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HabitType represents the type of habit
type HabitType int

const (
	HabitTypeBit HabitType = iota
	HabitTypeCount
	HabitTypeFloat
)

// Habit represents a single habit
type Habit struct {
	Name string
	Type HabitType
}

// HabitEntry represents a habit entry for a specific date
type HabitEntry struct {
	Date      time.Time
	Completed bool   // for bit type
	Value     string // for count/float type, "-" for skipped
}

// WeeklyCalendar is a bubbletea model for the weekly calendar view
type WeeklyCalendar struct {
	habits         []Habit
	entries        map[string]map[time.Time]HabitEntry // habitName -> date -> entry
	selectedDate   time.Time
	selectedHabit  int // index of selected habit
	viewStartDate  time.Time // First day of the visible week
	width          int
	height         int

	// Styles
	headerStyle       lipgloss.Style
	dateHeaderStyle   lipgloss.Style
	dayNameStyle      lipgloss.Style
	cellStyle         lipgloss.Style
	selectedCellStyle lipgloss.Style
	todayStyle        lipgloss.Style
	habitLabelStyle   lipgloss.Style
	selectedHabitLabelStyle lipgloss.Style
}

// NewWeeklyCalendar creates a new weekly calendar component
func NewWeeklyCalendar(habits []Habit) *WeeklyCalendar {
	now := time.Now()
	// Start the week on Monday (ISO standard)
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday becomes 7
	}
	weekStart := now.AddDate(0, 0, -(weekday - 1))

	return &WeeklyCalendar{
		habits:         habits,
		entries:        make(map[string]map[time.Time]HabitEntry),
		selectedDate:   now,
		selectedHabit:  0,
		viewStartDate:  weekStart,
		width:          80,
		height:         20,

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

// Init initializes the component
func (w *WeeklyCalendar) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (w *WeeklyCalendar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			// Move to previous day
			w.selectedDate = w.selectedDate.AddDate(0, 0, -1)
			w.adjustViewToSelection()
		case "l", "right":
			// Move to next day
			w.selectedDate = w.selectedDate.AddDate(0, 0, 1)
			w.adjustViewToSelection()
		case "H":
			// Move to previous week
			w.selectedDate = w.selectedDate.AddDate(0, 0, -7)
			w.viewStartDate = w.viewStartDate.AddDate(0, 0, -7)
		case "L":
			// Move to next week
			w.selectedDate = w.selectedDate.AddDate(0, 0, 7)
			w.viewStartDate = w.viewStartDate.AddDate(0, 0, 7)
		case "j", "down":
			// Move to next habit
			if w.selectedHabit < len(w.habits)-1 {
				w.selectedHabit++
			}
		case "k", "up":
			// Move to previous habit
			if w.selectedHabit > 0 {
				w.selectedHabit--
			}
		case "t":
			// Jump to today
			now := time.Now()
			w.selectedDate = now
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			w.viewStartDate = now.AddDate(0, 0, -(weekday - 1))
		}
	case tea.WindowSizeMsg:
		w.width = msg.Width
		w.height = msg.Height
	}
	return w, nil
}

// adjustViewToSelection adjusts the view to show the selected date
func (w *WeeklyCalendar) adjustViewToSelection() {
	// If selected date is before view start, shift view back
	if w.selectedDate.Before(w.viewStartDate) {
		weekday := int(w.selectedDate.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		w.viewStartDate = w.selectedDate.AddDate(0, 0, -(weekday - 1))
	}

	// If selected date is after view end, shift view forward
	viewEndDate := w.viewStartDate.AddDate(0, 0, 6)
	if w.selectedDate.After(viewEndDate) {
		weekday := int(w.selectedDate.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		w.viewStartDate = w.selectedDate.AddDate(0, 0, -(weekday - 1))
	}
}

// View renders the component
func (w *WeeklyCalendar) View() string {
	var sb strings.Builder

	// Calculate week number
	_, week := w.viewStartDate.ISOWeek()
	endDate := w.viewStartDate.AddDate(0, 0, 6)

	// Header with week range and week number
	header := fmt.Sprintf("Week: %s - %s, %d",
		w.viewStartDate.Format("Jan 02"),
		endDate.Format("Jan 02"),
		w.viewStartDate.Year(),
	)
	weekIndicator := fmt.Sprintf("◀ [%d/52] ▶", week)
	
	// Calculate spacing to fill the width
	headerLen := len(header)
	indicatorLen := len(weekIndicator)
	spacingLen := w.width - headerLen - indicatorLen - 4
	if spacingLen < 0 {
		spacingLen = 0
	}
	
	headerLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		w.headerStyle.Render(header),
		strings.Repeat(" ", spacingLen),
		w.headerStyle.Render(weekIndicator),
	)
	sb.WriteString(headerLine + "\n\n")

	// Calculate how many days can fit in the available width
	// Each day takes 8 characters, habit label takes 20 characters
	availableWidth := w.width - 20
	daysToShow := availableWidth / 8
	if daysToShow < 7 {
		daysToShow = 7 // minimum one week
	}

	// Day names row
	dayNames := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	dayNamesRow := w.habitLabelStyle.Render("")
	for i := 0; i < daysToShow; i++ {
		dayName := dayNames[i%7]
		dayNamesRow += w.dayNameStyle.Render(dayName)
	}
	sb.WriteString(dayNamesRow + "\n")

	// Dates row
	datesRow := w.habitLabelStyle.Render("")
	today := time.Now()
	for i := 0; i < daysToShow; i++ {
		date := w.viewStartDate.AddDate(0, 0, i)
		dateStr := date.Format("02")
		
		// Check if this date is today
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
		date := w.viewStartDate.AddDate(0, 0, i)
		isToday := date.Year() == today.Year() && 
			date.Month() == today.Month() && 
			date.Day() == today.Day()
		
		if isToday {
			todayIndicatorRow += w.cellStyle.Render("▼ Today")
		} else {
			todayIndicatorRow += w.cellStyle.Render("")
		}
	}
	sb.WriteString(todayIndicatorRow + "\n\n")

	// Habit rows - fill the remaining vertical space
	// Calculate how many rows we can show
	usedLines := 7 // header + spacing + day names + dates + today indicator + spacing
	// availableLines := w.height - usedLines - 2 // -2 for help text
	
	// Create habit rows
	habitRows := make([]string, len(w.habits))
	for idx, habit := range w.habits {
		var habitLabel string
		if idx == w.selectedHabit {
			habitLabel = w.selectedHabitLabelStyle.Render(habit.Name)
		} else {
			habitLabel = w.habitLabelStyle.Render(habit.Name)
		}
		row := habitLabel

		for i := 0; i < daysToShow; i++ {
			date := w.viewStartDate.AddDate(0, 0, i)
			cellValue := w.getCellValue(habit, date)
			
			// Check if this cell is selected
			isSelected := idx == w.selectedHabit && 
				date.Year() == w.selectedDate.Year() && 
				date.Month() == w.selectedDate.Month() && 
				date.Day() == w.selectedDate.Day()
			
			if isSelected {
				row += w.selectedCellStyle.Render(cellValue)
			} else {
				row += w.cellStyle.Render(cellValue)
			}
		}
		habitRows[idx] = row
	}

	// Write habit rows
	for _, row := range habitRows {
		sb.WriteString(row + "\n")
	}
	
	// Fill remaining vertical space with empty lines if needed
	linesUsed := usedLines + len(w.habits)
	for i := linesUsed; i < w.height-2; i++ {
		sb.WriteString("\n")
	}

	return sb.String()
}

// getCellValue returns the display value for a habit on a specific date
func (w *WeeklyCalendar) getCellValue(habit Habit, date time.Time) string {
	// Check if we have an entry for this habit on this date
	habitEntries, exists := w.entries[habit.Name]
	if !exists {
		return w.getDefaultValue(habit, date)
	}

	// Normalize date to remove time component
	dateKey := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	entry, exists := habitEntries[dateKey]
	if !exists {
		return w.getDefaultValue(habit, date)
	}

	switch habit.Type {
	case HabitTypeBit:
		if entry.Completed {
			return "✓"
		}
		return "-"
	case HabitTypeCount, HabitTypeFloat:
		if entry.Value == "-" {
			return "-"
		}
		return entry.Value
	}

	return "?"
}

// getDefaultValue returns the default display value based on date
func (w *WeeklyCalendar) getDefaultValue(habit Habit, date time.Time) string {
	today := time.Now()
	if date.After(today) {
		return "?"
	}
	return "-"
}

// SetEntry sets an entry for a habit on a specific date
func (w *WeeklyCalendar) SetEntry(habitName string, date time.Time, completed bool, value string) {
	if w.entries[habitName] == nil {
		w.entries[habitName] = make(map[time.Time]HabitEntry)
	}
	
	// Normalize date to remove time component
	dateKey := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	w.entries[habitName][dateKey] = HabitEntry{
		Date:      dateKey,
		Completed: completed,
		Value:     value,
	}
}

// Resize updates the component dimensions
func (w *WeeklyCalendar) Resize(width, height int) {
	w.width = width
	w.height = height
}
