package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "habits",
		Short: "A terminal-based habit tracking application",
		Long:  `Habits is a terminal-based habit tracking application built with Bubble Tea TUI framework.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// When no subcommand is provided, run TUI by default
			return runTUI(cmd, args)
		},
	}

	// Add persistent flags
	cmd.PersistentFlags().String("db", "./habits.db", "Database file path")

	// Add subcommands
	cmd.AddCommand(newTUICmd())
	cmd.AddCommand(newAddCmd())
	cmd.AddCommand(newDeleteCmd())
	cmd.AddCommand(newTrackUpCmd())
	cmd.AddCommand(newTrackDownCmd())
	cmd.AddCommand(newNextMonthCmd())
	cmd.AddCommand(newPrevMonthCmd())
	cmd.AddCommand(newWriteCmd())
	cmd.AddCommand(newExportCmd())

	return cmd
}

// Execute runs the CLI
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
