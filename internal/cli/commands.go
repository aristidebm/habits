package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"example.com/habits/internal/app"
	"github.com/spf13/cobra"
)

// newAddCmd creates the add command
func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <name> <type> [goal]",
		Short: "Add a new habit",
		Long:  `Add a new habit with the specified name and type. Types: bit, count, float. Goal is optional for count/float types.`,
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			dbPath, _ := cmd.Flags().GetString("db")

			application, err := app.NewApp(dbPath)
			if err != nil {
				return err
			}
			defer application.Close()

			if err := application.Migrate(); err != nil {
				return err
			}

			name := args[0]
			typeStr := args[1]

			var habitType app.HabitType
			switch typeStr {
			case "bit":
				habitType = app.HabitTypeBit
			case "count":
				habitType = app.HabitTypeCount
			case "float":
				habitType = app.HabitTypeFloat
			default:
				return fmt.Errorf("invalid habit type: %s (use: bit, count, float)", typeStr)
			}

			// Parse optional goal
			var goal float64
			if len(args) >= 3 {
				goalStr := args[2]
				var err error
				goal, err = strconv.ParseFloat(goalStr, 64)
				if err != nil {
					return fmt.Errorf("invalid goal value: %s", goalStr)
				}
			}

			habit, err := application.CreateHabit(context.Background(), name, habitType, goal)
			if err != nil {
				return fmt.Errorf("failed to create habit: %w", err)
			}

			fmt.Printf("Created habit: %s (ID: %d)\n", habit.Name, habit.ID)
			return nil
		},
	}

	return cmd
}

// newRenameCmd creates the rename command
func newRenameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename <old_name> <new_name>",
		Short: "Rename a habit",
		Long:  `Rename an existing habit to a new name.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dbPath, _ := cmd.Flags().GetString("db")

			application, err := app.NewApp(dbPath)
			if err != nil {
				return err
			}
			defer application.Close()

			if err := application.Migrate(); err != nil {
				return err
			}

			oldName := args[0]
			newName := args[1]

			habit, err := application.GetHabitByName(context.Background(), oldName)
			if err != nil {
				return fmt.Errorf("habit '%s' not found: %w", oldName, err)
			}

			if err := application.UpdateHabit(context.Background(), habit.ID, newName, habit.Type, habit.Goal); err != nil {
				return fmt.Errorf("failed to rename habit: %w", err)
			}

			fmt.Printf("Renamed habit: %s -> %s\n", oldName, newName)
			return nil
		},
	}

	return cmd
}

// newDeleteCmd creates the delete command
func newDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a habit",
		Long:  `Delete the habit with the specified name.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dbPath, _ := cmd.Flags().GetString("db")

			application, err := app.NewApp(dbPath)
			if err != nil {
				return err
			}
			defer application.Close()

			name := args[0]

			habit, err := application.GetHabitByName(context.Background(), name)
			if err != nil {
				return fmt.Errorf("habit '%s' not found: %w", name, err)
			}

			if err := application.DeleteHabit(context.Background(), habit.ID); err != nil {
				return fmt.Errorf("failed to delete habit: %w", err)
			}

			fmt.Printf("Deleted habit: %s\n", name)
			return nil
		},
	}

	return cmd
}

// newEditCmd creates the edit command
func newEditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit the configuration file",
		Long:  `Open the configuration file in your preferred editor.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := app.ConfigPath()

			// Ensure config exists
			if _, err := app.LoadConfig(""); err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Get editor from environment
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "nano" // fallback editor
			}

			// Launch editor
			editCmd := exec.Command(editor, configPath)
			editCmd.Stdin = os.Stdin
			editCmd.Stdout = os.Stdout
			editCmd.Stderr = os.Stderr

			if err := editCmd.Run(); err != nil {
				return fmt.Errorf("failed to launch editor: %w", err)
			}

			fmt.Printf("Configuration saved to %s\n", configPath)
			return nil
		},
	}

	return cmd
}

// newTrackUpCmd creates the track-up command
func newTrackUpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "track-up <habit_name>",
		Short: "Mark habit as done or increment value",
		Long:  `Mark a habit as done or increment its value. Date defaults to today.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dbPath, _ := cmd.Flags().GetString("db")

			application, err := app.NewApp(dbPath)
			if err != nil {
				return err
			}
			defer application.Close()

			if err := application.Migrate(); err != nil {
				return err
			}

			habitName := args[0]
			valueStr, _ := cmd.Flags().GetString("value")
			dateStr, _ := cmd.Flags().GetString("date")

			// Default values
			if valueStr == "" {
				valueStr = "1"
			}
			if dateStr == "" {
				dateStr = time.Now().Format("2006-01-02")
			} else {
				if _, err := time.Parse("2006-01-02", dateStr); err != nil {
					return fmt.Errorf("invalid date format (use YYYY-MM-DD): %w", err)
				}
			}

			habit, err := application.GetHabitByName(context.Background(), habitName)
			if err != nil {
				return fmt.Errorf("habit '%s' not found: %w", habitName, err)
			}

			var value float64
			switch habit.Type {
			case app.HabitTypeBit:
				value = 1
			case app.HabitTypeCount, app.HabitTypeFloat:
				value, err = strconv.ParseFloat(valueStr, 64)
				if err != nil {
					return fmt.Errorf("invalid value: %w", err)
				}
			}

			date, _ := time.Parse("2006-01-02", dateStr)
			if err := application.UpsertEntry(context.Background(), habit.ID, date, value); err != nil {
				return fmt.Errorf("failed to track habit: %w", err)
			}

			fmt.Printf("Tracked up: %s on %s\n", habitName, dateStr)
			return nil
		},
	}

	cmd.Flags().String("value", "", "Value to track (defaults to 1)")
	cmd.Flags().String("date", "", "Date to track (YYYY-MM-DD, defaults to today)")

	return cmd
}

