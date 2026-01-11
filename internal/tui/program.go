package tui

import (
	"fmt"
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
	calendarHabits := make([]calendar.Habit, len(application.GetHabits()))
	for i, h := range application.GetHabits() {
		calendarHabits[i] = calendar.Habit{
			ID:   h.ID,
			Name: h.Name,
			Type: calendar.HabitType(h.Type),
		}
	}

	cal := calendar.NewCalendar(calendarHabits)

	// Sync entries from app to calendar
	syncEntriesToCalendar(application, cal)

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

	// Execute command
	if err := p.app.AddHabit(name, habitType); err != nil {
		return command.Error(fmt.Sprintf("Error: %s", err))
	}

	// Reload calendar with new habits
	p.reloadCalendar()
	return command.Success(fmt.Sprintf("Added habit: %s", name))
}

func (p *Program) handleDeleteCommand(args []string) command.Result {
	// Validate arguments
	if len(args) < 1 {
		return command.Error("Usage: delete <name>")
	}

	name := args[0]

	// Execute command
	if err := p.app.DeleteHabit(name); err != nil {
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
	value := ""
	if len(args) > 1 {
		value = args[1]
	}

	// Execute command
	selectedDate := p.calendar.GetSelectedDate()
	if err := p.app.TrackUp(habitName, selectedDate, value); err != nil {
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

	// Execute command
	selectedDate := p.calendar.GetSelectedDate()
	if err := p.app.TrackDown(habitName, selectedDate); err != nil {
		return command.Error(fmt.Sprintf("Error: %s", err))
	}

	p.reloadCalendar()
	return command.Success(fmt.Sprintf("Tracked down: %s", habitName))
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
	// Convert app habits to calendar habits
	calendarHabits := make([]calendar.Habit, len(p.app.GetHabits()))
	for i, h := range p.app.GetHabits() {
		calendarHabits[i] = calendar.Habit{
			ID:   h.ID,
			Name: h.Name,
			Type: calendar.HabitType(h.Type),
		}
	}

	p.calendar = calendar.NewCalendar(calendarHabits)
	p.calendar.Resize(p.width, p.height-2)

	// Sync entries
	syncEntriesToCalendar(p.app, p.calendar)
}

// syncEntriesToCalendar syncs entries from app to calendar
func syncEntriesToCalendar(application *app.App, cal *calendar.Calendar) {
	for _, habit := range application.GetHabits() {
		// Get entries for the past year
		now := time.Now()
		start := now.AddDate(-1, 0, 0)

		for d := start; !d.After(now); d = d.AddDate(0, 0, 1) {
			entry, exists := application.GetEntry(habit.ID, d)
			if exists {
				cal.SetEntry(habit.Name, d, entry.Completed, entry.Value)
			}
		}
	}
}
