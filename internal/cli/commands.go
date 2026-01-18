package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"example.com/habits/internal/app"
	"github.com/spf13/cobra"
)

// newAddCmd creates the add command
func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <name> <type>",
		Short: "Add a new habit",
		Long:  `Add a new habit with the specified name and type. Types: bit, count, float.`,
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

			habit, err := application.CreateHabit(name, habitType, 0)
			if err != nil {
				return fmt.Errorf("failed to create habit: %w", err)
			}

			fmt.Printf("Created habit: %s (ID: %d)\n", habit.Name, habit.ID)
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

			if err := application.Migrate(); err != nil {
				return err
			}

			name := args[0]

			habit, err := application.GetHabitByName(name)
			if err != nil {
				return fmt.Errorf("habit '%s' not found: %w", name, err)
			}

			if err := application.DeleteHabit(habit.ID); err != nil {
				return fmt.Errorf("failed to delete habit: %w", err)
			}

			fmt.Printf("Deleted habit: %s\n", name)
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

			habit, err := application.GetHabitByName(habitName)
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
			if err := application.UpsertEntry(habit.ID, date, value); err != nil {
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

			habit, err := application.GetHabitByName(habitName)
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
			if err := application.UpsertEntry(habit.ID, date, value); err != nil {
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

			exportHabits, err := application.Export()
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
