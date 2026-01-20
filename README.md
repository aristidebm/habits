# Habits

A terminal-based habit tracking application written in Go, featuring a modern text-based user interface for managing and tracking daily habits.

## Overview

Habits is a fully functional terminal-based habit tracking application written in Go. It provides an intuitive calendar-based interface for tracking personal habits with support for boolean, integer, and decimal value tracking. The application features a modern text-based user interface built with the Bubble Tea framework, comprehensive theming support, and a powerful command system.

Perfect for tracking various habit types from simple completion tasks to quantitative measurements like reading pages, exercise reps, or water consumption.

## Features

### Core Features
- **Calendar-based tracking** - Visual calendar interface showing habit completion over time
- **Multiple habit types** - Support for boolean (bit), integer (count), and decimal (float) values
- **Flexible view modes** - Weekly and monthly calendar views with consistent theming
- **Vim-style command system** - Efficient keyboard-driven interface with command history
- **Real-time updates** - Immediate visual feedback for habit tracking
- **Date navigation** - Easy movement between days, months, and specific date selection

### Advanced Features
- **Note-taking system** - Attach detailed notes to any habit entry with full text editing
- **Customizable themes** - Built-in themes (Default, Dark, Light, Monochrome, Rosé Pine) with support for custom theme files
- **Configurable symbols** - Customize completion, missed, and untracked symbols per view mode
- **Database persistence** - SQLite-based storage with automatic migrations
- **Comprehensive logging** - Debug logging for troubleshooting and development
- **Unit test coverage** - Automated tests for core functionality

### Habit Management
- **Habit creation** - Add habits with custom names, types, and goals
- **Habit deletion** - Remove habits and associated data
- **Value tracking** - Increment/decrement values for count and float habits
- **Completion tracking** - Mark habits as complete/incomplete for bit habits
- **Goal setting** - Set target values for quantitative habits

### User Interface
- **Keyboard navigation** - Full keyboard control with vim-style bindings
- **Visual feedback** - Color-coded calendar cells with theme support
- **Command line interface** - Vim-style command system with tab completion
- **Status indicators** - Database path and selected date display
- **Responsive design** - Adapts to terminal width and height

### Technical Features
- **Clean architecture** - Well-organized codebase with clear separation of concerns
- **Error handling** - Comprehensive error handling with user-friendly messages
- **Configuration management** - TOML-based configuration with defaults and user overrides
- **Migration system** - Automatic database schema migrations with Goose
- **Cross-platform** - Works on Linux, macOS, and Windows

## Architecture

The application follows clean architecture principles with clear separation of concerns:

- `cmd/main.go` - Application entry point
- `internal/app/` - Core business logic and data models
- `internal/tui/` - Terminal user interface components
- `internal/cli/` - Command-line interface (prepared for future CLI mode)

## Installation

### Prerequisites

- Go 1.25.0 or later
- SQLite3 (usually pre-installed on most systems)

### Building from Source

```bash
git clone <repository-url>
cd habits
go mod tidy
go build ./...
```

### Running

```bash
./habits
# or
go run ./...
```

### Configuration

The application creates configuration files automatically on first run:
- Database: `~/.local/state/habits/habits.db`
- Config: `~/.config/habits/config.toml`
- Themes: `~/.config/habits/themes/`

## Usage

### Interface Navigation

The application uses a keyboard-driven interface with the following controls:

#### Global Navigation
- `:` - Open command line
- `q` - Quick quit
- `Tab` - Switch between weekly and monthly views

#### Calendar Navigation
- Arrow keys / `HJKL` - Navigate between dates and habits
- `j`/`k` - Move between habits (up/down)
- `h`/`l` - Move between dates (left/right)
- `H`/`L` - Jump to previous/next month

#### Habit Management
- `e` - Edit note for selected habit and date
- `Space` - Toggle completion for bit habits
- `+`/`-` - Increment/decrement values for count/float habits

### Commands

The application supports vim-style commands accessed by pressing `:`:

#### Habit Management
- `add <name> <type> [goal]` - Create new habit
  - Types: `bit`, `count`, `float`
  - Goal: Optional target value for count/float habits
  - Examples:
    - `add "Drink water" bit`
    - `add "Read pages" count 50`
    - `add "Exercise time" float 1.5`
