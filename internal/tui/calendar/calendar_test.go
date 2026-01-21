package calendar

import (
	"testing"
	"time"

	"example.com/habits/internal/app"
	tea "github.com/charmbracelet/bubbletea"
)

func TestCalendarViewMode(t *testing.T) {
	// Test view mode constants
	if ViewModeWeekly != 0 {
		t.Errorf("Expected ViewModeWeekly to be 0, got %d", ViewModeWeekly)
	}
	if ViewModeMonthly != 1 {
		t.Errorf("Expected ViewModeMonthly to be 1, got %d", ViewModeMonthly)
	}
}

func TestCalendarGetCompletedSymbol(t *testing.T) {
	// Test the default behavior
	if app.DefaultConfig().Views.Weekly.Completed != "🗸" {
		t.Errorf("Expected default weekly completed '🗸', got '%s'", app.DefaultConfig().Views.Weekly.Completed)
	}
}

func TestCalendarGetMissedSymbol(t *testing.T) {
	if app.DefaultConfig().Views.Weekly.Missed != "⛌" {
		t.Errorf("Expected default weekly missed '⛌', got '%s'", app.DefaultConfig().Views.Weekly.Missed)
	}
}

func TestCalendarGetUntrackedSymbol(t *testing.T) {
	if app.DefaultConfig().Views.Weekly.Untracked != "○" {
		t.Errorf("Expected default weekly untracked '○', got '%s'", app.DefaultConfig().Views.Weekly.Untracked)
	}
}

func TestHabitTypeString(t *testing.T) {
	if HabitTypeBit.String() != "bit" {
		t.Errorf("Expected HabitTypeBit string 'bit', got '%s'", HabitTypeBit.String())
	}
	if HabitTypeCount.String() != "count" {
		t.Errorf("Expected HabitTypeCount string 'count', got '%s'", HabitTypeCount.String())
	}
	if HabitTypeFloat.String() != "float" {
		t.Errorf("Expected HabitTypeFloat string 'float', got '%s'", HabitTypeFloat.String())
	}
}

func TestTimeNormalization(t *testing.T) {
	// Test that dates are normalized (time part set to 00:00:00)
	date := time.Date(2024, 1, 15, 14, 30, 45, 123456789, time.UTC)
	normalized := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	if date.Hour() != 14 {
		t.Errorf("Original date should have hour 14, got %d", date.Hour())
	}
	if normalized.Hour() != 0 {
		t.Errorf("Normalized date should have hour 0, got %d", normalized.Hour())
	}
	if normalized.Minute() != 0 {
		t.Errorf("Normalized date should have minute 0, got %d", normalized.Minute())
	}
	if normalized.Second() != 0 {
		t.Errorf("Normalized date should have second 0, got %d", normalized.Second())
	}
}

// Test message handling in calendar Update method
func TestCalendarMessageHandling(t *testing.T) {
	config := app.DefaultConfig()
	tm := app.NewThemeManager()
	theme, _ := tm.LoadTheme("default")
	styles := app.NewStyles(theme)
	habits := []Habit{
		{ID: 1, Name: "Exercise", Type: HabitTypeBit, Goal: 0},
		{ID: 2, Name: "Reading", Type: HabitTypeCount, Goal: 1},
		{ID: 3, Name: "Water", Type: HabitTypeCount, Goal: 8},
	}

	cal := NewCalendar(habits, config, styles, nil)

	// Test HabitDeletedMsg - should remove habit and adjust cursor
	initialLength := len(cal.habits)
	cal.selectedHabit = 1 // Select middle habit

	_, cmd := cal.Update(HabitDeletedMsg{HabitID: 2}) // Delete "Reading" habit
	if cmd != nil {
		t.Errorf("HabitDeletedMsg should not return a command")
	}

	if len(cal.habits) != initialLength-1 {
		t.Errorf("Expected %d habits after deletion, got %d", initialLength-1, len(cal.habits))
	}

	// Should still have "Exercise" and "Water"
	if len(cal.habits) >= 1 && cal.habits[0].Name != "Exercise" {
		t.Errorf("Expected first habit to be 'Exercise', got '%s'", cal.habits[0].Name)
	}

	// Cursor should adjust (was at index 1, now should be at index 1 which is "Water")
	if cal.selectedHabit != 1 {
		t.Errorf("Expected selected habit index to be 1, got %d", cal.selectedHabit)
	}

	// Test EntryDeletedMsg - should be handled without error
	_, cmd = cal.Update(EntryDeletedMsg{EntryIDs: []int{1, 2, 3}})
	if cmd != nil {
		t.Errorf("EntryDeletedMsg should not return a command")
	}
}

