package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"example.com/habits/internal/app"
	"example.com/habits/internal/tui"
)

func main() {
	// Initialize the core application
	application := app.NewApp()

	// Add some sample habits for testing
	application.AddHabit("Morning Run", app.HabitTypeBit)
	application.AddHabit("Read Pages", app.HabitTypeCount)
	application.AddHabit("Water (L)", app.HabitTypeFloat)
	application.AddHabit("Meditation", app.HabitTypeBit)

	// Create and run TUI program
	program := tui.NewProgram(application)
	p := tea.NewProgram(program, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
