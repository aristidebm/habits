package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"example.com/habits/internal/app"
	"example.com/habits/internal/tui"
	"github.com/spf13/cobra"
)

// runTUI launches the terminal user interface
func runTUI(cmd *cobra.Command, args []string) error {
	f, err := tea.LogToFile("/tmp/debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	// Get database path from flag or default
	dbPath, _ := cmd.Flags().GetString("db")

	// Initialize app
	application, err := app.NewApp(dbPath)
	if err != nil {
		return err
	}
	defer application.Close()

	// Migrate database
	if err := application.Migrate(); err != nil {
		return err
	}

	// Create and run TUI
	program := tui.NewProgram(application)
	p := tea.NewProgram(program, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}

	return nil
}

// newTUICmd creates the tui command
func newTUICmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "tui",
		Short: "Launch the terminal user interface",
		Long:  `Launch the interactive terminal user interface for managing habits.`,
		RunE:  runTUI,
	}

	return cmd
}