func TestCalendarCursorPreservation(t *testing.T) {
	config := app.DefaultConfig()
	tm := app.NewThemeManager()
	theme, _ := tm.LoadTheme("default")
	styles := app.NewStyles(theme)
	habits := []Habit{
		{ID: 1, Name: "Exercise", Type: HabitTypeBit, Goal: 0},
		{ID: 2, Name: "Reading", Type: HabitTypeCount, Goal: 1},
		{ID: 3, Name: "Water", Type: HabitTypeCount, Goal: 8},
	}

	cal := NewCalendar(habits, config, styles, nil)

	// Set cursor to middle habit
	cal.selectedHabit = 1
	selectedDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	cal.selectedDate = selectedDate

	// Delete the habit before cursor
	_, _ = cal.Update(HabitDeletedMsg{HabitID: 1}) // Delete "Exercise"

	// Cursor should now be at index 0 (previously index 1, but list shifted)
	if cal.selectedHabit != 0 {
		t.Errorf("Expected cursor to move to index 0 after deletion, got %d", cal.selectedHabit)
	}

	// Date should be preserved
	if !cal.selectedDate.Equal(selectedDate) {
		t.Errorf("Expected date to be preserved, got %v", cal.selectedDate)
	}
}

func TestCalendarReloadHabitsWithSelection(t *testing.T) {
	config := app.DefaultConfig()
	tm := app.NewThemeManager()
	theme, _ := tm.LoadTheme("default")
	styles := app.NewStyles(theme)

	// Initial habits
	oldHabits := []Habit{
		{ID: 1, Name: "Exercise", Type: HabitTypeBit, Goal: 0},
		{ID: 2, Name: "Reading", Type: HabitTypeCount, Goal: 1},
	}

	cal := NewCalendar(oldHabits, config, styles, nil)

	// Set initial state
	cal.selectedHabit = 1 // "Reading"
	selectedDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	cal.selectedDate = selectedDate

	// New habits (different order, "Reading" is now first)
	newHabits := []Habit{
		{ID: 2, Name: "Reading", Type: HabitTypeCount, Goal: 1}, // Same habit, different position
		{ID: 4, Name: "Writing", Type: HabitTypeBit, Goal: 0},   // New habit
	}

	// Reload with cursor preservation
	cal.ReloadHabitsWithSelection(newHabits, "Reading", selectedDate)

	// Should have 2 habits
	if len(cal.habits) != 2 {
		t.Errorf("Expected 2 habits, got %d", len(cal.habits))
	}

	// Cursor should be at index 0 ("Reading" is now first)
	if cal.selectedHabit != 0 {
		t.Errorf("Expected cursor at index 0, got %d", cal.selectedHabit)
	}

	// Date should be preserved
	if !cal.selectedDate.Equal(selectedDate) {
		t.Errorf("Expected date to be preserved")
	}

	// Test with non-existent habit name
	cal.ReloadHabitsWithSelection(newHabits, "NonExistent", selectedDate)
	// Should default to first habit
	if cal.selectedHabit != 0 {
		t.Errorf("Expected cursor to default to 0 for non-existent habit, got %d", cal.selectedHabit)
	}
}

func TestHabitEntryStatus(t *testing.T) {
	// Test status enum values
	if HabitEntryStatusActive != 0 {
		t.Errorf("Expected HabitEntryStatusActive to be 0, got %d", HabitEntryStatusActive)
	}
	if HabitEntryStatusPendingAdding != 1 {
		t.Errorf("Expected HabitEntryStatusPendingAdding to be 1, got %d", HabitEntryStatusPendingAdding)
	}
	if HabitEntryStatusPendingDeletion != 2 {
		t.Errorf("Expected HabitEntryStatusPendingDeletion to be 2, got %d", HabitEntryStatusPendingDeletion)
	}

	// Test String() method
	if HabitEntryStatusActive.String() != "active" {
		t.Errorf("Expected 'active', got '%s'", HabitEntryStatusActive.String())
	}
	if HabitEntryStatusPendingAdding.String() != "pending_adding" {
		t.Errorf("Expected 'pending_adding', got '%s'", HabitEntryStatusPendingAdding.String())
	}
	if HabitEntryStatusPendingDeletion.String() != "pending_deletion" {
		t.Errorf("Expected 'pending_deletion', got '%s'", HabitEntryStatusPendingDeletion.String())
	}
}

