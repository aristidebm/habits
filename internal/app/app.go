package app

import (
	"fmt"
	"time"

	_ "github.com/pressly/goose/v3"
)

type App struct {
	*Store
}

func NewApp(dbPath string) (*App, error) {
	db, err := OpenDB(dbPath)
	if err != nil {
		return nil, err
	}

	return &App{
		Store: &Store{db: db},
	}, nil
}

func (a *App) Migrate() error {
	return Migrate(a.Store.db, "migrations")
}

func (a *App) Export() ([]ExportHabit, error) {
	habits, err := a.ListHabits()
	if err != nil {
		return nil, fmt.Errorf("failed to list habits: %w", err)
	}

	exportHabits := []ExportHabit{}
	for _, habit := range habits {
		// Get all entries for this habit (no date limit)
		startDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Now().AddDate(10, 0, 0) // Far future
		entries, err := a.ListEntries(habit.ID, startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("failed to list entries for habit %s: %w", habit.Name, err)
		}

		exportEntries := []ExportEntry{}
		for _, entry := range entries {
			// Get notes for this entry
			notes, err := a.ListNotes(entry.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to list notes for entry %d: %w", entry.ID, err)
			}

			exportNotes := []ExportNote{}
			for _, note := range notes {
				exportNotes = append(exportNotes, ExportNote{
					Note: note.Note,
				})
			}

			exportEntries = append(exportEntries, ExportEntry{
				ID:    entry.ID,
				Date:  entry.Date.Format("2006-01-02"),
				Value: entry.Value,
				Notes: exportNotes,
			})
		}

		// Ensure entries is never nil
		if exportEntries == nil {
			exportEntries = []ExportEntry{}
		}

		exportHabits = append(exportHabits, ExportHabit{
			ID:      habit.ID,
			Name:    habit.Name,
			Goal:    habit.Goal,
			Kind:    habit.Type.String(),
			Entries: exportEntries,
		})
	}

	return exportHabits, nil
}

func (a *App) Close() error {
	return a.Store.Close()
}
