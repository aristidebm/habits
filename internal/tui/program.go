package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"example.com/habits/internal/app"
	"example.com/habits/internal/tui/calendar"
	"example.com/habits/internal/tui/command"
)

// NoteEditedMsg is sent when the user finishes editing a note
type NoteEditedMsg struct {
	Habit calendar.Habit
	Date  time.Time
	Note  string
	Error error
}

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
	// Ensure lipgloss is used
	_ = lipgloss.NewStyle()
	// Use default color profile
	// lipgloss.SetColorProfile(termenv.ANSI256)

	// Convert app habits to calendar habits
	calendarHabits := make([]calendar.Habit, len(application.GetHabits(context.Background())))
	for i, h := range application.GetHabits(context.Background()) {
		calendarHabits[i] = calendar.Habit{
			ID:      h.ID,
			Name:    h.Name,
			Type:    calendar.HabitType(h.Type),
			Goal:    h.Goal,
			Pending: false,
		}
	}

	// Create hasNoteFunc callback
	hasNoteFunc := func(habitID int, date time.Time) bool {
		// First check if there's a pending note
		// Note: We can't check pending notes here since we don't have habit name
		// The calendar handles pending notes in HasEntryNote

		// Check database for notes
		entry, err := application.GetEntry(context.Background(), habitID, date)
		if err != nil || entry == nil {
			return false
		}
		hasNote, err := application.HasNote(context.Background(), entry.ID)
		return err == nil && hasNote
	}

	cal := calendar.NewCalendar(calendarHabits, application.GetConfig(), application.GetStyles(), hasNoteFunc)

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
		Aliases:     []string{"a"},
		Description: "Add a new habit",
		Usage:       "add <name> <type> [goal]",
		Handler:     p.handleAddCommand,
	})

	p.commandLine.RegisterCommand(command.Command{
		Name:        "delete",
		Aliases:     []string{"d"},
		Description: "Delete a habit",
		Usage:       "delete <name>",
		Handler:     p.handleDeleteCommand,
	})

	p.commandLine.RegisterCommand(command.Command{
		Name:        "rename",
		Aliases:     []string{"r"},
		Description: "Rename a habit",
		Usage:       "rename <old_name> <new_name>",
		Handler:     p.handleRenameCommand,
	})

	p.commandLine.RegisterCommand(command.Command{
		Name:        "note",
		Aliases:     []string{"n"},
		Description: "Edit note for current habit entry",
		Usage:       "note",
		Handler:     p.handleNoteCommand,
	})

	p.commandLine.RegisterCommand(command.Command{
		Name:        "write",
		Aliases:     []string{"w"},
		Description: "Write pending habits to database",
		Usage:       "write",
		Handler:     p.handleWriteCommand,
	})

	p.commandLine.RegisterCommand(command.Command{
		Name:        "quit",
		Aliases:     []string{"q"},
		Description: "Quit the application",
		Usage:       "quit",
		Handler:     p.handleQuitCommand,
	})

	p.commandLine.RegisterCommand(command.Command{
		Name:        "wq",
		Description: "Write pending habits and quit",
		Usage:       "wq",
		Handler:     p.handleWQCommand,
	})

	p.commandLine.RegisterCommand(command.Command{
		Name:        "delete",
		Description: "Delete a habit",
		Usage:       "delete <name>",
		Handler:     p.handleDeleteCommand,
	})

	p.commandLine.RegisterCommand(command.Command{
		Name:        "rename",
		Description: "Rename a habit",
		Usage:       "rename <old_name> <new_name>",
		Handler:     p.handleRenameCommand,
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
		case "e":
			if !p.commandLine.IsVisible() {
				// Edit note for current habit entry
				habit := p.calendar.GetSelectedHabit()
				if habit != nil {
					selectedDate := p.calendar.GetSelectedDate()
					cmd := p.openNoteEditor(*habit, selectedDate)
					return p, cmd
				}
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

	case NoteEditedMsg:
		if msg.Error != nil {
			p.commandLine.SetError(fmt.Sprintf("Failed to edit note: %s", msg.Error))
		} else {
			// Store as pending note
			if p.calendar.PendingNotes[msg.Habit.Name] == nil {
				p.calendar.PendingNotes[msg.Habit.Name] = make(map[time.Time]string)
			}
			p.calendar.PendingNotes[msg.Habit.Name][msg.Date] = msg.Note

			// Update the entry's HasNote status
			if entry, exists := p.calendar.GetEntry(msg.Habit.Name, msg.Date); exists {
				entry.HasNote = msg.Note != ""
				p.calendar.SetEntryWithNote(msg.Habit.Name, msg.Date, entry.Completed, entry.Value, entry.Pending, entry.HasNote)
			}

			p.commandLine.SetSuccess("Note updated")
		}
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

	// Database path line (like Vim's file path)
	dbPath := p.app.Datasource
	dbPathLine := p.app.GetStyles().DBPath.Render(dbPath)

	// Add selected date on right side for monthly view
	if p.calendar.GetViewMode() == calendar.ViewModeMonthly {
		dateStr := p.app.GetStyles().DBPath.Render(p.calendar.GetSelectedDate().Format("02/Jan/06"))
		spacingLen := p.width - len(dbPath) - len(p.calendar.GetSelectedDate().Format("02/Jan/06")) - 2
		if spacingLen > 0 {
			spacing := strings.Repeat(" ", spacingLen)
			dbPathLine = lipgloss.JoinHorizontal(lipgloss.Left, dbPathLine, spacing, dateStr)
		} else {
			dbPathLine = lipgloss.JoinHorizontal(lipgloss.Left, dbPathLine, "  ", dateStr)
		}
	}

	view.WriteString(dbPathLine)
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
		return command.Error("Usage: add <name> <type> [goal]")
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

	// Parse optional goal
	var goal float64
	if len(args) >= 3 {
		_, err := fmt.Sscanf(args[2], "%f", &goal)
		if err != nil {
			return command.Error(fmt.Sprintf("Invalid goal value: %s", args[2]))
		}
	}

	// Add to pending habits (not database yet)
	p.calendar.AddPendingHabit(calendar.Habit{
		ID:   0,
		Name: name,
		Type: calendar.HabitType(habitType),
		Goal: goal,
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

func (p *Program) handleRenameCommand(args []string) command.Result {
	// Validate arguments
	if len(args) < 2 {
		return command.Error("Usage: rename <old_name> <new_name>")
	}

	oldName := args[0]
	newName := args[1]

	// Find habit by old name
	habits := p.app.GetHabits(context.Background())
	var habitID int
	var habitType app.HabitType
	var goal float64
	found := false
	for _, h := range habits {
		if h.Name == oldName {
			habitID = h.ID
			habitType = h.Type
			goal = h.Goal
			found = true
			break
		}
	}

	if !found {
		return command.Error(fmt.Sprintf("Habit '%s' not found", oldName))
	}

	// Check if new name already exists
	for _, h := range habits {
		if h.Name == newName {
			return command.Error(fmt.Sprintf("Habit '%s' already exists", newName))
		}
	}

	// Rename the habit
	if err := p.app.UpdateHabit(context.Background(), habitID, newName, habitType, goal); err != nil {
		return command.Error(fmt.Sprintf("Error: %s", err))
	}

	p.reloadCalendar()
	return command.Success(fmt.Sprintf("Renamed habit: %s -> %s", oldName, newName))
}

func (p *Program) handleNoteCommand(args []string) command.Result {
	selectedDate := p.calendar.GetSelectedDate()

	habit := p.calendar.GetSelectedHabit()
	if habit == nil {
		return command.Error("No habit selected")
	}

	return command.Result{
		Type: command.ResultSuccess,
		Cmd:  p.openNoteEditor(*habit, selectedDate),
	}
}

func (p *Program) handleTrackUpCommand(args []string) command.Result {
	// This is a placeholder - track commands are typically handled via key bindings
	return command.Success("Track up command")
}

func (p *Program) handleTrackDownCommand(args []string) command.Result {
	// This is a placeholder - track commands are typically handled via key bindings
	return command.Success("Track down command")
}

func (p *Program) openNoteEditor(habit calendar.Habit, date time.Time) tea.Cmd {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "habit_note_*.txt")
	if err != nil {
		return tea.Cmd(func() tea.Msg {
			return NoteEditedMsg{Habit: habit, Date: date, Error: fmt.Errorf("failed to create temp file: %w", err)}
		})
	}

	// Get existing note content
	existingNote := ""
	if p.calendar.PendingNotes != nil {
		if habitNotes, exists := p.calendar.PendingNotes[habit.Name]; exists {
			if note, hasNote := habitNotes[date]; hasNote {
				existingNote = note
			}
		}
	}

	// If no pending note, check database
	if existingNote == "" {
		entry, err := p.app.GetEntry(context.Background(), habit.ID, date)
		if err == nil && entry != nil {
			note, err := p.app.GetNote(context.Background(), entry.ID)
			if err == nil && note != nil {
				existingNote = note.Note
			}
		}
	}

	// Determine status
	status := "not tracked"
	entry, exists := p.calendar.GetEntry(habit.Name, date)
	if exists {
		switch habit.Type {
		case calendar.HabitTypeBit:
			if entry.Completed {
				status = "completed"
			} else {
				status = "not completed"
			}
		case calendar.HabitTypeCount, calendar.HabitTypeFloat:
			if entry.Value != "-" && entry.Value != "" {
				if habit.Goal > 0 {
					val, _ := strconv.ParseFloat(entry.Value, 64)
					if val >= habit.Goal {
						status = "goal met"
					} else {
						status = fmt.Sprintf("%.1f/%g", val, habit.Goal)
					}
				} else {
					status = entry.Value
				}
			} else {
				status = "skipped"
			}
		}
	}

	// Write context header and existing note to temp file
	header := fmt.Sprintf("# Habit: %s\n# Date: %s\n# Goal: %.0f\n# Status: %s\n\n",
		habit.Name,
		date.Format("2006-01-02"),
		habit.Goal,
		status,
	)

	_, err = tmpFile.WriteString(header)
	if err != nil {
		os.Remove(tmpFile.Name())
		return tea.Cmd(func() tea.Msg {
			return NoteEditedMsg{Habit: habit, Date: date, Error: fmt.Errorf("failed to write header: %w", err)}
		})
	}

	_, err = tmpFile.WriteString(existingNote)
	if err != nil {
		os.Remove(tmpFile.Name())
		return tea.Cmd(func() tea.Msg {
			return NoteEditedMsg{Habit: habit, Date: date, Error: fmt.Errorf("failed to write existing note: %w", err)}
		})
	}

	tmpFile.Close()

	// Launch external editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano" // fallback
	}

	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		defer os.Remove(tmpFile.Name()) // Clean up temp file

		if err != nil {
			// Editor exited with error - don't save note
			return NoteEditedMsg{Habit: habit, Date: date, Error: fmt.Errorf("editor exited with error: %w", err)}
		}

		// Read the edited content
		content, err := os.ReadFile(tmpFile.Name())
		if err != nil {
			return NoteEditedMsg{Habit: habit, Date: date, Error: fmt.Errorf("failed to read edited file: %w", err)}
		}

		// Extract note content (skip header lines)
		lines := strings.Split(string(content), "\n")
		noteContent := ""
		inNoteSection := false

		for _, line := range lines {
			if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
				// Skip header lines and empty lines at start
				if inNoteSection {
					noteContent += line + "\n"
				}
			} else {
				inNoteSection = true
				noteContent += line + "\n"
			}
		}

		// Trim whitespace
		noteContent = strings.TrimSpace(noteContent)

		return NoteEditedMsg{Habit: habit, Date: date, Note: noteContent}
	})
}