func TestCalendarEntryStatusDisplay(t *testing.T) {
	config := app.DefaultConfig()
	tm := app.NewThemeManager()
	theme, _ := tm.LoadTheme("default")
	styles := app.NewStyles(theme)
	habits := []Habit{
		{ID: 1, Name: "Exercise", Type: HabitTypeCount, Goal: 5},
	}

	cal := NewCalendar(habits, config, styles, nil)

	// Add entry with pending deletion status
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	cal.SetEntryWithStatus("Exercise", testDate, false, "-", false, HabitEntryStatusPendingDeletion, false)

	// Check that it displays as untracked
	displayValue := cal.GetCellValue(habits[0], testDate, ViewModeWeekly)
	expected := config.Views.Weekly.Untracked
	if displayValue != expected {
		t.Errorf("Expected pending deletion entry to show '%s', got '%s'", expected, displayValue)
	}

	// Test with monthly view too
	displayValue = cal.GetCellValue(habits[0], testDate, ViewModeMonthly)
	expected = config.Views.Monthly.Untracked
	if displayValue != expected {
		t.Errorf("Expected pending deletion entry to show '%s' in monthly view, got '%s'", expected, displayValue)
	}
}

func TestCalendarDecrementToDeletion(t *testing.T) {
	config := app.DefaultConfig()
	tm := app.NewThemeManager()
	theme, _ := tm.LoadTheme("default")
	styles := app.NewStyles(theme)
	habits := []Habit{
		{ID: 1, Name: "Exercise", Type: HabitTypeCount, Goal: 5},
	}

	cal := NewCalendar(habits, config, styles, nil)
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Set initial value of 2
	cal.SetEntry("Exercise", testDate, false, "2", false)

	// Decrement to 1
	cal.decrementEntry(habits[0], testDate)
	entry, exists := cal.GetEntry("Exercise", testDate)
	if !exists || entry.Value != "1" {
		t.Errorf("Expected value '1', got '%s'", entry.Value)
	}

	// Decrement to 0 - should become pending deletion
	cal.decrementEntry(habits[0], testDate)
	entry, exists = cal.GetEntry("Exercise", testDate)
	if !exists || entry.Value != "-" || entry.Status != HabitEntryStatusPendingDeletion {
		t.Errorf("Expected pending deletion status with value '-', got status %v value '%s'", entry.Status, entry.Value)
	}
}

func TestCalendarMessageTypes(t *testing.T) {
	// Test that message types implement tea.Msg interface
	var _ tea.Msg = HabitEditMsg{}
	var _ tea.Msg = EntryEditMsg{}
	var _ tea.Msg = NoteEditMsg{}
	var _ tea.Msg = HabitDeletedMsg{}
	var _ tea.Msg = EntryDeletedMsg{}

	// Test message construction
	habitMsg := HabitEditMsg{HabitID: 123}
	if habitMsg.HabitID != 123 {
		t.Errorf("Expected HabitID 123, got %d", habitMsg.HabitID)
	}

	entryMsg := EntryEditMsg{HabitID: 456, Date: time.Now()}
	if entryMsg.HabitID != 456 {
		t.Errorf("Expected HabitID 456, got %d", entryMsg.HabitID)
	}

	noteMsg := NoteEditMsg{HabitEntryID: 789}
	if noteMsg.HabitEntryID != 789 {
		t.Errorf("Expected HabitEntryID 789, got %d", noteMsg.HabitEntryID)
	}

	habitDelMsg := HabitDeletedMsg{HabitID: 101}
	if habitDelMsg.HabitID != 101 {
		t.Errorf("Expected HabitID 101, got %d", habitDelMsg.HabitID)
	}

	entryDelMsg := EntryDeletedMsg{EntryIDs: []int{1, 2, 3}}
	if len(entryDelMsg.EntryIDs) != 3 {
		t.Errorf("Expected 3 entry IDs, got %d", len(entryDelMsg.EntryIDs))
	}
}
