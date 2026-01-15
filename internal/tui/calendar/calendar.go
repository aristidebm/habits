package calendar

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
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

func (ht HabitType) String() string {
	switch ht {
	case HabitTypeBit:
		return "bit"
	case HabitTypeCount:
		return "count"
	case HabitTypeFloat:
		return "float"
	default:
		return ""
	}
}

// Habit represents a single habit
type Habit struct {
	ID   string
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
	habits        []Habit
	entries       map[string]map[time.Time]HabitEntry // habitName -> date -> entry
	pendingHabits []Habit                             // Habits added but not yet written to database

	// State
	viewMode      ViewMode
	selectedDate  time.Time
	selectedHabit int       // index of selected habit
	viewStartDate time.Time // For weekly view - first day of visible week
	viewMonth     time.Time // For monthly view - first day of visible month

	// UI
	width    int
	height   int
	viewport viewport.Model
	ready    bool

	// View renderers
	weeklyView  *WeeklyView
	monthlyView *MonthlyView
}

// NewCalendar creates a new calendar component
func NewCalendar(habits []Habit) *Calendar {
	now := time.Now()

	// Month start
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	cal := &Calendar{
		habits:        habits,
		entries:       make(map[string]map[time.Time]HabitEntry),
		viewMode:      ViewModeWeekly,
		selectedDate:  now,
		selectedHabit: 0,
		viewMonth:     monthStart,
		width:         80,
		height:        24,
	}

	// Calculate view start to center today horizontally (after width is set)
	cal.centerDateInView(now)

	// Initialize view renderers
	cal.weeklyView = NewWeeklyView(cal)
	cal.monthlyView = NewMonthlyView(cal)

	// Initialize viewport
	cal.viewport = viewport.New(80, 20)

	return cal
}

// Init initializes the component
func (c *Calendar) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (c *Calendar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// Cycle through view modes
			c.viewMode = (c.viewMode + 1) % 2 // Currently just weekly and monthly
			c.updateViewportContent()

		case "h", "left":
			c.handleLeftNavigation()
			c.updateViewportContent()

		case "l", "right":
			c.handleRightNavigation()
			c.updateViewportContent()

		case "H":
			c.handleLeftJump()
			c.updateViewportContent()

		case "L":
			c.handleRightJump()
			c.updateViewportContent()

		case "j", "down":
			// Move to next habit (both weekly and monthly)
			if c.selectedHabit < len(c.habits)-1 {
				c.selectedHabit++
				c.scrollToSelectedHabit()
				c.updateViewportContent()
			}

		case "k", "up":
			// Move to previous habit (both weekly and monthly)
			if c.selectedHabit > 0 {
				c.selectedHabit--
				c.scrollToSelectedHabit()
				c.updateViewportContent()
			}

		case "n":
			// Next habit
			if c.selectedHabit < len(c.habits)-1 {
				c.selectedHabit++
				c.scrollToSelectedHabit()
				c.updateViewportContent()
			}

		case "p":
			// Previous habit
			if c.selectedHabit > 0 {
				c.selectedHabit--
				c.scrollToSelectedHabit()
				c.updateViewportContent()
			}

		case "t":
			// Jump to today
			c.jumpToToday()
			c.updateViewportContent()

		default:
			// Pass other keys to viewport for scrolling
			c.viewport, cmd = c.viewport.Update(msg)
			return c, cmd
		}

	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		c.updateViewportSize()
	}

	return c, nil
}

// updateViewportSize updates viewport dimensions based on calendar size
func (c *Calendar) updateViewportSize() {
	// Reserve space for header (varies by view mode)
	headerHeight := 6 // Default for weekly
	if c.viewMode == ViewModeMonthly {
		headerHeight = 3 // Monthly view has less header
	}

	viewportHeight := c.height - headerHeight
	if viewportHeight < 1 {
		viewportHeight = 1
	}

	c.viewport.Width = c.width
	c.viewport.Height = viewportHeight
	c.ready = true

	// Re-center today in weekly view when terminal resizes
	if c.viewMode == ViewModeWeekly {
		c.centerDateInView(time.Now())
	}

	c.updateViewportContent()
}

// updateViewportContent updates the viewport with current view content
func (c *Calendar) updateViewportContent() {
	if !c.ready {
		return
	}

	var content string
	if c.viewMode == ViewModeWeekly {
		content = c.weeklyView.RenderContent()
	} else {
		content = c.monthlyView.RenderContent()
	}

	c.viewport.SetContent(content)
}

// scrollToSelectedHabit scrolls viewport to keep selected habit visible
func (c *Calendar) scrollToSelectedHabit() {
	if !c.ready {
		return
	}

	var selectedY, itemHeight int

	if c.viewMode == ViewModeWeekly {
		// In weekly view, each habit is one line
		itemHeight = 1
		selectedY = c.selectedHabit * itemHeight
	} else {
		// In monthly view, calculate based on card layout
		cardWidth := 7*4 + 4
		cardsPerRow := c.width / cardWidth
		if cardsPerRow < 1 {
			cardsPerRow = 1
		}

		// Each card is ~10 lines tall (name + blank + 6 weeks + blank line between rows)
		cardHeight := 10
		rowIndex := c.selectedHabit / cardsPerRow
		selectedY = rowIndex * cardHeight
		itemHeight = cardHeight
	}

	// Calculate the bottom of the selected item
	selectedBottom := selectedY + itemHeight

	// If selected item's bottom is below viewport, scroll down to show full item
	if selectedBottom > c.viewport.YOffset+c.viewport.Height {
		c.viewport.YOffset = selectedBottom - c.viewport.Height
		if c.viewport.YOffset < 0 {
			c.viewport.YOffset = 0
		}
	}

	// If selected item's top is above viewport, scroll up
	if selectedY < c.viewport.YOffset {
		c.viewport.YOffset = selectedY
	}
}