func (p *Program) handleWriteCommand(args []string) command.Result {
	pendingHabits := p.calendar.GetPendingHabits()
	pendingEntries := p.calendar.GetPendingEntries()

	// Count pending notes
	pendingNotesCount := 0
	for _, habitNotes := range p.calendar.PendingNotes {
		pendingNotesCount += len(habitNotes)
	}

	if len(pendingHabits) == 0 && len(pendingEntries) == 0 && pendingNotesCount == 0 {
		return command.Error("No pending habits, entries, or notes to write")
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
	}, func(habitEntryID int, note string) error {
		return p.app.UpsertNote(context.Background(), habitEntryID, note)
	}, func(habitID int, date time.Time) (int, error) {
		entry, err := p.app.GetEntry(context.Background(), habitID, date)
		if err != nil {
			return 0, err
		}
		if entry == nil {
			return 0, nil
		}
		return entry.ID, nil
	}); err != nil {
		return command.Error(fmt.Sprintf("Error: %s", err))
	}

	p.reloadCalendar()
	msg := fmt.Sprintf("Wrote %d habits, %d entries, and %d notes to database", len(pendingHabits), len(pendingEntries), pendingNotesCount)
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

func (p *Program) handleWQCommand(args []string) command.Result {
	// First write pending habits
	result := p.handleWriteCommand(args)
	if result.Type == command.ResultError {
		return result // Return the error if write failed
	}

	// Then quit
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
			Goal:    h.Goal,
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
			cal.SetEntryWithNote(habit.Name, entry.Date, completed, value, false, entry.HasNote)
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
