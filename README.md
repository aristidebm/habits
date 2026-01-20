# Habits

A terminal-based habit tracking application with calendar views, multiple habit types, and note-taking capabilities.

## Features

- **Calendar Views**: Weekly and monthly calendar interfaces
- **Habit Types**: Boolean (done/not done), count (integers), and float (decimals)
- **Note System**: Attach detailed notes to any habit entry
- **Themes**: Built-in themes (Default, Dark, Light, Monochrome, Rosé Pine)
- **Commands**: Vim-style command system for habit management
- **Persistence**: SQLite database with automatic migrations

## Installation

### Prerequisites
- Go 1.25.0 or later
- SQLite3

### Build from Source
```bash
git clone <repository-url>
cd habits
go mod tidy
go build ./...
```

### Run
```bash
./habits
```

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

### Navigation
- `Tab` - Switch between weekly/monthly views
- `HJKL` / Arrow keys - Navigate dates and habits
- `:` - Open command mode
- `q` - Quit

### Commands
- `add <name> <type>` - Create habit (types: bit, count, float)
- `delete <name>` - Remove habit
- `track-up <habit>` - Increment habit value
- `track-down <habit>` - Decrement habit value
- `next-month` / `prev-month` - Navigate months

### Habit Types
- **bit**: Boolean completion (done/not done)
- **count**: Integer values (e.g., pages read)
- **float**: Decimal values (e.g., hours exercised)

### Notes
- `e` - Edit notes for selected habit/date
- `Space` - Toggle completion for bit habits
- `+/-` - Adjust values for count/float habits

## Configuration

Configuration files are created automatically in:
- `~/.config/habits/config.toml` - Application settings
- `~/.config/habits/themes/` - Custom themes

Built-in themes: Default, Dark, Light, Monochrome, Rosé Pine

## Development

```bash
go test ./...    # Run tests
go build ./...   # Build application
go run ./...     # Run application
```

## Dependencies

- Go 1.25.0+
- SQLite3

## License

MIT