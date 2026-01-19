package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"example.com/habits/internal/app"
	"example.com/habits/internal/tui/calendar"
	"example.com/habits/internal/tui/command"
)

// Program is the main TUI program
type Program struct {
	app         *app.App
	calendar    *calendar.Calendar
	commandLine *command.CommandLine
	width       int
	height      int
}

// NewProgram creates a new TUI program
func NewProgram(application *app.App) *Program {
	// Convert app habits to calendar habits
	calendarHabits := make([]calendar.Habit, len(application.GetHabits(context.Background())))
	for i, h := range application.GetHabits(context.Background()) {
		calendarHabits[i] = calendar.Habit{
			ID:      h.ID,
			Name:    h.Name,
			Type:    calendar.HabitType(h.Type),
			Pending: false,
		}
	}

	cal := calendar.NewCalendar(calendarHabits, application.GetConfig())

	// Sync entries (but not pending habits - they have no entries yet)
	syncEntriesToCalendar(application, cal, true)

	cmdLine := command.NewCommandLine()

	p := &Program{
		app:         application,
		calendar:    cal,
		commandLine: cmdLine,
		width:       80,
		height:      24,
	}

	// Register commands
	p.registerCommands()

	return p
}

// registerCommands registers all available commands
func (p *Program) registerCommands() {
	p.commandLine.RegisterCommand(command.Command{
		Name:        "add",
		Description: "Add a new habit",
		Usage:       "add <name> <type>  (types: bit, count, float)",
		Handler:     p.handleAddCommand,
	})

	p.commandLine.RegisterCommand(command.Command{
		Name:        "delete",
		Description: "Delete a habit",
		Usage:       "delete <name>",
		Handler:     p.handleDeleteCommand,
	})

	p.commandLine.RegisterCommand(command.Command{
		Name:        "track-up",
		Description: "Mark habit as done or increment value",
		Usage:       "track-up <habit> [value]",
		Handler:     p.handleTrackUpCommand,
	})

	p.commandLine.RegisterCommand(command.Command{
		Name:        "track-down",
		Description: "Mark habit as not done or decrement value",
		Usage:       "track-down <habit>",
		Handler:     p.handleTrackDownCommand,
	})

	p.commandLine.RegisterCommand(command.Command{
		Name:        "next-month",
		Description: "Go to next month",
		Usage:       "next-month",
		Handler:     p.handleNextMonthCommand,
	})

	p.commandLine.RegisterCommand(command.Command{
		Name:        "prev-month",
		Description: "Go to previous month",
		Usage:       "prev-month",
		Handler:     p.handlePrevMonthCommand,
	})

	p.commandLine.RegisterCommand(command.Command{
		Name:        "quit",
		Description: "Quit the application",
		Usage:       "quit",
		Handler:     p.handleQuitCommand,
	})

	p.commandLine.RegisterCommand(command.Command{
		Name:        "write",
		Description: "Write all pending habits to database",
		Usage:       "write",
		Handler:     p.handleWriteCommand,
	})

	p.commandLine.RegisterCommand(command.Command{
		Name:        "export",
		Description: "Export all habits and entries to JSON file",
		Usage:       "export <path>",
		Handler:     p.handleExportCommand,
	})
}

// Init initializes the program
func (p *Program) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (p *Program) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// If command line is visible, let it handle the input
		if p.commandLine.IsVisible() {
			return p, p.commandLine.Update(msg)
		}

		// Handle global keys
		switch msg.String() {
		case ":":
			return p, p.commandLine.Show()
		case "q":
			if !p.commandLine.IsVisible() {
				return p, tea.Quit
			}
		}

		// Let calendar handle the input
		_, cmd := p.calendar.Update(msg)
		return p, cmd

	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		p.calendar.Resize(msg.Width, msg.Height-2) // -2 for command line
		p.commandLine.SetWidth(msg.Width)
		return p, nil
	}

	// If command line is visible, update it
	if p.commandLine.IsVisible() {
		return p, p.commandLine.Update(msg)
	}

	return p, nil
}

// View renders the program
func (p *Program) View() string {
	var view strings.Builder

	// Render calendar
	view.WriteString(p.calendar.View())
	view.WriteString("\n")

	// Render command line or status
	view.WriteString(p.commandLine.View())

	// Help text
	if !p.commandLine.IsVisible() && p.commandLine.View() == "" {
		helpText := "[:]cmd [q]quit"
		view.WriteString(helpText)
	}

	return view.String()
}

// Command handlers - each validates its own arguments

