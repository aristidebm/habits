package calendar

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"example.com/habits/internal/app"
)

// ViewMode represents different calendar view modes
type ViewMode int

const (
	ViewModeWeekly ViewMode = iota
	ViewModeMonthly
	ViewModeHeatmap
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
	ID      int
	Name    string
	Type    HabitType
	Pending bool
}

// HabitEntry represents a habit entry for a specific date
type HabitEntry struct {
	Date      time.Time
	Completed bool   // for bit type
	Value     string // for count/float type, "-" for skipped
	Pending   bool
}

// Calendar is the main calendar component that manages habits and their entries
type Calendar struct {
	// Data
	habits  []Habit
	entries map[string]map[time.Time]HabitEntry // habitName -> date -> entry

	// Configuration
	config *app.Config

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
	heatmapView *HeatmapView
}

// NewCalendar creates a new calendar component
func NewCalendar(habits []Habit, config *app.Config) *Calendar {
	now := time.Now()

	// Month start
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	cal := &Calendar{
		habits:        habits,
		entries:       make(map[string]map[time.Time]HabitEntry),
		config:        config,
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
	cal.heatmapView = NewHeatmapView(cal)

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
			c.viewMode = (c.viewMode + 1) % 3 // Weekly, Monthly, Heatmap
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

		case "enter":
			// Modify entry at selected position
			habit := c.GetSelectedHabit()
			if habit != nil {
				c.toggleOrIncrementEntry(*habit, c.selectedDate)
				c.updateViewportContent()
			}

		case "backspace":
			// Decrement entry for count/float habits
			habit := c.GetSelectedHabit()
			if habit != nil && (habit.Type == HabitTypeCount || habit.Type == HabitTypeFloat) {
				c.decrementEntry(*habit, c.selectedDate)
				c.updateViewportContent()
			}

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
	if c.viewMode == ViewModeMonthly || c.viewMode == ViewModeHeatmap {
		headerHeight = 3 // Monthly and heatmap views have less header
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
	switch c.viewMode {
	case ViewModeWeekly:
		content = c.weeklyView.RenderContent()
	case ViewModeMonthly:
		content = c.monthlyView.RenderContent()
	case ViewModeHeatmap:
		content = c.heatmapView.RenderContent()
	}

	c.viewport.SetContent(content)
}

// scrollToSelectedHabit scrolls viewport to keep selected habit visible
func (c *Calendar) scrollToSelectedHabit() {
	if !c.ready {
		return
	}

	var selectedY, itemHeight int

	switch c.viewMode {
	case ViewModeWeekly, ViewModeHeatmap:
		// In weekly and heatmap views, each habit is one line
		itemHeight = 1
		selectedY = c.selectedHabit * itemHeight
	case ViewModeMonthly:
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
	switch c.viewMode {
	case ViewModeWeekly:
		header = c.weeklyView.RenderHeader()
	case ViewModeMonthly:
		header = c.monthlyView.RenderHeader()
	case ViewModeHeatmap:
		header = c.heatmapView.RenderHeader()
	}

	return header + "\n" + c.viewport.View()
}

// SetEntry sets an entry for a habit on a specific date
func (c *Calendar) SetEntry(habitName string, date time.Time, completed bool, value string, pending bool) {
	if c.entries[habitName] == nil {
		c.entries[habitName] = make(map[time.Time]HabitEntry)
	}

	// Normalize date to remove time component
	dateKey := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	c.entries[habitName][dateKey] = HabitEntry{
		Date:      dateKey,
		Completed: completed,
		Value:     value,
		Pending:   pending,
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

// GetCellValue returns the display value for a habit on a specific date for the given view mode
func (c *Calendar) GetCellValue(habit Habit, date time.Time, viewMode ViewMode) string {
	entry, exists := c.GetEntry(habit.Name, date)
	if !exists {
		return c.getDefaultValueForView(habit, date, viewMode)
	}

	slog.Info("", "Habit", habit.Name, "entry", entry.Value, "Date", entry.Date, "Completed", entry.Completed)

	completedSymbol := c.getCompletedSymbol(viewMode)
	missedSymbol := c.getMissedSymbol(viewMode)

	switch habit.Type {
	case HabitTypeBit:
		if entry.Completed {
			return completedSymbol
		}
		return missedSymbol
	case HabitTypeCount, HabitTypeFloat:
		if entry.Value == "-" || entry.Value == "" || entry.Value == "0" {
			return missedSymbol // No value = not completed
		}
		return completedSymbol // Has value = completed
	}

	return "?"
}

// getDefaultValueForView returns the default display value based on date and view mode
func (c *Calendar) getDefaultValueForView(habit Habit, date time.Time, viewMode ViewMode) string {
	today := time.Now()
	if date.After(today) {
		return c.getUntrackedSymbol(viewMode) // Future date
	}
	return c.getMissedSymbol(viewMode) // Past date with no entry (missed)
}

// GetCellValue returns the display value for a habit on a specific date (legacy method for weekly view)
func (c *Calendar) GetCellValueWeekly(habit Habit, date time.Time) string {
	return c.GetCellValue(habit, date, ViewModeWeekly)
}

// getDefaultValue returns the default display value based on date
func (c *Calendar) getDefaultValue(habit Habit, date time.Time) string {
	today := time.Now()
	if date.After(today) {
		return c.config.Views.Weekly.Untracked // Future date
	}
	return c.config.Views.Weekly.Missed // Past date with no entry (missed)
}

// getCompletedSymbol returns the completed symbol for the given view mode
func (c *Calendar) getCompletedSymbol(viewMode ViewMode) string {
	switch viewMode {
	case ViewModeWeekly:
		return c.config.Views.Weekly.Completed
	case ViewModeMonthly:
		return c.config.Views.Monthly.Completed
	case ViewModeHeatmap:
		return c.config.Views.Heatmap.Completed
	default:
		return "●"
	}
}

// getMissedSymbol returns the missed symbol for the given view mode
func (c *Calendar) getMissedSymbol(viewMode ViewMode) string {
	switch viewMode {
	case ViewModeWeekly:
		return c.config.Views.Weekly.Missed
	case ViewModeMonthly:
		return c.config.Views.Monthly.Missed
	case ViewModeHeatmap:
		return c.config.Views.Heatmap.Missed
	default:
		return "○"
	}
}

// getUntrackedSymbol returns the untracked symbol for the given view mode
func (c *Calendar) getUntrackedSymbol(viewMode ViewMode) string {
	switch viewMode {
	case ViewModeWeekly:
		return c.config.Views.Weekly.Untracked
	case ViewModeMonthly:
		return c.config.Views.Monthly.Untracked
	case ViewModeHeatmap:
		return c.config.Views.Heatmap.Untracked
	default:
		return "·"
	}
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
	for i := range c.habits {
		if c.habits[i].Name == name {
			c.habits = append(c.habits[:i], c.habits[i+1:]...)
			break
		}
	}
}

// ClearPendingHabits clears all pending habits
func (c *Calendar) ClearPendingHabits() {
	var filtered []Habit
	for _, h := range c.habits {
		if !h.Pending {
			filtered = append(filtered, h)
		}
	}
	c.habits = filtered
}

// GetPendingHabits returns all pending habits
func (c *Calendar) GetPendingHabits() []Habit {
	var pending []Habit
	for _, h := range c.habits {
		if h.Pending {
			pending = append(pending, h)
		}
	}
	return pending
}

// AddPendingHabit adds a habit to the pending list
func (c *Calendar) AddPendingHabit(habit Habit) {
	for _, h := range c.habits {
		if h.Name == habit.Name {
			return
		}
	}
	habit.Pending = true
	c.habits = append(c.habits, habit)
}

// ToggleEntry toggles or increments an entry based on habit type
func (c *Calendar) toggleOrIncrementEntry(habit Habit, date time.Time) {
	switch habit.Type {
	case HabitTypeBit:
		entry, exists := c.GetEntry(habit.Name, date)
		if exists {
			c.SetEntry(habit.Name, date, !entry.Completed, entry.Value, true)
		} else {
			c.SetEntry(habit.Name, date, true, "-", true)
		}

	case HabitTypeCount, HabitTypeFloat:
		entry, exists := c.GetEntry(habit.Name, date)
		if exists && entry.Value != "-" && entry.Value != "" {
			var val int
			fmt.Sscanf(entry.Value, "%d", &val)
			val++
			c.SetEntry(habit.Name, date, entry.Completed, fmt.Sprintf("%d", val), true)
		} else {
			c.SetEntry(habit.Name, date, false, "1", true)
		}
	}
}

// DecrementEntry decrements an entry for count/float habits
func (c *Calendar) decrementEntry(habit Habit, date time.Time) {
	entry, exists := c.GetEntry(habit.Name, date)

	if exists && entry.Value != "-" && entry.Value != "" {
		var val int
		fmt.Sscanf(entry.Value, "%d", &val)
		if val > 1 {
			val--
			c.SetEntry(habit.Name, date, entry.Completed, fmt.Sprintf("%d", val), true)
		} else {
			c.SetEntry(habit.Name, date, false, "-", true)
		}
	}
}

// GetPendingEntries returns all pending entries
func (c *Calendar) GetPendingEntries() []struct {
	HabitName string
	Date      time.Time
	Completed bool
	Value     string
} {
	var pending []struct {
		HabitName string
		Date      time.Time
		Completed bool
		Value     string
	}

	for habitName, habitEntries := range c.entries {
		for _, entry := range habitEntries {
			if entry.Pending {
				pending = append(pending, struct {
					HabitName string
					Date      time.Time
					Completed bool
					Value     string
				}{
					HabitName: habitName,
					Date:      entry.Date,
					Completed: entry.Completed,
					Value:     entry.Value,
				})
			}
		}
	}

	return pending
}

// ClearPendingEntries clears pending flag from all entries
func (c *Calendar) ClearPendingEntries() {
	for habitName := range c.entries {
		for dateKey := range c.entries[habitName] {
			entry := c.entries[habitName][dateKey]
			entry.Pending = false
			c.entries[habitName][dateKey] = entry
		}
	}
}

// WritePendingHabits writes all pending habits and entries to database via Store
func (c *Calendar) WritePendingHabits(createHabitsBulk func(habits []struct {
	Name      string
	HabitType string
	Goal      float64
}) ([]int, error), upsertEntry func(habitID int, date time.Time, value float64) error) error {
	var pendingHabits []Habit
	for _, h := range c.habits {
		if h.Pending {
			pendingHabits = append(pendingHabits, h)
		}
	}

	if len(pendingHabits) == 0 && len(c.GetPendingEntries()) == 0 {
		return nil
	}

	if len(pendingHabits) > 0 {
		habitInputs := make([]struct {
			Name      string
			HabitType string
			Goal      float64
		}, len(pendingHabits))

		for i, h := range pendingHabits {
			habitInputs[i] = struct {
				Name      string
				HabitType string
				Goal      float64
			}{
				Name:      h.Name,
				HabitType: h.Type.String(),
				Goal:      0,
			}
		}

		if _, err := createHabitsBulk(habitInputs); err != nil {
			return fmt.Errorf("failed to write habits: %w", err)
		}

		for i := range c.habits {
			c.habits[i].Pending = false
		}
	}

	pendingEntries := c.GetPendingEntries()
	for _, entry := range pendingEntries {
		habitID := -1
		for _, h := range c.habits {
			if h.Name == entry.HabitName {
				habitID = h.ID
				break
			}
		}

		if habitID == -1 {
			continue
		}

		var value float64
		if entry.Completed {
			value = 1
		} else if entry.Value != "-" && entry.Value != "" {
			fmt.Sscanf(entry.Value, "%f", &value)
		}

		if err := upsertEntry(habitID, entry.Date, value); err != nil {
			return fmt.Errorf("failed to upsert entry for %s on %s: %w", entry.HabitName, entry.Date.Format("2006-01-02"), err)
		}
	}

	c.ClearPendingEntries()

	return nil
}