// newTrackDownCmd creates the track-down command
func newTrackDownCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "track-down <habit_name>",
		Short: "Mark habit as not done or decrement value",
		Long:  `Mark a habit as not done or decrement its value. Date defaults to today.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dbPath, _ := cmd.Flags().GetString("db")

			application, err := app.NewApp(dbPath)
			if err != nil {
				return err
			}
			defer application.Close()

			if err := application.Migrate(); err != nil {
				return err
			}

			habitName := args[0]
			dateStr, _ := cmd.Flags().GetString("date")

			// Default value
			if dateStr == "" {
				dateStr = time.Now().Format("2006-01-02")
			} else {
				if _, err := time.Parse("2006-01-02", dateStr); err != nil {
					return fmt.Errorf("invalid date format (use YYYY-MM-DD): %w", err)
				}
			}

			habit, err := application.GetHabitByName(context.Background(), habitName)
			if err != nil {
				return fmt.Errorf("habit '%s' not found: %w", habitName, err)
			}

			var value float64
			switch habit.Type {
			case app.HabitTypeBit:
				value = 0
			case app.HabitTypeCount, app.HabitTypeFloat:
				value = 0
			}

			date, _ := time.Parse("2006-01-02", dateStr)
			if err := application.UpsertEntry(context.Background(), habit.ID, date, value); err != nil {
				return fmt.Errorf("failed to track habit: %w", err)
			}

			fmt.Printf("Tracked down: %s on %s\n", habitName, dateStr)
			return nil
		},
	}

	cmd.Flags().String("date", "", "Date to track (YYYY-MM-DD, defaults to today)")

	return cmd
}

// newExportCmd creates the export command
func newExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export <path>",
		Short: "Export all habits and entries to JSON file",
		Long:  `Export all habits and their entries to a JSON file.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dbPath, _ := cmd.Flags().GetString("db")

			application, err := app.NewApp(dbPath)
			if err != nil {
				return err
			}
			defer application.Close()

			if err := application.Migrate(); err != nil {
				return err
			}

			path := args[0]

			exportHabits, err := application.Export(context.Background())
			if err != nil {
				return fmt.Errorf("failed to export data: %w", err)
			}

			// Write to JSON file (reuse TUI logic)
			file, err := os.Create(path)
			if err != nil {
				return fmt.Errorf("error creating file: %w", err)
			}
			defer file.Close()

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(exportHabits); err != nil {
				return fmt.Errorf("error writing JSON: %w", err)
			}

			fmt.Printf("Exported %d habits to %s\n", len(exportHabits), path)
			return nil
		},
	}

	return cmd
}

// newWriteCmd creates the write command
func newWriteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "write",
		Short: "Write all pending habits to database",
		Long:  `This command is primarily used by the TUI. Use add/delete commands instead.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("The write command is primarily used by the TUI. Use add/delete commands instead.")
			return nil
		},
	}

	return cmd
}

// newNextMonthCmd creates the next-month command
func newNextMonthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "next-month",
		Short: "Navigate to next month (TUI only)",
		Long:  `This command is only available in the TUI interface.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("The next-month command is only available in the TUI interface.")
			fmt.Println("Run 'habits tui' to launch the interactive interface.")
			return nil
		},
	}

	return cmd
}

// newPrevMonthCmd creates the prev-month command
func newPrevMonthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prev-month",
		Short: "Navigate to previous month (TUI only)",
		Long:  `This command is only available in the TUI interface.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("The prev-month command is only available in the TUI interface.")
			fmt.Println("Run 'habits tui' to launch the interactive interface.")
			return nil
		},
	}

	return cmd
}
