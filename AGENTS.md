# AGENTS.md - Agent Guide for Habits Codebase

This document provides essential information for agentic coding assistants working in this repository.

## Project Overview

Habits is a terminal-based habit tracking application written in Go using the Bubble Tea TUI framework. It follows clean architecture principles with separation between business logic (`internal/app/`) and UI components (`internal/tui/`).

## Build and Development Commands

### Running the Application
```bash
go run ./...
# or
make run
```

### Building
```bash
go build ./...
# Output binary: habits (or habits.exe on Windows)
```

### Formatting
```bash
go fmt ./...
# or
make format
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/app

# Run a single test
go test -run TestSpecificName ./internal/app

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...
```

### Linting
No explicit linter configured. Use `go vet ./...`

## Code Style Guidelines

### Imports
Group: standard library → third-party → internal (blank lines between). Use aliases only when necessary:
```go
import (
    "fmt"
    "time"
    tea "github.com/charmbracelet/bubbletea"
    "example.com/habits/internal/app"
)
```

### Naming Conventions
- Types: PascalCase (`HabitType`, `Calendar`)
- Exported funcs: PascalCase (`NewCalendar`, `TrackUp`)
- Unexported funcs: camelCase (`normalizeDate`)
- Enums: `iota` starting at 0
- Receivers: Single lowercase letter (`c *Calendar`, `a *App`)

### Constants/Enums
```go
type ViewMode int
const (
    ViewModeWeekly ViewMode = iota
    ViewModeMonthly
)
```

### Struct Design
Unexported fields, `New*` constructors, pointer receivers for mutation:
```go
type Calendar struct { habits []Habit; entries map[string]map[time.Time]HabitEntry }
func NewCalendar(habits []Habit) *Calendar {
    return &Calendar{habits: habits, entries: make(map[string]map[time.Time]HabitEntry)}
}
```

### Error Handling
Return `error`, use `fmt.Errorf` with context:
```go
func (a *App) DeleteHabit(name string) error {
    return fmt.Errorf("habit '%s' not found", name)
}
```

### Bubble Tea Pattern
```go
func (c *Calendar) Init() tea.Cmd { return nil }
func (c *Calendar) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return c, nil }
func (c *Calendar) View() string { return "content" }
```

### Date Handling
Always normalize: `time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)`

### Lipgloss Styling
```go
cellStyle: lipgloss.NewStyle().Align(lipgloss.Center).Width(8)
```

### Comments
Doc comments for exported types/funcs. Use TODO for pending.

### File Organization
- `cmd/main.go` - Entry point
- `internal/app/` - Business logic
- `internal/tui/` - Terminal UI

### Testing
`*_test.go` files, table-driven tests, `Test<FunctionName>` naming

## Architecture Notes

App layer: data. TUI layer: presentation. Command pattern for CLI. View pattern (WeeklyView/MonthlyView). State immutable through Update methods. Dependencies: Bubble Tea, Charmbracelet, Go 1.25.0+. Check existing patterns before adding deps.