func (p *Program) handleAddCommand(args []string) command.Result {
	// Validate arguments
	if len(args) < 2 {
		return command.Error("Usage: add <name> <type>")
	}

	name := args[0]
	typeStr := strings.ToLower(args[1])

	var habitType app.HabitType
	switch typeStr {
	case "bit":
		habitType = app.HabitTypeBit
	case "count":
		habitType = app.HabitTypeCount
	case "float":
		habitType = app.HabitTypeFloat
	default:
		return command.Error(fmt.Sprintf("Invalid habit type: %s (use: bit, count, float)", typeStr))
	}

	// Add to pending habits (not database yet)
	p.calendar.AddPendingHabit(calendar.Habit{
		ID:   0,
		Name: name,
		Type: calendar.HabitType(habitType),
	})

	// Reload calendar with new habits
	p.reloadCalendar()
	return command.Success(fmt.Sprintf("Added habit: %s (pending, use :write to save)", name))
}

func (p *Program) handleDeleteCommand(args []string) command.Result {
	// Validate arguments
	if len(args) < 1 {
		return command.Error("Usage: delete <name>")
	}

	name := args[0]

	// Check pending habits first
	pending := p.calendar.GetPendingHabits()
	var habitID int
	for _, h := range pending {
		if h.Name == name {
			p.calendar.RemovePendingHabit(name)
			p.reloadCalendar()
			return command.Success(fmt.Sprintf("Deleted pending habit: %s", name))
		}
	}

	// Find habit by name in database
	habits := p.app.GetHabits(context.Background())
	found := false
	for _, h := range habits {
		if h.Name == name {
			habitID = h.ID
			found = true
			break
		}
	}

	if !found {
		return command.Error(fmt.Sprintf("Habit '%s' not found", name))
	}

	// Execute command
	if err := p.app.DeleteHabit(context.Background(), habitID); err != nil {
		return command.Error(fmt.Sprintf("Error: %s", err))
	}

	p.reloadCalendar()
	return command.Success(fmt.Sprintf("Deleted habit: %s", name))
}

func (p *Program) handleTrackUpCommand(args []string) command.Result {
	// Validate arguments
	if len(args) < 1 {
		return command.Error("Usage: track-up <habit> [value]")
	}

	habitName := args[0]
	valueStr := "1"
	if len(args) > 1 {
		valueStr = args[1]
	}

	// Find habit by name
	habits := p.app.GetHabits(context.Background())
	var habitID int
	var habitType app.HabitType
	found := false
	for _, h := range habits {
		if h.Name == habitName {
			habitID = h.ID
			habitType = h.Type
			found = true
			break
		}
	}

	if !found {
		return command.Error(fmt.Sprintf("Habit '%s' not found", habitName))
	}

	// Parse value
	var value float64
	switch habitType {
	case app.HabitTypeBit:
		value = 1
	case app.HabitTypeCount, app.HabitTypeFloat:
		_, err := fmt.Sscanf(valueStr, "%f", &value)
		if err != nil {
			return command.Error(fmt.Sprintf("Invalid value: %s", valueStr))
		}
	}

	// Execute command
	selectedDate := p.calendar.GetSelectedDate()
	if err := p.app.UpsertEntry(context.Background(), habitID, selectedDate, value); err != nil {
		return command.Error(fmt.Sprintf("Error: %s", err))
	}

	p.reloadCalendar()
	return command.Success(fmt.Sprintf("Tracked up: %s", habitName))
}

func (p *Program) handleTrackDownCommand(args []string) command.Result {
	// Validate arguments
	if len(args) < 1 {
		return command.Error("Usage: track-down <habit>")
	}

	habitName := args[0]

	// Find habit by name
	habits := p.app.GetHabits(context.Background())
	var habitID int
	var habitType app.HabitType
	found := false
	for _, h := range habits {
		if h.Name == habitName {
			habitID = h.ID
			habitType = h.Type
			found = true
			break
		}
	}

	if !found {
		return command.Error(fmt.Sprintf("Habit '%s' not found", habitName))
	}

	// Execute command
	selectedDate := p.calendar.GetSelectedDate()
	var value float64
	switch habitType {
	case app.HabitTypeBit:
		value = 0
	case app.HabitTypeCount, app.HabitTypeFloat:
		value = 0
	}

	if err := p.app.UpsertEntry(context.Background(), habitID, selectedDate, value); err != nil {
		return command.Error(fmt.Sprintf("Error: %s", err))
	}

	p.reloadCalendar()
	return command.Success(fmt.Sprintf("Tracked down: %s", habitName))
}