- `delete <name>` - Remove habit and all associated data
- `track-up <habit> [value]` - Mark as complete or increment value
  - Value: Optional increment amount (default: 1 for count, 0.1 for float)
- `track-down <habit> [value]` - Mark as incomplete or decrement value

#### Navigation
- `next-month` - Navigate to next month
- `prev-month` - Navigate to previous month
- `today` - Jump to today's date

#### System
- `quit` - Exit application
- `help` - Show available commands

### Habit Types

1. **Bit** - Boolean completion (done/not done)
2. **Count** - Integer values (e.g., pages read, exercises completed)
3. **Float** - Decimal values (e.g., water consumption in liters, hours studied)

### Configuration

The application can be configured via `~/.config/habits/config.toml`:

```toml
[db]
  path = "/home/user/.local/state/habits/habits.db"

[views]
  [views.weekly]
    missed = "⛌"
    completed = "🗸"
    untracked = "○"
  [views.monthly]
    missed = "⛌"
    completed = "🗸"
    untracked = "○"

[theme]
  name = "rose-pine"
```

### Theming

Habits supports multiple built-in themes and custom theme files in `~/.config/habits/themes/`:

- **Default** - Classic color scheme
- **Dark** - Dark theme for low-light environments
- **Light** - Light theme for bright environments
- **Monochrome** - True monochrome with white and black
- **Rosé Pine** - Elegant purple-based theme

Custom themes can be created by adding `.toml` files to the themes directory with color definitions for all UI elements.

## Development

This project is a fully functional habit tracking application built using modern Go development practices and clean architecture principles. The application demonstrates comprehensive TUI development, database integration, and user experience design.

### Current Status

The application is fully functional with complete habit tracking capabilities:

- ✅ Database persistence with SQLite and automatic migrations
- ✅ Full command system with validation and error handling
- ✅ Note-taking system with text editing
- ✅ Theme system with multiple built-in themes
- ✅ Comprehensive logging for debugging
- ✅ Unit test coverage for core functionality
- ✅ Cross-platform compatibility

### Architecture Details

```
habits/
├── cmd/main.go              # Application entry point
├── internal/
│   ├── app/                 # Core business logic
│   │   ├── models.go        # Data structures
│   │   ├── store.go         # Database operations
│   │   ├── theme.go         # Theme management
│   │   ├── config.go        # Configuration handling
│   │   ├── styles.go        # UI styling
│   │   └── db.go            # Database initialization
│   └── tui/                 # Terminal user interface
│       ├── program.go       # Main TUI program
│       ├── calendar/        # Calendar views
│       │   ├── calendar.go  # Core calendar logic
│       │   ├── weekly_view.go
│       │   └── monthly_view.go
│       └── command/         # Command system
├── migrations/              # Database schema migrations
└── themes/                  # Custom theme files
```

### Key Technologies

- **Go 1.25.0+** - Modern Go with latest features
- **Bubble Tea** - Elegant TUI framework
- **SQLite** - Embedded database with migrations
- **Lipgloss** - Terminal styling and layout
- **Goose** - Database migration tool

### Testing

Run the test suite:
```bash
go test ./...
```

### Development Commands

```bash
# Build and run
go run ./...

# Build only
go build ./...

# Run tests
go test ./...

# Format code
go fmt ./...

# Vet code
go vet ./...

# Database migrations
goose -dir migrations sqlite3 ./habits.db up
```

## Dependencies

### Core Dependencies
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Modern TUI framework for Go
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling and layout
- [Bubbles](https://github.com/charmbracelet/bubbles) - Interactive TUI components
- [SQLite3 driver](https://github.com/mattn/go-sqlite3) - SQLite database connectivity
- [Goose](https://github.com/pressly/goose) - Database migration tool
- [TOML parser](https://github.com/BurntSushi/toml) - Configuration file parsing

### Development Dependencies
- Go 1.25.0+ - Modern Go features and standard library

## License

This project is open source. See LICENSE file for details.

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

### Development Setup

```bash
git clone <your-fork-url>
cd habits
go mod tidy
go test ./...
go run ./...
```

### Code Style

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add documentation for exported functions
- Include unit tests for new features
- Use structured logging with appropriate levels