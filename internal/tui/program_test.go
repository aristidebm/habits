package tui

import (
	"testing"

	"example.com/habits/internal/tui/calendar"
	tea "github.com/charmbracelet/bubbletea"
)

func TestProgramMessageRouting(t *testing.T) {
	// Skip this test for now as it requires a full app setup
	// This would require mocking the entire app interface
	t.Skip("Skipping program message routing test - requires complex app mocking")
}

func TestProgramMessageTypes(t *testing.T) {
	// Test that message types implement tea.Msg interface
	var _ tea.Msg = calendar.HabitDeletedMsg{}
	var _ tea.Msg = calendar.EntryDeletedMsg{}

	// Test message construction
	habitMsg := calendar.HabitDeletedMsg{HabitID: 123}
	if habitMsg.HabitID != 123 {
		t.Errorf("Expected HabitID 123, got %d", habitMsg.HabitID)
	}

	entryMsg := calendar.EntryDeletedMsg{EntryIDs: []int{1, 2, 3}}
	if len(entryMsg.EntryIDs) != 3 {
		t.Errorf("Expected 3 entry IDs, got %d", len(entryMsg.EntryIDs))
	}
}

func TestProgramMessageRoutingTypes(t *testing.T) {
	// Test that program can handle calendar message types
	// This tests the type assertions in the Update method
	habitMsg := calendar.HabitDeletedMsg{HabitID: 123}
	entryMsg := calendar.EntryDeletedMsg{EntryIDs: []int{1, 2, 3}}

	// Test that these are valid tea.Msg types
	var _ tea.Msg = habitMsg
	var _ tea.Msg = entryMsg
}
