package calendar

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// ViewMode represents different calendar view modes
type ViewMode int

const (
	ViewModeWeekly ViewMode = iota
	ViewModeMonthly
	// Future: ViewModeHeatmap, ViewModeYearly, etc.
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

// Calendar is the main calendar component that manages habits and their entries
type Calendar struct {
	// Data
	habits  []Habit
	entries map[string]map[time.Time]HabitEntry // habitName -> date -> entry

	// State
	viewMode      ViewMode
	selectedDate  time.Time
	selectedHabit int       // index of selected habit
	viewStartDate time.Time // For weekly view - first day of visible week
	viewMonth     time.Time // For monthly view - first day of visible month

	// UI
	width  int
	height int

	// View renderers
	weeklyView  *WeeklyView
	monthlyView *MonthlyView
}

// NewCalendar creates a new calendar component
func NewCalendar(habits []Habit) *Calendar {
	now := time.Now()

	// Calculate week start (Monday)
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday becomes 7
	}
	weekStart := now.AddDate(0, 0, -(weekday - 1))

	// Month start
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	cal := &Calendar{
		habits:        habits,
		entries:       make(map[string]map[time.Time]HabitEntry),
		viewMode:      ViewModeWeekly,
		selectedDate:  now,
		selectedHabit: 0,
		viewStartDate: weekStart,
		viewMonth:     monthStart,
		width:         80,
		height:        24,
	}

	// Initialize view renderers
	cal.weeklyView = NewWeeklyView(cal)
	cal.monthlyView = NewMonthlyView(cal)

	return cal
}

// Init initializes the component
func (c *Calendar) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (c *Calendar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// Cycle through view modes
			c.viewMode = (c.viewMode + 1) % 2 // Currently just weekly and monthly

		case "h", "left":
			c.handleLeftNavigation()

		case "l", "right":
			c.handleRightNavigation()

		case "H":
			c.handleLeftJump()

		case "L":
			c.handleRightJump()

		case "j", "down":
			// Move to next habit
			if c.selectedHabit < len(c.habits)-1 {
				c.selectedHabit++
			}

		case "k", "up":
			// Move to previous habit
			if c.selectedHabit > 0 {
				c.selectedHabit--
			}

		case "t":
			// Jump to today
			c.jumpToToday()
		}

	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
	}

	return c, nil
}

// handleLeftNavigation handles left navigation based on current view mode
func (c *Calendar) handleLeftNavigation() {
	switch c.viewMode {
	case ViewModeWeekly:
		// Move to previous day
		c.selectedDate = c.selectedDate.AddDate(0, 0, -1)
		c.adjustWeeklyViewToSelection()

	case ViewModeMonthly:
		// Move to previous month
		c.viewMonth = c.viewMonth.AddDate(0, -1, 0)
	}
}

// handleRightNavigation handles right navigation based on current view mode
func (c *Calendar) handleRightNavigation() {
	switch c.viewMode {
	case ViewModeWeekly:
		// Move to next day
		c.selectedDate = c.selectedDate.AddDate(0, 0, 1)
		c.adjustWeeklyViewToSelection()

	case ViewModeMonthly:
		// Move to next month
		c.viewMonth = c.viewMonth.AddDate(0, 1, 0)
	}
}

// handleLeftJump handles left jump based on current view mode
func (c *Calendar) handleLeftJump() {
	switch c.viewMode {
	case ViewModeWeekly:
		// Move to previous week
		c.selectedDate = c.selectedDate.AddDate(0, 0, -7)
		c.viewStartDate = c.viewStartDate.AddDate(0, 0, -7)

	case ViewModeMonthly:
		// Move to previous year
		c.viewMonth = c.viewMonth.AddDate(-1, 0, 0)
	}
}

// handleRightJump handles right jump based on current view mode
func (c *Calendar) handleRightJump() {
	switch c.viewMode {
	case ViewModeWeekly:
		// Move to next week
		c.selectedDate = c.selectedDate.AddDate(0, 0, 7)
		c.viewStartDate = c.viewStartDate.AddDate(0, 0, 7)

	case ViewModeMonthly:
		// Move to next year
		c.viewMonth = c.viewMonth.AddDate(1, 0, 0)
	}
}

// jumpToToday jumps to today's date and adjusts view accordingly
func (c *Calendar) jumpToToday() {
	now := time.Now()
	c.selectedDate = now

	// Adjust weekly view
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	c.viewStartDate = now.AddDate(0, 0, -(weekday - 1))

	// Adjust monthly view
	c.viewMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// adjustWeeklyViewToSelection adjusts the weekly view to show the selected date
func (c *Calendar) adjustWeeklyViewToSelection() {
	// If selected date is before view start, shift view back
	if c.selectedDate.Before(c.viewStartDate) {
		weekday := int(c.selectedDate.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		c.viewStartDate = c.selectedDate.AddDate(0, 0, -(weekday - 1))
	}

	// If selected date is after view end, shift view forward
	viewEndDate := c.viewStartDate.AddDate(0, 0, 6)
	if c.selectedDate.After(viewEndDate) {
		weekday := int(c.selectedDate.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		c.viewStartDate = c.selectedDate.AddDate(0, 0, -(weekday - 1))
	}
}

// View renders the component based on current view mode
func (c *Calendar) View() string {
	switch c.viewMode {
	case ViewModeWeekly:
		return c.weeklyView.Render()
	case ViewModeMonthly:
		return c.monthlyView.Render()
	default:
		return "Unknown view mode"
	}
}

// SetEntry sets an entry for a habit on a specific date
func (c *Calendar) SetEntry(habitName string, date time.Time, completed bool, value string) {
	if c.entries[habitName] == nil {
		c.entries[habitName] = make(map[time.Time]HabitEntry)
	}

	// Normalize date to remove time component
	dateKey := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	c.entries[habitName][dateKey] = HabitEntry{
		Date:      dateKey,
		Completed: completed,
		Value:     value,
	}
}

// GetEntry retrieves an entry for a habit on a specific date
func (c *Calendar) GetEntry(habitName string, date time.Time) (HabitEntry, bool) {
	habitEntries, exists := c.entries[habitName]
	if !exists {
		return HabitEntry{}, false
	}

	// Normalize date
	dateKey := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	entry, exists := habitEntries[dateKey]
	return entry, exists
}

// GetCellValue returns the display value for a habit on a specific date
func (c *Calendar) GetCellValue(habit Habit, date time.Time) string {
	entry, exists := c.GetEntry(habit.Name, date)
	if !exists {
		return c.getDefaultValue(habit, date)
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
func (c *Calendar) getDefaultValue(habit Habit, date time.Time) string {
	today := time.Now()
	if date.After(today) {
		return "?"
	}
	return "-"
}

// Resize updates the component dimensions
func (c *Calendar) Resize(width, height int) {
	c.width = width
	c.height = height
}