func (p *Program) handleWriteCommand(args []string) command.Result {
	pendingHabits := p.calendar.GetPendingHabits()
	pendingEntries := p.calendar.GetPendingEntries()
	if len(pendingHabits) == 0 && len(pendingEntries) == 0 {
		return command.Error("No pending habits or entries to write")
	}

	if err := p.calendar.WritePendingHabits(func(habits []struct {
		Name      string
		HabitType string
		Goal      float64
	}) ([]int, error) {
		habitInputs := make([]struct {
			Name      string
			HabitType app.HabitType
			Goal      float64
		}, len(habits))

		for i, h := range habits {
			var hType app.HabitType
			switch h.HabitType {
			case "bit":
				hType = app.HabitTypeBit
			case "count":
				hType = app.HabitTypeCount
			case "float":
				hType = app.HabitTypeFloat
			default:
				return nil, fmt.Errorf("invalid habit type: %s", h.HabitType)
			}
			habitInputs[i] = struct {
				Name      string
				HabitType app.HabitType
				Goal      float64
			}{
				Name:      h.Name,
				HabitType: hType,
				Goal:      h.Goal,
			}
		}

		createdHabits, err := p.app.CreateHabitsBulk(context.Background(), habitInputs)
		if err != nil {
			return nil, err
		}

		ids := make([]int, len(createdHabits))
		for i, h := range createdHabits {
			ids[i] = h.ID
		}

		return ids, nil
	}, func(habitID int, date time.Time, value float64) error {
		return p.app.UpsertEntry(context.Background(), habitID, date, value)
	}); err != nil {
		return command.Error(fmt.Sprintf("Error: %s", err))
	}

	p.reloadCalendar()
	msg := fmt.Sprintf("Wrote %d habits and %d entries to database", len(pendingHabits), len(pendingEntries))
	return command.Success(msg)
}

func (p *Program) handleNextMonthCommand(args []string) command.Result {
	// No arguments needed - validation implicit
	p.calendar.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'L'}})
	return command.Success("Moved to next month")
}

func (p *Program) handlePrevMonthCommand(args []string) command.Result {
	// No arguments needed - validation implicit
	p.calendar.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'H'}})
	return command.Success("Moved to previous month")
}

func (p *Program) handleQuitCommand(args []string) command.Result {
	// No arguments needed - validation implicit
	return command.Quit()
}

// reloadCalendar reloads the calendar with current app data
func (p *Program) reloadCalendar() {
	// Save pending habits before reload
	pending := p.calendar.GetPendingHabits()

	// Convert app habits to calendar habits
	calendarHabits := make([]calendar.Habit, len(p.app.GetHabits(context.Background())))
	for i, h := range p.app.GetHabits(context.Background()) {
		calendarHabits[i] = calendar.Habit{
			ID:      h.ID,
			Name:    h.Name,
			Type:    calendar.HabitType(h.Type),
			Pending: false,
		}
	}

	// Append pending habits to calendar habits list
	calendarHabits = append(calendarHabits, pending...)

	p.calendar.ReloadHabits(calendarHabits)
	p.calendar.Resize(p.width, p.height-2)

	// Sync entries (but not pending habits - they have no entries yet)
	syncEntriesToCalendar(p.app, p.calendar, true)
}

// syncEntriesToCalendar syncs entries from app to calendar
func syncEntriesToCalendar(application *app.App, cal *calendar.Calendar, skipPending bool) {
	for _, habit := range application.GetHabits(context.Background()) {
		// Get all entries (no date limit)
		start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Now().AddDate(10, 0, 0) // Far future

		slog.Info("Attempt to get entries with", "start", start, "end", end)

		entries, err := application.ListEntries(context.Background(), habit.ID, start, end)
		if err != nil {
			continue
		}

		slog.Info("Got entries from database", "count", len(entries))

		for _, entry := range entries {
			var completed bool
			var value string
			if habit.Type == app.HabitTypeBit {
				completed = entry.Value == 1
				value = ""
			} else {
				completed = entry.Value > 0
				value = fmt.Sprintf("%.0f", entry.Value)
				if value == "1.000000" {
					value = "1"
				}
			}
			cal.SetEntry(habit.Name, entry.Date, completed, value, false)
		}
	}
}

func (p *Program) handleExportCommand(args []string) command.Result {
	// Validate arguments
	if len(args) < 1 {
		return command.Error("Usage: export <path>")
	}

	path := args[0]

	// Export habits and entries
	exportHabits, err := p.app.Export(context.Background())
	if err != nil {
		return command.Error(fmt.Sprintf("Error exporting data: %s", err))
	}

	// Write to JSON file
	file, err := os.Create(path)
	if err != nil {
		return command.Error(fmt.Sprintf("Error creating file: %s", err))
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(exportHabits); err != nil {
		return command.Error(fmt.Sprintf("Error writing JSON: %s", err))
	}

	return command.Success(fmt.Sprintf("Exported %d habits to %s", len(exportHabits), path))
}
