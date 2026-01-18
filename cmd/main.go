package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"example.com/habits/internal/app"
	"example.com/habits/internal/tui"
)

func main() {
	dbPath := "./habits.db"

	// Initialize the core application with database
	application, err := app.NewApp(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer application.Close()

	// Run database migrations
	if err := application.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if os.Getenv("DEBUG") == "1" {
		f, err := tea.LogToFile("/tmp/habits-debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	// Create and run TUI program
	program := tui.NewProgram(application)
	p := tea.NewProgram(program, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