// handleLeftNavigation handles left navigation based on current view mode
func (c *Calendar) handleLeftNavigation() {
	switch c.viewMode {
	case ViewModeWeekly:
		// Move to previous day
		c.selectedDate = c.selectedDate.AddDate(0, 0, -1)
		c.adjustWeeklyViewToSelection()

	case ViewModeMonthly:
		// Move to previous day
		c.selectedDate = c.selectedDate.AddDate(0, 0, -1)
		// Adjust month view if we moved to previous month
		if c.selectedDate.Month() != c.viewMonth.Month() || c.selectedDate.Year() != c.viewMonth.Year() {
			c.viewMonth = time.Date(c.selectedDate.Year(), c.selectedDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		}
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
		// Move to next day
		c.selectedDate = c.selectedDate.AddDate(0, 0, 1)
		// Adjust month view if we moved to next month
		if c.selectedDate.Month() != c.viewMonth.Month() || c.selectedDate.Year() != c.viewMonth.Year() {
			c.viewMonth = time.Date(c.selectedDate.Year(), c.selectedDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		}
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
		// Move to previous month
		c.viewMonth = c.viewMonth.AddDate(0, -1, 0)
		// Keep selected date in the same day of month if possible
		c.selectedDate = c.selectedDate.AddDate(0, -1, 0)
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
		// Move to next month
		c.viewMonth = c.viewMonth.AddDate(0, 1, 0)
		// Keep selected date in the same day of month if possible
		c.selectedDate = c.selectedDate.AddDate(0, 1, 0)
	}
}

// jumpToToday jumps to today's date and adjusts view accordingly
func (c *Calendar) jumpToToday() {
	now := time.Now()
	c.selectedDate = now

	// Adjust weekly view - center today
	c.centerDateInView(now)

	// Adjust monthly view
	c.viewMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// adjustWeeklyViewToSelection adjusts the weekly view to show the selected date
func (c *Calendar) adjustWeeklyViewToSelection() {
	availableWidth := c.width - 20
	daysToShow := availableWidth / 8
	if daysToShow < 7 {
		daysToShow = 7
	}

	// Calculate view end
	viewEndDate := c.viewStartDate.AddDate(0, 0, daysToShow-1)

	// If selected date is before view start or after view end, center it
	if c.selectedDate.Before(c.viewStartDate) || c.selectedDate.After(viewEndDate) {
		c.centerDateInView(c.selectedDate)
	}
}

// centerDateInView centers the given date in the weekly view based on current width
func (c *Calendar) centerDateInView(date time.Time) {
	availableWidth := c.width - 20
	daysToShow := availableWidth / 8
	if daysToShow < 7 {
		daysToShow = 7
	}
	offset := daysToShow / 2
	c.viewStartDate = date.AddDate(0, 0, -offset)
}

// View renders the component based on current view mode
func (c *Calendar) View() string {
	if !c.ready {
		return "Initializing..."
	}

	var header string
	if c.viewMode == ViewModeWeekly {
		header = c.weeklyView.RenderHeader()
	} else {
		header = c.monthlyView.RenderHeader()
	}

	return header + "\n" + c.viewport.View()
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
	c.updateViewportSize()
}

// GetSelectedDate returns the currently selected date
func (c *Calendar) GetSelectedDate() time.Time {
	return c.selectedDate
}

// GetSelectedHabit returns the currently selected habit
func (c *Calendar) GetSelectedHabit() *Habit {
	if c.selectedHabit >= 0 && c.selectedHabit < len(c.habits) {
		return &c.habits[c.selectedHabit]
	}
	return nil
}

// ReloadHabits replaces the habits list
func (c *Calendar) ReloadHabits(habits []Habit) {
	c.habits = habits
	c.selectedHabit = 0
}

// RemovePendingHabit removes a habit from pending list
func (c *Calendar) RemovePendingHabit(name string) {
	for i, h := range c.pendingHabits {
		if h.Name == name {
			if i == 0 {
				c.pendingHabits = c.pendingHabits[1:]
			} else {
				c.pendingHabits = append(c.pendingHabits[:i], c.pendingHabits[i+1:]...)
			}
			break
		}
	}
}

// ClearPendingHabits clears all pending habits
func (c *Calendar) ClearPendingHabits() {
	c.pendingHabits = []Habit{}
}

// GetPendingHabits returns all pending habits
func (c *Calendar) GetPendingHabits() []Habit {
	return c.pendingHabits
}

// AddPendingHabit adds a habit to the pending list
func (c *Calendar) AddPendingHabit(habit Habit) {
	c.pendingHabits = append(c.pendingHabits, habit)
}

// WritePendingHabits writes all pending habits to database via Store
func (c *Calendar) WritePendingHabits(createHabit func(name string, habitType string, goal float64) (int, error)) error {
	var errs []error
	for _, h := range c.pendingHabits {
		id, err := createHabit(h.Name, h.Type.String(), 0)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to write habit '%s': %w", h.Name, err))
		} else {
			_ = id
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("failed to write %d habits: %v", len(errs), errs)
	}
	c.ClearPendingHabits()
	return nil
}
