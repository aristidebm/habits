# Habits

A terminal-based habit tracking application written in Go, featuring a modern text-based user interface for managing and tracking daily habits.

## Overview

Habits is a command-line application under development that provides an intuitive calendar-based interface for tracking personal habits with different data types. The application aims to support boolean, integer, and decimal value tracking, making it suitable for various habit types from simple completion tracking to quantitative measurements.

## Features

- **Calendar-based tracking** - Visual calendar interface showing habit completion over time
- **Multiple habit types** - Support for boolean (bit), integer (count), and decimal (float) values
- **Flexible view modes** - Weekly and monthly calendar views
- **Vim-style command system** - Efficient keyboard-driven interface
- **Real-time updates** - Immediate visual feedback for habit tracking
- **Date navigation** - Easy movement between months and specific date selection

*Note: This application is currently in development and not all features are fully functional.*

## Architecture

The application follows clean architecture principles with clear separation of concerns:

- `cmd/main.go` - Application entry point
- `internal/app/` - Core business logic and data models
- `internal/tui/` - Terminal user interface components
- `internal/cli/` - Command-line interface (prepared for future CLI mode)

## Installation

### Prerequisites

- Go 1.25.0 or later
- Make (optional, for build scripts)

### Building from Source

```bash
git clone <repository-url>
cd habits
go build ./...
```

### Running

```bash
make run
# or
go run ./...
```

## Usage

### Interface Navigation

The application uses a keyboard-driven interface with the following controls:

- `:` - Open command line
- `q` - Quick quit
- Arrow keys/HJKL - Navigate calendar and habits
- `L`/`H` - Navigate between months

### Commands

The application supports vim-style commands accessed by pressing `:`:

- `add <name> <type>` - Create new habit
  - Types: `bit`, `count`, `float`
  - Example: `add "pages read" count`
- `delete <name>` - Remove habit
- `track-up <habit> [value]` - Mark as complete or increment value
- `track-down <habit>` - Mark as incomplete or decrement value
- `next-month` - Navigate to next month
- `prev-month` - Navigate to previous month
- `quit` - Exit application

### Habit Types

1. **Bit** - Boolean completion (done/not done)
2. **Count** - Integer values (e.g., pages read, exercises completed)
3. **Float** - Decimal values (e.g., water consumption in liters, hours studied)

## Development

This project is an experimental application built entirely using the Claude Opus model. The development process focused on architectural decisions and design patterns, with the AI model handling the complete implementation. The project demonstrates modern Go development practices and clean architecture principles.

**Development Conversation:** The complete conversation that generated this application can be viewed at: https://claude.ai/share/ac346bb4-b57e-45af-b7cb-53dcdd35f54e

### Current Status

The application is currently in development with a functional terminal interface and basic habit tracking capabilities. While the core UI components are implemented, several features are still being refined and data persistence is not yet complete.

### Development Roadmap

**Immediate Priorities:**
- Complete database persistence implementation
- Fix command execution bugs
- Improve error handling and validation
- Complete habit management functionality

**Future Enhancements:**
- CLI interface mode
- Data import/export functionality
- Configuration management
- Additional view modes (heatmap, yearly views)
- Statistics and analytics
- Habit streaks and achievements

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Modern TUI framework for Go
- [Charmbracelet ecosystem](https://github.com/charmbracelet) - UI components and styling
- Go 1.25.0+ - Modern Go features and standard library

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]