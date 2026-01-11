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
	application.AddHabit("One", app.HabitTypeBit)
	application.AddHabit("Two", app.HabitTypeBit)
	application.AddHabit("Three", app.HabitTypeBit)
	application.AddHabit("Four", app.HabitTypeBit)
	application.AddHabit("Five", app.HabitTypeBit)
	application.AddHabit("Six", app.HabitTypeBit)
	application.AddHabit("Seven", app.HabitTypeBit)
	application.AddHabit("Eight", app.HabitTypeBit)
	application.AddHabit("Nine", app.HabitTypeBit)
	application.AddHabit("Ten", app.HabitTypeBit)
	application.AddHabit("Eleven", app.HabitTypeBit)
	application.AddHabit("Twelve", app.HabitTypeBit)
	application.AddHabit("Thirteen", app.HabitTypeBit)
	application.AddHabit("Fourteen", app.HabitTypeBit)
	application.AddHabit("Fifteen", app.HabitTypeBit)
	application.AddHabit("Sixten", app.HabitTypeBit)
	application.AddHabit("SevenTeen", app.HabitTypeBit)
	application.AddHabit("Eighteen", app.HabitTypeBit)
	application.AddHabit("NineTeen", app.HabitTypeBit)
	application.AddHabit("Twenty", app.HabitTypeBit)

	// Create and run TUI program
	program := tui.NewProgram(application)
	p := tea.NewProgram(program, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
