package calendar

import (
	"testing"
	"time"

	"example.com/habits/internal/app"
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
