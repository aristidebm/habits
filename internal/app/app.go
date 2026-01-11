package app

import (
	"fmt"
	"time"
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
	ID   string
	Name string
	Type HabitType
}

// HabitEntry represents a habit entry for a specific date
type HabitEntry struct {
	HabitID   string
	Date      time.Time
	Completed bool   // for bit type
	Value     string // for count/float type
}

// App is the core application that manages habits and entries
// This can be used by TUI, CLI, or Web interfaces
type App struct {
	// TODO: Add database connection here
	habits  []Habit
	entries map[string]map[time.Time]HabitEntry // habitID -> date -> entry
}

// NewApp creates a new application instance
func NewApp() *App {
	// TODO: Initialize database connection
	return &App{
		habits:  make([]Habit, 0),
		entries: make(map[string]map[time.Time]HabitEntry),
	}
}

// AddHabit adds a new habit
func (a *App) AddHabit(name string, habitType HabitType) error {
	// TODO: Store in database
	habit := Habit{
		ID:   fmt.Sprintf("habit_%d", len(a.habits)+1),
		Name: name,
		Type: habitType,
	}
	a.habits = append(a.habits, habit)
	return nil
}

// DeleteHabit deletes a habit by name
func (a *App) DeleteHabit(name string) error {
	// TODO: Delete from database
	for i, h := range a.habits {
		if h.Name == name {
			a.habits = append(a.habits[:i], a.habits[i+1:]...)
			delete(a.entries, h.ID)
			return nil
		}
	}
	return fmt.Errorf("habit '%s' not found", name)
}

// TrackUp increments or marks a habit as done
func (a *App) TrackUp(habitName string, date time.Time, value string) error {
	// TODO: Store in database
	habit := a.findHabitByName(habitName)
	if habit == nil {
		return fmt.Errorf("habit '%s' not found", habitName)
	}

	if a.entries[habit.ID] == nil {
		a.entries[habit.ID] = make(map[time.Time]HabitEntry)
	}

	dateKey := normalizeDate(date)

	switch habit.Type {
	case HabitTypeBit:
		a.entries[habit.ID][dateKey] = HabitEntry{
			HabitID:   habit.ID,
			Date:      dateKey,
			Completed: true,
			Value:     "",
		}
	case HabitTypeCount:
		// Parse current value and increment
		entry, exists := a.entries[habit.ID][dateKey]
		if !exists {
			entry = HabitEntry{
				HabitID: habit.ID,
				Date:    dateKey,
				Value:   "1",
			}
		} else {
			// TODO: Parse and increment
			entry.Value = value
		}
		a.entries[habit.ID][dateKey] = entry
	case HabitTypeFloat:
		a.entries[habit.ID][dateKey] = HabitEntry{
			HabitID: habit.ID,
			Date:    dateKey,
			Value:   value,
		}
	}

	return nil
}

// TrackDown decrements or marks a habit as not done
func (a *App) TrackDown(habitName string, date time.Time) error {
	// TODO: Store in database
	habit := a.findHabitByName(habitName)
	if habit == nil {
		return fmt.Errorf("habit '%s' not found", habitName)
	}

	if a.entries[habit.ID] == nil {
		return nil
	}

	dateKey := normalizeDate(date)

	switch habit.Type {
	case HabitTypeBit:
		a.entries[habit.ID][dateKey] = HabitEntry{
			HabitID:   habit.ID,
			Date:      dateKey,
			Completed: false,
			Value:     "",
		}
	case HabitTypeCount:
		// TODO: Decrement value
		delete(a.entries[habit.ID], dateKey)
	case HabitTypeFloat:
		delete(a.entries[habit.ID], dateKey)
	}

	return nil
}

// GetHabits returns all habits
func (a *App) GetHabits() []Habit {
	return a.habits
}

// GetEntry retrieves an entry for a habit on a specific date
func (a *App) GetEntry(habitID string, date time.Time) (HabitEntry, bool) {
	dateKey := normalizeDate(date)
	if a.entries[habitID] == nil {
		return HabitEntry{}, false
	}
	entry, exists := a.entries[habitID][dateKey]
	return entry, exists
}

// findHabitByName finds a habit by name
func (a *App) findHabitByName(name string) *Habit {
	for i := range a.habits {
		if a.habits[i].Name == name {
			return &a.habits[i]
		}
	}
	return nil
}

// normalizeDate removes time component from date
func normalizeDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}
